package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update one entity example
func Update(c *gin.Context) {
	c.String(http.StatusOK, "Updated an example with id \""+c.Param("id")+"\"!")
}
