package service

import (
	"context"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type UploadService interface {
	GetBookPresignedURL(ctx context.Context) (web.GetBookPresignedURLResponse, error)
}
