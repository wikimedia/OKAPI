package pages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"okapi-public-api/pkg/contenttype"
	"okapi-public-api/schema/v3"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httperr"
)

// Download http handler
func Download(storage storage.Getter, cType contenttype.ContentType) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbName := c.Param("project")

		if len(dbName) <= 1 || len(dbName) > 255 {
			httperr.BadRequest(c)
			return
		}

		name := strings.ReplaceAll(c.Param("name"), " ", "_")

		if len(name) <= 1 || len(name) > 1000 {
			httperr.BadRequest(c)
			return
		}

		rc, err := storage.Get(fmt.Sprintf("page/json/%s/%s.json", dbName, name))

		if err != nil {
			httperr.NotFound(c, err.Error())
			return
		}

		defer rc.Close()

		if cType == contenttype.JSON {
			data, err := ioutil.ReadAll(rc)

			if err != nil {
				httperr.InternalServerError(c, err.Error())
				return
			}

			c.Data(http.StatusOK, fmt.Sprintf("%s; charset=UTF-8", gin.MIMEJSON), data)
			return
		}

		page := new(schema.Page)

		if err := json.NewDecoder(rc).Decode(page); err != nil {
			httperr.BadRequest(c, err.Error())
			return
		}

		switch cType {
		case contenttype.HTML:
			c.Data(http.StatusOK, fmt.Sprintf("%s; charset=UTF-8", gin.MIMEHTML), []byte(page.ArticleBody.HTML))
		case contenttype.WText:
			c.Data(http.StatusOK, "text/wikitext; charset=UTF-8", []byte(page.ArticleBody.Wikitext))
		}
	}
}
