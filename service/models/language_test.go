package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageBeforeInsert(t *testing.T) {
	lang := new(Language)
	createdAt := lang.CreatedAt
	updatedAt := lang.UpdatedAt

	_, err := lang.BeforeInsert(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, createdAt, lang.CreatedAt)
	assert.NotEqual(t, updatedAt, lang.UpdatedAt)
}

func TestLanguagesBeforeUpdate(t *testing.T) {
	lang := new(Language)
	createdAt := lang.CreatedAt
	updatedAt := lang.UpdatedAt

	_, err := lang.BeforeUpdate(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, createdAt, lang.CreatedAt)
	assert.NotEqual(t, updatedAt, lang.UpdatedAt)
}
