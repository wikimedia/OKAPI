package namespaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSupported(t *testing.T) {
	assert := assert.New(t)

	for _, ns := range []string{"0", "6", "14"} {
		assert.True(IsSupported(ns))
	}

	for _, ns := range []string{"1", "2", "3", "4", "5", "7", "8", "9", "10", "13", "15"} {
		assert.False(IsSupported(ns))
	}
}
