package routes

import (
	"fmt"
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
	var params createURLParams
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

	export := models.Export{
		UserID:       user.ID,
		ResourceType: string(params.ResourceType),
		ResourceID:   params.ID,
	}

	resource := exports.Types[params.ResourceType]
	resourceExists, _ := models.DB().Model(resource).Where("id = ?", params.ID).Exists()

	if !resourceExists {
		c.AbortWithStatusJSON(http.StatusNotFound, exception.Message(fmt.Errorf("Resource does not exist")))
		return
	}

	if err = models.Save(&export); err != nil {
		c.JSON(http.StatusInternalServerError, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, &export)
	}
}
