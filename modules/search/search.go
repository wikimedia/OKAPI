package search

import (
	"net/http"

	"okapi/lib/module"
	"okapi/modules/search/routes"

	"github.com/gin-gonic/gin"
)

// Module instance of example module
var Module = module.Module{
	Path:       "/search",
	Middleware: []func() gin.HandlerFunc{},
	Routes: []module.Route{
		{
			Path:    "/pages",
			Method:  http.MethodPost,
			Handler: routes.Pages,
		},
		{
			Path:    "/projects",
			Method:  http.MethodPost,
			Handler: routes.Projects,
		},
		{
			Path:    "/options/:lang",
			Method:  http.MethodGet,
			Handler: routes.Options,
		},
	},
}
