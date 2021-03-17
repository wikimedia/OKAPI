package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExclude(t *testing.T) {
	assert := assert.New(t)
	file, err := os.Open("../specials.json")
	assert.NoError(err)
	defer file.Close()

	data := []struct {
		DbName string `json:"dbname"`
	}{}
	assert.NoError(json.NewDecoder(file).Decode(&data))

	for _, info := range data {
		assert.True(Exclude(info.DbName))
	}
}
