package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

func NewAuthorRepository(db *pgxpool.Pool) AuthorRepository {
	return &AuthorRepositoryImpl{
		DB: db,
	}
}

type AuthorRepositoryImpl struct {
	DB *pgxpool.Pool
}

func (repository *AuthorRepositoryImpl) Save(ctx context.Context, author domain.Author) (domain.Author, error) {
	// Open transaction
	tx, err := repository.DB.Begin(ctx)
	if err != nil {
		return domain.Author{}, nil
	}
	defer helper.CommitOrRollback(ctx, tx)

	sqlQuery := `
	INSERT INTO authors (id, full_name, nationality)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(
		ctx,
		sqlQuery,
		uuid.NewString(),
		author.FullName,
		author.Nationality,
	).Scan(
		&author.Id,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		return domain.Author{}, err
	}

	return author, nil
}

func (repository *AuthorRepositoryImpl) CheckByFullName(ctx context.Context, tx pgx.Tx, fullName string) error {
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

func (repository *AuthorRepositoryImpl) FindAll(ctx context.Context, tx pgx.Tx, fullName, nationality string) ([]domain.Author, error) {
	baseQuery := `
	SELECT id, full_name, nationality, created_at, updated_at
	FROM authors
	`

	// Slice to aggregate arguments and WHERE condition dynamically
	args := []interface{}{}
	conditions := []string{}
	argCount := 1

	if fullName != "" {
		conditions = append(conditions, fmt.Sprintf("full_name ILIKE $%d", argCount))
		args = append(args, "%"+fullName+"%")
		argCount++
	}
	if nationality != "" {
		conditions = append(conditions, fmt.Sprintf("nationality ILIKE $%d", argCount))
		args = append(args, "%"+nationality+"%")
		argCount++
	}

	sqlQuery := baseQuery
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := tx.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := make([]domain.Author, 0)

	for rows.Next() {
		var author domain.Author

		err := rows.Scan(
			&author.Id,
			&author.FullName,
			&author.Nationality,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (repository *AuthorRepositoryImpl) FindById(ctx context.Context, tx pgx.Tx, authorId string) (domain.Author, error) {
	sqlQuery := `
	SELECT full_name, nationality, created_at, updated_at
	FROM authors
	WHERE id = $1
	`

	author := domain.Author{
		Id: authorId,
	}

	err := tx.QueryRow(ctx, sqlQuery, authorId).Scan(
		&author.FullName,
		&author.Nationality,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		return domain.Author{}, err
	}

	return author, nil
}

func (repository *AuthorRepositoryImpl) Update(ctx context.Context, tx pgx.Tx, authorId string, author domain.Author) (domain.Author, error) {
	sqlQuery := `
	UPDATE authors
	SET full_name = $1, nationality = $2, updated_at = $3
	WHERE id = $4
	RETURNING created_at
	`

	updatedAt := time.Now()

	err := tx.QueryRow(
		ctx,
		sqlQuery,
		author.FullName,
		author.Nationality,
		updatedAt,
		authorId,
	).Scan(
		&author.CreatedAt,
	)
	if err != nil {
		return domain.Author{}, err
	}

	return author, nil
}

func (repository *AuthorRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, authorId string) error {
	sqlQuery := `
	DELETE FROM authors
	WHERE id = $1
	`

	_, err := tx.Exec(ctx, sqlQuery, authorId)
	if err != nil {
		return err
	}

	return nil
}
