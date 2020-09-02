package middleware

import (
	"fmt"
	"net/http"
	"time"

	"okapi/helpers/logger"

	"github.com/gin-gonic/gin"
)

// Log log middleware for gin
func Log() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		shortMessage := fmt.Sprintf("Api request to '%s'", param.Path)
		fullMessage := fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		info := map[string]interface{}{
			"_client_ip":     param.ClientIP,
			"_method":        param.Method,
			"_path":          param.Path,
			"_status":        param.StatusCode,
			"_latency":       param.Latency,
			"_user_agent":    param.Request.UserAgent(),
			"_error_message": param.ErrorMessage,
		}

		if param.StatusCode >= http.StatusBadRequest {
			logger.API.Error(shortMessage, fullMessage, info)
		} else {
			logger.API.Info(shortMessage, fullMessage, info)
		}

		return fullMessage
	})
}
