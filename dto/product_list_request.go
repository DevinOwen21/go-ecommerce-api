package dto

type ProductListRequest struct {
	Pagination PaginationRequest
	Search     string
	MinPrice   *float64
	MaxPrice   *float64
	SortBy     string
	Order      string
	Category   string
}
