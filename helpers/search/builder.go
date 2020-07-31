package search

import (
	"github.com/go-pg/pg/v9/orm"
	"okapi/models"
)

// Builder search class builder
type Builder struct {
	Search *Search
}

// Model set search model
func (builder *Builder) Model(model interface{}) *Builder {
	if model != nil {
		builder.Search.Model = model
		builder.Search.Query = models.DB().Model(model)
	}

	return builder
}

// Filters set search filters
func (builder *Builder) Filters(filters *map[Field]func(query *orm.Query)) *Builder {
	if filters != nil {
		builder.Search.Filters = filters
	}

	return builder
}

// Params set search params
func (builder *Builder) Params(params *map[Field]interface{}) *Builder {
	if params != nil {
		builder.Search.Params = params
	}

	return builder
}

// Request set search request
func (builder *Builder) Request(request *Request) *Builder {
	if request != nil {
		builder.Search.Request = request
	}

	return builder
}

// Build create new search instance
func (builder *Builder) Build() *Search {
	return builder.Search
}

// New create new search builder
func New() *Builder {
	builder := new(Builder)
	builder.Search = new(Search)
	builder.Search.Filters = &map[Field]func(query *orm.Query){}
	builder.Search.Params = &map[Field]interface{}{}
	builder.Search.Request = &Request{}
	builder.Search.Response = &Response{}
	return builder
}
