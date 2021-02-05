package elastic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElastic(t *testing.T) {
	defer func() { client = nil }()
	assert := assert.New(t)
	assert.NoError(Init())
	assert.NotNil(Client())
	assert.NotNil(client)
	assert.Equal(ErrDuplicateClient, Init())
}
