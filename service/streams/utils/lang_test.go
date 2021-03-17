package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLang(t *testing.T) {
	assert := assert.New(t)

	for _, testCase := range []struct {
		data   string
		result string
	}{
		{
			"en.wikipedia.org",
			"en",
		},
		{
			"af",
			"af",
		},
		{
			"af.wikibooks.org",
			"af",
		},
	} {
		assert.Equal(testCase.result, Lang(testCase.data))
	}
}
