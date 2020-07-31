package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	column := db.Column{
		TableName: "projects",
		Name:      "schedule",
		Type:      "jsonb not null default '{}'::jsonb",
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(column.Add())
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(column.Drop())
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200720102312_add_fields_to_projects_table", up, down, opts)
}
