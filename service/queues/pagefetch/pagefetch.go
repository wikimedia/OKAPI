package pagefetch

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"okapi-data-service/lib/env"
	"okapi-data-service/models"
	"okapi-data-service/pkg/page"
	"okapi-data-service/pkg/producer"
	"okapi-data-service/pkg/worker"
	"okapi-data-service/schema/v3"
	"okapi-data-service/server/pages/fetch"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/mediawiki-api-client"
)

// Name redis key for the queue
const Name string = "queue/pagefetch"

// ErrPageNotFound page was not found
var ErrPageNotFound = errors.New("page not found")

// Hold the list of clients
var clients sync.Map

// Data item of the queue
type Data struct {
	Title     string         `json:"title"`
	Revision  int            `json:"revision"`
	DbName    string         `json:"db_name"`
	Lang      string         `json:"lang"`
	SiteURL   string         `json:"site_url"`
	Namespace int            `json:"namespace"`
	Scores    *schema.Scores `json:"scores,omitempty"`
	Editor    *schema.Editor `json:"editor,omitempty"`
}

func Worker(fetcher fetch.FetcherFactory, store fetch.Storage, repo fetch.Repo, producer producer.Producer) worker.Worker {
	return func(ctx context.Context, payload []byte) error {
		data := new(Data)

		if err := json.Unmarshal(payload, data); err != nil {
			return err
		}

		proj := new(models.Project)
		pquery := func(q *orm.Query) *orm.Query {
			return q.
				ColumnExpr("project.*, language.local_name as language__local_name, language.code as language__code").
				Join("left join languages as language").
				JoinOn("project.lang = language.code").
				Where("db_name = ?", data.DbName)
		}

		if err := repo.Find(ctx, proj, pquery); err != nil {
			return err
		}

		ns := new(models.Namespace)
		nquery := func(q *orm.Query) *orm.Query {
			return q.Where("lang = ? and id = ?", proj.Lang, data.Namespace)
		}

		if err := repo.Find(ctx, ns, nquery); err != nil {
			return err
		}

		info, _ := clients.LoadOrStore(data.SiteURL, mediawiki.
			NewBuilder(data.SiteURL).
			HTTPClient(&http.Client{Timeout: time.Second * 30}).
			Headers(map[string]string{
				"User-Agent": env.MediawikiAPIUserAgent,
			}).
			Build())
		cl := info.(*mediawiki.Client)
		worker := fetcher.Create(
			&page.Factory{
				Project:   proj,
				Language:  proj.Language,
				Namespace: ns,
			},
			store,
			cl,
			repo)

		tmCtx, cancel := context.WithTimeout(ctx, time.Second*120)
		defer cancel()

		pages, errs, err := worker.Fetch(tmCtx, data.Title)

		if err != nil {
			return err
		}

		if err := errs[data.Title]; err != nil {
			return err
		}

		page, ok := pages[data.Title]

		if !ok {
			return ErrPageNotFound
		}

		if page.Version != nil {
			if page.Version.Identifier == data.Revision {
				page.Version.Editor = data.Editor

				if data.Scores != nil {
					page.Version.Scores = data.Scores
				}
			} else if page.Version.Editor.Identifier != 0 {
				user, err := cl.User(ctx, page.Version.Editor.Identifier)

				if err == nil {
					page.Version.Editor.DateStarted = &user.Registration
					page.Version.Editor.EditCount = user.EditCount
					page.Version.Editor.Groups = user.Groups

					if !page.Version.Editor.IsAnonymous {
						for _, group := range user.Groups {
							if group == "bot" {
								page.Version.Editor.IsBot = true
								break
							}
						}
					}
				}
			}
		}

		value, err := json.Marshal(page)

		if err != nil {
			return err
		}

		key, err := json.Marshal(schema.PageKey{
			Name:     data.Title,
			IsPartOf: data.DbName,
		})

		if err != nil {
			return err
		}

		producer.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &schema.TopicPageUpdate, Partition: 0},
			Key:            key,
			Value:          value,
		}

		return nil
	}
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}
