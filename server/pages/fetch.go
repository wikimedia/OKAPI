package pages

import (
	"context"
	"errors"
	"log"
	"math"
	"okapi-data-service/models"
	pb "okapi-data-service/server/pages/protos"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/mediawiki-api-client"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
)

var errBatchSizeToBig = errors.New("batch size should be no more that 50")

type fetchRepo interface {
	repository.SelectOrCreator
	repository.Finder
	repository.Updater
}

// Fetch get page titles from the dumps and add the to the sotage and database
func Fetch(ctx context.Context, req *pb.FetchRequest, repo fetchRepo, dumps *dumps.Client) (*pb.FetchResponse, error) {
	res := new(pb.FetchResponse)
	proj := new(models.Project)
	err := repo.Find(ctx, proj, func(q *orm.Query) *orm.Query {
		return q.
			Where("db_name = ?", req.DbName)
	})

	if err != nil {
		return nil, err
	}

	if req.Batch > 50 {
		return nil, errBatchSizeToBig
	}

	if req.Batch <= 0 {
		req.Batch = 50
	}

	titles, err := dumps.PageTitles(ctx, req.DbName, time.Now())

	if err != nil {
		return nil, err
	}

	length := len(titles)
	batches := int(math.Ceil(float64(length) / float64(req.Batch)))
	mwiki := mediawiki.NewClient(proj.SiteURL)
	jobs := make(chan []string, batches)
	errs := make(chan []error, batches)

	for i := 0; i < int(req.Workers); i++ {
		go func() {
			for titles := range jobs {
				errs <- fetchWorker(ctx, titles, proj, repo, mwiki)
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
		for _, err := range <-errs {
			if err != nil {
				log.Println(err)
				res.Errors++
			} else {
				res.Redirects--
			}
		}
	}

	return res, nil
}

func fetchWorker(ctx context.Context, titles []string, proj *models.Project, repo fetchRepo, mwiki *mediawiki.Client) []error {
	errs := make([]error, 0)
	data, err := mwiki.PagesData(ctx, titles...)

	if err != nil {
		errs = append(errs, err)
	}

	for title, meta := range data {
		page := &models.Page{
			Title:   title,
			QID:     meta.Pageprops.WikibaseItem,
			PID:     meta.PageID,
			NsID:    meta.Ns,
			Lang:    proj.Lang,
			DbName:  proj.DbName,
			SiteURL: proj.SiteURL,
		}
		query := func(q *orm.Query) *orm.Query {
			return q.
				Where("title = ? and db_name = ?", page.Title, page.DbName)
		}

		if len(meta.Revisions) <= 0 {
			log.Printf("db_name: %s, title: %s, err: revisions not found\n", meta.Title, page.DbName)
			continue
		}

		page.SetRevision(meta.LastRevID, meta.Revisions[0].Timestamp)
		created, err := repo.SelectOrCreate(ctx, page, query)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if !created {
			page.SetRevision(meta.LastRevID, meta.Revisions[0].Timestamp)
			_, err := repo.Update(ctx, page, query)
			errs = append(errs, err)
			continue
		}

		errs = append(errs, nil)
	}

	return errs
}
