package params

import "github.com/elangreza/e-commerce/gen"

type ProductResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       *gen.Money `json:"price"`
	ImageUrl    string     `json:"image_url"`
	Stock       int64      `json:"stock"`
}

type ListProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	TotalPages int64             `json:"total_pages"`
}

type GetProductRequest struct {
	ProductID string `json:"product_id"`
}

type GetProductsRequest struct {
	ProductIDs []string `json:"product_id"`
	WithStock  bool     `json:"withStock"`
}

type GetProductResponse struct {
	Product *ProductResponse `json:"product"`
}

type GetProductsResponse struct {
	Products []ProductResponse `json:"products"`
}
