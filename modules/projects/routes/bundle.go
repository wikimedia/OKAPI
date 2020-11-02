package routes

import (
	"net/http"
	project_bundle "okapi/events/project/bundle"
	"okapi/helpers/exception"
	"okapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gookit/event"
)

// Bundle trigger bundle job for project
func Bundle(c *gin.Context) {
	var err error
	model := models.Project{}
	model.ID, err = strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	err = models.DB().Model(&model).Column("db_name").WherePK().Select()

	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	go event.Fire(project_bundle.Name, map[string]interface{}{
		"payload": project_bundle.Payload{
			DBName: model.DBName,
		},
	})

	c.Status(http.StatusNoContent)
}
