package entity

import (
	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
)

type Cart struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Items  []CartItem
}

type CartItem struct {
	ProductID string
	Quantity  int64
	Price     *gen.Money
	Stock     int64
}
