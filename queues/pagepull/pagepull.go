package pagepull

import (
	"context"
	"encoding/json"
	"errors"
	"okapi-data-service/models"
	"okapi-data-service/pkg/worker"
	"okapi-data-service/server/pages/content"

	"github.com/go-pg/pg/v10/orm"

	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/mediawiki-api-client"
)

// ErrPageNotFound page was not found
var ErrPageNotFound = errors.New("page not found")

// ErrRevisionsNotFound no revisions list in API response
var ErrRevisionsNotFound = errors.New("revisions not found")

// Name redis key for the queue
const Name string = "queue/pagepull"

// Probability for model
type Probability struct {
	True  float64 `json:"true"`
	False float64 `json:"false"`
}

// Damaging model for prediction whether revision is vandalism
type Damaging struct {
	Probability *Probability
}

// Models model scores
type Models struct {
	Damaging Damaging
}

// Data item of the queue
type Data struct {
	Title   string `json:"title"`
	DbName  string `json:"db_name"`
	Lang    string `json:"lang"`
	SiteURL string `json:"site_url"`
	Models  Models `json:"models"`
}

// Repo all needed interfaces for pagepull worker
type Repo interface {
	repository.SelectOrCreator
	repository.Updater
}

// Worker processing function
func Worker(repo Repo, storages *content.Storage) worker.Worker {
	return func(ctx context.Context, payload []byte) error {
		data := new(Data)

		if err := json.Unmarshal(payload, data); err != nil {
			return err
		}

		mwiki := mediawiki.NewClient(data.SiteURL)
		pages, err := mwiki.PagesData(ctx, data.Title)

		if err != nil {
			return err
		}

		meta, ok := pages[data.Title]

		if !ok {
			return ErrPageNotFound
		}

		if len(meta.Revisions) <= 0 {
			return ErrRevisionsNotFound
		}

		page := &models.Page{
			Title:   data.Title,
			QID:     meta.Pageprops.WikibaseItem,
			PID:     meta.PageID,
			NsID:    meta.Ns,
			Lang:    data.Lang,
			DbName:  data.DbName,
			SiteURL: data.SiteURL,
		}
		query := func(q *orm.Query) *orm.Query {
			return q.
				Where("title = ? and db_name = ?", data.Title, data.DbName)
		}

		page.SetRevision(meta.LastRevID, meta.Revisions[0].Timestamp)

		if _, err := repo.SelectOrCreate(ctx, page, func(q *orm.Query) *orm.Query {
			return query(q).Column("id", "created_at", "updated_at", "revision", "revisions", "lang")
		}); err != nil {
			return err
		}

		page.SetRevision(meta.LastRevID, meta.Revisions[0].Timestamp)

		if err := content.Pull(ctx, page, storages, mwiki); err != nil {
			return err
		}

		_, err = repo.Update(ctx, page, query)

		return err
	}
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}
