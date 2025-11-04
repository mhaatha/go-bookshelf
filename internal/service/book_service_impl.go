package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mhaatha/go-bookshelf/internal/config"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/minio/minio-go/v7"
)

func NewBookService(uow UnitOfWork, authorService AuthorService, validate *validator.Validate, minioClient *minio.Client, cfg *config.Config) BookService {
	return &BookServiceImpl{
		UoW:           uow,
		AuthorService: authorService,
		Validate:      validate,
		MiniIOClient:  minioClient,
		Config:        cfg,
	}
}

type BookServiceImpl struct {
	UoW           UnitOfWork
	AuthorService AuthorService
	Validate      *validator.Validate
	MiniIOClient  *minio.Client
	Config        *config.Config
}

func (service *BookServiceImpl) CreateNewBook(ctx context.Context, request web.CreateBookRequest) (web.CreateBookResponse, error) {
	// Validate request body
	err := service.Validate.Struct(request)
	if err != nil {
		return web.CreateBookResponse{}, err
	}

	// Open transaction
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return web.CreateBookResponse{}, err
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

	// It creates a new instance of BookRepository
	bookRepo := tx.GetBookRepository()

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
	err = bookRepo.CheckByNameAndAuthorId(ctx, request.Name, request.AuthorId)
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

	book := domain.Book{
		Name:          request.Name,
		TotalPage:     request.TotalPage,
		AuthorId:      request.AuthorId,
		PhotoKey:      request.PhotoKey,
		Status:        request.Status,
		CompletedDate: request.CompletedDate,
	}

	// Call repository
	book, err = bookRepo.Save(ctx, book)
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
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return []web.GetBookResponse{}, err
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

	// It creates a new instance of BookRepository
	bookRepo := tx.GetBookRepository()

	// Call repository
	books, err := bookRepo.FindAll(ctx, queries.Name, queries.Status, queries.AuthorName)
	if err != nil {
		return []web.GetBookResponse{}, err
	}

	// No records return []
	if len(books) == 0 {
		return []web.GetBookResponse{}, nil
	}

	booksWithURL := []domain.BookWithURL{}

	for _, book := range books {
		presignedURL, err := service.MiniIOClient.PresignedGetObject(ctx, service.Config.BookBucket, book.PhotoKey, 24*time.Hour, nil)
		if err != nil {
			return []web.GetBookResponse{}, err
		}

		booksWithURL = append(booksWithURL, domain.BookWithURL{
			Id:            book.Id,
			Name:          book.Name,
			TotalPage:     book.TotalPage,
			AuthorId:      book.AuthorId,
			PhotoURL:      presignedURL.String(),
			Status:        book.Status,
			CompletedDate: book.CompletedDate,
			CreatedAt:     book.CreatedAt,
			UpdatedAt:     book.UpdatedAt,
		})
	}

	return helper.ToGetBooksResponse(booksWithURL), nil
}

func (service *BookServiceImpl) GetBookById(ctx context.Context, pathValues web.PathParamsGetBook) (web.GetBookResponse, error) {
	// Validate path params
	err := service.Validate.Struct(pathValues)
	if err != nil {
		return web.GetBookResponse{}, err
	}

	// Open transcation
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return web.GetBookResponse{}, err
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

	// It creates a new instance of BookRepository
	bookRepo := tx.GetBookRepository()

	// Call repository
	book, err := bookRepo.FindById(ctx, pathValues.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "id",
				Message: fmt.Sprintf("book with id '%s' is not found", pathValues.Id),
			})

			return web.GetBookResponse{}, appError.NewAppError(
				http.StatusNotFound,
				errAggregate,
				fmt.Errorf("book with id '%s' is not found", pathValues.Id),
			)
		}
		return web.GetBookResponse{}, err
	}

	// Create presigned URL for GET object
	presignedURL, err := service.MiniIOClient.PresignedGetObject(ctx, service.Config.BookBucket, book.PhotoKey, 24*time.Hour, nil)
	if err != nil {
		return web.GetBookResponse{}, err
	}

	bookWithURL := domain.BookWithURL{
		Id:            book.Id,
		Name:          book.Name,
		TotalPage:     book.TotalPage,
		AuthorId:      book.AuthorId,
		PhotoURL:      presignedURL.String(),
		Status:        book.Status,
		CompletedDate: book.CompletedDate,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}

	return helper.ToGetBookResponse(bookWithURL), nil
}

