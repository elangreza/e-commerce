package constanta

type OrderStatus string

// 'pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'failed'
const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusFailed    OrderStatus = "FAILED"
)
