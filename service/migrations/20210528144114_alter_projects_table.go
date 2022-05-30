package main

import (
	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := "projects"
	columns := []pgmigrations.Column{
		{
			Table: table,
			Name:  "delay",
			Type:  "smallint not null default 48",
		},
		{
			Table: table,
			Name:  "threshold",
			Type:  "real not null default 0.6",
		},
		{
			Table: table,
			Name:  "html_size",
			Type:  "double precision",
		},
		{
			Table: table,
			Name:  "wikitext_size",
			Type:  "double precision",
		},
		{
			Table: table,
			Name:  "json_size",
			Type:  "double precision",
		},
		{
			Table: table,
			Name:  "html_path",
			Type:  "varchar(255)",
		},
		{
			Table: table,
			Name:  "wikitext_path",
			Type:  "varchar(255)",
		},
		{
			Table: table,
			Name:  "json_path",
			Type:  "varchar(255)",
		},
		{
			Table: table,
			Name:  "html_at",
			Type:  "timestamp with time zone",
		},
		{
			Table: table,
			Name:  "wikitext_at",
			Type:  "timestamp with time zone",
		},
		{
			Table: table,
			Name:  "json_at",
			Type:  "timestamp with time zone",
		},
	}

	up := func(db orm.DB) error {
		for _, colum := range columns {
			if _, err := db.Exec(colum.Drop()); err != nil {
				return err
			}
		}

		return nil
	}

	down := func(db orm.DB) error {
		for _, colum := range columns {
			if _, err := db.Exec(colum.Add()); err != nil {
				return err
			}
		}

		return nil
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210528144114_alter_projects_table", up, down, opts)
}
