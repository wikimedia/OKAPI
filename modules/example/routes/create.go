package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create example create handler
func Create(c *gin.Context) {
	c.String(http.StatusCreated, "Create an example!")
}
