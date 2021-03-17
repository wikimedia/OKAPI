package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"okapi-data-service/models"
	pb "okapi-data-service/server/search/protos"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/go-pg/pg/v10/orm"
)

const lang string = "en"

type option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type fields struct {
	LangName []option `json:"lang_name"`
	SiteName []option `json:"site_name"`
	SiteCode []option `json:"site_code"`
	NsID     []option `json:"ns_id"`
	Lang     []option `json:"lang"`
}

type group struct {
	Model interface{}
	Query func(q *orm.Query) *orm.Query
}

// Aggregate create dataset for UI filters and upload it to the storage
func Aggregate(ctx context.Context, req *pb.AggregateRequest, repo repository.Finder, store storage.Putter) (*pb.AggregateResponse, error) {
	fields := new(fields)
	groups := map[interface{}]group{
		&fields.LangName: {
			&models.Language{},
			func(q *orm.Query) *orm.Query {
				return q.
					ColumnExpr("name as value, name as label").
					GroupExpr("value, label").
					OrderExpr("label asc")
			},
		},
		&fields.SiteName: {
			&models.Project{},
			func(q *orm.Query) *orm.Query {
				return q.
					ColumnExpr("site_name as value, site_name as label").
					GroupExpr("value, label").
					OrderExpr("label asc")
			},
		},
		&fields.SiteCode: {
			&models.Project{},
			func(q *orm.Query) *orm.Query {
				return q.
					ColumnExpr("site_code as value, site_name as label").
					Where("lang = ?", lang).
					GroupExpr("value, label").
					OrderExpr("label asc")
			},
		},
		&fields.Lang: {
			&models.Language{},
			func(q *orm.Query) *orm.Query {
				return q.
					ColumnExpr("code as value, local_name as label").
					GroupExpr("value, label").
					OrderExpr("label asc")
			},
		},
		&fields.NsID: {
			&models.Namespace{},
			func(q *orm.Query) *orm.Query {
				return q.
					ColumnExpr("id as value, title as label").
					Where("lang = ?", lang).
					GroupExpr("value, label").
					OrderExpr("label asc")
			},
		},
	}

	for field, group := range groups {
		if err := repo.Find(ctx, group.Model, group.Query, field); err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(fields)

	if err != nil {
		return nil, err
	}

	if err := store.Put(fmt.Sprintf("options/%s.json", lang), bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return new(pb.AggregateResponse), nil
}
