package pages

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"okapi-data-service/models"
	"okapi-data-service/pkg/page"
	"okapi-data-service/server/pages/fetch"
	pb "okapi-data-service/server/pages/protos"
	"strings"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/mediawiki-api-client"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
)

var errBatchSizeToBig = errors.New("batch size should be no more that 50")
var errMinNumberOfWorkers = errors.New("min number of workers is 1")

type fetchRepo interface {
	repository.Creator
	repository.Finder
	repository.Updater
}

type fetchStorage interface {
	storage.Putter
	storage.Deleter
}

// Fetch get page titles from the dumps and add the to the storage and database
func Fetch(ctx context.Context, req *pb.FetchRequest, repo fetchRepo, mwdump *dumps.Client, store fetchStorage, fetcher fetch.FetcherFactory) (*pb.FetchResponse, error) {
	res := new(pb.FetchResponse)
	proj := new(models.Project)
	err := repo.Find(ctx, proj, func(q *orm.Query) *orm.Query {
		return q.
			ColumnExpr("project.*, language.local_name as language__local_name, language.code as language__code").
			Join("left join languages as language").
			JoinOn("project.lang = language.code").
			Where("db_name = ?", req.DbName)
	})

	if err != nil {
		return nil, err
	}

	ns := new(models.Namespace)
	err = repo.Find(ctx, ns, func(q *orm.Query) *orm.Query {
		return q.Where("id = ? and lang = ?", req.Ns, proj.Lang)
	})

	if err != nil {
		return nil, err
	}

	if req.Batch > 50 {
		return nil, errBatchSizeToBig
	}

	if req.Workers <= 0 {
		return nil, errMinNumberOfWorkers
	}

	if req.Batch <= 0 {
		req.Batch = 50
	}

	titles := []string{}

	if req.Failed {
		pages := []*models.Page{}

		err := repo.Find(ctx, &pages, func(q *orm.Query) *orm.Query {
			return q.
				Column("title").
				Where("db_name = ? and ns_id = ? and failed = true", req.DbName, req.Ns)
		})

		if err != nil {
			return nil, err
		}

		for _, page := range pages {
			titles = append(titles, page.Title)
		}
	} else {
		filter := func(p *dumps.Page) {
			if p.Ns == int(req.Ns) {
				titles = append(titles, p.Title)
			}
		}

		if req.Ns == 0 {
			if err := mwdump.PageTitles(ctx, req.DbName, time.Now().UTC(), filter); err != nil {
				if err := mwdump.PageTitles(ctx, req.DbName, time.Now().UTC().Add(-24*time.Hour), filter); err != nil {
					return nil, err
				}
			}
		} else {
			date := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)

			if time.Now().UTC().Day() > 20 {
				date = time.Date(time.Now().Year(), time.Now().Month(), 20, 0, 0, 0, 0, time.UTC)
			}

			if err := mwdump.PageTitlesNs(ctx, req.DbName, date, filter); err != nil {
				return nil, err
			}
		}
	}

	length := len(titles)
	batches := int(math.Ceil(float64(length) / float64(req.Batch)))
	mwiki := mediawiki.NewClient(proj.SiteURL)
	jobs := make(chan []string, batches)
	errs := make(chan map[string]error, batches)
	worker := fetcher.Create(
		&page.Factory{
			Project:   proj,
			Language:  proj.Language,
			Namespace: ns,
		},
		store,
		mwiki,
		repo)

	for i := 0; i < int(req.Workers); i++ {
		go func() {
			for titles := range jobs {
				_, fErrs, err := worker.Fetch(ctx, titles...)

				if err != nil {
					log.Println(strings.Replace(err.Error(), "\n", " ", -1))
				}

				errs <- fErrs
			}
		}()
	}

	for i := 1; i <= batches; i++ {
		start, end := (i-1)*int(req.Batch), i*int(req.Batch)

		if end > length {
			jobs <- titles[start:]
		} else {
			jobs <- titles[start:end]
		}
	}

	close(jobs)
	res.Total = int32(length)
	res.Redirects = int32(length)

	for i := 1; i <= batches; i++ {
		for title, err := range <-errs {
			if err != nil {
				log.Println(fmt.Sprintf("title: %s err: %v", title, err))
				res.Errors++
			} else {
				res.Redirects--
			}
		}
	}

	return res, nil
}
