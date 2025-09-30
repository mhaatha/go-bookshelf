package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

type BookRepository interface {
	Save(ctx context.Context, tx pgx.Tx, book domain.Book) (domain.Book, error)
	CheckByNameAndAuthorId(ctx context.Context, tx pgx.Tx, name, authorId string) error
	FindAll(ctx context.Context, tx pgx.Tx, name, status, author_name string) ([]domain.Book, error)
	FindById(ctx context.Context, tx pgx.Tx, bookId string) (domain.Book, error)
}
