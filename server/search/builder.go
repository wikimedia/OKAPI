package search

import (
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// NewBuilder create new instance of server by custom params
func NewBuilder() *Builder {
	return &Builder{
		new(Server),
	}
}

// Builder create new namespaces server
type Builder struct {
	srv *Server
}

// Repository set data repository for server
func (bu *Builder) Repository(repo repository.Repository) *Builder {
	bu.srv.repo = repo
	return bu
}

// Storage set new storage for the server
func (bu *Builder) Storage(store storage.Storage) *Builder {
	bu.srv.store = store
	return bu
}

// Build create new server instance
func (bu *Builder) Build() *Server {
	return bu.srv
}
