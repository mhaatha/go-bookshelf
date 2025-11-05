package config

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// Regex for name, only a-z, A-Z, ., ', and -
	nameRegex = regexp.MustCompile(`^[a-zA-Z .'-]+$`)
)

func ValidatorInit() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	// Register custom validation
	validate.RegisterValidation("validName", validName)
	validate.RegisterValidation("bookStatus", bookStatus)
	validate.RegisterValidation("validPhotoKey", validPhotoKey)

	return validate
}

func validName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return nameRegex.MatchString(value)
}

func bookStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()

	if status != "completed" && status != "reading" && status != "plan_to_read" {
		return false
	}

	return true
}

func validPhotoKey(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" {
		return false
	}

	return strings.HasSuffix(strings.ToLower(value), ".jpg")
}
