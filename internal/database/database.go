package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/go-bookshelf/internal/config"
)

var DB *pgxpool.Pool

func ConnectDB(cfg *config.Config) error {
	dbPool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		return err
	}

	DB = dbPool

	// Check DB connection
	err = DB.Ping(context.Background())
	if err != nil {
		return err
	}

	slog.Info("connected to the database")
	return nil
}
