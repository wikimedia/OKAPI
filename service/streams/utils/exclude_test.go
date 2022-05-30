package utils

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExclude(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		DbName string `json:"dbname"`
	}{}

	assert.NoError(json.NewDecoder(strings.NewReader(exclude)).Decode(&data))

	for _, info := range data {
		assert.True(Exclude(info.DbName))
	}
}
