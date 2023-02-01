package projects

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"okapi-data-service/models"
	"okapi-data-service/schema/v3"
	pb "okapi-data-service/server/projects/protos"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

type aggrStore interface {
	storage.Putter
	storage.Getter
}

// Aggregate generate list of projects for API to serve
func Aggregate(ctx context.Context, _ *pb.AggregateRequest, repo repository.Finder, store aggrStore) (*pb.AggregateResponse, error) {
	res := new(pb.AggregateResponse)
	pointer := 0
	schemas := []*schema.Project{}
	exports := map[int][]*schema.Project{}

	for _, ns := range namespaces {
		exports[ns] = []*schema.Project{}
	}

	for {
		projects := make([]models.Project, 0)
		err := repo.Find(ctx, &projects, func(q *orm.Query) *orm.Query {
			return q.
				ColumnExpr("project.*, language.name as language__name, language.local_name as language__local_name, language.code as language__code").
				Join("left join languages as language").
				JoinOn("language.code = project.lang").
				Where("project.id > ? and active = true", pointer).
				Limit(1000).
				Order("project.id asc")
		})

		if err != nil {
			return nil, err
		}

		if len(projects) <= 0 {
			break
		}

		for _, proj := range projects {
			// collect an actual metadata from the storage by namespace
			for nsID := range exports {
				mrc, err := store.Get(fmt.Sprintf("export/%s/%s_%d.json", proj.DbName, proj.DbName, nsID))

				if err != nil {
					log.Println(err)
					continue
				}

				meta := new(schema.Project)

				if err := json.NewDecoder(mrc).Decode(meta); err != nil {
					log.Println(err)
					continue
				}

				exports[nsID] = append(exports[nsID], meta)

				_ = mrc.Close()
			}

			// collect initial list of projects with base metadata
			schemas = append(schemas, &schema.Project{
				Name:       proj.SiteName,
				Identifier: proj.DbName,
				URL:        proj.SiteURL,
				InLanguage: &schema.Language{
					Name:       proj.Language.LocalName,
					Identifier: proj.Language.Code,
				},
			})

			res.Total++
		}

		pointer = projects[len(projects)-1].ID
	}

	for nsID, meta := range exports {
		if len(meta) == 0 {
			continue
		}

		data, err := json.Marshal(meta)

		if err != nil {
			log.Println(err)
			continue
		}

		if err := store.Put(fmt.Sprintf("public/exports_%d.json", nsID), bytes.NewReader(data)); err != nil {
			log.Println(err)
		}
	}

	data, err := json.Marshal(schemas)

	if err != nil {
		return res, err
	}

	return res, store.Put("public/projects.json", bytes.NewReader(data))
}
