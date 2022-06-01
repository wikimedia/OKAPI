package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/casbin/casbin/v2"
	ginswagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "okapi-streams/docs"
	"okapi-streams/lib/auth"
	"okapi-streams/lib/env"
	"okapi-streams/lib/redis"
	"okapi-streams/modules"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/gin-toolkit/httphandler"
	"github.com/protsack-stephan/gin-toolkit/httpmod"
	"github.com/protsack-stephan/gin-toolkit/httpmw"
)

const port = ":4040"

// @title Wikimedia Enterprise Realtime API
// @version 1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	setup := []func() error{
		env.Init,
		auth.Init,
		redis.Init,
	}

	for _, init := range setup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	gin.SetMode(env.APIMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(httpmw.LogFormatter))
	router.Use(httpmw.CORS())
	router.NoRoute(httpmw.NotFound())

	e, err := casbin.NewEnforcer(env.AccessModelPath, env.AccessPolicyPath)

	if err != nil {
		log.Panic(err)
	}

	router.GET("/v1/status/", httphandler.Status(map[string]httphandler.StatusCheck{}))
	router.GET("/v1/docs/*any", ginswagger.WrapHandler(swaggerFiles.Handler))
	router.Use(httpmw.IpCognitoAuth(
		cognitoidentityprovider.New(auth.Session()),
		redis.Client(),
		env.CognitoClientID,
		env.IpRange,
		time.Second*30,
	))
	router.Use(httpmw.RBAC(httpmw.CasbinRBACAuthorizer(e)))

	if err := httpmod.Init(router, modules.Init()); err != nil {
		log.Panic(err)
	}

	if err := router.Run(port); err != nil {
		log.Panic(err)
	}
}
