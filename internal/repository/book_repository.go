package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

type BookRepository interface {
	Save(ctx context.Context, tx pgx.Tx, book domain.Book) (domain.Book, error)
	CheckByNameAndAuthorId(ctx context.Context, tx pgx.Tx, name, authorId string) error
}
