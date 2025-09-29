package handler

import "net/http"

type UploadHandler interface {
	GetBookPresignedURL(w http.ResponseWriter, r *http.Request)
}
