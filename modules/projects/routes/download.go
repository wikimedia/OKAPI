package routes

import (
	"fmt"
	"net/http"
	"okapi/helpers/exception"
	"okapi/lib/storage"
	"okapi/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Download project dump
func Download(c *gin.Context) {
	var err error
	model := models.Project{}
	model.ID, err = strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	err = models.DB().Model(&model).WherePK().Select()

	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	if len(model.Path) <= 0 {
		c.JSON(http.StatusNotFound,
			exception.Message(fmt.Errorf("nothing to download for this project")))
		return
	}

	url, err := storage.
		Remote.
		Client().
		Link(model.Path, 1*time.Minute)

	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}
