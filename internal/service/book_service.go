package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type BookService interface {
	CreateNewBook(ctx context.Context, request web.CreateBookRequest) (web.CreateBookResponse, error)
	GetAllBooks(ctx context.Context, queries web.QueryParamsGetBooks) ([]web.GetBookResponse, error)
	GetBookById(ctx context.Context, pathValues web.PathParamsGetBook) (web.GetBookResponse, error)
	UpdateBookById(ctx context.Context, pathValues web.PathParamsUpdateBook, request web.UpdateBookRequest) (web.UpdateBookResponse, error)
	DeleteBookById(ctx context.Context, pathValues web.PathParamsDeleteBook) error
}
