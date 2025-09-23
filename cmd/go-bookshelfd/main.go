package main

import (
	"log/slog"
	"net/http"

	"github.com/mhaatha/go-bookshelf/internal/config"
	"github.com/mhaatha/go-bookshelf/internal/database"
	"github.com/mhaatha/go-bookshelf/internal/handler"
	"github.com/mhaatha/go-bookshelf/internal/repository"
	"github.com/mhaatha/go-bookshelf/internal/router"
	"github.com/mhaatha/go-bookshelf/internal/service"
)

func main() {
	// Log init
	config.LogInit()

	// Config init
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("godotenv fails to load .env file", "err", err)
	}

	// Database init
	err = database.ConnectDB(cfg)
	if err != nil {
		slog.Error("failed connect to database", "err", err)
	}
	defer database.DB.Close()

	// Main router
	mux := http.NewServeMux()

	// Author resources
	authorRepository := repository.NewAuthorRepository()
	authorService := service.NewAuthorService(authorRepository, database.DB)
	authorHandler := handler.NewAuthorHandler(authorService)

	// Author Router
	router.AuthorRouter(authorHandler, mux)

	// Server
	server := http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	slog.Info("starting server on :" + cfg.AppPort)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start the server", "err", err)
	}
}
