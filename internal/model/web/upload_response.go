package web

type GetBookPresignedURLResponse struct {
	URL      string      `json:"url"`
	FormData interface{} `json:"form_data"`
}
