package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create project create handler
func Create(c *gin.Context) {
	c.String(http.StatusOK, "Create a project!")
}
