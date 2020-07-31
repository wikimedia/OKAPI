package pages

import (
	"net/http"

	"okapi/lib/module"
	"okapi/middleware"
	"okapi/modules/pages/routes"

	"github.com/gin-gonic/gin"
)

// Module projects module instance
var Module = module.Module{
	Path: "/pages",
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
		{
			Path:    "/:id/download",
			Method:  http.MethodGet,
			Handler: routes.Download,
		},
	},
}
