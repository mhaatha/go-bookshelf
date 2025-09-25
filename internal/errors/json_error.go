package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

func RequestJSONErrorHandler(w http.ResponseWriter, err error) {
	// Handle JSON syntax error
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		slog.Error("invalid JSON syntax", "err", err)

		helper.WriteToResponseBody(w, http.StatusInternalServerError, web.WebFailedResponse{
			Errors: "Invalid JSON payload",
		})
		return
	}

	// Handle invalid JSON field type
	var jsonTypeErr *json.UnmarshalTypeError
	if errors.As(err, &jsonTypeErr) {
		slog.Error("invalid JSON type", "err", err)

		helper.WriteToResponseBody(w, http.StatusInternalServerError, web.WebFailedResponse{
			Errors: fmt.Sprintf("Invalid JSON type for field: %v", jsonTypeErr.Field),
		})
		return
	}

	// Unexpected error
	slog.Error("failed to read JSON from request body", "err", err)

	helper.WriteToResponseBody(w, http.StatusInternalServerError, web.WebFailedResponse{
		Errors: http.StatusText(http.StatusInternalServerError),
	})
}
