package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
		uuid.NewString(),
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

func (repository *AuthorRepositoryImpl) GetByFullName(ctx context.Context, tx pgx.Tx, fullName string) error {
	sqlQuery := `
	SELECT 1 FROM authors 
	WHERE full_name = $1
	`

	var exists int
	err := tx.QueryRow(ctx, sqlQuery, fullName).Scan(&exists)
	if exists == 1 {
		return fmt.Errorf("author %v is already exists", fullName)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	return nil
}
