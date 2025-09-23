package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/repository"
)

func NewAuthorService(authorRepository repository.AuthorRepository, db *pgxpool.Pool, validate *validator.Validate) AuthorService {
	return &AuthorServiceImpl{
		AuthorRepository: authorRepository,
		DB:               db,
		Validate:         validate,
	}
}

type AuthorServiceImpl struct {
	AuthorRepository repository.AuthorRepository
	DB               *pgxpool.Pool
	Validate         *validator.Validate
}

func (service *AuthorServiceImpl) CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error) {
	// Validate request body
	err := service.Validate.Struct(request)
	if err != nil {
		return web.CreateAuthorResponse{}, err
	}

	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return web.CreateAuthorResponse{}, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	// Check if full_name already exists
	err = service.AuthorRepository.GetByFullName(ctx, tx, request.FullName)
	if err != nil {
		// !!! Will be change, will be using custom error !!!
		return web.CreateAuthorResponse{}, err
	}

	author := domain.Author{
		Id:          uuid.NewString(),
		FullName:    request.FullName,
		Nationality: request.Nationality,
	}

	// Call repository
	author, err = service.AuthorRepository.Save(ctx, tx, author)
	if err != nil {
		return web.CreateAuthorResponse{}, err
	}

	return helper.ToCreateAuthorResponse(author), nil
}
