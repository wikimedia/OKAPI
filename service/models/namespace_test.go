package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespaceBeforeInsert(t *testing.T) {
	ns := new(Namespace)
	createdAt := ns.CreatedAt
	updatedAt := ns.UpdatedAt

	_, err := ns.BeforeInsert(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, createdAt, ns.CreatedAt)
	assert.NotEqual(t, updatedAt, ns.UpdatedAt)
}

func TestNamespacesBeforeUpdate(t *testing.T) {
	ns := new(Namespace)
	createdAt := ns.CreatedAt
	updatedAt := ns.UpdatedAt

	_, err := ns.BeforeUpdate(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, createdAt, ns.CreatedAt)
	assert.NotEqual(t, updatedAt, ns.UpdatedAt)
}
