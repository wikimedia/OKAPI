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

// List http handler
// @Summary Returns list of all available exports for namespace
// @Tags exports
// @Description Includes identifiers, file sizes and other relevant metadata.
// @ID v1-exports-list
// @Security ApiKeyAuth
// @Param namespace path number true "Pages namespace (currently supported 0, 6, 14)"
// @Success 200 {object} []schema.Project
// @Failure 404 {object} httperr.Error
// @Router /v1/exports/meta/{namespace} [get]
func List(storage storage.Getter, group string) gin.HandlerFunc {
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

		path := fmt.Sprintf("public/exports_%s.json", ns)

		for _, role := range user.GetGroups() {
			if role == group {
				path = fmt.Sprintf("public/exports_%s_%s.json", group, ns)
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
