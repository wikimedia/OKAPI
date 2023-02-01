package projects

import (
	"fmt"
	"net/http"
	"okapi-public-api/lib/aws"
	"okapi-public-api/lib/env"
	"okapi-public-api/lib/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/gin-toolkit/httpmod"
	"github.com/protsack-stephan/gin-toolkit/httpmw"
)

// Init for projects endpoints
func Init() httpmod.Module {
	ses := aws.Session()
	expire := time.Second * time.Duration(env.ProjectsExpire)
	cmdable := redis.Client()
	store := s3.NewStorage(ses, env.AWSBucket)

	return httpmod.Module{
		Path:       "/v1/projects",
		Middleware: []gin.HandlerFunc{},
		Routes: []httpmod.Route{
			{
				Method: http.MethodGet,
				Handler: httpmw.Cache(&httpmw.CacheParams{
					Cache:       cmdable,
					Expire:      expire,
					Handle:      List(store),
					ContentType: fmt.Sprintf("%s; charset=UTF-8", gin.MIMEJSON),
				}),
			},
		},
	}
}
