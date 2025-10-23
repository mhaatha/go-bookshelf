package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/repository"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetAuthorRepository() repository.AuthorRepository
	GetBookRepository() repository.BookRepository
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
}
