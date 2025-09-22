package handler

import "net/http"

type AuthorHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}
