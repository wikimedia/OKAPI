package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPageBeforeInsert(t *testing.T) {
	page := new(Page)
	createdAt := page.CreatedAt
	updatedAt := page.UpdatedAt

	_, err := page.BeforeInsert(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, createdAt, page.CreatedAt)
	assert.NotEqual(t, updatedAt, page.UpdatedAt)
}

func TestPagesBeforeUpdate(t *testing.T) {
	page := new(Page)
	createdAt := page.CreatedAt
	updatedAt := page.UpdatedAt

	_, err := page.BeforeUpdate(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, createdAt, page.CreatedAt)
	assert.NotEqual(t, updatedAt, page.UpdatedAt)
}

func TestPageSerRevision(t *testing.T) {
	page := new(Page)
	testCases := []struct {
		rev int
		dt  time.Time
	}{
		{
			1,
			time.Now().Add(1 * time.Hour),
		},
		{
			2,
			time.Now().Add(2 * time.Hour),
		},
		{
			3,
			time.Now().Add(3 * time.Hour),
		},
		{
			4,
			time.Now().Add(4 * time.Hour),
		},
		{
			5,
			time.Now().Add(5 * time.Hour),
		},
		{
			6,
			time.Now().Add(6 * time.Hour),
		},
		{
			7,
			time.Now().Add(7 * time.Hour),
		},
	}

	for _, testCase := range testCases {
		page.SetRevision(testCase.rev, testCase.dt)
		assert.Equal(t, testCase.rev, page.Revision)
		assert.Equal(t, testCase.rev, page.Revisions[0])
		assert.Equal(t, testCase.dt, page.RevisionDt)
	}

	for i := len(testCases) - 1; i > 0; i-- {
		assert.Equal(t, page.Revisions[revisions-i], testCases[i].rev)
	}
}
