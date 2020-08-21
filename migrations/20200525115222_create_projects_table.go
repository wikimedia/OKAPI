package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "projects"
	columns := []db.Column{
		{
			Name: "id",
			Type: "SERIAL PRIMARY KEY",
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
			Name: "db_name",
			Type: "varchar(255) not null unique",
		},
		{
			Name: "schedule",
			Type: "jsonb not null default '{}'::jsonb",
		},
		{
			Name: "threshold",
			Type: "jsonb not null default '{\"damaging\": 0.6}'::jsonb",
		},
		{
			Name: "lang_name",
			Type: "varchar(255) not null",
		},
		{
			Name: "lang_local_name",
			Type: "varchar(255) not null",
		},
		{
			Name: "updates",
			Type: "int not null default 0",
		},
		{
			Name: "time_delay",
			Type: "smallint not null default 48",
		},
		{
			Name: "lang",
			Type: "varchar(25) not null",
		},
		{
			Name: "dir",
			Type: "varchar(3) not null",
		},
		{
			Name: "size",
			Type: "double precision",
		},
		{
			Name: "path",
			Type: "varchar(255)",
		},
		{
			Name: "active",
			Type: "bool not null",
		},
		{
			Name: "dumped_at",
			Type: "timestamp with time zone",
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
			Columns: []string{
				"db_name",
			},
		},
	}
	table := db.Table{
		Name:    tableName,
		Columns: columns,
		Indexes: indexes,
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

	migrations.Register("20200525115222_create_projects_table", up, down, opts)
}
