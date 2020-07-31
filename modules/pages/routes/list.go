package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9/orm"
	"okapi/helpers/filter"
	"okapi/helpers/search"
	"okapi/models"
)

// List list all entities example
func List(c *gin.Context) {
	request := search.Request{}
	c.BindQuery(&request)

	model := []models.Page{}

	params := map[search.Field]interface{}{
		"title":      c.Query("title"),
		"lang":       c.Query("lang"),
		"project_id": c.Query("project_id"),
	}

	filters := map[search.Field]func(query *orm.Query){
		"title":      filter.Like("title", params["title"].(string)),
		"lang":       filter.Like("lang", params["lang"].(string)),
		"project_id": filter.Equal("project_id", params["project_id"].(string)),
	}

	search := search.
		New().
		Model(&model).
		Filters(&filters).
		Request(&request).
		Params(&params).
		Build()

	search.Query.Relation("Project")

	c.JSON(http.StatusOK, search.Run())
}
