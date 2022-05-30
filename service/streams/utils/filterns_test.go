package utils

import (
	"okapi-data-service/schema/v3"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNs(t *testing.T) {
	assert := assert.New(t)

	for _, ns := range []int{schema.NamespaceArticle, schema.NamespaceFile, schema.NamespaceCategory, schema.NamespaceTemplate} {
		assert.True(FilterNs(ns))
	}

	for _, ns := range []int{1, 2, 3, 4, 5, 7, 8, 9, 11, 12, 13, 15, 16} {
		assert.False(FilterNs(ns))
	}
}
