package models

import (
	"context"

	"github.com/go-pg/pg/v9"
)

// User struct for "users" table representation
type User struct {
	baseModel
	Email    string `pg:"type:varchar(255),unique,notnull" json:"email"`
	Password string `pg:"type:varchar(255),notnull" json:"-"`
}

var _ pg.BeforeUpdateHook = (*User)(nil)

// BeforeUpdate model hook
func (user *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	user.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*User)(nil)

// BeforeInsert model hook
func (user *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	user.OnInsert()
	return ctx, nil
}
