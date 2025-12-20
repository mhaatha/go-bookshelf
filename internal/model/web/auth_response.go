package web

type CreateUserResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
}
