package helper

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(err error) []map[string]string {
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		out := make([]map[string]string, 0, len(validateErrs))
		for _, e := range validateErrs {
			msg := ""
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", e.Field())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
			case "validName":
				msg = fmt.Sprintf("%s must not contain numbers and symbols", e.Field())
			case "alpha":
				msg = fmt.Sprintf("%s must not contain numbers and symbols", e.Field())
			default:
				msg = fmt.Sprintf("%s is invalid", e.Field())
			}
			out = append(out, map[string]string{
				"field":   e.Field(),
				"message": msg,
			})
		}
		return out
	}
	return nil
}
