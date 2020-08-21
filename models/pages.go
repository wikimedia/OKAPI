package models

import (
	"context"
	"okapi/lib/ores"

	"github.com/go-pg/pg/v9"
)

// Page struct for "pages" table representation
type Page struct {
	baseModel
	Title     string      `pg:"type:varchar(750),notnull" json:"title"`
	ProjectID int         `pg:",notnull" json:"project_id"`
	Path      string      `pg:"type:varchar(1000)" json:"path"`
	SiteURL   string      `pg:"type:varchar(255),notnull" json:"site_url"`
	Revision  int         `pg:",use_zero" json:"revision"`
	Revisions [6]int      `pg:",array,notnull" json:"revisions"`
	Updates   int         `pg:",use_zero" json:"updates"`
	Scores    ores.Scores `pg:"type:jsonb,notnull" json:"scores,omitempty"`
	Project   *Project    `json:"project"`
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
