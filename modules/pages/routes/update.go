package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update one entity example
func Update(c *gin.Context) {
	c.String(http.StatusOK, "Updated a page with id \""+c.Param("id")+"\"!")
}
