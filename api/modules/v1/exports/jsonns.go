package exports

import (
	"okapi-public-api/lib/env"
	"okapi-public-api/pkg/contenttype"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// JSONNS http handler
// @Summary Returns today’s tar.gz file with JSON export of entire project corpus in specified namespace
// @Tags exports
// @Description Full project export of current revisions updated daily at 12:00 UTC. The archive contains JSON files for each article including revision Wikitext, HTML, and relevant metadata.
// @ID v1-exports-download-ns
// @Security ApiKeyAuth
// @Param namespace path number true "Pages namespace (currently supported 0, 6, 14)"
// @Param project path string true "Project identifier"
// @Success 307 string nil "Redirects to the direct download URL"
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Router /v1/exports/download/{namespace}/{project} [get]
func JSONNS(store storage.Storage) gin.HandlerFunc {
	return Download(store, contenttype.JSON, env.Group)
}
