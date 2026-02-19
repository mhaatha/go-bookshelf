package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

func NewAuthService(uow UnitOfWork, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UoW:      uow,
		Validate: validate,
	}
}

type AuthServiceImpl struct {
	UoW      UnitOfWork
	Validate *validator.Validate
}

func (service *AuthServiceImpl) CreateNewUser(ctx context.Context, request web.CreateUserRequest) (web.CreateUserResponse, error) {
	return web.CreateUserResponse{}, nil
}

func (service *AuthServiceImpl) LoginExistingUser(ctx context.Context, request web.LoginRequest) (web.LoginResponse, error) {
	return web.LoginResponse{}, nil
}
