package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{}
}

type BookRepositoryImpl struct{}

func (repository *BookRepositoryImpl) Save(ctx context.Context, tx pgx.Tx, book domain.Book) (domain.Book, error) {
	sqlQuery := `
	INSERT INTO books (id, name, total_page, author_id, photo_url, status, date_complete)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(
		ctx,
		sqlQuery,
		uuid.NewString(),
		book.Name,
		book.TotalPage,
		book.AuthorId,
		book.PhotoURL,
		book.Status,
		book.CompletedDate,
	).Scan(
		&book.Id,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return domain.Book{}, err
	}

	return book, nil
}

func (repository *BookRepositoryImpl) CheckByNameAndAuthorId(ctx context.Context, tx pgx.Tx, name, authorId string) error {
	sqlQuery := `
	SELECT 1 FROM books
	WHERE name = $1 AND author_id = $2
	`

	var exists int
	err := tx.QueryRow(ctx, sqlQuery, name, authorId).Scan(&exists)
	if exists == 1 {
		return fmt.Errorf("book with name %v and author_id '%v' is already exists", name, authorId)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
	}
	return err
}
