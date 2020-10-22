package routes

import (
	"io/ioutil"
	"net/http"
	"okapi/helpers/exception"
	"okapi/lib/storage"

	"github.com/gin-gonic/gin"
)

// Options options for pages and projects
func Options(c *gin.Context) {
	readCloser, err := storage.Local.Client().Get("options/" + c.Param("lang") + "/options.json")

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	body, err := ioutil.ReadAll(readCloser)

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	c.Writer.Write(body)
}
