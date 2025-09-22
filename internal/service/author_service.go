package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type AuthorService interface {
	CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error)
}
