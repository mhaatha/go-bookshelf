package handler

import (
	"log/slog"
	"net/http"

	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/helper"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
	"github.com/mhaatha/go-bookshelf/internal/service"
)

func NewUploadHandler(uploadService service.UploadService) UploadHandler {
	return &UploadHandlerImpl{
		UploadService: uploadService,
	}
}

type UploadHandlerImpl struct {
	UploadService service.UploadService
}

func (handler *UploadHandlerImpl) GetBookPresignedURL(w http.ResponseWriter, r *http.Request) {
	// Call the service
	presignedURLResponse, err := handler.UploadService.GetBookPresignedURL(r.Context())
	if err != nil {
		appError.ResponseServiceErrorHandler(w, err, "failed to get presigned url")
		return
	}

	// Log the info
	slog.Info("request handled",
		"method", r.Method,
		"endpoint", r.URL,
		"status", http.StatusOK,
	)

	// Write and send the response
	helper.WriteToResponseBody(w, http.StatusOK, web.WebSuccessResponse{
		Message: "Success get presigned URL",
		Data: web.GetBookPresignedURLResponse{
			URL:      presignedURLResponse.URL,
			FormData: presignedURLResponse.FormData,
		},
	})
}
