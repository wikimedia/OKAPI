package projects

import (
	"net/http"
	"okapi/models/permissions"
	"okapi/models/roles"

	"okapi/middleware"

	"okapi/lib/module"
	"okapi/modules/projects/routes"

	"github.com/gin-gonic/gin"
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
			Middleware: []func() gin.HandlerFunc{
				middleware.Permissions(
					middleware.PermissionsList{permissions.ProjectCreate},
				),
			},
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
			Middleware: []func() gin.HandlerFunc{
				middleware.Permissions(
					middleware.PermissionsList{permissions.ProjectDelete},
				),
			},
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
			Middleware: []func() gin.HandlerFunc{
				middleware.DownloadRestrictions(map[roles.Type]int{
					roles.Client:     5,
					roles.Subscriber: -1,
					roles.Admin:      -1,
				}),
			},
		},
		{
			Path:    "/:id/bundle",
			Method:  http.MethodPost,
			Handler: routes.Bundle,
			Middleware: []func() gin.HandlerFunc{
				middleware.Permissions(
					middleware.PermissionsList{permissions.ProjectBundle},
				),
			},
		},
	},
}
