package helper

import (
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

func ToCreateAuthorResponse(author domain.Author) web.CreateAuthorResponse {
	return web.CreateAuthorResponse{
		Id:          author.Id,
		FullName:    author.FullName,
		Nationality: author.Nationality,
		CreatedAt:   author.CreatedAt,
		UpdatedAt:   author.UpdatedAt,
	}
}
