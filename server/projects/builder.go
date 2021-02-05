package projects

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/mediawiki-api-client"
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

// MWiki set mediawiki url
func (bu *Builder) MWiki(mWiki *mediawiki.Client) *Builder {
	bu.srv.mWiki = mWiki
	return bu
}

// Repository assign storage repository to the server
func (bu *Builder) Repository(repo repository.Repository) *Builder {
	bu.srv.repo = repo
	return bu
}

// Elastic set elasticsearch client
func (bu *Builder) Elastic(elastic *elasticsearch.Client) *Builder {
	bu.srv.elastic = elastic
	return bu
}

// Build create new server instance
func (bu *Builder) Build() *Server {
	return bu.srv
}
