package revisioncreate

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

const revisioncreateTestExpire = time.Hour * 24
const revisioncreateTestQueueName = "queue/pagepull"
const revisioncreateTestName = "stream/revisioncreate"
const revisioncreateTestTitle = "ninja"
const revisioncreateTestDbName = "ninjas"
const revisioncreateSiteURL = "en.wikipedia.org"
const revisioncreateLang = "en"

type revisioncreateRedisMock struct {
	mock.Mock
	redis.Client
}

func (r *revisioncreateRedisMock) RPush(_ context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := r.Called(key, values[0])
	cmd := new(redis.IntCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (r *revisioncreateRedisMock) Set(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := r.Called(key, value, expiration)
	cmd := new(redis.StatusCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func TestRevisioncreate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	assert := assert.New(t)
	ctx := context.Background()

	date := time.Now().Add(24 * time.Hour)
	evt := new(eventstream.RevisionCreate)
	evt.Data.PageTitle = revisioncreateTestTitle
	evt.Data.Database = revisioncreateTestDbName
	evt.Data.Meta.Domain = revisioncreateSiteURL
	evt.Data.Meta.Dt = date

	data, err := json.Marshal(&pagepull.Data{
		Title:   revisioncreateTestTitle,
		DbName:  revisioncreateTestDbName,
		SiteURL: utils.SiteURL(revisioncreateSiteURL),
		Lang:    revisioncreateLang,
	})
	assert.NoError(err)

	t.Run("revisioncreate success", func(t *testing.T) {
		cmdable := new(revisioncreateRedisMock)
		cmdable.On("RPush", revisioncreateTestQueueName, data).Return(nil)
		cmdable.On("Set", revisioncreateTestName, date, revisioncreateTestExpire).Return(nil)

		Handler(ctx, cmdable, revisioncreateTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisioncreateTestQueueName, data)
		cmdable.AssertCalled(t, "Set", revisioncreateTestName, date, revisioncreateTestExpire)
	})

	t.Run("revisioncreate push error", func(t *testing.T) {
		cmdable := new(revisioncreateRedisMock)
		cmdable.On("RPush", revisioncreateTestQueueName, data).Return(errors.New("redis not available"))
		cmdable.On("Set", revisioncreateTestName, date, revisioncreateTestExpire).Return(nil)

		Handler(ctx, cmdable, revisioncreateTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisioncreateTestQueueName, data)
		cmdable.AssertNotCalled(t, "Set", revisioncreateTestName, date, revisioncreateTestExpire)
	})

	t.Run("revisioncreate set error", func(t *testing.T) {
		cmdable := new(revisioncreateRedisMock)
		cmdable.On("RPush", revisioncreateTestQueueName, data).Return(nil)
		cmdable.On("Set", revisioncreateTestName, date, revisioncreateTestExpire).Return(errors.New("offline"))

		Handler(ctx, cmdable, revisioncreateTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", revisioncreateTestQueueName, data)
		cmdable.AssertCalled(t, "Set", revisioncreateTestName, date, revisioncreateTestExpire)
	})
}
