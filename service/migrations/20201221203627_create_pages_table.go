package main

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	pgmigrations "github.com/protsack-stephan/go-pg-migrations-helper"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	table := pgmigrations.Table{
		Name: "pages",
		Constraints: map[pgmigrations.Constraint][]string{
			pgmigrations.ConstraintPrimaryKey: {
				pgmigrations.Columns([]string{
					"title",
					"db_name",
				}),
			},
		},
		Columns: []pgmigrations.Column{
			{
				Name: "id",
				Type: "bigserial not null",
			},
			{
				Name: "title",
				Type: "varchar(750) not null",
			},
			{
				Name: "ns_id",
				Type: "int not null default 0",
			},
			{
				Name: "pid",
				Type: "bigint not null",
			},
			{
				Name: "qid",
				Type: "varchar(500)",
			},
			{
				Name: "db_name",
				Type: fmt.Sprintf("varchar(255) not null references projects(db_name) on update %s", pgmigrations.ActionCascade),
			},
			{
				Name: "site_url",
				Type: "varchar(255)",
			},
			{
				Name: "lang",
				Type: fmt.Sprintf("varchar(25) not null references languages(code) on update %s", pgmigrations.ActionCascade),
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
				Name: "revision_dt",
				Type: "timestamp with time zone not null",
			},
			{
				Name: "html_path",
				Type: "varchar(1000)",
			},
			{
				Name: "wikitext_path",
				Type: "varchar(1000)",
			},
			{
				Name: "json_path",
				Type: "varchar(1000)",
			},
			{
				Name: "updated_at",
				Type: "timestamp with time zone not null",
			},
			{
				Name: "created_at",
				Type: "timestamp with time zone not null",
			},
		},
		Indexes: []pgmigrations.Index{
			{
				Table:   "pages",
				Columns: []string{"id"},
			},
			{
				Table:   "pages",
				Columns: []string{"title"},
			},
		},
		Partition: &pgmigrations.Partition{
			Columns: []string{"db_name"},
			By:      pgmigrations.PartitionByList,
		},
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(table.Create())

		if err != nil {
			return err
		}

		_, err = db.Exec(fmt.Sprintf("create table %s_default partition of %s default", table.Name, table.Name))

		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(table.Drop())
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20201221203627_create_pages_table", up, down, opts)
}
