package params

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"image_url"`
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
