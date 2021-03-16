package content

import (
	"context"
	"fmt"
	"okapi-data-service/models"
	"okapi-data-service/schema/v1"

	"github.com/protsack-stephan/mediawiki-api-client"
)

// Storer intrface wrapper for Storage
type Storer interface {
	Deleter
	Puller
}

// Deleter interface for delete method
type Deleter interface {
	Delete(ctx context.Context, page *models.Page) error
}

// Puller interface for pull method
type Puller interface {
	Pull(ctx context.Context, page *models.Page, mwiki *mediawiki.Client) (*schema.Page, error)
}

// License default license for the pages
const License = "CC BY-SA"

// NewStructured create new structured content instance
func NewStructured(page *models.Page) *schema.Page {
	return &schema.Page{
		Title: page.Title,
		URL: schema.PageURL{
			Canonical: fmt.Sprintf("%s/wiki/%s", page.SiteURL, page.Title),
		},
		License:      []string{License},
		DbName:       page.DbName,
		PID:          page.PID,
		QID:          page.QID,
		InLanguage:   page.Lang,
		Revision:     page.Revision,
		DateModified: page.RevisionDt,
	}
}
