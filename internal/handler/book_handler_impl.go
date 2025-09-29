package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/service"
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

	fmt.Println(bookRequest)

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
