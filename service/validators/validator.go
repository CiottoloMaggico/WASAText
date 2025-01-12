package validators

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func NewAppValidator() error {
	if err := Validate.RegisterValidation("emoji", ValidateEmoji); err != nil {
		return err
	}
	if err := Validate.RegisterValidation("image", ValidateImage); err != nil {
		return err
	}
	if err := Validate.RegisterValidation("formula", ValidateFormula); err != nil {
		return err
	}

	return nil
}
