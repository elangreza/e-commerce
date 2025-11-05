package entity

import (
	"github/elangreza/e-commerce/pkg/converter"

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
	Price     converter.Money
	Stock     int64
}
