package diffs

import (
	"okapi-public-api/lib/auth"
	"okapi-public-api/lib/aws"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffs(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(aws.Init())
	assert.NoError(auth.Init())
	assert.NotZero(len(Init().Routes))
}
