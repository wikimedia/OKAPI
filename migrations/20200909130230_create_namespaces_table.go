package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "namespaces"
	columns := []db.Column{
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
			Type: "varchar(25) not null",
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
			Columns:   []string{"lang"},
		},
		{
			TableName: tableName,
			Columns:   []string{"lang", "id"},
		},
	}
	table := db.Table{
		Name:       tableName,
		PrimaryKey: []string{"id", "lang"},
		Columns:    columns,
		Indexes:    indexes,
	}
	column := db.Column{
		TableName: "pages",
		Name:      "ns_id",
		Type:      "int not null default 0",
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(table.Create())

		if err != nil {
			return err
		}

		_, err = db.Exec(column.Add())

		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(column.Drop())

		if err != nil {
			return err
		}

		_, err = db.Exec(table.Drop())

		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200909130230_create_namespaces_table", up, down, opts)
}
