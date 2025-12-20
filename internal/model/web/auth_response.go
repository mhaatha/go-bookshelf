package web

import "time"

type CreateUserResponse struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
}
