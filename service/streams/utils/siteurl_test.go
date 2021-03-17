package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteURL(t *testing.T) {
	assert := assert.New(t)

	for _, testCase := range []struct {
		data   string
		result string
	}{
		{
			"en.wikipedia.org",
			"https://en.wikipedia.org",
		},
		{
			"af.wikibooks.org",
			"https://af.wikibooks.org",
		},
	} {
		assert.Equal(testCase.result, SiteURL(testCase.data))
	}
}
