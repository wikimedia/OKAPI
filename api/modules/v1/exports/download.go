package exports

import (
	"fmt"
	"net/http"
	"okapi-public-api/pkg/contenttype"
	"okapi-public-api/pkg/namespaces"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httperr"
	"github.com/protsack-stephan/gin-toolkit/httpmw"
)

const defaultNs = "0"

type downloadStorage interface {
	storage.Linker
	storage.Stater
}

// Download http handler creator
func Download(storage downloadStorage, cType contenttype.ContentType, group string) gin.HandlerFunc {
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

		dbName := c.Param("project")

		if len(dbName) <= 1 || len(dbName) > 255 {
			httperr.BadRequest(c)
			return
		}

		ns := c.Param("namespace")

		if len(ns) > 0 && !namespaces.IsSupported(ns) {
			httperr.BadRequest(c, fmt.Sprintf("Namespace '%s' not supported!", ns))
			return
		}

		if len(ns) == 0 {
			ns = defaultNs
		}

		path := fmt.Sprintf("export/%s/%s_%s_%s.tar.gz", dbName, dbName, cType, ns)

		for _, role := range user.GetGroups() {
			if role == group {
				path = fmt.Sprintf("export/%s/%s_%s_%s_%s.tar.gz", dbName, dbName, group, cType, ns)
			}
		}

		if _, err := storage.Stat(path); err != nil {
			httperr.NotFound(c, fmt.Sprintf("Export for '%s' not found!", dbName))
			return
		}

		url, err := storage.Link(
			path,
			10*time.Second,
		)

		if err != nil {
			httperr.InternalServerError(c, err.Error())
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
