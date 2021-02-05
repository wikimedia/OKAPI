package search

import (
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/stretchr/testify/assert"
)

var builderTestRepo = &repository.Mock{}
var builderTestStore = &storage.Mock{}

func TestBuilder(t *testing.T) {
	srv := NewBuilder().
		Repository(builderTestRepo).
		Storage(builderTestStore).
		Build()

	assert.NotNil(t, srv)
	assert.Equal(t, srv.repo, builderTestRepo)
	assert.Equal(t, srv.store, builderTestStore)
}
