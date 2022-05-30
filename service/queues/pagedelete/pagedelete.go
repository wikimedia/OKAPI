package pagedelete

import (
	"context"
	"encoding/json"
	"errors"
	"okapi-data-service/models"
	"okapi-data-service/pkg/index"
	"okapi-data-service/pkg/producer"
	"okapi-data-service/pkg/worker"
	"okapi-data-service/schema/v3"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Name redis key for the queue
const Name string = "queue/pagedelete"

// Data item of the queue
type Data struct {
	Title  string         `json:"title"`
	DbName string         `json:"db_name"`
	Editor *schema.Editor `json:"editor,omitempty"`
}

// Repo all the needed repositories to call delete
type Repo interface {
	repository.Deleter
	repository.Finder
}

// Storage all the needed storages to call delete
type Storage interface {
	storage.Getter
	storage.Deleter
}

// Worker processing function
func Worker(repo Repo, storage Storage, producer producer.Producer, elastic *elasticsearch.Client) worker.Worker {
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

		rc, err := storage.Get(page.Path)

		if err != nil {
			return err
		}

		evt := new(schema.Page)
		err = json.NewDecoder(rc).Decode(evt)
		_ = rc.Close()

		if err != nil {
			return err
		}

		if err := storage.Delete(page.Path); err != nil {
			return err
		}

		key, err := json.Marshal(schema.PageKey{
			Name:     data.Title,
			IsPartOf: data.DbName,
		})

		if err != nil {
			return err
		}

		if evt.Version != nil {
			evt.Version.Editor = data.Editor
		}

		evt.ArticleBody = nil
		msg, err := json.Marshal(evt)

		if err != nil {
			return err
		}

		producer.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &schema.TopicPageDelete, Partition: 0},
			Key:            key,
			Value:          msg,
		}

		res, err := elastic.Delete(index.Page, strconv.Itoa(page.ID))

		if err != nil {
			return err
		}

		if res.IsError() {
			return errors.New(res.String())
		}

		return nil
	}
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}
