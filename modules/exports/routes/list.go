package routes

import (
	"net/http"
	"okapi/helpers/exception"
	user_helper "okapi/helpers/user"
	"okapi/models"
	"okapi/models/exports"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type listResponse struct {
	Pages    []models.Page    `json:"page"`
	Projects []models.Project `json:"project"`
}

// List exports list grouped by resource type
func List(c *gin.Context) {
	user, err := user_helper.FromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exception.Message(err))
		return
	}

	userExports := []models.Export{}
	ids := map[exports.Resource][]int{
		exports.Page:    []int{},
		exports.Project: []int{},
	}

	models.DB().
		Model(&userExports).
		Column("resource_type", "resource_id").
		Where("user_id = ?", user.ID).
		Select()

	for _, userExport := range userExports {
		exportType := exports.Resource(userExport.ResourceType)

		if _, ok := ids[exports.Resource(userExport.ResourceType)]; ok {
			ids[exportType] = append(ids[exportType], userExport.ResourceID)
		}
	}

	res := listResponse{
		Projects: []models.Project{},
		Pages:    []models.Page{},
	}

	models.DB().Model(&res.Projects).Where("id in (?)", pg.In(ids[exports.Project])).Select()
	models.DB().Model(&res.Pages).Where("id in (?)", pg.In(ids[exports.Page])).Select()

	c.JSON(http.StatusOK, res)
}
