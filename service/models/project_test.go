package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectBeforeInsert(t *testing.T) {
	project := new(Project)
	createdAt := project.CreatedAt
	updatedAt := project.UpdatedAt

	_, err := project.BeforeInsert(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, createdAt, project.CreatedAt)
	assert.NotEqual(t, updatedAt, project.UpdatedAt)
}

func TestProjectsBeforeUpdate(t *testing.T) {
	project := new(Project)
	createdAt := project.CreatedAt
	updatedAt := project.UpdatedAt

	_, err := project.BeforeUpdate(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, createdAt, project.CreatedAt)
	assert.NotEqual(t, updatedAt, project.UpdatedAt)
}
