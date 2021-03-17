package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

const revisions = 6

// Page database table representation
type Page struct {
	ID           int            `json:"id"`
	Title        string         `pg:"type:varchar(750),notnull" json:"title"`
	NsID         int            `pg:",use_zero" json:"ns_id"`
	QID          string         `pg:"type:varchar(500)"`
	PID          int            `pg:"type:bigint" json:"pid"`
	Revision     int            `pg:",use_zero" json:"revision"`
	Revisions    [revisions]int `pg:",array,notnull" json:"revisions"`
	RevisionDt   time.Time      `pg:"type:timestamp" json:"revision_dt"`
	Lang         string         `pg:"type:varchar(25),notnull" json:"lang"`
	DbName       string         `pg:"type:varchar(255),notnull" json:"db_name"`
	SiteURL      string         `pg:"type:varchar(1000),notnull" json:"site_url"`
	HTMLPath     string         `pg:"type:varchar(1000)" json:"html_path"`
	WikitextPath string         `pg:"type:varchar(1000)" json:"wikitext_path"`
	JSONPath     string         `pg:"type:varchar(1000)" json:"json_path"`
	Language     *Language      `pg:"rel:has-one" json:"language,omitempty"`
	Project      *Project       `pg:"rel:has-one" json:"project,omitempty"`
	timestamp
}

// SetRevision set new default revision
func (page *Page) SetRevision(rev int, dt time.Time) {
	if page.Revision != rev && rev > 0 {
		page.Revision = rev
		page.RevisionDt = dt

		for i, rev := range append([]int{rev}, page.Revisions[:]...) {
			if i < revisions {
				page.Revisions[i] = rev
			}
		}
	}
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
