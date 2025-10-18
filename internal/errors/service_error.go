package errors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

func ResponseServiceErrorHandler(w http.ResponseWriter, err error, message string) {
	// Validation error
	validationErrs := TranslateValidationErrors(err)
	if validationErrs != nil {
		slog.Error("validation error", "err", err)

		helper.WriteToResponseBody(w, http.StatusBadRequest, web.WebFailedResponse{
			Errors: validationErrs,
		})
		return
	}

	// Custom error
	var customErr *AppError
	if errors.As(err, &customErr) {
		slog.Error(message, "err", err)

		helper.WriteToResponseBody(w, customErr.StatusCode, web.WebFailedResponse{
			Errors: customErr.ErrAggregate,
		})
		return
	}

	// Unexpected error
	slog.Error(message, "err", err)

	helper.WriteToResponseBody(w, http.StatusInternalServerError, web.WebFailedResponse{
		Errors: http.StatusText(http.StatusInternalServerError),
	})
}
