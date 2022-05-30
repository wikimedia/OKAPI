package main

import (
	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := "pages"
	dropCols := []pgmigrations.Column{
		{
			Table: table,
			Name:  "revisions",
			Type:  "int[6] not null default '{}'",
		},
		{
			Table: table,
			Name:  "html_path",
			Type:  "varchar(1000)",
		},
		{
			Table: table,
			Name:  "wikitext_path",
			Type:  "varchar(1000)",
		},
	}
	rnFromCol := pgmigrations.Column{
		Table: table,
		Name:  "json_path",
	}
	rnToCol := pgmigrations.Column{
		Table: table,
		Name:  "path",
	}

	up := func(db orm.DB) error {
		for _, column := range dropCols {
			if _, err := db.Exec(column.Drop()); err != nil {
				return err
			}
		}

		_, err := db.Exec(rnFromCol.Rename(rnToCol.Name))
		return err
	}

	down := func(db orm.DB) error {
		for _, column := range dropCols {
			if _, err := db.Exec(column.Add()); err != nil {
				return err
			}
		}

		_, err := db.Exec(rnToCol.Rename(rnFromCol.Name))
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210531073108_alter_pages_table", up, down, opts)
}
