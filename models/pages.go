package models

import (
	"context"
	page_index "okapi/indexes/page"
	"okapi/lib/ores"

	"github.com/go-pg/pg/v10"
)

// Page struct for "pages" table representation
type Page struct {
	baseModel
	Title        string      `pg:"type:varchar(750),notnull" json:"title"`
	ProjectID    int         `pg:",notnull" json:"project_id"`
	NsID         int         `pg:",use_zero" json:"ns_id"`
	Path         string      `pg:"type:varchar(1000)" json:"path"`
	WikitextPath string      `pg:"type:varchar(1000)" json:"wikitext_path"`
	SiteURL      string      `pg:"type:varchar(255),notnull" json:"site_url"`
	Revision     int         `pg:",use_zero" json:"revision"`
	Revisions    [6]int      `pg:",array,notnull" json:"revisions"`
	Updates      int         `pg:",use_zero" json:"updates"`
	Scores       ores.Scores `pg:"type:jsonb,notnull" json:"scores,omitempty"`
	Project      *Project    `pg:"rel:has-one" json:"project"`
}

// SetRevision add new revision to the list
func (page *Page) SetRevision(rev int) {
	if page.Revisions[0] != rev {
		page.Revision = rev
		revs := page.Revisions

		for i := 1; i < len(page.Revisions); i++ {
			page.Revisions[i] = revs[i-1]
		}

		page.Revisions[0] = rev
		page.Updates++
	}
}

// SetScore add new revision score to the list
func (page *Page) SetScore(rev int, model ores.Model, score ores.Score) {
	newScore := make(ores.Scores)

	if page.Scores == nil {
		page.Scores = make(ores.Scores)
	}

	if _, exists := page.Scores[rev]; !exists {
		page.Scores[rev] = map[ores.Model]ores.Score{}
	}

	page.Scores[rev][model] = score

	for _, revision := range page.Revisions {
		if revision > 0 {
			if _, exists := page.Scores[revision]; exists {
				newScore[revision] = page.Scores[revision]
			}
		}
	}

	page.Scores = newScore
}

// Index get indexed data structure
func (page *Page) Index() *page_index.Index {
	index := &page_index.Index{
		ID:        page.ID,
		Title:     page.Title,
		NsID:      page.NsID,
		SiteURL:   page.SiteURL,
		UpdatedAt: page.UpdatedAt,
	}

	if page.Project == nil || page.Project.ID <= 0 {
		page.Project = &Project{}
		db.Model(page.Project).
			Where("id = ?", page.ProjectID).
			Select()
	}

	index.SiteCode = page.Project.SiteCode
	index.Lang = page.Project.Lang
	index.LangName = page.Project.LangName
	index.LangLocalName = page.Project.LangLocalName
	index.ProjectID = page.ProjectID
	return index
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

var _ pg.AfterInsertHook = (*Page)(nil)

// AfterInsert model hook
func (page *Page) AfterInsert(ctx context.Context) error {
	page.Index().Update()
	return nil
}

var _ pg.AfterUpdateHook = (*Page)(nil)

// AfterUpdate model hook
func (page *Page) AfterUpdate(ctx context.Context) error {
	page.Index().Update()
	return nil
}

var _ pg.AfterDeleteHook = (*Page)(nil)

// AfterDelete model hook
func (page *Page) AfterDelete(ctx context.Context) error {
	page_index.Delete(page.ID)
	return nil
}
