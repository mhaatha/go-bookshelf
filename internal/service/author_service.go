package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type AuthorService interface {
	CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error)
	GetAllAuthors(ctx context.Context, queris web.QueryParamsGetAuthors) ([]web.GetAuthorResponse, error)
	GetAuthorById(ctx context.Context, pathValues web.PathParamsGetAuthor) (web.GetAuthorResponse, error)
	UpdateAuthorById(ctx context.Context, pathValues web.PathParamsGetAuthor, request web.UpdateAuthorRequest) (web.UpdateAuthorResponse, error)
	DeleteAuthorById(ctx context.Context, pathValues web.PathParamsGetAuthor) error
}
