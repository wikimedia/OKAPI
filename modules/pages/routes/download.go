package routes

import (
	"net/http"
	"strconv"

	"okapi/helpers/exception"
	"okapi/lib/storage"
	"okapi/models"

	"github.com/gin-gonic/gin"
)

// Download sends page's HTML file
func Download(c *gin.Context) {
	var err error
	model := models.Page{}
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

	filePath, err := storage.Local.Client().Link(model.Path, 0)
	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	c.File(filePath)
}
