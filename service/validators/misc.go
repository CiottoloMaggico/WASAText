package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/tmdvs/Go-Emoji-Utils"
)

func ValidateEmoji(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	if _, err := emoji.LookupEmoji(fieldValue); err != nil {
		return false
	}
	return true
}
