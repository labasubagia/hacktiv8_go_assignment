package helpers

import (
	"final_project/dto"

	"github.com/go-playground/validator/v10"
)

func Validate(input any) (fieldErr []dto.FieldValidationError) {
	validate := validator.New()
	errors := validate.Struct(input)
	if errors == nil {
		return
	}
	for _, err := range errors.(validator.ValidationErrors) {
		fieldErr = append(fieldErr, dto.FieldValidationError{
			Field:   err.Field(),
			Message: "invalid",
		})
	}
	return
}
