package projects

import (
	"okapi-public-api/lib/auth"
	"okapi-public-api/lib/aws"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjects(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(auth.Init())
	assert.NoError(aws.Init())
	assert.NotZero(len(Init().Routes))
}
