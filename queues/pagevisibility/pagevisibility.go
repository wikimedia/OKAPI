package pagevisibility

import (
	"context"
	"okapi-data-service/pkg/worker"

	"github.com/go-redis/redis/v8"
)

// Name redis key for the queue
const Name string = "queue/pagevisibility"

// Data item of the queue
type Data struct {
	Title    string `json:"title"`
	DbName   string `json:"db_name"`
	Revision int    `json:"revision"`
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}

// Worker processing function
func Worker() worker.Worker {
	return func(ctx context.Context, data []byte) error {
		return nil
	}
}
