package repository

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

type BookRepository interface {
	Save(ctx context.Context, book domain.Book) (domain.Book, error)
	CheckByNameAndAuthorId(ctx context.Context, name, authorId string) error
	FindAll(ctx context.Context, name, status, author_name string) ([]domain.Book, error)
	FindById(ctx context.Context, bookId string) (domain.Book, error)
	Update(ctx context.Context, bookId string, book domain.Book) (domain.Book, error)
	Delete(ctx context.Context, bookId string) error
}
