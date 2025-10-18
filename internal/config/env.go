package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvProduction  Environment   = "production"
	ShutdownPeriod time.Duration = 10 * time.Second
)

type Config struct {
	AppEnv  string
	AppPort string

	DBURL string

	MinIOEndpoint        string
	MinIOAccessKeyId     string
	MinIOSecretAccessKey string

	BookBucket string
}

func LoadConfig() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv != string(EnvProduction) {
		err := godotenv.Load("../../.env")
		if err != nil {
			return &Config{}, err
		}
	}

	slog.Info("env loaded successfully")

	return &Config{
		AppEnv:               appEnv,
		DBURL:                os.Getenv("DB_URL"),
		AppPort:              os.Getenv("APP_PORT"),
		MinIOEndpoint:        os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKeyId:     os.Getenv("MINIO_ACCESS_KEY_ID"),
		MinIOSecretAccessKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
		BookBucket:           os.Getenv("BOOK_BUCKET"),
	}, nil
}
