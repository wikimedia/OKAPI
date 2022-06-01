package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const daysmapTestLength = 10

func TestDaysMap(t *testing.T) {
	assert := assert.New(t)
	days := DaysMap(daysmapTestLength)

	for i := 0; i < daysmapTestLength; i++ {
		assert.Contains(days, time.Now().Add(time.Duration(i)*-24*time.Hour).UTC().Format(DateFormat))
	}

	assert.Len(days, daysmapTestLength)
}
