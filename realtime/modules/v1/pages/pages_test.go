package pages

import (
	"okapi-streams/lib/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPages(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(auth.Init())
	assert.NotZero(len(Init().Routes))
}
