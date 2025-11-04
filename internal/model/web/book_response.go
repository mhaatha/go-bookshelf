package web

import "time"

type CreateBookResponse struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	TotalPage     int       `json:"total_page"`
	AuthorId      string    `json:"author_id"`
	PhotoKey      string    `json:"photo_key"`
	Status        string    `json:"status"`
	CompletedDate string    `json:"completed_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GetBookResponse struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	TotalPage     int       `json:"total_page"`
	AuthorId      string    `json:"author_id"`
	PhotoURL      string    `json:"photo_url"`
	Status        string    `json:"status"`
	CompletedDate string    `json:"completed_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateBookResponse struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	TotalPage     int       `json:"total_page"`
	AuthorId      string    `json:"author_id"`
	PhotoKey      string    `json:"photo_key"`
	Status        string    `json:"status"`
	CompletedDate string    `json:"completed_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
