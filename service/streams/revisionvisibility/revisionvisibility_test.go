package revisionvisibility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"okapi-data-service/queues/pagevisibility"
	"okapi-data-service/schema/v3"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const revisionvisibilityTestExpire = time.Hour * 24
const revisionvisibilityTestQueueName = "queue/pagevisibility"
const revisionvisibilityTestName = "stream/revisionvisibility"
const revisionvisibilityTestTitle = "ninja"
const revisionvisibilityTestDbName = "ninjas"
const revisionvisibilityTestRev = 1
const revisionvisibilityTestSiteURL = "en.wikipedia.org"
const revisionvisibilityTestLang = "en"
const revisionvisibilityTestTextVisible = true
const revisionvisibilityTestCommentVisible = true
const revisionvisibilityTestUserVisible = true
const revisionvisibilityTestPageID = 100
const revisionvisibilityTestUserID = 10
const revisionvisibilityTestUserText = "unknown"
const revisionvisibilityTestUserEditCount = 100
const revisionvisibilityTestUserIsBot = false

var revisionvisibilityTestUserRegistrationDt = time.Now()
var revisionvisibilityTestUserGroups = []string{"bot", "admin"}
var revisionvisibilityRevDt = time.Now()

type revisionvisibilityRedisMock struct {
	mock.Mock
	redis.Client
}

func (r *revisionvisibilityRedisMock) RPush(_ context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := r.Called(key, values[0])
	cmd := new(redis.IntCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (r *revisionvisibilityRedisMock) Set(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := r.Called(key, value, expiration)
	cmd := new(redis.StatusCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func TestRevisionvisibility(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	assert := assert.New(t)
	ctx := context.Background()

	date := time.Now().Add(24 * time.Hour)
	evt := new(eventstream.RevisionVisibilityChange)
	evt.Data.PageID = revisionvisibilityTestPageID
	evt.Data.PageTitle = revisionvisibilityTestTitle
	evt.Data.Database = revisionvisibilityTestDbName
	evt.Data.RevID = revisionvisibilityTestRev
	evt.Data.Meta.Dt = date
	evt.Data.Meta.Domain = revisionvisibilityTestSiteURL
	evt.Data.Visibility.Text = revisionvisibilityTestTextVisible
	evt.Data.Visibility.Comment = revisionvisibilityTestCommentVisible
	evt.Data.Visibility.User = revisionvisibilityTestUserVisible
	evt.Data.RevTimestamp = revisionvisibilityRevDt
	evt.Data.Performer.UserID = revisionvisibilityTestUserID
	evt.Data.Performer.UserText = revisionvisibilityTestUserText
	evt.Data.Performer.UserEditCount = revisionvisibilityTestUserEditCount
	evt.Data.Performer.UserGroups = revisionvisibilityTestUserGroups
	evt.Data.Performer.UserIsBot = revisionvisibilityTestUserIsBot
	evt.Data.Performer.UserRegistrationDt = revisionvisibilityTestUserRegistrationDt

	qData := new(pagevisibility.Data)
	qData.ID = revisionvisibilityTestPageID
	qData.Title = revisionvisibilityTestTitle
	qData.Revision = revisionvisibilityTestRev
	qData.DbName = revisionvisibilityTestDbName
	qData.RevisionDt = revisionvisibilityRevDt
	qData.Lang = revisionvisibilityTestLang
	qData.SiteURL = fmt.Sprintf("https://%s", revisionvisibilityTestSiteURL)
	qData.Visible = revisionvisibilityTestTextVisible
	qData.Visibility.Text = revisionvisibilityTestTextVisible
	qData.Visibility.Comment = revisionvisibilityTestCommentVisible
	qData.Visibility.User = revisionvisibilityTestUserVisible
	qData.Editor = &schema.Editor{
		Identifier:  revisionvisibilityTestUserID,
		Name:        revisionvisibilityTestUserText,
		EditCount:   revisionvisibilityTestUserEditCount,
		Groups:      revisionvisibilityTestUserGroups,
		IsBot:       revisionvisibilityTestUserIsBot,
		DateStarted: &revisionvisibilityTestUserRegistrationDt,
	}
	data, err := json.Marshal(qData)
	assert.NoError(err)

	t.Run("revisionvisibility success", func(t *testing.T) {
		cmdable := new(revisionvisibilityRedisMock)
		cmdable.On("RPush", revisionvisibilityTestQueueName, data).Return(nil)
		cmdable.On("Set", revisionvisibilityTestName, date, revisionvisibilityTestExpire).Return(nil)

		Handler(ctx, cmdable, revisionvisibilityTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisionvisibilityTestQueueName, data)
		cmdable.AssertCalled(t, "Set", revisionvisibilityTestName, date, revisionvisibilityTestExpire)
	})

	t.Run("revisionvisibility push error", func(t *testing.T) {
		cmdable := new(revisionvisibilityRedisMock)
		cmdable.On("RPush", revisionvisibilityTestQueueName, data).Return(errors.New("redis not available"))
		cmdable.On("Set", revisionvisibilityTestName, date, revisionvisibilityTestExpire).Return(nil)

		Handler(ctx, cmdable, revisionvisibilityTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisionvisibilityTestQueueName, data)
		cmdable.AssertNotCalled(t, "Set", revisionvisibilityTestName, date, revisionvisibilityTestExpire)
	})

	t.Run("revisionvisibility set error", func(t *testing.T) {
		cmdable := new(revisionvisibilityRedisMock)
		cmdable.On("RPush", revisionvisibilityTestQueueName, data).Return(nil)
		cmdable.On("Set", revisionvisibilityTestName, date, revisionvisibilityTestExpire).Return(errors.New("offline"))

		Handler(ctx, cmdable, revisionvisibilityTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisionvisibilityTestQueueName, data)
		cmdable.AssertCalled(t, "Set", revisionvisibilityTestName, date, revisionvisibilityTestExpire)
	})
}
