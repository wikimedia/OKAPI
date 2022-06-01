package exports

import (
	"net/http"
	"okapi-public-api/lib/aws"
	"okapi-public-api/lib/env"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

// Init for projects endpoints
func Init() httpmod.Module {
	store := s3.NewStorage(aws.Session(), env.AWSBucket)

	return httpmod.Module{
		Path:       "/v1/exports",
		Middleware: []gin.HandlerFunc{},
		Routes: []httpmod.Route{
			// Deprecated enpoint
			{
				Path:    "/json/:project",
				Method:  http.MethodGet,
				Handler: JSON(store),
			},
			{
				Path:    "/download/:namespace/:project",
				Method:  http.MethodGet,
				Handler: JSONNS(store),
			},
			{
				Path:    "/meta/:namespace",
				Method:  http.MethodGet,
				Handler: List(store, env.Group),
			},
			{
				Path:    "/meta/:namespace/:project",
				Method:  http.MethodGet,
				Handler: Detail(store, env.Group),
			},
		},
	}
}
