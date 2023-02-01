package main

import (
	"context"
	"fmt"
	"log"
	"okapi-public-api/lib/auth"
	"okapi-public-api/lib/aws"
	"okapi-public-api/lib/env"
	"okapi-public-api/lib/redis"
	"os"
	"time"

	"okapi-public-api/modules"

	cog "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/gin-toolkit/httphandler"
	"github.com/protsack-stephan/gin-toolkit/httpmod"
	"github.com/protsack-stephan/gin-toolkit/httpmw"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "okapi-public-api/docs"

	ginswagger "github.com/protsack-stephan/gin-swagger"

	"github.com/casbin/casbin/v2"
)

var startup = []func() error{
	env.Init,
	redis.Init,
	aws.Init,
	auth.Init,
}

// @title Wikimedia Enterprise API
// @version 1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	for _, init := range startup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	gin.SetMode(env.APIMode)
	_ = os.Setenv("TZ", "UTC")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(httpmw.LogFormatter))
	router.Use(httpmw.CORS())
	router.NoRoute(httpmw.NotFound())

	router.GET("/v1/docs/*any", httpmw.Limit(5), ginswagger.WrapHandler(swaggerFiles.Handler, func(c *ginswagger.Config) {
		c.Title = "Wikimedia Enterprise API"
	}))

	cmd := redis.Client()

	router.GET("/v1/status/", httpmw.Limit(5), httphandler.Status(map[string]httphandler.StatusCheck{
		"redis": func(ctx context.Context) error {
			return cmd.Ping(ctx).Err()
		},
	}))

	enf, err := casbin.NewEnforcer(env.AccessModelPath, env.AccessPolicyPath)

	if err != nil {
		log.Panic(err)
	}

	user := new(httpmw.CognitoUser)
	user.SetUsername(env.IpCognitoUsername)
	user.SetGroups([]string{env.IpCognitoUsergroup})

	router.Use(httpmw.IpCognitoAuth(&httpmw.IpCognitoParams{
		Srv:      cog.New(auth.Session()),
		Cache:    cmd,
		ClientID: env.CognitoClientID,
		IpRange:  env.IpRange,
		Expire:   time.Minute * 5,
		User:     user,
	}))
	router.Use(httpmw.RBAC(httpmw.CasbinRBACAuthorizer(enf)))

	for group, limit := range env.QPSLimitPerGroup {
		router.Use(httpmw.LimitPerUser(cmd, limit, fmt.Sprintf("qps_%s", group), time.Second*1, group))
	}

	if len(env.IpRange) > 0 {
		router.Use(httpmw.Limit(env.IpRangeRequestsLimit, env.IpRange))
	}

	if err := httpmod.Init(router, modules.Init()); err != nil {
		log.Panic(err)
	}

	if err := router.Run(fmt.Sprintf(":%s", env.APIPort)); err != nil {
		log.Panic(err)
	}
}
