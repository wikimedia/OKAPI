package models

import (
	"context"
	"okapi/models/permissions"
	"okapi/models/roles"

	"github.com/go-pg/pg/v9"
)

// Role struct for "roles" table representation
type Role struct {
	baseModel
	ID          roles.Type         `pg:",pk" json:"id"`
	Description string             `pg:"type:varchar(255),notnull" json:"description"`
	Permissions []permissions.Type `pg:"type:varchar(255)[],notnull" json:"permissions"`
}

// IsUpdate check if it's update or insert
func (role *Role) IsUpdate() bool {
	return role.ID != ""
}

var _ pg.BeforeUpdateHook = (*Role)(nil)

// BeforeUpdate model hook
func (role *Role) BeforeUpdate(ctx context.Context) (context.Context, error) {
	role.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Role)(nil)

// BeforeInsert model hook
func (role *Role) BeforeInsert(ctx context.Context) (context.Context, error) {
	role.OnInsert()
	return ctx, nil
}
