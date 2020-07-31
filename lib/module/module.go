package module

import "github.com/gin-gonic/gin"

// Module struct to represent module
type Module struct {
	Path       string
	Middleware []func() gin.HandlerFunc
	Routes     []Route
}
