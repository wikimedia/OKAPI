package diffs

import (
	"okapi-public-api/pkg/contenttype"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// JSONNS http handler
// @Summary Returns tar.gz file with a specified date’s project revisions in JSON for specified namespace
// @Tags diffs
// @Description Hourly updated bundle of revised pages starting at 00:00 UTC each day.
// @ID v1-diffs-json-ns
// @Security ApiKeyAuth
// @Param date path string true "Date of the diff in YYYY-MM-DD"
// @Param project path string true "Project identifier"
// @Param namespace path number true "Pages namespace (currently supported 0, 6, 14)"
// @Success 307 string nil "Redirects to the direct download URL"
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Router /v1/diffs/download/{date}/{namespace}/{project} [get]
func JSONNS(store storage.Storage) gin.HandlerFunc {
	return Download(store, contenttype.JSON)
}
