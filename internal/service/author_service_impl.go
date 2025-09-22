package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/repository"
)

func NewAuthorService(authorRepository repository.AuthorRepository, db *pgxpool.Pool) AuthorService {
	return &AuthorServiceImpl{
		AuthorRepository: authorRepository,
		DB:               db,
	}
}

type AuthorServiceImpl struct {
	AuthorRepository repository.AuthorRepository
	DB               *pgxpool.Pool
}

func (service *AuthorServiceImpl) CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error) {
	// Open transaction
	tx, err := service.DB.Begin(ctx)
	if err != nil {
		return web.CreateAuthorResponse{}, nil
	}
	defer helper.CommitOrRollback(ctx, tx)

	author := domain.Author{
		Id:          uuid.NewString(),
		FullName:    request.FullName,
		Nationality: request.Nationality,
	}

	// Call repository
	author, err = service.AuthorRepository.Save(ctx, tx, author)
	if err != nil {
		return web.CreateAuthorResponse{}, nil
	}

	return helper.ToCreateAuthorResponse(author), nil
}
