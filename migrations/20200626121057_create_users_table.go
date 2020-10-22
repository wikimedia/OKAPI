package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "users"
	columns := []db.Column{
		{
			Name: "id",
			Type: "SERIAL PRIMARY KEY",
		},
		{
			Name: "email",
			Type: "varchar(255) not null unique",
		},
		{
			Name: "password",
			Type: "varchar(255) not null",
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
			Columns:   []string{"email"},
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

	migrations.Register("20200626121057_create_users_table", up, down, opts)
}
