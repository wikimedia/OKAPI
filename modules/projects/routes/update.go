package routes

import (
	"fmt"
	"net/http"
	"okapi/helpers/exception"
	"okapi/lib/ores"
	"okapi/models"

	"github.com/gin-gonic/gin"
)

type projectParams struct {
	Threshold map[ores.Model]float64 `form:"threshold"`
}

// Update one entity example
func Update(c *gin.Context) {
	var params projectParams
	projectId := c.Param("id")

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))

		return
	}

	var project models.Project

	err := models.DB().Model(&project).Where("id = ?", projectId).Select()

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			exception.Message(fmt.Errorf("project with id %s does not exist", projectId)),
		)

		return
	}

	if len(params.Threshold) > 0 {
		for modelName, thresholdVal := range params.Threshold {
			project.Threshold[modelName] = thresholdVal
		}
	}

	if err := models.Save(&project); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, &project)
	}
}
