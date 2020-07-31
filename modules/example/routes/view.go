package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// View view one entity example
func View(c *gin.Context) {
	c.String(http.StatusOK, "View an example with id \""+c.Param("id")+"\"!")
}
