package web

type CreateAuthorRequest struct {
	FullName    string `json:"full_name" validate:"min=3,max=255"`
	Nationality string `json:"nationality" validate:"min=3,max=255"`
}
