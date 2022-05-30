package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/protsack-stephan/mediawiki-api-client"
)

func NewPage(title string, data *mediawiki.PageData, proj *Project, prev *Page) *Page {
	page := &Page{
		Title:      title,
		QID:        data.Pageprops.WikibaseItem,
		PID:        data.PageID,
		NsID:       data.Ns,
		Lang:       proj.Lang,
		Revision:   data.LastRevID,
		RevisionDt: time.Now(),
		DbName:     proj.DbName,
		SiteURL:    proj.SiteURL,
	}

	if len(data.Revisions) > 0 {
		page.RevisionDt = data.Revisions[0].Timestamp
	}

	if prev != nil {
		page.ID = prev.ID
		page.CreatedAt = prev.CreatedAt
		page.UpdatedAt = prev.UpdatedAt
	}

	return page
}

// Page database table representation
type Page struct {
	ID         int        `json:"id"`
	Title      string     `pg:"type:varchar(750),notnull" json:"title"`
	NsID       int        `pg:",use_zero" json:"ns_id"`
	QID        string     `pg:"type:varchar(500)"`
	PID        int        `pg:"type:bigint" json:"pid"`
	Revision   int        `pg:",use_zero" json:"revision"`
	RevisionDt time.Time  `pg:"type:timestamp" json:"revision_dt"`
	Failed     bool       `pg:",use_zero" json:"failed"`
	Lang       string     `pg:"type:varchar(25),notnull" json:"lang"`
	DbName     string     `pg:"type:varchar(255),notnull" json:"db_name"`
	SiteURL    string     `pg:"type:varchar(1000),notnull" json:"site_url"`
	Path       string     `pg:"type:varchar(1000)" json:"path"`
	Language   *Language  `pg:"rel:has-one" json:"language,omitempty"`
	Project    *Project   `pg:"rel:has-one" json:"project,omitempty"`
	Namespace  *Namespace `pg:"rel:has-one" json:"namespace,omitempty"`
	timestamp
}

// SetRevision set new default revision
func (page *Page) SetRevision(rev int, dt time.Time) {
	page.Revision = rev
	page.RevisionDt = dt
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
