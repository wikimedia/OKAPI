package fetch

import (
	"okapi-data-service/pkg/page"

	"github.com/protsack-stephan/mediawiki-api-client"
)

// Factory for fetch worker
type Factory struct{}

// Create create new fetch worker
func (f Factory) Create(fact *page.Factory, store Storage, mwiki *mediawiki.Client, repo Repo) Fetcher {
	return &Worker{
		fact,
		store,
		mwiki,
		repo,
	}
}
