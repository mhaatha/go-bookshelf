package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvProduction Environment = "production"
)

type Config struct {
	DBURL string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") != string(EnvProduction) {
		err := godotenv.Load("../../.env")
		if err != nil {
			return &Config{}, err
		}
	}

	return &Config{
		DBURL: os.Getenv("DB_URL"),
	}, nil
}
