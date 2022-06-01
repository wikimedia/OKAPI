package exports

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"okapi-public-api/pkg/namespaces"

	_ "okapi-public-api/schema/v3"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httperr"
	"github.com/protsack-stephan/gin-toolkit/httpmw"
)

// Detail http handler
// @Summary Returns export metadata for namespace
// @Tags exports
// @Description Includes identifiers, file sizes and other relevant metadata.
// @ID v1-exports-detail
// @Security ApiKeyAuth
// @Param namespace path number true "Pages namespace (currently supported 0, 6, 14)"
// @Param project path string true "Project identifier"
// @Success 200 {object} schema.Project
// @Failure 404 {object} httperr.Error
// @Router /v1/exports/meta/{namespace}/{project} [get]
func Detail(storage storage.Getter, group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uRaw, ok := c.Get("user")

		if !ok {
			httperr.Unauthorized(c, "User not found!")
			return
		}

		user, ok := uRaw.(*httpmw.CognitoUser)

		if !ok {
			httperr.InternalServerError(c, "Unknown user type!")
			return
		}

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

		path := fmt.Sprintf("export/%s/%s_%s.json", dbName, dbName, ns)

		for _, role := range user.GetGroups() {
			if role == group {
				path = fmt.Sprintf("export/%s/%s_%s_%s.json", dbName, dbName, group, ns)
			}
		}

		rc, err := storage.Get(path)

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
