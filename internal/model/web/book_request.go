package web

type CreateBookRequest struct {
	Name          string `json:"name" validate:"required,min=3,max=255"`
	TotalPage     int    `json:"total_page" validate:"required,number,min=1,max=12000"`
	AuthorId      string `json:"author_id" validate:"required,uuid"`
	PhotoKey      string `json:"photo_key" validate:"required,min=3,max=255"`
	Status        string `json:"status" validate:"required,bookStatus"`
	CompletedDate string `json:"completed_date" validate:"omitempty,datetime=2006-01-02"`
}

type QueryParamsGetBooks struct {
	Status     string `json:"status" validate:"omitempty,bookStatus"`
	Name       string `json:"name" validate:"omitempty,min=3,max=255"`
	AuthorName string `json:"author_name" validate:"omitempty,min=3,max=255,validName"`
}

type PathParamsGetBook struct {
	Id string `json:"id" validate:"omitempty,uuid"`
}

type PathParamsUpdateBook struct {
	Id string `json:"id" validate:"omitempty,uuid"`
}

type UpdateBookRequest struct {
	Name          string `json:"name" validate:"required,min=3,max=255"`
	TotalPage     int    `json:"total_page" validate:"required,number,min=1,max=12000"`
	AuthorId      string `json:"author_id" validate:"required,uuid"`
	PhotoKey      string `json:"photo_key" validate:"required,min=3,max=255"`
	Status        string `json:"status" validate:"required,bookStatus"`
	CompletedDate string `json:"completed_date" validate:"omitempty,datetime=2006-01-02"`
}

type PathParamsDeleteBook struct {
	Id string `json:"id" validate:"omitempty,uuid"`
}
