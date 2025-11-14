package entity

import (
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/constanta"
	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID               `json:"id" db:"id"`
	Status        constanta.PaymentStatus `json:"status" db:"status"`
	TotalAmount   *gen.Money              `json:"total_amount" db:"total_amount"`
	TransactionID string                  `json:"transaction_id" db:"transaction_id"` // Add this field
	OrderID       string                  `json:"order_id" db:"order_id"`             // Link back to the order
	CreatedAt     time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at" db:"updated_at"`
}
