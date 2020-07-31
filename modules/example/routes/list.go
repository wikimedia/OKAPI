package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// List list all entities example
func List(c *gin.Context) {
	c.String(http.StatusOK, "List all examples!")
}
