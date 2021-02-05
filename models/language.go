package models

import (
	"context"

	"github.com/go-pg/pg/v10"
)

// Language database table representation
type Language struct {
	ID        int    `pg:",pk" json:"id"`
	Code      string `pg:"type:varchar(25),unique,notnull" json:"code"`
	Dir       string `pg:"type:varchar(255),notnull" json:"dir"`
	Name      string `pg:"type:varchar(255),notnull" json:"name"`
	LocalName string `pg:"type:varchar(255),notnull" json:"local_name"`
	timestamp
}

var _ pg.BeforeUpdateHook = (*Page)(nil)

// BeforeUpdate model hook
func (lang *Language) BeforeUpdate(ctx context.Context) (context.Context, error) {
	lang.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Page)(nil)

// BeforeInsert model hook
func (lang *Language) BeforeInsert(ctx context.Context) (context.Context, error) {
	lang.OnInsert()
	return ctx, nil
}
