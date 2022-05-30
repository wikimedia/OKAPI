package pages

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
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

// JSONStorage set new json storage for the server
func (bu *Builder) JSONStorage(store storage.Storage) *Builder {
	bu.srv.jsonStore = store
	return bu
}

// GenStorage set new general storage for the server
func (bu *Builder) GenStorage(store storage.Storage) *Builder {
	bu.srv.genStore = store
	return bu
}

// RemoteStorage set new remote storage
func (bu *Builder) RemoteStorage(store storage.Storage) *Builder {
	bu.srv.remoteStore = store
	return bu
}

// Repository set new storage repository for the server
func (bu *Builder) Repository(repo repository.Repository) *Builder {
	bu.srv.repo = repo
	return bu
}

// Dumps set new wikimedia dumps client
func (bu *Builder) Dumps(dumps *dumps.Client) *Builder {
	bu.srv.dumps = dumps
	return bu
}

// Elastic set elasticsearch client
func (bu *Builder) Elastic(elastic *elasticsearch.Client) *Builder {
	bu.srv.elastic = elastic
	return bu
}

// Build create new server instance with custom params
func (bu *Builder) Build() *Server {
	return bu.srv
}
