package search

import (
	"github.com/go-pg/pg/v10/orm"
)

// Field it's field that will be used to filter
type Field string

// Search base struct
type Search struct {
	Model    interface{}
	Filters  *map[Field]func(query *orm.Query)
	Params   *map[Field]interface{}
	Query    *orm.Query
	Response *Response
	Request  *Request
}

// Run run the search
func (search *Search) Run() *Response {
	search.Response.Limit = search.Request.Limit
	search.Response.Page = search.Request.Page
	search.Response.Items = search.Model

	if search.Response.Limit <= 0 {
		search.Response.Limit = 100
	}

	filters := *search.Filters
	for name, value := range *search.Params {
		if value != nil {
			if filter, exists := filters[name]; exists {
				filter(search.Query)
			}
		}
	}

	sort := search.Request.Sort
	order := search.Request.Order
	_, hasField := filters[Field(sort)]
	if (sort == "id" || hasField) && (order == ASC || order == DESC) {
		search.Query.Order(sort + " " + string(order))
	}

	search.Response.Total, _ = search.Query.
		Limit(search.Response.Limit).
		Offset(search.Response.Limit * search.Response.Page).
		SelectAndCountEstimate(100000)
	search.Response.Pages = search.Response.Total / search.Response.Limit

	return search.Response
}
