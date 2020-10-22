package captcha

import (
	"net/http"
	"okapi/modules/captcha/routes"

	"okapi/lib/module"
)

// Module instance of captcha module
var Module = module.Module{
	Path: "/captcha",
	Routes: []module.Route{
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: routes.Get,
		},
		{
			Path:    "/:id",
			Method:  http.MethodGet,
			Handler: routes.Show,
		},
	},
}
