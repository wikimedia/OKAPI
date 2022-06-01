package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	defer func() { client = nil }()
	assert := assert.New(t)
	assert.NoError(Init())
	assert.NotNil(client)
	assert.NotNil(Client())
	assert.NoError(Init())
	assert.Equal(client, Client())
}
