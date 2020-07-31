package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete one entity
func Delete(c *gin.Context) {
	c.String(http.StatusNoContent, "Deleted an example with id \""+c.Param("id")+"\"!")
}
