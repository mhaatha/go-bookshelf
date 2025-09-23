package web

type CreateAuthorRequest struct {
	FullName    string `json:"full_name" validate:"required,min=3,max=255,validName"`
	Nationality string `json:"nationality" validate:"required,min=3,max=255,alpha"`
}
