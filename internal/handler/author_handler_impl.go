package handler

import (
	"log/slog"
	"net/http"

	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/service"
)

const (
	queryFullName    = "full_name"
	queryNationality = "nationality"

	wildcardId = "id"
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
		appError.RequestJSONErrorHandler(w, err)
		return
	}

	// Call the service
	authorResponse, err := handler.AuthorService.CreateNewAuthor(r.Context(), authorRequest)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to create new author")
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

func (handler *AuthorHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get query params if any
	queries := web.QueryParamsGetAuthors{
		FullName:    r.URL.Query().Get(queryFullName),
		Nationality: r.URL.Query().Get(queryNationality),
	}

	// Call the service
	authorsResponse, err := handler.AuthorService.GetAllAuthors(r.Context(), queries)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to get authors")
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusOK,
	)

	// Write and send the response
	helper.WriteToResponseBody(w, http.StatusOK, web.WebSuccessResponse{
		Message: "Success get all authors",
		Data:    authorsResponse,
	})
}

func (handler *AuthorHandlerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	// Get path values if any
	pathValue := web.PathParamsGetAuthor{
		Id: r.PathValue(wildcardId),
	}

	// Call the service
	authorResponse, err := handler.AuthorService.GetAuthorById(r.Context(), pathValue)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to get author by id")
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusOK,
	)

	// Write and send the response
	helper.WriteToResponseBody(w, http.StatusOK, web.WebSuccessResponse{
		Message: "Success get author",
		Data:    authorResponse,
	})
}

func (handler *AuthorHandlerImpl) UpdateById(w http.ResponseWriter, r *http.Request) {
	// Get path values if any
	pathValue := web.PathParamsUpdateAuthor{
		Id: r.PathValue(wildcardId),
	}

	// Get request body and write it to authorRequest
	authorRequest := web.UpdateAuthorRequest{}
	err := helper.ReadFromRequestBody(r, &authorRequest)
	if err != nil {
		appError.RequestJSONErrorHandler(w, err)
		return
	}

	// Call the service
	authorResponse, err := handler.AuthorService.UpdateAuthorById(r.Context(), pathValue, authorRequest)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to update author by id")
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusOK,
	)

	// Write and send the response
	helper.WriteToResponseBody(w, http.StatusOK, web.WebSuccessResponse{
		Message: "Author updated successfully",
		Data:    authorResponse,
	})
}

func (handler *AuthorHandlerImpl) DeleteById(w http.ResponseWriter, r *http.Request) {
	// Get path values if any
	pathValue := web.PathParamsDeleteAuthor{
		Id: r.PathValue(wildcardId),
	}

	// Call the service
	err := handler.AuthorService.DeleteAuthorById(r.Context(), pathValue)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to delete author by id")
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusNoContent,
	)

	// Set to 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
