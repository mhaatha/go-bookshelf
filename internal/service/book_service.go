package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type BookService interface {
	CreateNewBook(ctx context.Context, request web.CreateBookRequest) (web.CreateBookResponse, error)
}
