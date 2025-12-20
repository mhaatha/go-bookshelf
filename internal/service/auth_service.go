package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type AuthService interface {
	CreateNewUser(ctx context.Context, request web.CreateUserRequest) (web.CreateUserResponse, error)
	LoginExistingUser(ctx context.Context, request web.LoginRequest) (web.LoginResponse, error)
}
