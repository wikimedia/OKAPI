package pagedelete

import (
	"context"
	"encoding/json"
	"okapi-data-service/models"
	"okapi-data-service/pkg/topics"
	"okapi-data-service/pkg/worker"
	"okapi-data-service/schema/v1"
	"okapi-data-service/server/pages/content"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
)

// Name redis key for the queue
const Name string = "queue/pagedelete"

// Data item of the queue
type Data struct {
	Title  string `json:"title"`
	DbName string `json:"db_name"`
}

// Repo all the needed repositories to call delete
type Repo interface {
	repository.Deleter
	repository.Finder
}

// Worker processing function
func Worker(repo Repo, stores content.Deleter, producer topics.Producer) worker.Worker {
	return func(ctx context.Context, payload []byte) error {
		data := new(Data)

		if err := json.Unmarshal(payload, data); err != nil {
			return err
		}

		page := new(models.Page)
		query := func(q *orm.Query) *orm.Query {
			return q.Where("title = ? and db_name = ?", data.Title, data.DbName)
		}

		if err := repo.Find(ctx, page, query); err != nil {
			return err
		}

		if _, err := repo.Delete(ctx, page, query); err != nil {
			return err
		}

		if err := stores.Delete(ctx, page); err != nil {
			return err
		}

		msg := content.NewStructured(page)
		value, err := json.Marshal(msg)

		if err != nil {
			return err
		}

		key, err := json.Marshal(schema.PageKey{
			Title:  data.Title,
			DbName: data.DbName,
		})

		if err != nil {
			return err
		}

		return producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topics.PageDelete, Partition: 0},
			Key:            key,
			Value:          value,
		}, nil)
	}
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}
