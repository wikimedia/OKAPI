package namespaces

import (
	"okapi-public-api/schema/v3"
	"strconv"
)

// Supported list of available namespaces
var Supported = map[string]string{
	strconv.Itoa(schema.NamespaceArticle):  "Article",
	strconv.Itoa(schema.NamespaceFile):     "File",
	strconv.Itoa(schema.NamespaceCategory): "Category",
}

// IsSupported check if we support particular namespace
func IsSupported(ns string) bool {
	_, ok := Supported[ns]
	return ok
}
