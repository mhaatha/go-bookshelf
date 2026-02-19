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

	// Regex for password, min 8 chars, at least one uppercase and lowercase, and at least one digit
	upperRe = regexp.MustCompile(`[A-Z]`)
	lowerRe = regexp.MustCompile(`[a-z]`)
	digitRe = regexp.MustCompile(`[0-9]`)
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
	validate.RegisterValidation("validPassword", validPassword)

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

func validPassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) < 8 {
		return false
	}
	if !upperRe.MatchString(value) {
		return false
	}
	if !lowerRe.MatchString(value) {
		return false
	}
	if !digitRe.MatchString(value) {
		return false
	}

	return true
}
