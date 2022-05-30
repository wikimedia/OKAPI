package projects

import (
	"bytes"
	"context"
	"encoding/json"
	"okapi-data-service/models"
	"okapi-data-service/pkg/index"
	pb "okapi-data-service/server/projects/protos"
	"runtime"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
)

// Index list all of the projects in elasticsearch
func Index(ctx context.Context, req *pb.IndexRequest, elastic *elasticsearch.Client, repo repository.Finder) (*pb.IndexResponse, error) {
	res := new(pb.IndexResponse)
	pointer := 0

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      index.Project,
		Client:     elastic,
		NumWorkers: runtime.NumCPU(),
	})

	if err != nil {
		return nil, err
	}

	for {
		projects := make([]models.Project, 0)
		err := repo.Find(ctx, &projects, func(q *orm.Query) *orm.Query {
			return q.
				ColumnExpr("project.*, language.name as language__name, language.local_name as language__local_name").
				Join("left join languages as language").
				JoinOn("language.code = project.lang").
				Where("project.id > ?", pointer).
				Limit(1000).
				Order("project.id asc")
		})

		if err != nil {
			return nil, err
		}

		if len(projects) <= 0 {
			break
		}

		for i, proj := range projects {
			body, err := json.Marshal(index.DocProject{
				ID:            proj.ID,
				DbName:        proj.DbName,
				SiteName:      proj.SiteName,
				SiteCode:      proj.SiteCode,
				SiteURL:       proj.SiteURL,
				Lang:          proj.Lang,
				LangName:      proj.Language.Name,
				LangLocalName: proj.Language.LocalName,
				Active:        proj.Active,
				UpdatedAt:     proj.UpdatedAt,
			})

			if err != nil {
				return nil, err
			}

			res.Total++
			err = bi.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(proj.ID),
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

			if len(projects)-1 >= i {
				pointer = proj.ID
			}
		}
	}

	if err := bi.Close(ctx); err != nil {
		return nil, err
	}

	return res, nil
}
