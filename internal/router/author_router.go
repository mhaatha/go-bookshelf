package router

import (
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/handler"
)

func AuthorRouter(handler handler.AuthorHandler, mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/authors", handler.Create)
	mux.HandleFunc("GET /api/v1/authors", handler.GetAll)
	mux.HandleFunc("GET /api/v1/authors/{id}", handler.GetById)
	mux.HandleFunc("PUT /api/v1/authors/{id}", handler.UpdateById)
}
