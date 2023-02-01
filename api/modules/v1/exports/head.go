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

// Head http handler
// @Summary Returns the headers with file size and other data before the download
// @Tags exports
// @Description Includes accept-ranges, Last-Modified, Content-Length, ETag, Cache-Control, Content-Disposition, Content-Encoding, Content-Type and Expires headers
// @ID v1-exports-head
// @Security ApiKeyAuth
// @Param date path string true "Date of the diff in YYYY-MM-DD"
// @Param project path string true "Project identifier"
// @Param namespace path number true "Pages namespace (currently supported 0, 6, 14)"
// @Success 200
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Router /v1/exports/download/{namespace}/{project} [head]
func Head(storage storage.Stater, cType contenttype.ContentType, group string) gin.HandlerFunc {
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

		if len(ns) <= 0 {
			ns = defaultNs
		} else if !namespaces.IsSupported(ns) {
			httperr.BadRequest(c)
			return
		}

		path := fmt.Sprintf("export/%s/%s_%s_%s.tar.gz", dbName, dbName, cType, ns)

		for _, role := range user.GetGroups() {
			if role == group {
				path = fmt.Sprintf("export/%s/%s_%s_%s_%s.tar.gz", dbName, dbName, group, cType, ns)
			}
		}

		stat, err := storage.Stat(path)

		if err != nil {
			httperr.NotFound(c)
			return
		}

		c.Header("accept-ranges", stat.AcceptRanges())
		c.Header("Last-Modified", stat.LastModified().Format(time.RFC1123))
		c.Header("Content-Length", fmt.Sprint(stat.Size()))
		c.Header("ETag", stat.ETag())
		c.Header("Cache-Control", stat.CacheControl())
		c.Header("Content-Disposition", stat.ContentDisposition())
		c.Header("Content-Encoding", stat.ContentEncoding())
		c.Header("Content-Type", stat.ContentType())
		c.Header("Expires", stat.Expires())
		c.Status(http.StatusOK)
	}
}
