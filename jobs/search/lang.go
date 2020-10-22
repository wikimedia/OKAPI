package search

import "okapi/models"

func getLangs(lang string) (interface{}, error) {
	options := []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}{}
	err := models.DB().
		Model(&models.Project{}).
		ColumnExpr("lang as value, lang_local_name as label").
		GroupExpr("value, label").
		OrderExpr("value asc").
		Select(&options)

	return options, err
}
