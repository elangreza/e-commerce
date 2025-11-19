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
	ID        uuid.UUID
	CartID    uuid.UUID
	ProductID string
	Quantity  int64
	Name      string
	Price     *gen.Money
	// ActualStock is used to compare the stock and qty in cart
	// will be used when user is getting the cart
	// let's say get from DB in order is 3
	// get stock from actual product service is 2
	// so the result is deficit 1 qty
	// must be appeared warning message in FE
	ActualStock int64
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

func (c *Cart) GetGenCart() *gen.Cart {
	res := &gen.Cart{
		Id:    c.ID.String(),
		Items: []*gen.CartItem{},
	}

	if len(c.Items) == 0 {
		return res
	}

	for _, items := range c.Items {
		res.Items = append(res.Items, &gen.CartItem{
			ProductId:   items.ProductID,
			Quantity:    items.Quantity,
			Name:        items.ProductID,
			Price:       items.Price,
			ActualStock: items.ActualStock,
		})
	}

	return res
}
