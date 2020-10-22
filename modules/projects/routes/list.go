package routes

import (
	"net/http"

	"okapi/models"

	"okapi/helpers/filter"
	"okapi/helpers/search"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
)

// List list all entities example
func List(c *gin.Context) {
	request := search.Request{}
	c.BindQuery(&request)

	model := []models.Project{}
	params := map[search.Field]interface{}{}
	filters := map[search.Field]func(query *orm.Query){}
	columns := map[search.Field]func(column string, param string) func(*orm.Query){
		"db_name":         filter.Like,
		"lang_local_name": filter.Like,
		"lang_name":       filter.Like,
		"size":            filter.Equal,
		"site_name":       filter.Equal,
		"updates":         filter.Equal,
	}

	for name, filter := range columns {
		params[name] = c.Query(string(name))
		switch params[name].(type) {
		case string:
			filters[name] = filter(string(name), params[name].(string))
		}
	}

	sch := search.
		New().
		Model(&model).
		Filters(&filters).
		Request(&request).
		Params(&params).
		Build()

	sch.Query.Where("active = ?", true)

	c.JSON(http.StatusOK, sch.Run())
}
