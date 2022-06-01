package namespaces

import (
	"net/http"
	"okapi-public-api/pkg/namespaces"
	"okapi-public-api/schema/v3"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/gin-toolkit/httperr"
)

// List http handler
// @Summary Returns list of available namespaces
// @Tags namespaces
// @Description Includes name and identifier.
// @ID v1-namespace-list
// @Security ApiKeyAuth
// @Success 200 {object} []schema.Namespace
// @Failure 500 {object} httperr.Error
// @Router /v1/namespaces [get]
func List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := []schema.Namespace{}

		for id, name := range namespaces.Supported {
			id, err := strconv.Atoi(id)

			if err != nil {
				httperr.InternalServerError(c, err.Error())
				return
			}

			ns = append(ns, schema.Namespace{Name: name, Identifier: id})
		}

		c.JSON(http.StatusOK, ns)
	}
}
