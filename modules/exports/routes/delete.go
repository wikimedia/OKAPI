package routes

import (
	"net/http"
	"okapi/helpers/exception"
	"okapi/helpers/success"
	user_helper "okapi/helpers/user"
	"okapi/models"
	"strconv"

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
	var urlParams deleteURLParams
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

	var bodyParams deleteBodyParams
	err = c.ShouldBindJSON(&bodyParams)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	var exportIDs []string

	for _, exportID := range bodyParams.ResourceIDs {
		exportIDs = append(exportIDs, strconv.Itoa(exportID))
	}

	sql := "delete from exports where "
	sql += "user_id = ? and "
	sql += "resource_type = ? and "
	sql += "resource_id in (?)"

	_, err = models.DB().Exec(sql, strconv.Itoa(user.ID), urlParams.ResourceType, pg.In(exportIDs))

	if err != nil {
		c.JSON(http.StatusInternalServerError, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, success.Message("Exports have been deleted"))
	}
}
