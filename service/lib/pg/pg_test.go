package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPG(t *testing.T) {
	defer func() { conn = nil }()
	assert.NoError(t, Init())
	assert.NotNil(t, conn)
	assert.NotNil(t, Conn())
	assert.Equal(t, ErrDuplicateConn, Init())
}
