package search

import "okapi/models"

func getSiteCodes(lang string) (interface{}, error) {
	options := []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}{}
	err := models.DB().
		Model(&models.Project{}).
		ColumnExpr("site_code as value, site_name as label").
		Where("lang = ?", lang).
		GroupExpr("value, label").
		OrderExpr("value asc").
		Select(&options)

	return options, err
}
