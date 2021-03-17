package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNs(t *testing.T) {
	assert := assert.New(t)
	assert.True(FilterNs(0))
	assert.False(FilterNs(1))
}