func (service *BookServiceImpl) UpdateBookById(ctx context.Context, pathValues web.PathParamsUpdateBook, request web.UpdateBookRequest) (web.UpdateBookResponse, error) {
	// Validate path params
	err := service.Validate.Struct(pathValues)
	if err != nil {
		return web.UpdateBookResponse{}, err
	}

	// Validate request body
	err = service.Validate.Struct(request)
	if err != nil {
		return web.UpdateBookResponse{}, err
	}

	// Open transaction
	tx, err := service.UoW.Begin(ctx)
	if err != nil {
		return web.UpdateBookResponse{}, err
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

	// It creates a new instance of BookRepository
	bookRepo := tx.GetBookRepository()

	// Check if id is exists
	_, err = bookRepo.FindById(ctx, pathValues.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "id",
				Message: fmt.Sprintf("book with id '%s' is not found", pathValues.Id),
			})

			// If id is not found, return earlier
			return web.UpdateBookResponse{}, appError.NewAppError(
				http.StatusNotFound,
				errAggregate,
				fmt.Errorf("book with id '%v' is not found", pathValues.Id),
			)
		} else {
			return web.UpdateBookResponse{}, err
		}
	}

	// Check if author_id exists
	_, err = service.AuthorService.GetAuthorById(ctx, web.PathParamsGetAuthor{Id: request.AuthorId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "author_id",
				Message: fmt.Sprintf("author with id '%v' is not found", request.AuthorId),
			})
		} else {
			return web.UpdateBookResponse{}, err
		}
	}

	// Check if there is a book with the same name and the same author_id
	err = bookRepo.CheckByNameAndAuthorId(ctx, request.Name, request.AuthorId)
	if err != nil {
		errAggregate = append(errAggregate, appError.ErrAggregate{
			Field:   "name",
			Message: fmt.Sprintf("%v with author_id '%v' is already exists", request.Name, request.AuthorId),
		})
	}

	if len(errAggregate) != 0 {
		return web.UpdateBookResponse{}, appError.NewAppError(
			http.StatusBadRequest,
			errAggregate,
			nil,
		)
	}

	book := domain.Book{
		Id:            pathValues.Id,
		Name:          request.Name,
		TotalPage:     request.TotalPage,
		AuthorId:      request.AuthorId,
		PhotoKey:      request.PhotoKey,
		Status:        request.Status,
		CompletedDate: request.CompletedDate,
	}

	// Call repository
	book, err = bookRepo.Update(ctx, pathValues.Id, book)
	if err != nil {
		return web.UpdateBookResponse{}, err
	}

	return helper.ToUpdateBookResponse(book), nil
}

func (service *BookServiceImpl) DeleteBookById(ctx context.Context, pathValues web.PathParamsDeleteBook) error {
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

	// It creates a new instance of BookRepository
	bookRepo := tx.GetBookRepository()

	// Check if id is exists
	_, err = bookRepo.FindById(ctx, pathValues.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errAggregate = append(errAggregate, appError.ErrAggregate{
				Field:   "id",
				Message: fmt.Sprintf("book with id '%s' is not found", pathValues.Id),
			})

			// if id not found, return earlier
			return appError.NewAppError(
				http.StatusNotFound,
				errAggregate,
				fmt.Errorf("book with id '%v' is not found", pathValues.Id),
			)
		} else {
			return err
		}
	}

	// Call repository
	err = bookRepo.Delete(ctx, pathValues.Id)
	if err != nil {
		return err
	}

	return nil
}
