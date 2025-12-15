package params

import (
	errs "github.com/elangreza/e-commerce/api/internal/error"
	"github.com/google/uuid"
)

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}

func (a *AddToCartRequest) Validate() error {
	if a.ProductID == "" {
		return errs.ValidationError{Message: "product_id is required"}
	}

	if a.Quantity < 1 {
		return errs.ValidationError{Message: "quantity must be positive"}
	}

	return nil
}

type (
	GetCartItemsResponse struct {
		ProductID string `json:"product_id"`
		Quantity  int64  `json:"quantity"`
	}

	GetCartResponse struct {
		CartID string                 `json:"cart_id"`
		Items  []GetCartItemsResponse `json:"items"`
	}
)

type (
	CreateOrderRequest struct {
		IdempotencyKey string `json:"idempotency_key"`
	}

	CreateOrderItemsResponse struct {
		ProductID    string `json:"product_id"`
		Quantity     int64  `json:"quantity"`
		Name         string `json:"name"`
		PricePerUnit *Money `json:"price_per_unit"`
	}

	OrderResponse struct {
		OrderID       string                 `json:"order_id"`
		TotalAmount   *Money                 `json:"total_amount"`
		Status        string                 `json:"status"`
		Items         []GetCartItemsResponse `json:"items,omitempty"`
		TransactionID string                 `json:"transaction_id"`
	}
)

func (a *CreateOrderRequest) Validate() error {
	if a.IdempotencyKey == "" {
		return errs.ValidationError{Message: "idempotency_key is required"}
	}

	_, err := uuid.Parse(a.IdempotencyKey)
	if err != nil {
		return errs.ValidationError{Message: "not valid idempotency_key"}
	}

	return nil
}

type (
	GetOrderListRequest struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status"`
	}

	GetOrderListResponse struct {
		OrderList []OrderResponse `json:"order_list"`
	}
)
