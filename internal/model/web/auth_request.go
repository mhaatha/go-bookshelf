package web

type CreateUserRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=255,validName"`
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,min=8,validPassword"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,min=8,validPassword"`
}
