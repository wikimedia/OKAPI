package routes

import (
	"io/ioutil"
	"net/http"
	"okapi/helpers/exception"
	"okapi/indexes/page"
	"okapi/lib/elastic"

	"github.com/gin-gonic/gin"
)

// Pages search in pages index
func Pages(c *gin.Context) {
	elastic := elastic.Client()

	res, err := elastic.Search(
		elastic.Search.WithIndex(page.Name),
		elastic.Search.WithBody(c.Request.Body),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	if res.IsError() {
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}

	c.Writer.Write(body)
}
