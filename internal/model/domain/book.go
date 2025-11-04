package domain

import "time"

type Book struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	TotalPage     int       `json:"total_page"`
	AuthorId      string    `json:"author_id"`
	PhotoKey      string    `json:"photo_key,omitempty"`
	Status        string    `json:"status"`
	CompletedDate string    `json:"completed_date,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookWithURL struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	TotalPage     int       `json:"total_page"`
	AuthorId      string    `json:"author_id"`
	PhotoURL      string    `json:"photo_url,omitempty"`
	Status        string    `json:"status"`
	CompletedDate string    `json:"completed_date,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
