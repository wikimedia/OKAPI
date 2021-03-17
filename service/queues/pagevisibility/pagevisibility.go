package pagevisibility

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
	"github.com/protsack-stephan/mediawiki-api-client"
)

// Name redis key for the queue
const Name string = "queue/pagevisibility"

// Data item of the queue
type Data struct {
	Title    string `json:"title"`
	DbName   string `json:"db_name"`
	Revision int    `json:"revision"`
	Visible  bool   `json:"visible"`
	Lang     string `json:"lang"`
	SiteURL  string `json:"site_url"`
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}

// Worker processing function
func Worker(repo repository.Finder, storage content.Storer, producer topics.Producer) worker.Worker {
	return func(ctx context.Context, payload []byte) error {
		data := new(Data)

		if err := json.Unmarshal(payload, data); err != nil {
			return err
		}

		page := new(models.Page)
		query := func(q *orm.Query) *orm.Query {
			return q.Where("title = ? and db_name = ? and revision = ?", data.Title, data.DbName, data.Revision)
		}

		if err := repo.Find(ctx, page, query); err != nil {
			return err
		}

		var cont *schema.Page

		key, err := json.Marshal(schema.PageKey{
			Title:  data.Title,
			DbName: data.DbName,
		})

		if err != nil {
			return err
		}

		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topics.PageVisibility, Partition: 0},
			Key:            key,
		}

		if data.Visible {
			if cont, err = storage.Pull(ctx, page, mediawiki.NewClient(data.SiteURL)); err != nil {
				return err
			}

			cont.Visible = &data.Visible

			if msg.Value, err = json.Marshal(cont); err != nil {
				return err
			}

			return producer.Produce(msg, nil)
		}

		if err := storage.Delete(ctx, page); err != nil {
			return err
		}

		cont = content.NewStructured(page)
		cont.Visible = &data.Visible

		if msg.Value, err = json.Marshal(cont); err != nil {
			return err
		}

		return producer.Produce(msg, nil)
	}
}
