package models

type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ResponseError struct {
	Message string      `json:"message"`
	Error   interface{} `json:"errors,omitempty"`
}

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalCount   int64 `json:"total_count"`
	TotalPages   int   `json:"total_pages"`
	FirstPage    int   `json:"first_page"`
	NextPage     int   `json:"next_page"`
	LastPage     int   `json:"last_page"`
	CurrentCount int   `json:"current_count"`
}
