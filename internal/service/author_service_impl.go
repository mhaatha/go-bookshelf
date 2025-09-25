package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
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

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// Check if full_name already exists
	err = service.AuthorRepository.GetByFullName(ctx, tx, request.FullName)
	if err != nil {
		errAggregate = append(errAggregate, appError.ErrAggregate{
			Field:   "full_name",
			Message: fmt.Sprintf("author %s is already exists", request.FullName),
		})
	}

	if len(errAggregate) != 0 {
		return web.CreateAuthorResponse{}, appError.NewAppError(
			http.StatusBadRequest,
			errAggregate,
			nil,
		)
	}

	author := domain.Author{
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
