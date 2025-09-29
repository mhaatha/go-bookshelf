package helper

import (
	"github.com/mhaatha/go-bookshelf/internal/model/domain"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

func ToCreateBookResponse(book domain.Book) web.CreateBookResponse {
	return web.CreateBookResponse{
		Id:            book.Id,
		Name:          book.Name,
		TotalPage:     book.TotalPage,
		AuthorId:      book.AuthorId,
		PhotoURL:      book.PhotoURL,
		Status:        book.Status,
		CompletedDate: book.CompletedDate,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}
}
