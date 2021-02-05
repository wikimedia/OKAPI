package main

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := pgmigrations.Table{
		Name: "projects",
		Columns: []pgmigrations.Column{
			{
				Name: "id",
				Type: "serial primary key",
			},
			{
				Name: "db_name",
				Type: "varchar(255) not null unique",
			},
			{
				Name: "site_name",
				Type: "varchar(255) not null",
			},
			{
				Name: "site_code",
				Type: "varchar(255) not null",
			},
			{
				Name: "site_url",
				Type: "varchar(255) not null",
			},
			{
				Name: "lang",
				Type: fmt.Sprintf("varchar(25) not null references languages(code) on update %s", pgmigrations.ActionCascade),
			},
			{
				Name: "active",
				Type: "bool not null",
			},
			{
				Name: "delay",
				Type: "smallint not null default 48",
			},
			{
				Name: "threshold",
				Type: "real not null default 0.6",
			},
			{
				Name: "html_size",
				Type: "double precision",
			},
			{
				Name: "wikitext_size",
				Type: "double precision",
			},
			{
				Name: "json_size",
				Type: "double precision",
			},
			{
				Name: "html_path",
				Type: "varchar(255)",
			},
			{
				Name: "wikitext_path",
				Type: "varchar(255)",
			},
			{
				Name: "json_path",
				Type: "varchar(255)",
			},
			{
				Name: "html_at",
				Type: "timestamp with time zone",
			},
			{
				Name: "wikitext_at",
				Type: "timestamp with time zone",
			},
			{
				Name: "json_at",
				Type: "timestamp with time zone",
			},
			{
				Name: "updated_at",
				Type: "timestamp with time zone not null",
			},
			{
				Name: "created_at",
				Type: "timestamp with time zone not null",
			},
		},
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(table.Create())
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(table.Drop())
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20201221203529_create_projects_table", up, down, opts)
}
