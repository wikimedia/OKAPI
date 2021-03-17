package pagemove

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"okapi-data-service/queues/pagedelete"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagemoveTestExpire = time.Hour * 24
const pagemoveTestQueueName = "queue/pagedelete"
const pagemoveTestName = "stream/pagemove"
const pagemoveTestTitle = "ninja"
const pagemoveTestDbName = "ninjas"

type pagemoveRedisMock struct {
	mock.Mock
	redis.Client
}

func (r *pagemoveRedisMock) RPush(_ context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := r.Called(key, values[0])
	cmd := new(redis.IntCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (r *pagemoveRedisMock) Set(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := r.Called(key, value, expiration)
	cmd := new(redis.StatusCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func TestPagemove(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	assert := assert.New(t)
	ctx := context.Background()

	date := time.Now().Add(24 * time.Hour)
	evt := new(eventstream.PageMove)
	evt.Data.PageTitle = pagemoveTestTitle
	evt.Data.Database = pagemoveTestDbName
	evt.Data.Meta.Dt = date

	data, err := json.Marshal(&pagedelete.Data{
		Title:  pagemoveTestTitle,
		DbName: pagemoveTestDbName,
	})
	assert.NoError(err)

	t.Run("pagemove success", func(t *testing.T) {
		cmdable := new(pagemoveRedisMock)
		cmdable.On("RPush", pagemoveTestQueueName, data).Return(nil)
		cmdable.On("Set", pagemoveTestName, date, pagemoveTestExpire).Return(nil)

		Handler(ctx, cmdable, pagemoveTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", pagemoveTestQueueName, data)
		cmdable.AssertCalled(t, "Set", pagemoveTestName, date, pagemoveTestExpire)
	})

	t.Run("pagemove push error", func(t *testing.T) {
		cmdable := new(pagemoveRedisMock)
		cmdable.On("RPush", pagemoveTestQueueName, data).Return(errors.New("redis not available"))
		cmdable.On("Set", pagemoveTestName, date, pagemoveTestExpire).Return(nil)

		Handler(ctx, cmdable, pagemoveTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", pagemoveTestQueueName, data)
		cmdable.AssertNotCalled(t, "Set", pagemoveTestName, date, pagemoveTestExpire)
	})

	t.Run("pagemove set error", func(t *testing.T) {
		cmdable := new(pagemoveRedisMock)
		cmdable.On("RPush", pagemoveTestQueueName, data).Return(nil)
		cmdable.On("Set", pagemoveTestName, date, pagemoveTestExpire).Return(errors.New("offline"))

		Handler(ctx, cmdable, pagemoveTestExpire)(evt)
		cmdable.AssertCalled(t, "RPush", pagemoveTestQueueName, data)
		cmdable.AssertCalled(t, "Set", pagemoveTestName, date, pagemoveTestExpire)
	})
}
