package handler

import (
	"errors"
	"log/slog"
	"net/http"

	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/service"
)

func NewAuthorHandler(authorService service.AuthorService) AuthorHandler {
	return &AuthorHandlerImpl{
		AuthorService: authorService,
	}
}

type AuthorHandlerImpl struct {
	AuthorService service.AuthorService
}

func (handler *AuthorHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	// Get request body and write it to authorRequest
	authorRequest := web.CreateAuthorRequest{}
	err := helper.ReadFromRequestBody(r, &authorRequest)
	if err != nil {
		slog.Error("failed read JSON from request body", "err", err)
		return
	}

	// Call the service
	authorResponse, err := handler.AuthorService.CreateNewAuthor(r.Context(), authorRequest)
	if err != nil {
		// Check if it's a validation error
		validationErrs := appError.TranslateValidationErrors(err)

		if validationErrs != nil {
			slog.Error("validation error", "err", err)

			helper.WriteToResponseBody(w, http.StatusBadRequest, web.WebFailedResponse{
				Errors: validationErrs,
			})
			return
		}

		// Check if it's a custom error
		var customErr *appError.AppError

		if errors.As(err, &customErr) {
			slog.Error("failed to create author", "err", err)

			helper.WriteToResponseBody(w, customErr.StatusCode, web.WebFailedResponse{
				Errors: customErr.ErrAggregate,
			})
			return
		}

		// If it's not a validation error, then it's an unexpected error
		slog.Error("error when calling CreateNewAuthor", "err", err)

		helper.WriteToResponseBody(w, http.StatusInternalServerError, web.WebFailedResponse{
			Errors: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusCreated,
	)

	// Write and send the response
	helper.WriteToResponseBody(w, http.StatusCreated, web.WebSuccessResponse{
		Message: "Author created successfully",
		Data:    authorResponse,
	})
}
