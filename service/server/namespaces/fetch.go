package namespaces

import (
	"context"
	"okapi-data-service/models"
	pb "okapi-data-service/server/namespaces/protos"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/mediawiki-api-client"

	"github.com/go-pg/pg/v10/orm"
)

type fetchRepo interface {
	repository.SelectOrCreator
	repository.Finder
}

// Fetch get all namespaces from mediawiki and add them to the database
func Fetch(ctx context.Context, req *pb.FetchRequest, repo fetchRepo) (*pb.FetchResponse, error) {
	pointer := 0
	defaults := map[int]string{
		0: "Article",
	}

	for {
		projects := make([]models.Project, 0)
		err := repo.Find(ctx, &projects, func(q *orm.Query) *orm.Query {
			return q.
				Where("project.id > ?", pointer).
				Limit(1000)
		})

		if err != nil {
			return nil, err
		}

		if len(projects) <= 0 {
			break
		}

		for i, proj := range projects {
			namespaces, err := mediawiki.
				NewClient(proj.SiteURL).
				Namespaces(ctx)

			if err != nil {
				return nil, err
			}

			for _, ns := range namespaces {
				model := &models.Namespace{
					ID:    ns.ID,
					Lang:  proj.Lang,
					Title: ns.Name,
				}

				if title, ok := defaults[ns.ID]; len(model.Title) <= 0 && ok {
					model.Title = title
				}

				_, err := repo.SelectOrCreate(ctx, model, func(q *orm.Query) *orm.Query {
					return q.Where("id = ? and lang = ?", ns.ID, proj.Lang)
				})

				if err != nil {
					return nil, err
				}
			}

			if len(projects)-1 >= i {
				pointer = proj.ID
			}
		}
	}

	return new(pb.FetchResponse), nil
}
