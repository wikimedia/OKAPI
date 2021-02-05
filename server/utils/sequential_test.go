package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const seqTestKey = "test/sequence"

var seqTestSleep = time.Millisecond * 100
var seqTestWait = time.Millisecond * 50

func TestSequential(t *testing.T) {
	seq := new(Sequential)
	callback := func() error { return nil }

	go func() {
		_ = seq.Once(seqTestKey, func() error {
			time.Sleep(seqTestSleep)
			return nil
		})
	}()

	time.Sleep(seqTestWait)
	assert.Error(t, seq.Once(seqTestKey, callback))
	assert.Error(t, seq.Once(seqTestKey, callback))

	time.Sleep(seqTestWait)
	assert.NoError(t, seq.Once(seqTestKey, callback))
}
