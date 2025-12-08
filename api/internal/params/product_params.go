package params

type ListProductsRequest struct {
	Search string `json:"search"`
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sort_by"`
}

type Money struct {
	Units        int64  `json:"units,omitempty"`
	CurrencyCode string `json:"currency_code,omitempty"`
}

type ListProductsResponse struct {
	Products   []*Product `json:"products,omitempty"`
	Total      int64      `json:"total,omitempty"`
	TotalPages int64      `json:"total_pages,omitempty"`
}

type GetProductsDetail struct {
	Ids       []string `json:"ids"`
	WithStock bool     `json:"with_stock"`
}

type WareHouse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type Product struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	ImageUrl       string      `json:"image_url"`
	Price          *Money      `json:"price"`
	Stock          int64       `json:"stock"`
	ShopID         int64       `json:"shop_id"`
	ShopName       string      `json:"shop_name"`
	ShopWareHouses []WareHouse `json:"shop_ware_houses,omitempty"`
}
