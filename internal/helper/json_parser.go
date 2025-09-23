package helper

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

// ReadFromRequestBody reads the request body and stores it in the result parameter
func ReadFromRequestBody(r *http.Request, result interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}

	return nil
}

// WriteToResponseBody writes the result parameter to the response body
func WriteToResponseBody(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		slog.Error("failed to encode response body", "err", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		fallback := web.WebFailedResponse{
			Errors: http.StatusText(http.StatusInternalServerError),
		}

		_ = json.NewEncoder(w).Encode(fallback)
	}
}
