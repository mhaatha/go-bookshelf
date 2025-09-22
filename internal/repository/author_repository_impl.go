package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

func NewAuthorRepository() AuthorRepository {
	return &AuthorRepositoryImpl{}
}

type AuthorRepositoryImpl struct{}

func (repository *AuthorRepositoryImpl) Save(ctx context.Context, tx pgx.Tx, author domain.Author) (domain.Author, error) {
	sqlQuery := `
	INSERT INTO authors (id, full_name, nationality)
	VALUES ($1, $2, $3)
	RETURNING created_at, updated_at
	`

	err := tx.QueryRow(
		ctx,
		sqlQuery,
		author.Id,
		author.FullName,
		author.Nationality,
	).Scan(
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		return domain.Author{}, err
	}

	return author, nil
}
