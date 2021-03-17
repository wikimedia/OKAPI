package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	defer func() { ses = nil }()
	assert := assert.New(t)
	assert.NoError(Init())
	assert.NotNil(ses)
	assert.NotNil(Session())
	assert.Equal(ErrDuplicateSession, Init())
}
