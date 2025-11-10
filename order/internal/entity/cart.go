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

func (c *Cart) GetProductIDs() []string {
	if len(c.Items) == 0 {
		return nil
	}

	res := []string{}
	for _, item := range c.Items {
		res = append(res, item.ProductID)
	}

	return res
}
