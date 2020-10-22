package routes

import (
	"net/http"

	"okapi/helpers/filter"
	"okapi/helpers/search"
	"okapi/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
)

// List list all entities example
func List(c *gin.Context) {
	request := search.Request{}
	c.BindQuery(&request)

	model := []models.Page{}
	params := map[search.Field]interface{}{}
	filters := map[search.Field]func(query *orm.Query){}
	columns := map[search.Field]func(column string, param string) func(*orm.Query){
		"title":      filter.Like,
		"lang":       filter.Like,
		"project_id": filter.Equal,
		"updates":    filter.Equal,
	}

	for name, filter := range columns {
		params[name] = c.Query(string(name))
		switch params[name].(type) {
		case string:
			filters[name] = filter(string(name), params[name].(string))
		}
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
