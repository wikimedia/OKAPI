package routes

import (
	"net/http"
	"okapi/helpers/exception"
	user_helper "okapi/helpers/user"
	"okapi/models"
	"okapi/models/exports"

	"github.com/gin-gonic/gin"
)

type createURLParams struct {
	ID           int              `uri:"id" binding:"required"`
	ResourceType exports.Resource `uri:"resource_type" binding:"required,export_type"`
}

// Create create project or page export
func Create(c *gin.Context) {
	params := createURLParams{}
	err := c.ShouldBindUri(&params)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	user, err := user_helper.FromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exception.Message(err))
		return
	}

	resource := exports.Types[params.ResourceType]
	export := models.Export{
		UserID:       user.ID,
		ResourceType: string(params.ResourceType),
		ResourceID:   params.ID,
	}

	if err = models.DB().Model(resource).Where("id = ?", params.ID).Select(); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, exception.Message(err))
		return
	}

	if err = models.Save(&export); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	c.JSON(http.StatusOK, resource)
}
