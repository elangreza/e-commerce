package constanta

import (
	"database/sql/driver"
	"fmt"
)

type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "PENDING"
	OrderStatusStockReserved OrderStatus = "STOCK_RESERVED"
	OrderStatusCompleted     OrderStatus = "COMPLETED" // when customer pays the bills
	OrderStatusCancelled     OrderStatus = "CANCELLED" // when the warehouse is inactive, cancelled all order
	OrderStatusFailed        OrderStatus = "FAILED"    // when customer exceeded the expiry order
)

// return string
func (ps OrderStatus) String() string {
	switch ps {
	case OrderStatusPending:
		return "PENDING"
	case OrderStatusStockReserved:
		return "STOCK_RESERVED"
	case OrderStatusCompleted:
		return "COMPLETED"
	case OrderStatusCancelled:
		return "CANCELLED"
	case OrderStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
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
