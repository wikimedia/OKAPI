package main

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := "pages"

	column := pgmigrations.Column{
		Table: table,
		Name:  "failed",
		Type:  "bool not null default false",
	}

	up := func(db orm.DB) error {
		if _, err := db.Exec(column.Add()); err != nil {
			return err
		}

		_, err := db.Exec(fmt.Sprintf("CREATE INDEX idx_%s_%s_notunique_btree ON pages (%s) WHERE %s = TRUE;", table, column.Name, column.Name, column.Name))

		return err
	}

	down := func(db orm.DB) error {
		if _, err := db.Exec(fmt.Sprintf("DROP INDEX idx_%s_%s_notunique_btree;", table, column.Name)); err != nil {
			return err
		}

		_, err := db.Exec(column.Drop())
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210428110201_alter_pages_table", up, down, opts)
}
