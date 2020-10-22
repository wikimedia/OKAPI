package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	tableName := "users"
	usernameColumn := db.Column{
		TableName: tableName,
		Name:      "username",
	}
	index := db.Index{
		TableName: tableName,
		Type:      "unique",
		Columns:   []string{"username"},
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(usernameColumn.Set("not null") + index.Create())
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(usernameColumn.Set("default null") + index.Drop())
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20201009124313_make_users_username_as_main", up, down, opts)
}
