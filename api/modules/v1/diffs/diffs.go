package diffs

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
		Path:       "/v1/diffs",
		Middleware: []gin.HandlerFunc{},
		Routes: []httpmod.Route{
			{
				Path:    "/json/:date/:project", // deprecated.
				Method:  http.MethodGet,
				Handler: JSON(store),
			},
			{
				Path:    "/download/:date/:namespace/:project",
				Method:  http.MethodGet,
				Handler: JSONNS(store),
			},
			{
				Path:    "/download/:date/:namespace/:project",
				Method:  http.MethodHead,
				Handler: Head(store),
			},
			{
				Path:    "/meta/:date/:namespace",
				Method:  http.MethodGet,
				Handler: List(store),
			},
			{
				Path:    "/meta/:date/:namespace/:project",
				Method:  http.MethodGet,
				Handler: Detail(store),
			},
		},
	}
}
