package pages

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"github.com/stretchr/testify/assert"
)

var builderTestRemoteStore = new(storage.Mock)
var builderTestJSONStore = new(storage.Mock)
var builderTestGenStore = new(storage.Mock)
var builderTestRepo = new(repository.Mock)
var builderTestDumps = new(dumps.Client)
var builderTestElastic = new(elasticsearch.Client)

func TestBuilder(t *testing.T) {
	client := NewBuilder().
		RemoteStorage(builderTestRemoteStore).
		JSONStorage(builderTestJSONStore).
		GenStorage(builderTestGenStore).
		Repository(builderTestRepo).
		Dumps(builderTestDumps).
		Elastic(builderTestElastic).
		Build()

	assert := assert.New(t)
	assert.Equal(builderTestRemoteStore, client.remoteStore)
	assert.Equal(builderTestJSONStore, client.jsonStore)
	assert.Equal(builderTestGenStore, client.genStore)
	assert.Equal(builderTestRepo, client.repo)
	assert.Equal(builderTestDumps, client.dumps)
	assert.Equal(builderTestElastic, client.elastic)
}
