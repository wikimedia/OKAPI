package dumps

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"okapi/lib/module"
	"okapi/middleware"
	"okapi/modules/dumps/routes"
)

// Module instance of dumps module
var Module = module.Module{
	Path: "/dumps",
	Middleware: []func() gin.HandlerFunc{
		middleware.JWT().MiddlewareFunc,
	},
	Routes: []module.Route{
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: routes.View,
		},
		{
			Path:    "/url",
			Method:  http.MethodGet,
			Handler: routes.URL,
		},
	},
}
