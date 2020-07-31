package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	"okapi/lib/db"
)

func init() {
	tableName := "projects"
	columns := []db.Column{
		{
			Name: "id",
			Type: "SERIAL PRIMARY KEY",
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
			Name: "code",
			Type: "varchar(25) not null",
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
	table := db.Table{
		Name:    tableName,
		Columns: columns,
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
