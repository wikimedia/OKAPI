package pages

import (
	"context"
	"log"
	"net/http"
	"okapi-data-service/models"

	"okapi-data-service/server/pages/content"
	pb "okapi-data-service/server/pages/protos"
	"sync"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"

	"github.com/protsack-stephan/mediawiki-api-client"

	"github.com/go-pg/pg/v10/orm"
)

const pullDefLimit = 10000

type pullRepo interface {
	repository.Finder
	repository.Updater
}

// Pull get all of the pages from db and store html/wikitext from API calls on the hard drive
func Pull(ctx context.Context, req *pb.PullRequest, repo pullRepo, storage *content.Storage) (*pb.PullResponse, error) {
	proj := new(models.Project)
	err := repo.Find(ctx, proj, func(q *orm.Query) *orm.Query {
		return q.
			Where("db_name = ?", req.DbName)
	})

	if err != nil {
		return nil, err
	}

	res := new(pb.PullResponse)
	jobWg, errWg := new(sync.WaitGroup), new(sync.WaitGroup)
	jobs := make(chan models.Page, int(req.Workers))
	errs := make(chan error, int(req.Workers))
	mwiki := mediawiki.
		NewBuilder(proj.SiteURL).
		HTTPClient(&http.Client{Timeout: time.Second * 60}).
		Build()

	if req.Limit <= 0 {
		req.Limit = pullDefLimit
	}

	errWg.Add(1)
	go func() {
		defer errWg.Done()
		for err := range errs {
			res.Total++

			if err != nil {
				log.Println(err)
				res.Errors++
			}
		}
	}()

	for i := 0; i < int(req.Workers); i++ {
		jobWg.Add(1)
		go func() {
			defer jobWg.Done()
			for page := range jobs {
				err := content.Pull(ctx, &page, storage, mwiki)

				if err != nil {
					errs <- err
				} else {
					_, err := repo.Update(ctx, &page, func(q *orm.Query) *orm.Query {
						return q.Where("title = ? and db_name = ?", page.Title, page.DbName)
					}, "wikitext_path", "html_path", "json_path", "updated_at")
					errs <- err
				}
			}
		}()
	}

	pointer := 0

	for {
		pages := make([]models.Page, 0)
		err := repo.Find(ctx, &pages, func(q *orm.Query) *orm.Query {
			return q.
				Where("db_name = ? and id > ?", proj.DbName, pointer).
				Limit(int(req.Limit)).
				Order("id asc")
		})

		if err != nil {
			return nil, err
		}

		if len(pages) <= 0 {
			break
		}

		for i, page := range pages {
			jobs <- page

			if len(pages)-1 >= i {
				pointer = page.ID
			}
		}
	}

	close(jobs)
	jobWg.Wait()
	close(errs)
	errWg.Wait()

	return res, nil
}
