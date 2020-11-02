package routes

import (
	"net/http"
	"okapi/helpers/exception"
	"okapi/helpers/success"
	user_helper "okapi/helpers/user"
	"okapi/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type deleteURLParams struct {
	ResourceType string `uri:"resource_type" binding:"required,export_type"`
}

type deleteBodyParams struct {
	ResourceIDs []int `json:"resourceIds" binding:"required"`
}

// Delete delete exports by id or array of ids
func Delete(c *gin.Context) {
	urlParams := deleteURLParams{}
	err := c.ShouldBindUri(&urlParams)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	user, err := user_helper.FromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exception.Message(err))
		return
	}

	bodyParams := deleteBodyParams{}
	err = c.ShouldBindJSON(&bodyParams)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	exportIDs := []int{}
	for _, exportID := range bodyParams.ResourceIDs {
		exportIDs = append(exportIDs, exportID)
	}

	_, err = models.DB().
		Model(&models.Export{}).
		Where(
			"user_id = ? and resource_type = ? and resource_id in (?)",
			user.ID,
			urlParams.ResourceType,
			pg.In(exportIDs)).
		Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, success.Message("Exports have been deleted"))
	}
}
