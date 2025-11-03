package entity

import (
	"time"

	"github.com/elangreza/e-commerce/order/internal/constanta"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID          uuid.UUID             `json:"id" db:"id"`
	UserID      string                `json:"user_id" db:"user_id"` // can be uuid
	Status      constanta.OrderStatus `json:"status" db:"status"`
	TotalAmount decimal.Decimal       `json:"total_amount" db:"total_amount"`
	Currency    string                `json:"currency" db:"currency"`
	CreatedAt   time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at" db:"updated_at"`
	ShippedAt   *time.Time            `json:"shipped_at,omitempty" db:"shipped_at"`
	CancelledAt *time.Time            `json:"cancelled_at,omitempty" db:"cancelled_at"`
}

type OrderItem struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	OrderID      uuid.UUID       `json:"order_id" db:"order_id"`
	ProductID    string          `json:"product_id" db:"product_id"`
	Quantity     int64           `json:"quantity" db:"quantity"`
	Name         string          `json:"name" db:"name"`
	PricePerUnit decimal.Decimal `json:"price_per_unit" db:"price_per_unit"`
	TotalPrice   decimal.Decimal `json:"total_price" db:"total_price"`
}
