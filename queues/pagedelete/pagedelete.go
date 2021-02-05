package pagedelete

import (
	"context"
	"encoding/json"
	"okapi-data-service/models"
	"okapi-data-service/pkg/worker"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
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

// Storages all necessary sotrages
type Storages struct {
	HTML  storage.Deleter
	WText storage.Deleter
	JSON  storage.Deleter
}

// Worker processing function
func Worker(repo Repo, stores *Storages) worker.Worker {
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

		var result error

		if err := stores.HTML.Delete(page.HTMLPath); err != nil {
			result = err
		}

		if err := stores.WText.Delete(page.WikitextPath); err != nil {
			result = err
		}

		if err := stores.JSON.Delete(page.JSONPath); err != nil {
			result = err
		}

		return result
	}
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}
