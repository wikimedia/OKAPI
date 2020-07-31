package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"okapi/middleware"

	"okapi/lib/module"
	"okapi/modules/auth/routes"
)

// Module instance of login module
var Module = module.Module{
	Path: "/auth",
	Routes: []module.Route{
		{
			Path:    "/login",
			Method:  http.MethodPost,
			Handler: middleware.JWT().LoginHandler,
		},
		{
			Path:   "/create",
			Method: http.MethodPost,
			Middleware: []func() gin.HandlerFunc{
				middleware.Auth,
			},
			Handler: routes.Create,
		},
	},
}
