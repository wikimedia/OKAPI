package revisionscore

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"okapi-data-service/queues/pagepull"
	"okapi-data-service/streams/utils"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const revisionscoreTestExpire = time.Hour * 24
const revisionscoreTestQueueName = "queue/pagepull"
const revisionscoreTestName = "stream/revisionscore"
const revisionscoreTestTitle = "ninja"
const revisionscoreTestDbName = "enwiki"
const revisionscoreTestTrue = 0.2
const revisionscoreTestFalse = 0.8
const revisionscoreSiteURL = "en.wikipedia.org"
const revisionscoreLang = "en"

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
	evt.Data.Database = revisionscoreTestDbName
	evt.Data.Meta.Dt = date
	evt.Data.Meta.Domain = revisionscoreSiteURL
	evt.Data.Scores.Damaging.Probability.True = revisionscoreTestTrue
	evt.Data.Scores.Damaging.Probability.False = revisionscoreTestFalse

	data, err := json.Marshal(&pagepull.Data{
		Title:   revisionscoreTestTitle,
		DbName:  revisionscoreTestDbName,
		SiteURL: utils.SiteURL(revisionscoreSiteURL),
		Lang:    revisionscoreLang,
		Models: pagepull.Models{
			Damaging: pagepull.Damaging{
				Probability: &pagepull.Probability{
					True:  revisionscoreTestTrue,
					False: revisionscoreTestFalse,
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
