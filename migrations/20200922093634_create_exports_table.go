package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "exports"
	columns := []db.Column{
		{
			Name: "user_id",
			Type: "bigint not null",
		},
		{
			Name: "resource_type",
			Type: "varchar(25) not null",
		},
		{
			Name: "resource_id",
			Type: "bigint not null",
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
	foreignKeys := []db.ForeignKey{
		{
			ParentTable: "users",
			TableName:   tableName,
			Name:        "pages_users_fk",
			Property:    "user_id",
			References:  "id",
			OnDelete:    db.Cascade,
		},
	}
	table := db.Table{
		Name:        tableName,
		Columns:     columns,
		ForeignKeys: foreignKeys,
		PrimaryKey:  []string{"user_id", "resource_type", "resource_id"},
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

	migrations.Register("20200922093634_create_exports_table", up, down, opts)
}
