package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const formatTestDir = "2021-02-16"
const formatTestDbName = "enwiki"
const formatTestContentType = "html"
const formatTestTitle = "Earth"
const formatTestFileType = "html"

func TestFormat(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("page/%s/%s/%s/%s.%s", formatTestDir, formatTestDbName, formatTestContentType, formatTestTitle, formatTestFileType), Format(formatTestDir, formatTestDbName, formatTestContentType, formatTestTitle, formatTestFileType))

}
