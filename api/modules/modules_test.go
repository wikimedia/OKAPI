package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModules(t *testing.T) {
	assert.NotZero(t, len(Init()))
}
