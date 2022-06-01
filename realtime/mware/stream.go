package mware

import "github.com/gin-gonic/gin"

// Stream set streaming headers
func Stream() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream; charset=utf-8")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Next()
	}
}
