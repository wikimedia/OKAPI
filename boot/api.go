package boot

import (
	"github.com/gin-gonic/gin"
	"okapi/helpers/logger"
	"okapi/lib/env"
	"okapi/modules"
)

// API function to start API server
func API() {
	router := gin.New()

	if env.Context.APIMode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	modules.Init(router)
	err := router.Run(":" + env.Context.APIPort)
	if err != nil {
		logger.SYSTEM.Panic(logger.Message{
			ShortMessage: "API: api failed to start",
			FullMessage:  err.Error(),
		})
	}
}
