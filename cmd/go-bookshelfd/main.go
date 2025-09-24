package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		os.Exit(1)
	}

	// Validator init
	validate := config.ValidatorInit()

	// Database init
	db, err := database.ConnectDB(cfg)
	if err != nil {
		slog.Error("failed connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Main router
	mux := http.NewServeMux()

	// Author resources
	authorRepository := repository.NewAuthorRepository()
	authorService := service.NewAuthorService(authorRepository, db, validate)
	authorHandler := handler.NewAuthorHandler(authorService)

	// Author Router
	router.AuthorRouter(authorHandler, mux)

	// Server
	server := http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	if cfg.AppEnv != string(config.EnvProduction) {
		go func() {
			slog.Info("starting server on :" + cfg.AppPort)
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				slog.Error("HTTP server error", "err", err)
				os.Exit(1)
			}
			slog.Info("stopped serving new connections")
		}()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		<-signalChan

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("error while shut the server down", "err", err)
			os.Exit(1)
		}

		slog.Info("server shut down gracefully")
	} else {
		slog.Info("starting server on :" + cfg.AppPort)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "err", err)
			os.Exit(1)
		}
	}
}
