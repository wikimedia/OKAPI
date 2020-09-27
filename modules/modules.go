package modules

import (
	"okapi/lib/module"
	"okapi/middleware"
	"okapi/modules/auth"
	"okapi/modules/example"
	"okapi/modules/pages"
	"okapi/modules/projects"
	"okapi/modules/search"

	"github.com/gin-gonic/gin"
)

var modules = []module.Module{
	example.Module,
	projects.Module,
	pages.Module,
	auth.Module,
	search.Module,
}

// Init initialize all modules
func Init(router *gin.Engine) {
	router.Use(middleware.CORS())
	router.Use(middleware.Log())

	for _, module := range modules {
		group := router.Group(module.Path)

		for _, middleware := range module.Middleware {
			group.Use(middleware())
		}

		for _, route := range module.Routes {
			handlers := make([]gin.HandlerFunc, 0)

			if route.Middleware != nil {
				for _, routeMiddleware := range route.Middleware {
					handlers = append(handlers, routeMiddleware())
				}
			}

			group.Handle(route.Method, route.Path, append(handlers, route.Handler)...)
		}
	}
}
