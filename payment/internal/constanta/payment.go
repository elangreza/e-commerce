package constanta

import (
	"database/sql/driver"
	"fmt"
)

type PaymentStatus string

const (
	// waiting payment to be done
	WAITING PaymentStatus = "WAITING"
	// successfully paid
	PAID PaymentStatus = "PAID"
	// if payment amount is less or more the actual amount
	FAILED PaymentStatus = "FAILED"
	// if payment is expired
	CANCELLED PaymentStatus = "CANCELLED"
)

// Implement driver.Valuer interface for writing to database
func (ps PaymentStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// Implement sql.Scanner interface for reading from database
func (ps *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		*ps = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*ps = PaymentStatus(v)
	case []byte:
		*ps = PaymentStatus(v)
	case int64:
		*ps = PaymentStatus(fmt.Sprintf("%d", v))
	default:
		return fmt.Errorf("cannot scan %T into PaymentStatus", value)
	}

	return nil
}
