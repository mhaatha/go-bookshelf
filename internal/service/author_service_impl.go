package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

func NewAuthorService(uow UnitOfWork, validate *validator.Validate) AuthorService {
	return &AuthorServiceImpl{
		UoW:      uow,
		Validate: validate,
	}
}

type AuthorServiceImpl struct {
	UoW      UnitOfWork
	Validate *validator.Validate
}

func (service *AuthorServiceImpl) CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error) {
	// Validate request body
	err := service.Validate.Struct(request)
	if err != nil {
		return web.CreateAuthorResponse{}, err
	}

	// Open transaction
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return web.CreateAuthorResponse{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// It creates a new instance of AuthorRepository
	authorRepo := tx.GetAuthorRepository()

	// Check if full_name already exists
	err = authorRepo.CheckByFullName(ctx, request.FullName)
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
	author, err = authorRepo.Save(ctx, author)
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
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// It creates a new instance of AuthorRepository
	authorRepo := tx.GetAuthorRepository()

	// Call repository
	authors, err := authorRepo.FindAll(ctx, queries.FullName, queries.Nationality)
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

	// Open transaction
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return web.GetAuthorResponse{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// It creates a new instance of AuthorRepository
	authorRepo := tx.GetAuthorRepository()

	// Call repository
	author, err := authorRepo.FindById(ctx, pathValues.Id)
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
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return web.UpdateAuthorResponse{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// It creates a new instance of AuthorRepository
	authorRepo := tx.GetAuthorRepository()

	// Check if id is exists
	author, err := authorRepo.FindById(ctx, pathValues.Id)
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
	err = authorRepo.CheckByFullName(ctx, request.FullName)
	if err != nil {
		if author.FullName != request.FullName {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "full_name",
				Message: fmt.Sprintf("author %s is already exists", request.FullName),
			})
		}
	}

	if len(errAggregate) != 0 {
		return web.UpdateAuthorResponse{}, appError.NewAppError(
			http.StatusBadRequest,
			errAggregate,
			nil,
		)
	}

	author = domain.Author{
		Id:          pathValues.Id,
		FullName:    request.FullName,
		Nationality: request.Nationality,
	}

	// Call repository
	author, err = authorRepo.Update(ctx, pathValues.Id, author)
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
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// It creates a new instance of AuthorRepository
	authorRepo := tx.GetAuthorRepository()

	// Check if id is exists
	_, err = authorRepo.FindById(ctx, pathValues.Id)
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
	err = authorRepo.Delete(ctx, pathValues.Id)
	if err != nil {
		return err
	}

	return nil
}
