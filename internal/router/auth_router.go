package router

import (
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/handler"
)

func AuthRouter(handler handler.AuthHandler, mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", handler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
}
