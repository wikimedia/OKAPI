package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errServerOffline = errors.New("server offline")

type workerRedisMock struct {
	mock.Mock
	redis.Client
}

func (r *workerRedisMock) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := r.Called(key, values[0])
	cmd := new(redis.IntCmd)
	cmd.SetErr(args.Error(0))
	return cmd
}

func TestWorker(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	name := "test/queue"
	info := struct {
		Name string
	}{
		"ninja",
	}

	t.Run("enqueue success", func(t *testing.T) {
		data, err := json.Marshal(info)
		assert.NoError(err)

		store := new(workerRedisMock)
		store.On("RPush", string(name), data).Return(nil)

		assert.NoError(Enqueue(ctx, name, store, &info))
	})

	t.Run("enqueue error", func(t *testing.T) {
		data, err := json.Marshal(info)
		assert.NoError(err)

		store := new(workerRedisMock)
		store.On("RPush", string(name), data).Return(errServerOffline)

		assert.Equal(Enqueue(ctx, name, store, &info), errServerOffline)
	})
}
