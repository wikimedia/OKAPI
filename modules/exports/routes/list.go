package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"net/http"
	"okapi/helpers/exception"
	user_helper "okapi/helpers/user"
	"okapi/models"
	"okapi/models/exports"
)

// List exports list grouped by resource type
func List(c *gin.Context) {
	user, err := user_helper.FromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exception.Message(err))
		return
	}

	userExports := []models.Export{}

	models.DB().
		Model(&userExports).
		Column("resource_type", "resource_id").
		Where("user_id = ?", user.ID).
		Select()

	projectIDs := []int{}
	pageIDs := []int{}

	for _, userExport := range userExports {
		switch exports.Resource(userExport.ResourceType) {
		case exports.Project:
			projectIDs = append(projectIDs, userExport.ResourceID)
		case exports.Page:
			pageIDs = append(pageIDs, userExport.ResourceID)
		}
	}

	projects := []models.Project{}
	pages := []models.Page{}

	models.DB().Model(&projects).Where("id in (?)", pg.In(projectIDs)).Select()
	models.DB().Model(&pages).Where("id in (?)", pg.In(pageIDs)).Select()

	res := map[exports.Resource]interface{}{
		exports.Project: projects,
		exports.Page:    pages,
	}

	c.JSON(http.StatusOK, res)
}
