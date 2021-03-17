package namespaces

import (
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"

	"github.com/stretchr/testify/assert"
)

var builderTestRepo = &repository.Mock{}

func TestBuilder(t *testing.T) {
	srv := NewBuilder().
		Repository(builderTestRepo).
		Build()

	assert.NotNil(t, srv)
	assert.Equal(t, srv.repo, builderTestRepo)
}
