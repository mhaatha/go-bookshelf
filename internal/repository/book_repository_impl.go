package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
	INSERT INTO books (id, name, total_page, author_id, photo_key, status, completed_date)
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
		book.PhotoKey,
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

func (repository *BookRepositoryImpl) FindAll(ctx context.Context, tx pgx.Tx, name, status, author_name string) ([]domain.Book, error) {
	baseQuery := `
	SELECT b.id, b.name, b.total_page, b.author_id, b.photo_key, 
       	   b.status, b.completed_date, b.created_at, b.updated_at 
	FROM books b
	JOIN authors a ON b.author_id = a.id
	`

	// Slice to aggregate arguments and WHERE condition dynamically
	args := []interface{}{}
	conditions := []string{}
	argCount := 1

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("b.name ILIKE $%d", argCount))
		args = append(args, "%"+name+"%")
		argCount++
	}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("b.status = $%d", argCount))
		args = append(args, status)
		argCount++
	}

	if author_name != "" {
		conditions = append(conditions, fmt.Sprintf("a.full_name ILIKE $%d", argCount))
		args = append(args, "%"+author_name+"%")
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

	books := make([]domain.Book, 0)

	for rows.Next() {
		var book domain.Book

		err := rows.Scan(
			&book.Id,
			&book.Name,
			&book.TotalPage,
			&book.AuthorId,
			&book.PhotoKey,
			&book.Status,
			&book.CompletedDate,
			&book.CreatedAt,
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (repository *BookRepositoryImpl) FindById(ctx context.Context, tx pgx.Tx, bookId string) (domain.Book, error) {
	sqlQuery := `
	SELECT name, total_page, author_id, photo_key, status, completed_date, created_at, updated_at
	FROM books
	WHERE id = $1
	`

	book := domain.Book{
		Id: bookId,
	}

	err := tx.QueryRow(ctx, sqlQuery, bookId).Scan(
		&book.Name,
		&book.TotalPage,
		&book.AuthorId,
		&book.PhotoKey,
		&book.Status,
		&book.CompletedDate,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return domain.Book{}, err
	}

	return book, nil
}

func (repository *BookRepositoryImpl) Update(ctx context.Context, tx pgx.Tx, bookId string, book domain.Book) (domain.Book, error) {
	sqlQuery := `
	UPDATE books
	SET name = $1, total_page = $2, author_id = $3, photo_key = $4, status = $5, completed_date = $6, updated_at = $7
	WHERE id = $8
	RETURNING created_at
	`

	updatedAt := time.Now()

	err := tx.QueryRow(
		ctx,
		sqlQuery,
		book.Name,
		book.TotalPage,
		book.AuthorId,
		book.PhotoKey,
		book.Status,
		book.CompletedDate,
		updatedAt,
		bookId,
	).Scan(
		&book.CreatedAt,
	)
	if err != nil {
		return domain.Book{}, err
	}

	return book, nil
}

func (repository *BookRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, bookId string) error {
	sqlQuery := `
	DELETE FROM books
	WHERE id = $1
	`

	_, err := tx.Exec(ctx, sqlQuery, bookId)
	if err != nil {
		return err
	}

	return nil
}
