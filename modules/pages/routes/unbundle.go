package routes

import (
	"net/http"
	page_unbundle "okapi/events/page/unbundle"
	"okapi/helpers/exception"
	"okapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gookit/event"
)

// Unbundle remove page from bundle
func Unbundle(c *gin.Context) {
	var err error
	model := models.Page{}

	model.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	err = models.DB().
		Model(&model).
		Relation("Project").
		WherePK().
		Select()

	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	go event.Fire(page_unbundle.Name, map[string]interface{}{
		"payload": page_unbundle.Payload{
			Title:  model.Title,
			DBName: model.Project.DBName,
		},
	})

	c.Status(http.StatusNoContent)
}
