package handler

import (
	"log/slog"
	"net/http"

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

	authorResponse, err := handler.AuthorService.CreateNewAuthor(r.Context(), authorRequest)
	if err != nil {
		slog.Error("error when calling CreateNewAuthor", "err", err)
		return
	}

	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusCreated,
	)

	// Write and send the response
	w.WriteHeader(http.StatusCreated)
	helper.WriteToResponseBody(w, web.WebSuccessResponse{
		Message: "Author created successfully",
		Data:    authorResponse,
	})
}
