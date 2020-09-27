package main

import (
	"okapi/lib/db"
	"okapi/models/roles"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	roleIdColumn := db.Column{
		TableName: "users",
		Name:      "role_id",
		Type:      "varchar(255) not null default '" + string(roles.Client) + "'",
	}
	usernameColumn := db.Column{
		TableName: "users",
		Name:      "username",
		Type:      "varchar(255) unique",
	}
	foreignKey := db.ForeignKey{
		ParentTable: "roles",
		TableName:   "users",
		Name:        "roles_users_fk",
		Property:    "role_id",
		References:  "id",
		OnDelete:    db.Cascade,
	}

	up := func(db orm.DB) error {
		_, err := db.Exec(
			roleIdColumn.Add() + usernameColumn.Add() + foreignKey.Create(),
		)

		if err != nil {
			return err
		}

		_, err = db.Exec("update users SET username = lower(substr(email,1,position('@' in email) - 1))")

		if err != nil {
			return err
		}

		_, err = db.Exec(usernameColumn.Set("not null"))

		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(foreignKey.Drop() + roleIdColumn.Drop() + usernameColumn.Drop())

		return err
	}

	migrations.Register("20200829220147_add_role_to_users", up, down, migrations.MigrationOptions{})
}
