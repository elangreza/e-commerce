package constanta

import (
	"database/sql/driver"
	"fmt"
)

type OrderStatus string

// 'pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'failed'
const (
	OrderStatusPending       OrderStatus = "PENDING"
	OrderStatusStockReserved OrderStatus = "STOCK_RESERVED"
	OrderStatusCompleted     OrderStatus = "COMPLETED"
	OrderStatusCancelled     OrderStatus = "CANCELLED"
	OrderStatusFailed        OrderStatus = "FAILED"
)

// return string
func (ps OrderStatus) String() string {
	return string(ps)
}

// Implement driver.Valuer interface for writing to database
func (ps OrderStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// Implement sql.Scanner interface for reading from database
func (ps *OrderStatus) Scan(value interface{}) error {
	if value == nil {
		*ps = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*ps = OrderStatus(v)
	case []byte:
		*ps = OrderStatus(v)
	case int64:
		*ps = OrderStatus(fmt.Sprintf("%d", v))
	default:
		return fmt.Errorf("cannot scan %T into PaymentStatus", value)
	}

	return nil
}
