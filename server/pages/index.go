package pages

import (
	"bytes"
	"context"
	"encoding/json"
	"okapi-data-service/models"
	"okapi-data-service/pkg/index"
	pb "okapi-data-service/server/pages/protos"
	"runtime"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
)

// Index all the pages for search
func Index(ctx context.Context, req *pb.IndexRequest, elastic *elasticsearch.Client, repo repository.Finder) (*pb.IndexResponse, error) {
	res := new(pb.IndexResponse)
	pointer := 0

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      index.Page,
		Client:     elastic,
		NumWorkers: runtime.NumCPU(),
	})

	if err != nil {
		return nil, err
	}

	for {
		pages := make([]models.Page, 0)
		err := repo.Find(ctx, &pages, func(q *orm.Query) *orm.Query {
			return q.
				ColumnExpr("page.*, project.site_code as project__site_code, language.name as language__name, language.local_name as language__local_name").
				Join("left join languages as language").
				JoinOn("language.code = page.lang").
				Join("left join projects as project").
				JoinOn("project.db_name = page.db_name").
				Where("page.id > ?", pointer).
				Limit(10000).
				Order("page.id asc")
		})

		if err != nil {
			return nil, err
		}

		if len(pages) <= 0 {
			break
		}

		for i, page := range pages {
			body, err := json.Marshal(index.DocPage{
				ID:            page.ID,
				Title:         page.Title,
				NsID:          page.NsID,
				DbName:        page.DbName,
				Lang:          page.Lang,
				LangName:      page.Language.Name,
				LangLocalName: page.Language.LocalName,
				SiteCode:      page.Project.SiteCode,
				SiteURL:       page.SiteURL,
				UpdatedAt:     page.UpdatedAt,
			})

			if err != nil {
				return nil, err
			}

			res.Total++
			err = bi.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(page.ID),
				Body:       bytes.NewReader(body),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, r esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						res.Errors++
					}
				},
			})

			if err != nil {
				return nil, err
			}

			if len(pages)-1 >= i {
				pointer = page.ID
			}
		}
	}

	if err := bi.Close(ctx); err != nil {
		return nil, err
	}

	return res, nil
}
