package search

import "okapi/models"

func getSiteNames(lang string) (interface{}, error) {
	options := []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}{}
	err := models.DB().
		Model(&models.Project{}).
		ColumnExpr("site_name as value, site_name as label").
		GroupExpr("value, label").
		OrderExpr("value asc").
		Select(&options)

	return options, err
}
