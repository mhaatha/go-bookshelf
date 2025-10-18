package config

import (
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinIOInit(cfg *Config) (*minio.Client, error) {
	endpoint := cfg.MinIOEndpoint
	accessKeyId := cfg.MinIOAccessKeyId
	secretAccessKey := cfg.MinIOSecretAccessKey

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	if err != nil {
		slog.Error("failed to create MinIO client", "err", err)
		return nil, err
	}

	slog.Info("minio server connected successfully")

	return minioClient, nil
}
