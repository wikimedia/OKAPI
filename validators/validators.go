package validators

import (
	"fmt"
	"okapi/validators/export_type"
	"okapi/validators/threshold"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Init function to initialize all the validators
func Init() error {
	validators := map[string]validator.Func{
		"threshold":   threshold.Validator,
		"export_type": export_type.Validator,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return fmt.Errorf("model validations failed to bind")
	}

	for name, fn := range validators {
		err := v.RegisterValidation(name, fn)

		if err != nil {
			return err
		}
	}

	return nil
}
