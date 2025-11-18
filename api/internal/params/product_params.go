package params

type ListProductsRequest struct {
	Search string `json:"search"`
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sort_by"`
}

type ListProductsResponse struct {
	Search string `json:"search"`
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sort_by"`
}
