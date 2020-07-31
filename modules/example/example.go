package example

import (
	"net/http"

	"okapi/lib/module"
	"okapi/modules/example/routes"

	"github.com/gin-gonic/gin"
)

// Module instance of example module
var Module = module.Module{
	Path:       "/example",
	Middleware: []func() gin.HandlerFunc{},
	Routes: []module.Route{
		{
			Path:    "",
			Method:  http.MethodPost,
			Handler: routes.Create,
		},
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: routes.List,
		},
		{
			Path:    "/:id",
			Method:  http.MethodGet,
			Handler: routes.View,
		},
		{
			Path:    "/:id",
			Method:  http.MethodDelete,
			Handler: routes.Delete,
		},
		{
			Path:    "/:id",
			Method:  http.MethodPut,
			Handler: routes.Update,
		},
	},
}
