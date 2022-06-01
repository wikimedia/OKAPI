package diffs

import (
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var builderTestRemoteStore = new(storage.Mock)
var builderTestLocalStore = new(storage.Mock)

func TestBuilder(t *testing.T) {
	client := NewBuilder().
		RemoteStorage(builderTestRemoteStore).
		LocalStorage(builderTestLocalStore).
		Build()

	assert := assert.New(t)
	assert.Equal(builderTestRemoteStore, client.remoteStore)
	assert.Equal(builderTestLocalStore, client.localStore)
}
