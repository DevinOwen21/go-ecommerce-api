package dto

type ProductListRequest struct {
	Pagination PaginationRequest
	Search     string
	MinPrice   *int
	MaxPrice   *int
	SortBy     string
	Order      string
	Category   string
}
