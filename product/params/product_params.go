package params

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Picture     string  `json:"image_url"`
}

type ListProductsRequest struct {
	Search string `json:"search"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
	SortBy string `json:"sort_by"`
}

func (r *ListProductsRequest) Validate() error {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 10
	}
	if r.SortBy == "" {
		r.SortBy = "updated_at"
	}
	return nil
}

type ListProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	TotalPages int64             `json:"total_pages"`
}

type GetProductRequest struct {
	ProductID string `json:"product_id"`
}

type GetProductResponse struct {
	Product *ProductResponse `json:"product"`
}
