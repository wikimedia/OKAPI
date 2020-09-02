package boot

import (
	"okapi/helpers/logger"
	"okapi/lib/env"
	"okapi/modules"

	"github.com/gin-gonic/gin"
)

// API function to start API server
func API() {
	router := gin.New()

	if env.Context.APIMode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	modules.Init(router)

	if err := router.Run(":" + env.Context.APIPort); err != nil {
		logger.System.Panic("API: api failed to start", err.Error())
	}
}
