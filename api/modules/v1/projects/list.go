package projects

import (
	"io/ioutil"
	"net/http"

	_ "okapi-public-api/schema/v3"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httperr"
)

// List http handler
// @Summary Returns list of all downloadable Wikimedia projects
// @Tags projects
// @Description Includes project identifier, file sizes, and other relevant metadata.
// @ID v1-projects-list
// @Security ApiKeyAuth
// @Success 200 {object} []schema.Project
// @Failure 404 {object} httperr.Error
// @Router /v1/projects [get]
func List(storage storage.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		rc, err := storage.Get("public/projects.json")

		if err != nil {
			httperr.NotFound(c, err.Error())
			return
		}

		defer rc.Close()
		data, err := ioutil.ReadAll(rc)

		if err != nil {
			httperr.InternalServerError(c, err.Error())
			return
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}
