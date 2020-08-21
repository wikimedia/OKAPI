package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "pages"
	columns := []db.Column{
		{
			Name: "id",
			Type: "SERIAL PRIMARY KEY",
		},
		{
			Name: "title",
			Type: "varchar(750) not null",
		},
		{
			Name: "project_id",
			Type: "bigint not null",
		},
		{
			Name: "revision",
			Type: "int not null",
		},
		{
			Name: "revisions",
			Type: "int[6] not null default '{}'",
		},
		{
			Name: "updates",
			Type: "int not null default 0",
		},
		{
			Name: "scores",
			Type: "jsonb default '{}'::jsonb",
		},
		{
			Name: "path",
			Type: "varchar(1000)",
		},
		{
			Name: "site_url",
			Type: "varchar(255)",
		},
		{
			Name: "updated_at",
			Type: "timestamp with time zone",
		},
		{
			Name: "created_at",
			Type: "timestamp with time zone",
		},
	}
	indexes := []db.Index{
		{
			TableName: tableName,
			Columns:   []string{"title"},
		},
		{
			TableName: tableName,
			Columns:   []string{"project_id"},
		},
		{
			Type:      "UNIQUE",
			TableName: tableName,
			Columns:   []string{"title", "project_id"},
		},
	}
	foreignKeys := []db.ForeignKey{
		{
			ParentTable: "projects",
			TableName:   tableName,
			Name:        "pages_projects_fk",
			Property:    "project_id",
			References:  "id",
			OnDelete:    db.Cascade,
		},
	}
	table := db.Table{
		Name:        tableName,
		Columns:     columns,
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
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

	migrations.Register("20200525121350_create_pages_table", up, down, opts)
}
