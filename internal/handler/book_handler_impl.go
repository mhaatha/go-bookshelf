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
	queryName       = "name"
	queryStatus     = "status"
	queryAuthorName = "author_name"
)

func NewBookHandler(bookService service.BookService) BookHandler {
	return &BookHandlerImpl{
		BookService: bookService,
	}
}

type BookHandlerImpl struct {
	BookService service.BookService
}

func (handler *BookHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	// Get request body and write it to bookRequest
	bookRequest := web.CreateBookRequest{}
	err := helper.ReadFromRequestBody(r, &bookRequest)
	if err != nil {
		appError.RequestJSONErrorHandler(w, err)
		return
	}

	// Call the service
	bookResponse, err := handler.BookService.CreateNewBook(r.Context(), bookRequest)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to create new book")
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
		Message: "Book created successfully",
		Data:    bookResponse,
	})
}

func (handler *BookHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get query params if any
	queries := web.QueryParamsGetBooks{
		Name:       r.URL.Query().Get(queryName),
		Status:     r.URL.Query().Get(queryStatus),
		AuthorName: r.URL.Query().Get(queryAuthorName),
	}

	// Call the service
	authorsResponse, err := handler.BookService.GetAllBooks(r.Context(), queries)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to get books")
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
		Message: "Success get all books",
		Data:    authorsResponse,
	})
}

func (handler *BookHandlerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	// Get path values if any
	pathValue := web.PathParamsGetBook{
		Id: r.PathValue(wildcardId),
	}

	// Call the service
	bookResponse, err := handler.BookService.GetBookById(r.Context(), pathValue)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to get book by id")
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
		Message: "Success get book",
		Data:    bookResponse,
	})
}

func (handler *BookHandlerImpl) UpdateById(w http.ResponseWriter, r *http.Request) {
	// Get path values if any
	pathValue := web.PathParamsUpdateBook{
		Id: r.PathValue(wildcardId),
	}

	// Get request body and write it to bookRequest
	bookRequest := web.UpdateBookRequest{}
	err := helper.ReadFromRequestBody(r, &bookRequest)
	if err != nil {
		appError.RequestJSONErrorHandler(w, err)
		return
	}

	// Call the service
	bookResponse, err := handler.BookService.UpdateBookById(r.Context(), pathValue, bookRequest)
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to update book by id")
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusOK,
	)

	// Write and send the response
	helper.WriteToResponseBody(w, http.StatusCreated, web.WebSuccessResponse{
		Message: "Book updated successfully",
		Data:    bookResponse,
	})
}
