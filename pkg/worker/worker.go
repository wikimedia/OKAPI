package worker

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

// Worker processing unit of the queue
type Worker func(context.Context, []byte) error

// Enqueue add message to the queue
func Enqueue(ctx context.Context, name string, store redis.Cmdable, model interface{}) error {
	data, err := json.Marshal(model)

	if err != nil {
		return err
	}

	return store.RPush(ctx, name, data).Err()
}
