package exports

import (
	"net/http"
	"okapi/modules/exports/routes"

	"okapi/lib/module"
	"okapi/middleware"

	"github.com/gin-gonic/gin"
)

// Module projects module instance
var Module = module.Module{
	Path: "/exports",
	Middleware: []func() gin.HandlerFunc{
		middleware.JWT().MiddlewareFunc,
	},
	Routes: []module.Route{
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: routes.List,
		},
		{
			Path:    "/:resource_type/delete",
			Method:  http.MethodPut,
			Handler: routes.Delete,
		},
		{
			Path:    "/:resource_type/:id",
			Method:  http.MethodPost,
			Handler: routes.Create,
		},
	},
}
