package export_type

import (
	"okapi/models/exports"

	"github.com/go-playground/validator/v10"
)

// Validator function for export resource type
func Validator(fl validator.FieldLevel) bool {
	exists := true
	var exportType exports.Resource

	switch fl.Field().Interface().(type) {
	case string:
		exportType = exports.Resource(fl.Field().Interface().(string))
	case exports.Resource:
		exportType = fl.Field().Interface().(exports.Resource)
	}

	if len(exportType) > 0 {
		_, exists = exports.Types[exportType]
	}

	return exists
}
