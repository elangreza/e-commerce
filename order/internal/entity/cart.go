package entity

import "github.com/google/uuid"

type Cart struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Items  []CartItem
}

type CartItem struct {
	ProductID string
	Quantity  int64
	Price     float64
	Stock     int64
}
