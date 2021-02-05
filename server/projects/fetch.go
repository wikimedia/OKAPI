package projects

import (
	"context"
	"fmt"
	"okapi-data-service/models"

	pb "okapi-data-service/server/projects/protos"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/mediawiki-api-client"
)

const partitionQuery = `CREATE TABLE IF NOT EXISTS  "pages_%s" PARTITION OF pages FOR VALUES in ('%s')`

type fetchRepo interface {
	repository.SelectOrCreator
	repository.Executor
}

// Fetch get all the projects from the api
func Fetch(ctx context.Context, req *pb.FetchRequest, mWiki *mediawiki.Client, repo fetchRepo) (*pb.FetchResponse, error) {
	sites, err := mWiki.Sitematrix(ctx)
	defaults := map[string]string{
		"shy": "Shawiya",
	}

	if err != nil {
		return nil, err
	}

	for _, proj := range sites.Projects {
		lang := models.Language{
			Name:      proj.Name,
			LocalName: proj.Localname,
			Code:      proj.Code,
			Dir:       proj.Dir,
		}

		if name, ok := defaults[proj.Code]; len(lang.Name) <= 0 && ok {
			lang.Name = name
		}

		if localName, ok := defaults[proj.Code]; len(lang.LocalName) <= 0 && ok {
			lang.LocalName = localName
		}

		_, err := repo.SelectOrCreate(ctx, &lang, func(q *orm.Query) *orm.Query {
			return q.Where("code = ?", proj.Code)
		})

		if err != nil {
			return nil, err
		}

		for _, site := range proj.Site {
			project := models.Project{
				DbName:   site.DBName,
				SiteName: site.Sitename,
				SiteCode: site.Code,
				SiteURL:  site.URL,
				Lang:     proj.Code,
				Active:   !site.Closed,
			}

			_, err := repo.SelectOrCreate(ctx, &project, func(q *orm.Query) *orm.Query {
				return q.Where("db_name = ?", project.DbName)
			})

			if err != nil {
				return nil, err
			}

			_, err = repo.Exec(ctx, fmt.Sprintf(partitionQuery, project.DbName, project.DbName))

			if err != nil {
				return nil, err
			}
		}
	}

	return new(pb.FetchResponse), nil
}
