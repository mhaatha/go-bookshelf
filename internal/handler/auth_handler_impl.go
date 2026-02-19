package handler

import (
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/service"
)

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
	}
}

type AuthHandlerImpl struct {
	AuthService service.AuthService
}

func (handler *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {}

func (handler *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {}
