package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Picture     string    `json:"image_url"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type ListProductRequest struct {
	Search string `json:"search"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
	SortBy string `json:"sort_by"`
}
