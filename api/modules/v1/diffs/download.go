package diffs

import (
	"fmt"
	"net/http"
	"okapi-public-api/pkg/contenttype"
	"okapi-public-api/pkg/namespaces"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httperr"
)

const dateFormat = "2006-01-02"
const defaultNs = "0"

type downloadStorage interface {
	storage.Linker
	storage.Stater
}

// Download http handler creator
func Download(storage downloadStorage, cType contenttype.ContentType) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")

		if len(date) <= 0 {
			httperr.BadRequest(c)
			return
		}

		if _, err := time.Parse(dateFormat, date); err != nil {
			httperr.BadRequest(c)
			return
		}

		dbName := c.Param("project")

		if len(dbName) <= 1 || len(dbName) > 255 {
			httperr.BadRequest(c)
			return
		}

		ns := c.Param("namespace")

		if len(ns) <= 0 {
			ns = defaultNs
		} else if !namespaces.IsSupported(ns) {
			httperr.BadRequest(c, fmt.Sprintf("Namespace '%s' not supported!", ns))
			return
		}

		path := fmt.Sprintf("diff/%s/%s/%s_%s_%s.tar.gz", date, dbName, dbName, cType, ns)

		if _, err := storage.Stat(path); err != nil {
			httperr.NotFound(c, fmt.Sprintf("Diff from '%s' for '%s' not found!", date, dbName))
			return
		}

		url, err := storage.Link(
			path,
			10*time.Second,
		)

		if err != nil {
			httperr.NotFound(c, err.Error())
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
