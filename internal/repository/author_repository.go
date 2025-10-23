package repository

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/domain"
)

type AuthorRepository interface {
	Save(ctx context.Context, author domain.Author) (domain.Author, error)
	CheckByFullName(ctx context.Context, fullName string) error
	FindAll(ctx context.Context, fullName, nationality string) ([]domain.Author, error)
	FindById(ctx context.Context, authorId string) (domain.Author, error)
	Update(ctx context.Context, authorId string, author domain.Author) (domain.Author, error)
	Delete(ctx context.Context, authorId string) error
}
