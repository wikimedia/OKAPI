package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV1(t *testing.T) {
	assert.NotZero(t, len(Init()))
}
