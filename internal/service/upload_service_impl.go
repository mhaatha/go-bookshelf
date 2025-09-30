package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mhaatha/go-bookshelf/internal/config"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/minio/minio-go/v7"
)

func NewUploadService(minioClient *minio.Client, cfg *config.Config) UploadService {
	return &UploadServiceImpl{
		MinIOClient: minioClient,
		Config:      cfg,
	}
}

type UploadServiceImpl struct {
	MinIOClient *minio.Client
	Config      *config.Config
}

func (service *UploadServiceImpl) GetBookPresignedURL(ctx context.Context) (web.GetBookPresignedURLResponse, error) {
	// Initialize policy condition config
	policy := minio.NewPostPolicy()

	policy.SetBucket(service.Config.BookBucket)
	policy.SetKey(uuid.NewString() + ".jpg")
	policy.SetContentLengthRange(1024, 5*1024*1024) // 1KB - 5MB
	policy.SetContentType("image/jpeg")
	policy.SetExpires(time.Now().UTC().Add(5 * time.Minute))

	// Get the POST form key/value object:
	url, formData, err := service.MinIOClient.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		return web.GetBookPresignedURLResponse{}, err
	}

	return web.GetBookPresignedURLResponse{
		URL:      url.String(),
		FormData: formData,
	}, nil
}
