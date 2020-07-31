package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var methods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPost,
	http.MethodPut,
	http.MethodOptions,
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Content-Type", "application/json")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
