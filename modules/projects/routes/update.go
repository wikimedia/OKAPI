package routes

import (
	"net/http"
	"okapi/helpers/exception"
	"okapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Update a project record
func Update(c *gin.Context) {
	var err error
	project := models.Project{}

	if project.ID, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	_, err = models.DB().
		Model(&project).
		Column("threshold", "time_delay", "updated_at").
		WherePK().
		Update()

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, &project)
	}
}
