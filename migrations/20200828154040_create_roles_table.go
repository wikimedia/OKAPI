package main

import (
	"okapi/lib/db"
	"okapi/models"
	"okapi/models/permissions"
	"okapi/models/roles"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

type roleMeta struct {
	Description string
	Permissions []permissions.Type
}

func init() {
	rolesData := map[roles.Type]roleMeta{
		roles.Admin: {
			Description: "Admin user with all available permissions.",
			Permissions: []permissions.Type{
				permissions.ProjectDelete,
				permissions.ProjectCreate,
				permissions.ProjectBundle,
			},
		},
		roles.Client: {
			Description: "Client user with limited amount of downloads.",
			Permissions: []permissions.Type{},
		},
		roles.Subscriber: {
			Description: "Client with unlimited amount of downloads and additional perks.",
			Permissions: []permissions.Type{},
		},
	}
	tableName := "roles"
	columns := []db.Column{
		{
			Name: "id",
			Type: "varchar(255) PRIMARY KEY",
		},
		{
			Name: "description",
			Type: "varchar(255) not null",
		},
		{
			Name: "permissions",
			Type: "varchar(255)[] not null default '{}'",
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
		if _, err := db.Exec(table.Create()); err != nil {
			return err
		}

		for roleType, meta := range rolesData {
			role := models.Role{
				ID:          roleType,
				Description: meta.Description,
				Permissions: meta.Permissions,
			}
			err := db.Insert(&role)

			if err != nil {
				return err
			}
		}

		return nil
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(table.Drop())

		return err
	}

	migrations.Register("20200828154040_create_roles_table", up, down, migrations.MigrationOptions{})
}
