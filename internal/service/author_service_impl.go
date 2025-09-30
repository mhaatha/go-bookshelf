package service

import (
	"context"
	"database/sql"
	"errors"
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
	err = service.AuthorRepository.CheckByFullName(ctx, tx, request.FullName)
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

func (service *AuthorServiceImpl) GetAllAuthors(ctx context.Context, queries web.QueryParamsGetAuthors) ([]web.GetAuthorResponse, error) {
	// Validate queries
	err := service.Validate.Struct(queries)
	if err != nil {
		return []web.GetAuthorResponse{}, err
	}

	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return []web.GetAuthorResponse{}, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	// Call repository
	authors, err := service.AuthorRepository.FindAll(ctx, tx, queries.FullName, queries.Nationality)
	if err != nil {
		return []web.GetAuthorResponse{}, err
	}

	// No records return []
	if len(authors) == 0 {
		return []web.GetAuthorResponse{}, nil
	}

	return helper.ToGetAuthorsResponse(authors), nil
}

func (service *AuthorServiceImpl) GetAuthorById(ctx context.Context, pathValues web.PathParamsGetAuthor) (web.GetAuthorResponse, error) {
	// Validate path params
	err := service.Validate.Struct(pathValues)
	if err != nil {
		return web.GetAuthorResponse{}, err
	}

	// Open transcation
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return web.GetAuthorResponse{}, err
	}

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// Call repository
	author, err := service.AuthorRepository.FindById(ctx, tx, pathValues.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "id",
				Message: fmt.Sprintf("author with id '%s' is not found", pathValues.Id),
			})

			return web.GetAuthorResponse{}, appError.NewAppError(
				http.StatusNotFound,
				errAggregate,
				fmt.Errorf("author with id '%s' is not found", pathValues.Id),
			)
		}
		return web.GetAuthorResponse{}, err
	}

	return helper.ToGetAuthorResponse(author), nil
}

func (service *AuthorServiceImpl) UpdateAuthorById(ctx context.Context, pathValues web.PathParamsUpdateAuthor, request web.UpdateAuthorRequest) (web.UpdateAuthorResponse, error) {
	// Validate path params
	err := service.Validate.Struct(pathValues)
	if err != nil {
		return web.UpdateAuthorResponse{}, err
	}

	// Validate request body
	err = service.Validate.Struct(request)
	if err != nil {
		return web.UpdateAuthorResponse{}, err
	}

	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return web.UpdateAuthorResponse{}, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// Check if id is exists
	_, err = service.AuthorRepository.FindById(ctx, tx, pathValues.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "id",
				Message: fmt.Sprintf("author with id '%s' is not found", pathValues.Id),
			})

			// If id is not found, return earlier
			return web.UpdateAuthorResponse{}, appError.NewAppError(
				http.StatusNotFound,
				errAggregate,
				fmt.Errorf("author with id '%v' is not found", pathValues.Id),
			)
		} else {
			return web.UpdateAuthorResponse{}, err
		}
	}

	// Check if full_name already exists
	err = service.AuthorRepository.CheckByFullName(ctx, tx, request.FullName)
	if err != nil {
		errAggregate = append(errAggregate, appError.ErrAggregate{
			Field:   "full_name",
			Message: fmt.Sprintf("author %s is already exists", request.FullName),
		})
	}

	if len(errAggregate) != 0 {
		return web.UpdateAuthorResponse{}, appError.NewAppError(
			http.StatusBadRequest,
			errAggregate,
			nil,
		)
	}

	author := domain.Author{
		Id:          pathValues.Id,
		FullName:    request.FullName,
		Nationality: request.Nationality,
	}

	// Call repository
	author, err = service.AuthorRepository.Update(ctx, tx, pathValues.Id, author)
	if err != nil {
		return web.UpdateAuthorResponse{}, err
	}

	return helper.ToUpdateAuthorResponse(author), nil
}

func (service *AuthorServiceImpl) DeleteAuthorById(ctx context.Context, pathValues web.PathParamsDeleteAuthor) error {
	// Validate path params
	err := service.Validate.Struct(pathValues)
	if err != nil {
		return err
	}

	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx)

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// Check if id is exists
	_, err = service.AuthorRepository.FindById(ctx, tx, pathValues.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "id",
				Message: fmt.Sprintf("author with id '%s' is not found", pathValues.Id),
			})

			// if id not found, return earlier
			return appError.NewAppError(
				http.StatusNotFound,
				errAggregate,
				fmt.Errorf("author with id '%v' is not found", pathValues.Id),
			)
		} else {
			return err
		}
	}

	// Call repository
	err = service.AuthorRepository.Delete(ctx, tx, pathValues.Id)
	if err != nil {
		return err
	}

	return nil
}
