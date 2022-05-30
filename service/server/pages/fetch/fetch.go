package fetch

import (
	"context"
	"okapi-data-service/pkg/page"
	"okapi-data-service/schema/v3"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/mediawiki-api-client"
)

// Repo fetch repository
type Repo interface {
	repository.Updater
	repository.Finder
	repository.Creator
}

// Storage interface for communication with storage
type Storage interface {
	storage.Putter
	storage.Deleter
}

// FetcherFactory interface for creating new fetch bulk processor
type FetcherFactory interface {
	Create(*page.Factory, Storage, *mediawiki.Client, Repo) Fetcher
}

// Fetcher interface for bulk page fetching worker
type Fetcher interface {
	Fetch(ctx context.Context, titles ...string) (map[string]*schema.Page, map[string]error, error)
}
