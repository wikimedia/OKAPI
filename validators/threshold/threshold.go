package threshold

import (
	"okapi/lib/ores"

	"github.com/go-playground/validator/v10"
)

type bounds struct {
	min float64
	max float64
}

var models = map[ores.Model]bounds{
	ores.Damaging: {0, 1},
}

// Validator function to validate ORES threshold
func Validator(fl validator.FieldLevel) bool {
	threshold, ok := fl.Field().Interface().(map[ores.Model]float64)

	if ok {
		for modelName, value := range threshold {
			if bounds, ok := models[modelName]; !ok || (value > bounds.max || value < bounds.min) {
				return false
			}
		}
	}

	return true
}
