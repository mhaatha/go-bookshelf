package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// MinIO init
	minioClient, err := config.MinIOInit(cfg)
	if err != nil {
		slog.Error("failed to initialize MinIO client", "err", err)
		os.Exit(1)
	}

	// Main router
	mux := http.NewServeMux()

	// Author resources
	authorRepository := repository.NewAuthorRepository(db)
	authorService := service.NewAuthorService(authorRepository, db, validate)
	authorHandler := handler.NewAuthorHandler(authorService)

	// Author router
	router.AuthorRouter(authorHandler, mux)

	// Upload resources
	uploadService := service.NewUploadService(minioClient, cfg)
	uploadHandler := handler.NewUploadHandler(uploadService)

	// Upload router
	router.UploadRouter(uploadHandler, mux)

	// Book resources
	bookRepository := repository.NewBookRepository()
	bookService := service.NewBookService(bookRepository, authorService, db, validate, minioClient, cfg)
	bookhandler := handler.NewBookHandler(bookService)

	// Book router
	router.BookRouter(bookhandler, mux)

	// Server
	server := http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	if cfg.AppEnv == string(config.EnvProduction) {
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

		// Wait for signal
		<-signalChan

		// 10 seconds shut down period
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), config.ShutdownPeriod)
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
