package router

import (
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/handler"
)

func BookRouter(handler handler.BookHandler, mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/books", handler.Create)
	mux.HandleFunc("GET /api/v1/books", handler.GetAll)
	mux.HandleFunc("GET /api/v1/books/{id}", handler.GetById)
}
