package pages

import (
	"okapi-public-api/pkg/contenttype"

	_ "okapi-public-api/schema/v3"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Meta http handler
// @Summary Returns JSON structure of page
// @Tags pages
// @Description Most current revision of a page.
// @ID v1-page-data
// @Security ApiKeyAuth
// @Param project path string true "Project identifier"
// @Param name path string true "Page name"
// @Success 200 {object} schema.Page
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Router /v1/pages/meta/{project}/{name} [get]
func Meta(store storage.Storage) gin.HandlerFunc {
	return Download(store, contenttype.JSON)
}
