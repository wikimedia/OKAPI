package revisionscore

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"okapi-data-service/queues/pagefetch"
	"okapi-data-service/schema/v3"
	"okapi-data-service/streams/utils"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	ores "github.com/protsack-stephan/mediawiki-ores-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const revisionscoreTestExpire = time.Hour * 24
const revisionscoreTestQueueName = "queue/pagefetch"
const revisionscoreTestName = "stream/revisionscore"
const revisionscoreTestTitle = "ninja"
const revisionscoreTestDbName = "enwiki"
const revisionscoreTestDamagingPrediction = false
const revisionscoreTestDamagingTrue = 0.2
const revisionscoreTestDamagingFalse = 0.8
const revisionscoreTestGoodFaithPrediction = true
const revisionscoreTestGoodFaithTrue = 0.8
const revisionscoreTestGoodFaithFalse = 0.2
const revisionscoreSiteURL = "en.wikipedia.org"
const revisionscoreLang = "en"
const revisionscoreRevID = 112
const revisionscoreNamespace = 6
const revisionscoreTestUserID = 10
const revisionscoreTestUserText = "unknown"
const revisionscoreTestUserEditCount = 100
const revisionscoreTestUserIsBot = false

var revisionscoreTestUserRegistrationDt = time.Now()
var revisionscoreTestUserGroups = []string{"bot", "admin"}

type revisionscoreRedisMock struct {
	mock.Mock
	redis.Client
}

func (r *revisionscoreRedisMock) RPush(_ context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := r.Called(key, values[0])
	cmd := new(redis.IntCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (r *revisionscoreRedisMock) Set(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := r.Called(key, value, expiration)
	cmd := new(redis.StatusCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func TestRevisionscore(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	assert := assert.New(t)
	ctx := context.Background()

	date := time.Now().Add(24 * time.Hour)
	evt := new(eventstream.RevisionScore)
	evt.Data.PageTitle = revisionscoreTestTitle
	evt.Data.PageNamespace = revisionscoreNamespace
	evt.Data.RevID = revisionscoreRevID
	evt.Data.Database = revisionscoreTestDbName
	evt.Data.Meta.Dt = date
	evt.Data.Meta.Domain = revisionscoreSiteURL
	evt.Data.Scores.Damaging.Probability.True = revisionscoreTestDamagingTrue
	evt.Data.Scores.Damaging.Probability.False = revisionscoreTestDamagingFalse
	evt.Data.Scores.Goodfaith.Probability.True = revisionscoreTestGoodFaithTrue
	evt.Data.Scores.Goodfaith.Probability.False = revisionscoreTestGoodFaithFalse
	evt.Data.Performer.UserID = revisionscoreTestUserID
	evt.Data.Performer.UserText = revisionscoreTestUserText
	evt.Data.Performer.UserEditCount = revisionscoreTestUserEditCount
	evt.Data.Performer.UserGroups = revisionscoreTestUserGroups
	evt.Data.Performer.UserIsBot = revisionscoreTestUserIsBot
	evt.Data.Performer.UserRegistrationDt = revisionscoreTestUserRegistrationDt

	data, err := json.Marshal(&pagefetch.Data{
		Title:     revisionscoreTestTitle,
		DbName:    revisionscoreTestDbName,
		SiteURL:   utils.SiteURL(revisionscoreSiteURL),
		Lang:      revisionscoreLang,
		Namespace: revisionscoreNamespace,
		Revision:  revisionscoreRevID,
		Editor: &schema.Editor{
			Identifier:  revisionscoreTestUserID,
			Name:        revisionscoreTestUserText,
			EditCount:   revisionscoreTestUserEditCount,
			Groups:      revisionscoreTestUserGroups,
			IsBot:       revisionscoreTestUserIsBot,
			DateStarted: &revisionscoreTestUserRegistrationDt,
		},
		Scores: &schema.Scores{
			Damaging: &ores.ScoreDamaging{
				Prediction: revisionscoreTestDamagingPrediction,
				Probability: struct {
					False float64 `json:"false"`
					True  float64 `json:"true"`
				}{
					True:  revisionscoreTestDamagingTrue,
					False: revisionscoreTestDamagingFalse,
				},
			},
			GoodFaith: &ores.ScoreGoodFaith{
				Prediction: revisionscoreTestGoodFaithPrediction,
				Probability: struct {
					False float64 `json:"false"`
					True  float64 `json:"true"`
				}{
					True:  revisionscoreTestGoodFaithTrue,
					False: revisionscoreTestGoodFaithFalse,
				},
			},
		},
	})
	assert.NoError(err)

	t.Run("revisionscore success", func(t *testing.T) {
		cmdable := new(revisionscoreRedisMock)
		cmdable.On("RPush", revisionscoreTestQueueName, data).Return(nil)
		cmdable.On("Set", revisionscoreTestName, date, revisionscoreTestExpire).Return(nil)

		Handler(ctx, cmdable, revisionscoreTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisionscoreTestQueueName, data)
		cmdable.AssertCalled(t, "Set", revisionscoreTestName, date, revisionscoreTestExpire)
	})

	t.Run("revisionscore push error", func(t *testing.T) {
		cmdable := new(revisionscoreRedisMock)
		cmdable.On("RPush", revisionscoreTestQueueName, data).Return(errors.New("redis not available"))
		cmdable.On("Set", revisionscoreTestName, date, revisionscoreTestExpire).Return(nil)

		Handler(ctx, cmdable, revisionscoreTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisionscoreTestQueueName, data)
		cmdable.AssertNotCalled(t, "Set", revisionscoreTestName, date, revisionscoreTestExpire)
	})

	t.Run("revisionscore set error", func(t *testing.T) {
		cmdable := new(revisionscoreRedisMock)
		cmdable.On("RPush", revisionscoreTestQueueName, data).Return(nil)
		cmdable.On("Set", revisionscoreTestName, date, revisionscoreTestExpire).Return(errors.New("offline"))

		Handler(ctx, cmdable, revisionscoreTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisionscoreTestQueueName, data)
		cmdable.AssertCalled(t, "Set", revisionscoreTestName, date, revisionscoreTestExpire)
	})
}
