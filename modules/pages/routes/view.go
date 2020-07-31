package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"okapi/helpers/exception"
	"okapi/models"
)

// View view one entity example
func View(c *gin.Context) {
	var err error

	model := models.Page{}
	model.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	err = models.DB().Model(&model).WherePK().Select()
	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	c.JSON(http.StatusOK, model)
}
