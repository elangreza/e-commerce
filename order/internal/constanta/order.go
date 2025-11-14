package constanta

type OrderStatus string

// 'pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'failed'
const (
	OrderStatusPending       OrderStatus = "PENDING"
	OrderStatusStockReserved OrderStatus = "STOCK_RESERVED"
	OrderStatusCompleted     OrderStatus = "COMPLETED"
	OrderStatusCancelled     OrderStatus = "CANCELLED"
	OrderStatusFailed        OrderStatus = "FAILED"
)
