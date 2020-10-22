package main

import (
	"okapi/lib/db"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	wikitextPath := db.Column{
		TableName: "pages",
		Name:      "wikitext_path",
		Type:      "varchar(1000)",
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(wikitextPath.Add())
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(wikitextPath.Drop())
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20201007124722_add_wikitext_path_to_pages_table", up, down, opts)
}
