package router

import (
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/handler"
)

func UploadRouter(handler handler.UploadHandler, mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/upload/books/presigned-url", handler.GetBookPresignedURL)
}
