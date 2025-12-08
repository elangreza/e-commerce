package params

import errs "github.com/elangreza/e-commerce/api/internal/error"

type SetWarehouseStatusRequest struct {
	WarehouseID int64 `json:"warehouse_id"`
	IsActive    bool  `json:"is_active"`
}

type TransferStockBetweenWarehouseRequest struct {
	FromWarehouseId int64  `json:"from_warehouse_id"`
	ToWarehouseId   int64  `json:"to_warehouse_id"`
	ProductId       string `json:"product_id"`
	Quantity        int64  `json:"quantity"`
}

func (rur *TransferStockBetweenWarehouseRequest) Validate() error {
	if rur.FromWarehouseId < 1 {
		return errs.ValidationError{Message: "from_warehouse_id must be larger than 0"}
	}
	if rur.ToWarehouseId < 1 {
		return errs.ValidationError{Message: "to_warehouse_id must be larger than 0"}
	}
	if rur.ToWarehouseId == rur.FromWarehouseId {
		return errs.ValidationError{Message: "to_warehouse_id and from_warehouse_id cannot be same"}
	}

	if rur.ProductId == "" {
		return errs.ValidationError{Message: "product is required"}
	}

	if rur.Quantity < 1 {
		return errs.ValidationError{Message: "quantity must be larger than 0"}
	}

	return nil
}
