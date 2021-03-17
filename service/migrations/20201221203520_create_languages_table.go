package main

import (
	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := pgmigrations.Table{
		Name: "languages",
		Columns: []pgmigrations.Column{
			{
				Name: "id",
				Type: "serial primary key",
			},
			{
				Name: "code",
				Type: "varchar(25) not null unique",
			},
			{
				Name: "name",
				Type: "varchar(255) not null",
			},
			{
				Name: "local_name",
				Type: "varchar(255) not null",
			},
			{
				Name: "dir",
				Type: "varchar(3) not null",
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

	migrations.Register("20201221203520_create_languages_table", up, down, opts)
}
