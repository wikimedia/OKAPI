package models

import (
	"context"

	"github.com/go-pg/pg/v9"
)

// Page struct for "pages" table representation
type Page struct {
	baseModel
	Title     string   `pg:"type:varchar(750),notnull" json:"title"`
	ProjectID int      `pg:",notnull" json:"project_id"`
	Revision  int      `pg:",notnull" json:"revision"`
	Path      string   `pg:"type:varchar(1000)" json:"path"`
	TID       string   `pg:"type:varchar(255),notnull" json:"tid"`
	Lang      string   `pg:"type:varchar(3),notnull" json:"lang"`
	SiteURL   string   `pg:"type:varchar(255),notnull" json:"site_url"`
	Project   *Project `json:"project"`
}

var _ pg.BeforeUpdateHook = (*Page)(nil)

// BeforeUpdate model hook
func (page *Page) BeforeUpdate(ctx context.Context) (context.Context, error) {
	page.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Page)(nil)

// BeforeInsert model hook
func (page *Page) BeforeInsert(ctx context.Context) (context.Context, error) {
	page.OnInsert()
	return ctx, nil
}
