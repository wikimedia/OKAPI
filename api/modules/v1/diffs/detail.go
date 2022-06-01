package diffs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"okapi-public-api/pkg/contenttype"
	"okapi-public-api/pkg/namespaces"

	_ "okapi-public-api/schema/v3"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httperr"
)

// Detail http handler
// @Summary Returns a day diff metadata for namespace
// @Tags diffs
// @Description Includes identifiers, file sizes and other relevant metadata.
// @ID v1-diffs-detail
// @Security ApiKeyAuth
// @Param date path string true "A datetime of diff (YYYY-MM-DD)"
// @Param namespace path number true "Pages namespace (currently supported 0, 6, 14)"
// @Param project path string true "Project identifier"
// @Success 200 {object} schema.Project
// @Failure 404 {object} httperr.Error
// @Router /v1/diffs/meta/{date}/{namespace}/{project} [get]
func Detail(storage storage.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		ns := c.Param("namespace")

		if len(ns) > 0 && !namespaces.IsSupported(ns) {
			httperr.BadRequest(c, fmt.Sprintf("Namespace '%s' not supported!", ns))
			return
		}

		dbName := c.Param("project")

		if len(dbName) <= 1 || len(dbName) > 255 {
			httperr.BadRequest(c)
			return
		}

		rc, err := storage.Get(fmt.Sprintf("diff/%s/%s/%s_%s_%s.json", date, dbName, dbName, contenttype.JSON, ns))

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
