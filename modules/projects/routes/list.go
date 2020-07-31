package routes

import (
	"net/http"

	"okapi/models"

	"okapi/helpers/filter"
	"okapi/helpers/search"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9/orm"
)

// List list all entities example
func List(c *gin.Context) {
	request := search.Request{}
	c.BindQuery(&request)

	model := []models.Project{}

	params := map[search.Field]interface{}{
		"db_name":    c.Query("db_name"),
		"local_name": c.Query("local_name"),
		"name":       c.Query("name"),
		"size":       c.Query("size"),
	}

	filters := map[search.Field]func(query *orm.Query){
		"db_name":    filter.Like("db_name", params["db_name"].(string)),
		"local_name": filter.Like("local_name", params["local_name"].(string)),
		"name":       filter.Like("name", params["name"].(string)),
		"size":       filter.Equal("size", params["size"].(string)),
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
