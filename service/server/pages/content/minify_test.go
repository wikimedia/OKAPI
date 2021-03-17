package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinify(t *testing.T) {
	for _, testCase := range []struct {
		input  string
		result string
	}{
		{
			`hello
			world`,
			"hello\t\t\tworld",
		},
		{
			"hello \nworld",
			"hello world",
		},
		{
			"hello world\n",
			"hello world",
		},
	} {
		assert.Equal(t, testCase.result, minify(testCase.input))
	}
}
