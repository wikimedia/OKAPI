package main

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := pgmigrations.Table{
		Name: "namespaces",
		Constraints: map[pgmigrations.Constraint][]string{
			pgmigrations.ConstraintPrimaryKey: {
				pgmigrations.Columns([]string{
					"id",
					"lang",
				}),
			},
		},
		Columns: []pgmigrations.Column{
			{
				Name: "id",
				Type: "int not null",
			},
			{
				Name: "title",
				Type: "varchar(255) not null",
			},
			{
				Name: "lang",
				Type: fmt.Sprintf("varchar(25) not null references languages(code) on update %s", pgmigrations.ActionCascade),
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

	migrations.Register("20201221203541_create_namespaces_table", up, down, opts)
}
