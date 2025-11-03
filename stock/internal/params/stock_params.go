package params

import "github.com/google/uuid"

type ReserveStock struct {
	Stocks []Stock `json:"stocks"`
}

type ReleaseStock struct {
	ReservedStockIDs []int64 `json:"reserved_stock_ids"`
}

type ConfirmStock struct {
	ReservedStockIDs []int64 `json:"reserved_stock_ids"`
}

type Stock struct {
	ProductID uuid.UUID `json:"id"`
	Quantity  int64     `json:"quantity"`
}
