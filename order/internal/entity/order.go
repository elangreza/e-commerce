package entity

import (
	"github.com/elangreza/e-commerce/order/internal/constanta"
	"github.com/google/uuid"
)

type Order struct {
	ID     uuid.UUID
	UserID string
	Cart   Cart
	Total  int64
	Status constanta.OrderStatus
}

type OrderItem struct {
	ProductID string
	Quantity  int64
	Price     float64
	Stock     int64
}
