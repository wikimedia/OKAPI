package utils

import "okapi-data-service/schema/v3"

var filterNs = map[int]bool{
	schema.NamespaceArticle:  true,
	schema.NamespaceFile:     true,
	schema.NamespaceCategory: true,
	schema.NamespaceTemplate: true,
}

// FilterNs check whether page in allowed namespace
func FilterNs(ns int) bool {
	_, ok := filterNs[ns]
	return ok
}
