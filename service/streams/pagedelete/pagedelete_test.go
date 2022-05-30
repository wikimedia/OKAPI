package pagedelete

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"okapi-data-service/queues/pagedelete"
	"okapi-data-service/schema/v3"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagedeleteTestExpire = time.Hour * 24
const pagedeleteTestQueueName = "queue/pagedelete"
const pagedeleteTestName = "stream/pagedelete"
const pagedeleteTestTitle = "ninja"
const pagedeleteTestDbName = "ninjas"
const pagedeleteTestUserID = 10
const pagedeleteTestUserText = "unknown"
const pagedeleteTestUserEditCount = 100
const pagedeleteTestUserIsBot = false

var pagedeleteTestUserRegistrationDt = time.Now()
var pagedeleteTestUserGroups = []string{"bot", "admin"}

type pagedeleteRedisMock struct {
	mock.Mock
	redis.Client
}

func (r *pagedeleteRedisMock) RPush(_ context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := r.Called(key, values[0])
	cmd := new(redis.IntCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (r *pagedeleteRedisMock) Set(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := r.Called(key, value, expiration)
	cmd := new(redis.StatusCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func TestPagedelete(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	assert := assert.New(t)
	ctx := context.Background()

	date := time.Now().Add(24 * time.Hour)
	evt := new(eventstream.PageDelete)
	evt.Data.PageTitle = pagedeleteTestTitle
	evt.Data.Database = pagedeleteTestDbName
	evt.Data.Meta.Dt = date
	evt.Data.Performer.UserID = pagedeleteTestUserID
	evt.Data.Performer.UserText = pagedeleteTestUserText
	evt.Data.Performer.UserEditCount = pagedeleteTestUserEditCount
	evt.Data.Performer.UserGroups = pagedeleteTestUserGroups
	evt.Data.Performer.UserIsBot = pagedeleteTestUserIsBot
	evt.Data.Performer.UserRegistrationDt = pagedeleteTestUserRegistrationDt

	data, err := json.Marshal(&pagedelete.Data{
		Title:  pagedeleteTestTitle,
		DbName: pagedeleteTestDbName,
		Editor: &schema.Editor{
			Identifier:  pagedeleteTestUserID,
			Name:        pagedeleteTestUserText,
			EditCount:   pagedeleteTestUserEditCount,
			Groups:      pagedeleteTestUserGroups,
			IsBot:       pagedeleteTestUserIsBot,
			DateStarted: &pagedeleteTestUserRegistrationDt,
		},
	})
	assert.NoError(err)

	t.Run("pagedelete success", func(t *testing.T) {
		cmdable := new(pagedeleteRedisMock)
		cmdable.On("RPush", pagedeleteTestQueueName, data).Return(nil)
		cmdable.On("Set", pagedeleteTestName, date, pagedeleteTestExpire).Return(nil)

		Handler(ctx, cmdable, pagedeleteTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", pagedeleteTestQueueName, data)
		cmdable.AssertCalled(t, "Set", pagedeleteTestName, date, pagedeleteTestExpire)
	})

	t.Run("pagedelete push error", func(t *testing.T) {
		cmdable := new(pagedeleteRedisMock)
		cmdable.On("RPush", pagedeleteTestQueueName, data).Return(errors.New("redis not available"))
		cmdable.On("Set", pagedeleteTestName, date, pagedeleteTestExpire).Return(nil)

		Handler(ctx, cmdable, pagedeleteTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", pagedeleteTestQueueName, data)
		cmdable.AssertNotCalled(t, "Set", pagedeleteTestName, date, pagedeleteTestExpire)
	})

	t.Run("pagedelete set error", func(t *testing.T) {
		cmdable := new(pagedeleteRedisMock)
		cmdable.On("RPush", pagedeleteTestQueueName, data).Return(nil)
		cmdable.On("Set", pagedeleteTestName, date, pagedeleteTestExpire).Return(errors.New("offline"))

		Handler(ctx, cmdable, pagedeleteTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", pagedeleteTestQueueName, data)
		cmdable.AssertCalled(t, "Set", pagedeleteTestName, date, pagedeleteTestExpire)
	})
}
