package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"okapi/helpers/logger"
)

// Log log middleware for gin
func Log() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		formatted := fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
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

		message := logger.Message{
			ShortMessage: fmt.Sprintf("Api request to '%s'", param.Path),
			FullMessage:  formatted,
			Params: map[string]interface{}{
				"_client_ip":     param.ClientIP,
				"_method":        param.Method,
				"_path":          param.Path,
				"_status":        param.StatusCode,
				"_latency":       param.Latency,
				"_user_agent":    param.Request.UserAgent(),
				"_error_message": param.ErrorMessage,
			},
		}

		if param.StatusCode >= http.StatusBadRequest {
			message.Level = logger.ERROR
		} else {
			message.Level = logger.INFO
		}

		logger.API.Send(message)
		return formatted
	})
}
