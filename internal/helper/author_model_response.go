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

func ToGetAuthorResponse(author domain.Author) web.GetAuthorResponse {
	return web.GetAuthorResponse{
		Id:          author.Id,
		FullName:    author.FullName,
		Nationality: author.Nationality,
		CreatedAt:   author.CreatedAt,
		UpdatedAt:   author.UpdatedAt,
	}
}

func ToGetAuthorsResponse(authors []domain.Author) []web.GetAuthorResponse {
	var authorResponses []web.GetAuthorResponse
	for _, author := range authors {
		authorResponses = append(authorResponses, ToGetAuthorResponse(author))
	}
	return authorResponses
}

func ToUpdateAuthorResponse(author domain.Author) web.UpdateAuthorResponse {
	return web.UpdateAuthorResponse{
		Id:          author.Id,
		FullName:    author.FullName,
		Nationality: author.Nationality,
		CreatedAt:   author.CreatedAt,
		UpdatedAt:   author.UpdatedAt,
	}
}
