package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"okapi-data-service/models"
	"okapi-data-service/pkg/page"
	"okapi-data-service/schema/v3"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/mediawiki-api-client"
)

const batchRequests = 10

type pageError struct {
	title string
	err   error
}

type response struct {
	data []byte
	err  error
}

type Worker struct {
	fact  *page.Factory
	store Storage
	mwiki *mediawiki.Client
	repo  Repo
}

func (w Worker) GetPagesData(ctx context.Context, titles []string) (map[string]mediawiki.PageData, error) {
	return w.mwiki.PagesData(ctx, titles)
}

func (w Worker) GetPagesModel(ctx context.Context, titles []string) (map[string]*models.Page, error) {
	pages := []*models.Page{}
	err := w.repo.Find(ctx, &pages, func(q *orm.Query) *orm.Query {
		return q.
			Where("db_name = ?", w.fact.Project.DbName).
			WhereIn("title in (?)", titles)
	})

	models := map[string]*models.Page{}

	for _, page := range pages {
		models[page.Title] = page
	}

	return models, err
}

func (w Worker) GetPagesHTML(ctx context.Context, titles []string) map[string]*response {
	workers := int(math.Ceil(float64(len(titles)) / batchRequests))
	data := make(chan map[string]*response, workers)

	for i := 0; i < workers; i++ {
		from, to := i*batchRequests, ((i * batchRequests) + batchRequests)

		if to > len(titles) {
			to = len(titles)
		}

		go func(titles []string) {
			resps := map[string]*response{}

			for _, title := range titles {
				res := new(response)
				res.data, res.err = w.mwiki.PageHTML(ctx, title)
				resps[title] = res
			}

			data <- resps
		}(titles[from:to])
	}

	resps := map[string]*response{}

	for i := 0; i < workers; i++ {
		for title, res := range <-data {
			resps[title] = res
		}
	}

	return resps
}

func (w Worker) UpdatePages(ctx context.Context, pages []*models.Page) map[string]error {
	errs := map[string]error{}

	for _, page := range pages {
		_, err := w.repo.Update(ctx, page, func(q *orm.Query) *orm.Query {
			return q.Where("db_name = ? and title = ?", page.DbName, page.Title)
		})

		errs[page.Title] = err
	}

	return errs
}

func (w Worker) CreatePages(ctx context.Context, pages []*models.Page) map[string]error {
	errs := map[string]error{}

	if len(pages) <= 0 {
		return errs
	}

	_, err := w.repo.Create(ctx, &pages)

	for _, page := range pages {
		errs[page.Title] = err
	}

	return errs
}

func (w Worker) SaveSchemas(ctx context.Context, pages map[string]*schema.Page) map[string]*response {
	semaphore := make(chan int, len(pages))
	resps := make(chan map[string]*response, len(pages))

	for title, page := range pages {
		semaphore <- 1
		go func(title string, page *schema.Page) {
			path := fmt.Sprintf("json/%s/%s.json", w.fact.Project.DbName, title)
			data, err := json.Marshal(page)

			if err != nil {
				resps <- map[string]*response{
					title: {
						err:  err,
						data: []byte(path),
					},
				}
			} else {
				resps <- map[string]*response{
					title: {
						err:  w.store.Put(path, bytes.NewReader(data)),
						data: []byte(path),
					},
				}
			}
			<-semaphore
		}(title, page)
	}

	data := map[string]*response{}

	for i := 0; i < len(pages); i++ {
		for title, resp := range <-resps {
			data[title] = resp
		}
	}

	return data
}

// Fetch bulk download and update data
func (w Worker) Fetch(ctx context.Context, titles ...string) (map[string]*schema.Page, map[string]error, error) {
	reqs := make(chan error, 3)
	htmls := map[string]*response{}

	go func() {
		htmls = w.GetPagesHTML(ctx, titles)
		reqs <- nil
	}()

	records := make(map[string]*models.Page)

	go func() {
		data, err := w.GetPagesModel(ctx, titles)
		records = data
		reqs <- err
	}()

	pages := make(map[string]mediawiki.PageData)

	go func() {
		data, err := w.GetPagesData(ctx, titles)
		pages = data
		reqs <- err
	}()

	for i := 0; i < 3; i++ {
		if err := <-reqs; err != nil {
			return nil, nil, err
		}
	}

	errs := map[string]error{}
	schemas := map[string]*schema.Page{}

	for title, pdata := range pages {
		res, ok := htmls[title]

		if ok && res.err == nil {
			schemas[title] = w.fact.Create(&pdata, string(res.data)) // #nosec G601
		} else {
			errs[title] = res.err
		}
	}

	updates, creates := []*models.Page{}, []*models.Page{}
	paths := w.SaveSchemas(ctx, schemas)

	for title, pdata := range pages {
		record, isUpdate := records[title]
		page := models.NewPage(title, &pdata, w.fact.Project, record) // #nosec G601

		if path, ok := paths[title]; ok {
			page.Failed = path.err != nil
			page.Path = string(path.data)
		} else {
			page.Failed = true
		}

		if isUpdate {
			updates = append(updates, page)
		} else {
			creates = append(creates, page)
		}
	}

	pLen := len(creates) + len(updates)
	pErrs := make(chan pageError, pLen)

	go func() {
		for title, err := range w.CreatePages(ctx, creates) {
			pErrs <- pageError{
				title,
				err,
			}
		}
	}()

	go func() {
		for title, err := range w.UpdatePages(ctx, updates) {
			pErrs <- pageError{
				title,
				err,
			}
		}
	}()

	for i := 0; i < pLen; i++ {
		pErr := <-pErrs
		errs[pErr.title] = pErr.err
	}

	return schemas, errs, nil
}
