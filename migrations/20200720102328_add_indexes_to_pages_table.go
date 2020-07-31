package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "pages"
	indexes := []db.Index{
		{
			TableName: tableName,
			Columns:   []string{"title"},
		},
		{
			TableName: tableName,
			Columns:   []string{"title", "project_id"},
		},
	}

	up := func(db orm.DB) error {
		for _, index := range indexes {
			_, err := db.Exec(index.Create())

			if err != nil {
				return err
			}
		}

		return nil
	}

	down := func(db orm.DB) error {
		for _, index := range indexes {
			_, err := db.Exec(index.Drop())

			if err != nil {
				return err
			}
		}

		return nil
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200720102328_add_indexes_to_pages_table", up, down, opts)
}
