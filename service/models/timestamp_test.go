package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type timestampTestModel struct {
	timestamp
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestTimestampOnUpdate(t *testing.T) {
	model := timestampTestModel{
		ID:   10,
		Name: "time",
	}
	createdAt := model.CreatedAt
	updatedAt := model.UpdatedAt

	model.OnUpdate()
	assert.Equal(t, createdAt, model.CreatedAt)
	assert.NotEqual(t, updatedAt, model.UpdatedAt)
}

func TestTimestampOnInsert(t *testing.T) {
	model := timestampTestModel{
		ID:   10,
		Name: "time",
	}
	createdAt := model.CreatedAt
	updatedAt := model.UpdatedAt

	model.OnInsert()
	assert.NotEqual(t, createdAt, model.CreatedAt)
	assert.NotEqual(t, updatedAt, model.UpdatedAt)
}
