package dto

type PaginationRequestDTO struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type PaginationResultDTO[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}
