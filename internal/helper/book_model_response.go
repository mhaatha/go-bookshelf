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
		PhotoKey:      book.PhotoKey,
		Status:        book.Status,
		CompletedDate: book.CompletedDate,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}
}

func ToGetBookResponse(book domain.BookWithURL) web.GetBookResponse {
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

func ToGetBooksResponse(books []domain.BookWithURL) []web.GetBookResponse {
	var bookResponses []web.GetBookResponse
	for _, book := range books {
		bookResponses = append(bookResponses, ToGetBookResponse(book))
	}
	return bookResponses
}

func ToUpdateBookResponse(book domain.Book) web.UpdateBookResponse {
	return web.UpdateBookResponse{
		Id:            book.Id,
		Name:          book.Name,
		TotalPage:     book.TotalPage,
		AuthorId:      book.AuthorId,
		PhotoKey:      book.PhotoKey,
		Status:        book.Status,
		CompletedDate: book.CompletedDate,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}
}
