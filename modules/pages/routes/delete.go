package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete one entity
func Delete(c *gin.Context) {
	c.String(http.StatusOK, "Deleted a page with id \""+c.Param("id")+"\"!")
}
