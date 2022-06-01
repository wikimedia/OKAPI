package diffs

import (
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Builder create new server with custom params
type Builder struct {
	srv *Server
}

// NewBuilder initialize new server builder
func NewBuilder() *Builder {
	return &Builder{
		new(Server),
	}
}

// LocalStorage set new local file storage for the server
func (bu *Builder) LocalStorage(store storage.Storage) *Builder {
	bu.srv.localStore = store
	return bu
}

// RemoteStorage set new remote storage
func (bu *Builder) RemoteStorage(store storage.Storage) *Builder {
	bu.srv.remoteStore = store
	return bu
}

// Build create new server instance with custom params
func (bu *Builder) Build() *Server {
	return bu.srv
}
