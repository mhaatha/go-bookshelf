package errors

import (
	"fmt"
	"strings"
)

type ErrAggregate struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	StatusCode   int
	ErrAggregate []ErrAggregate
	Err          error
}

func (ae *AppError) Error() string {
	if len(ae.ErrAggregate) != 0 {
		errs := []string{}
		for _, e := range ae.ErrAggregate {
			errs = append(errs, fmt.Sprintf("field=%s: %s", e.Field, e.Message))
		}

		return fmt.Sprintf("(%v) errors: %s", ae.StatusCode, strings.Join(errs, "; "))
	}

	return fmt.Sprintf("(%v) - cause=%v", ae.StatusCode, ae.Err)
}

func NewAppError(statusCode int, errAggregate []ErrAggregate, err error) *AppError {
	return &AppError{
		StatusCode:   statusCode,
		ErrAggregate: errAggregate,
		Err:          err,
	}
}
