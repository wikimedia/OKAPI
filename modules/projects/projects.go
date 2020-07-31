package projects

import (
	"net/http"

	"okapi/middleware"

	"github.com/gin-gonic/gin"
	"okapi/lib/module"
	"okapi/modules/projects/routes"
)

// Module projects module instance
var Module = module.Module{
	Path: "/projects",
	Middleware: []func() gin.HandlerFunc{
		middleware.JWT().MiddlewareFunc,
	},
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
