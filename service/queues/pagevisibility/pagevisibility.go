package pagevisibility

import (
	"context"
	"encoding/json"
	"fmt"
	"okapi-data-service/models"
	"okapi-data-service/pkg/producer"
	"okapi-data-service/pkg/worker"
	"okapi-data-service/schema/v3"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Name redis key for the queue
const Name string = "queue/pagevisibility"

// Data item of the queue
type Data struct {
	ID         int            `json:"id"`
	Title      string         `json:"title"`
	DbName     string         `json:"db_name"`
	Revision   int            `json:"revision"`
	RevisionDt time.Time      `json:"revision_dt"`
	Visible    bool           `json:"visible"`
	Lang       string         `json:"lang"`
	SiteURL    string         `json:"site_url"`
	Namespace  int            `json:"namespace"`
	Editor     *schema.Editor `json:"editor,omitempty"`
	Visibility struct {
		Text    bool `json:"text"`
		User    bool `json:"user"`
		Comment bool `json:"comment"`
	} `json:"visibility"`
}

type Storage interface {
	storage.Deleter
	storage.Getter
}

// Enqueue add data to the worker queue
func Enqueue(ctx context.Context, store redis.Cmdable, data *Data) error {
	return worker.Enqueue(ctx, Name, store, data)
}

// Worker processing function
func Worker(repo repository.Finder, storage Storage, producer producer.Producer) worker.Worker {
	return func(ctx context.Context, payload []byte) error {
		data := new(Data)

		if err := json.Unmarshal(payload, data); err != nil {
			return err
		}

		path := fmt.Sprintf("%s/%s.json", data.DbName, data.Title)
		page := new(schema.Page)

		if prc, err := storage.Get(path); err == nil {
			defer prc.Close()

			if err := json.NewDecoder(prc).Decode(page); err != nil {
				return err
			}
		}

		resps := make(chan error, 2)

		go func() {
			if len(page.Name) == 0 {
				proj := new(models.Project)
				pQuery := func(q *orm.Query) *orm.Query {
					return q.
						ColumnExpr("project.*, language.local_name as language__local_name, language.code as language__code").
						Join("left join languages as language").
						JoinOn("project.lang = language.code").
						Where("db_name = ?", data.DbName)
				}

				if err := repo.Find(ctx, proj, pQuery); err != nil {
					resps <- err
					return
				}

				ns := new(models.Namespace)
				nQuery := func(q *orm.Query) *orm.Query {
					return q.Where("lang = ? and id = ?", proj.Lang, data.Namespace)
				}

				if err := repo.Find(ctx, ns, nQuery); err != nil {
					resps <- err
					return
				}

				page.Name = data.Title
				page.Identifier = data.ID
				page.URL = fmt.Sprintf("%s/wiki/%s", data.SiteURL, data.Title)
				page.Version = &schema.Version{
					Identifier: data.Revision,
				}

				license := schema.NewLicense()

				// Custom license for wikinews projects.
				if proj.SiteCode == "wikinews" {
					license = &schema.License{
						Name:       "Attribution 2.5 Generic",
						Identifier: "CC BY 2.5",
						URL:        "https://creativecommons.org/licenses/by/2.5/",
					}
				}

				page.License = append(page.License, license)
				page.DateModified = &data.RevisionDt
				page.Namespace = &schema.Namespace{
					Name:       ns.Title,
					Identifier: ns.ID,
				}
				page.IsPartOf = &schema.Project{
					Identifier: proj.DbName,
					Name:       proj.SiteName,
				}

				if proj.Language != nil {
					page.InLanguage = &schema.Language{
						Name:       proj.Language.LocalName,
						Identifier: proj.Language.Code,
					}
				}
			} else {
				page.ArticleBody = nil
				page.MainEntity = nil
			}

			if page.Version != nil {
				page.Version.Editor = data.Editor
			}

			page.Visibility = new(schema.Visibility)
			page.Visibility.Comment = data.Visibility.Comment
			page.Visibility.Text = data.Visibility.Text
			page.Visibility.User = data.Visibility.User
			value, err := json.Marshal(page)

			if err != nil {
				resps <- err
				return
			}

			key, err := json.Marshal(schema.PageKey{
				Name:     data.Title,
				IsPartOf: data.DbName,
			})

			if err != nil {
				resps <- err
				return
			}

			producer.ProduceChannel() <- &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &schema.TopicPageVisibility,
					Partition: 0,
				},
				Key:   key,
				Value: value,
			}

			resps <- nil
		}()

		go func() {
			if page.Version != nil && page.Version.Identifier == data.Revision && (!data.Visibility.Text || !data.Visibility.User || !data.Visibility.Comment) {
				resps <- storage.Delete(path)
			} else {
				resps <- nil
			}
		}()

		errs := []error{}

		for i := 0; i < 2; i++ {
			if err := <-resps; err != nil {
				errs = append(errs, err)
			}
		}

		for _, err := range errs {
			if err != nil {
				return err
			}
		}

		return nil
	}
}
