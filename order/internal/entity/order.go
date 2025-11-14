package entity

import (
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/constanta"
	"github.com/google/uuid"
)

type Order struct {
	IdempotencyKey uuid.UUID             `json:"idempotency_key" db:"idempotency_key"`
	ID             uuid.UUID             `json:"id" db:"id"`
	UserID         uuid.UUID             `json:"user_id" db:"user_id"` // can be uuid
	Status         constanta.OrderStatus `json:"status" db:"status"`
	TotalAmount    *gen.Money            `json:"total_amount" db:"total_amount"`
	// TransactionID is available after payment is processed, and successfully created
	TransactionID string     `json:"transaction_id" db:"transaction_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	ShippedAt     *time.Time `json:"shipped_at,omitempty" db:"shipped_at"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`

	Items []OrderItem
}

type OrderItem struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	OrderID           uuid.UUID  `json:"order_id" db:"order_id"`
	ProductID         string     `json:"product_id" db:"product_id"`
	Name              string     `json:"name" db:"name"`
	PricePerUnit      *gen.Money `json:"price_per_unit" db:"price_per_unit"`
	Quantity          int64      `json:"quantity" db:"quantity"`
	TotalPricePerUnit *gen.Money `json:"total_price" db:"total_price"`
}
