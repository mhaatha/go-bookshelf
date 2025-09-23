package config

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidatorInit() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	return validate
}
