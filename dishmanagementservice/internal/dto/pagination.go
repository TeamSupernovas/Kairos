package dto

type PaginatedResponse[T any] struct {
	TotalCount		int		`json:"totalCount"`
	TotalPages		int		`json:"totalPages"`
	Page			int		`json:"page"`
	PageSize		int		`json:"pageSize"`
	Data[]			[]T		`json:"data"`
}