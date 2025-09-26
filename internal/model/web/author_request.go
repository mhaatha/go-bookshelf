package web

type CreateAuthorRequest struct {
	FullName    string `json:"full_name" validate:"required,min=3,max=255,validName"`
	Nationality string `json:"nationality" validate:"required,min=3,max=255,alpha"`
}

type QueryParamsGetAuthors struct {
	FullName    string `json:"full_name" validate:"omitempty,min=3,max=255,validName"`
	Nationality string `json:"nationality" validate:"omitempty,min=3,max=255,alpha"`
}

type PathParamsGetAuthor struct {
	Id string `json:"id" validate:"omitempty,uuid"`
}
