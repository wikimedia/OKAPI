package models

import (
	"context"

	"github.com/go-pg/pg/v10"
)

// Project database table representation
type Project struct {
	ID       int       `pg:",pk" json:"id"`
	DbName   string    `pg:"type:varchar(255),unique,notnull" json:"db_name"`
	SiteName string    `pg:"type:varchar(255),notnull" json:"site_name"`
	SiteCode string    `pg:"type:varchar(255),notnull" json:"site_code"`
	SiteURL  string    `pg:"type:varchar(255),notnull" json:"site_url"`
	Lang     string    `pg:"type:varchar(25),notnull" json:"lang"`
	Active   bool      `pg:",use_zero,notnull" json:"active"`
	Language *Language `pg:"rel:has-one" json:"language,omitempty"`
	timestamp
}

var _ pg.BeforeUpdateHook = (*Project)(nil)

// BeforeUpdate model hook
func (proj *Project) BeforeUpdate(ctx context.Context) (context.Context, error) {
	proj.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Project)(nil)

// BeforeInsert model hook
func (proj *Project) BeforeInsert(ctx context.Context) (context.Context, error) {
	proj.OnInsert()
	return ctx, nil
}
