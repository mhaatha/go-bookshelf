package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/repository"
)

func NewBookService(bookRepository repository.BookRepository, authorService AuthorService, db *pgxpool.Pool, validate *validator.Validate) BookService {
	return &BookServiceImpl{
		BookRepository: bookRepository,
		AuthorService:  authorService,
		DB:             db,
		Validate:       validate,
	}
}

type BookServiceImpl struct {
	BookRepository repository.BookRepository
	AuthorService  AuthorService
	DB             *pgxpool.Pool
	Validate       *validator.Validate
}

func (service *BookServiceImpl) CreateNewBook(ctx context.Context, request web.CreateBookRequest) (web.CreateBookResponse, error) {
	// Validate request body
	err := service.Validate.Struct(request)
	if err != nil {
		return web.CreateBookResponse{}, err
	}

	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return web.CreateBookResponse{}, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	// errAggregate aggregates errors from user bad request
	errAggregate := []appError.ErrAggregate{}

	// Check if author_id exists
	_, err = service.AuthorService.GetAuthorById(ctx, web.PathParamsGetAuthor{Id: request.AuthorId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "author_id",
				Message: fmt.Sprintf("author with id '%v' is not found", request.AuthorId),
			})
		} else {
			return web.CreateBookResponse{}, err
		}
	}

	// Check if there is a book with the same name and the same author_id
	err = service.BookRepository.CheckByNameAndAuthorId(ctx, tx, request.Name, request.AuthorId)
	if err != nil {
		errAggregate = append(errAggregate, appError.ErrAggregate{
			Field:   "name",
			Message: fmt.Sprintf("%v with author_id '%v' is already exists", request.Name, request.AuthorId),
		})
	}

	if len(errAggregate) != 0 {
		return web.CreateBookResponse{}, appError.NewAppError(
			http.StatusBadRequest,
			errAggregate,
			nil,
		)
	}

	// Parse completed_date to time.Time manually
	t, err := time.Parse("2006-01-02", request.CompletedDate)
	if err != nil {
		return web.CreateBookResponse{}, err
	}

	book := domain.Book{
		Name:          request.Name,
		TotalPage:     request.TotalPage,
		AuthorId:      request.AuthorId,
		PhotoURL:      request.PhotoURL,
		Status:        request.Status,
		CompletedDate: t,
	}

	// Call repository
	book, err = service.BookRepository.Save(ctx, tx, book)
	if err != nil {
		return web.CreateBookResponse{}, err
	}

	return helper.ToCreateBookResponse(book), nil
}

func (service *BookServiceImpl) GetAllBooks(ctx context.Context, queries web.QueryParamsGetBooks) ([]web.GetBookResponse, error) {
	// Validate queries
	err := service.Validate.Struct(queries)
	if err != nil {
		return []web.GetBookResponse{}, err
	}

	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return []web.GetBookResponse{}, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	// Call repository
	books, err := service.BookRepository.FindAll(ctx, tx, queries.Name, queries.Status, queries.AuthorName)
	if err != nil {
		return []web.GetBookResponse{}, err
	}

	// No records return []
	if len(books) == 0 {
		return []web.GetBookResponse{}, nil
	}

	return helper.ToGetBooksResponse(books), nil
}
