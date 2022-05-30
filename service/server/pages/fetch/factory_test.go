package fetch

import (
	"okapi-data-service/pkg/page"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	assert := assert.New(t)
	pfact := new(page.Factory)
	store := new(storage.Mock)
	mwiki := new(mediawiki.Client)
	repo := new(repository.Mock)

	factory := new(Factory)
	fetcher := factory.Create(pfact, store, mwiki, repo).(*Worker)
	assert.NotNil(fetcher)
	assert.Equal(fetcher.fact, pfact)
	assert.Equal(fetcher.store, store)
	assert.Equal(fetcher.mwiki, mwiki)
	assert.Equal(fetcher.repo, repo)
}
