package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

type AuthorRepository interface {
	Save(ctx context.Context, author domain.Author) (domain.Author, error)
	CheckByFullName(ctx context.Context, tx pgx.Tx, fullName string) error
	FindAll(ctx context.Context, tx pgx.Tx, fullName, nationality string) ([]domain.Author, error)
	FindById(ctx context.Context, tx pgx.Tx, authorId string) (domain.Author, error)
	Update(ctx context.Context, tx pgx.Tx, authorId string, author domain.Author) (domain.Author, error)
	Delete(ctx context.Context, tx pgx.Tx, authorId string) error
}
