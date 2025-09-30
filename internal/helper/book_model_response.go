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

func ToGetBookResponse(book domain.Book) web.GetBookResponse {
	return web.GetBookResponse{
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

func ToGetBooksResponse(books []domain.Book) []web.GetBookResponse {
	var bookResponses []web.GetBookResponse
	for _, book := range books {
		bookResponses = append(bookResponses, ToGetBookResponse(book))
	}
	return bookResponses
}
