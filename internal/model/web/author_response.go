package web

import "time"

type CreateAuthorResponse struct {
	Id          string    `json:"id"`
	FullName    string    `json:"full_name"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetAuthorResponse struct {
	Id          string    `json:"id"`
	FullName    string    `json:"full_name"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
