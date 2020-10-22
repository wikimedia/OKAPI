package search

import (
	"okapi/models"
)

func getNamespaces(lang string) (interface{}, error) {
	options := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{}
	err := models.DB().
		Model(&models.Namespace{}).
		ColumnExpr("id as value, title as label").
		Where("lang = ?", lang).
		GroupExpr("value, label").
		OrderExpr("value asc").
		Select(&options)

	return options, err
}
