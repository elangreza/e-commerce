package entity

type Cart struct {
	Items []CartItem
}

type CartItem struct {
	ProductID string
	Quantity  int64
	Price     int64
}
