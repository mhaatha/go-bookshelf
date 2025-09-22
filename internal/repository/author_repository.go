package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

type AuthorRepository interface {
	Save(ctx context.Context, tx pgx.Tx, author domain.Author) (domain.Author, error)
}
