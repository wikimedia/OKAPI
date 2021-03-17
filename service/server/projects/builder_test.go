package projects

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
)

var builderTestElastic = &elasticsearch.Client{}
var builderTestMWiki = &mediawiki.Client{}
var builderTestRepo = &repository.Mock{}

func TestBuilder(t *testing.T) {
	srv := NewBuilder().
		MWiki(builderTestMWiki).
		Repository(builderTestRepo).
		Elastic(builderTestElastic).
		Build()

	assert.NotNil(t, srv)
	assert.Equal(t, builderTestMWiki, srv.mWiki)
	assert.Equal(t, builderTestElastic, srv.elastic)
	assert.Equal(t, builderTestRepo, srv.repo)
}
