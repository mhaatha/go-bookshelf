package config

import (
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
	DBURL   string
	AppPort string
}

func LoadConfig() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv != string(EnvProduction) {
		err := godotenv.Load("../../.env")
		if err != nil {
			return &Config{}, err
		}
	}

	return &Config{
		AppEnv:  appEnv,
		DBURL:   os.Getenv("DB_URL"),
		AppPort: os.Getenv("APP_PORT"),
	}, nil
}
