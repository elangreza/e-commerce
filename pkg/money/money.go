// pkg/money/money.go
package money

import (
	"fmt"
	"strings"

	"github.com/elangreza/e-commerce/gen" // adjust import as needed
	"github.com/shopspring/decimal"
)

// FromProto converts gen.Money to internal representation (just returns copy)
// You could add validation here.
func FromProto(m *gen.Money) (*gen.Money, error) {
	if m == nil {
		return nil, fmt.Errorf("money is nil")
	}
	if err := ValidateCurrency(m.CurrencyCode); err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}
	if m.Units < 0 {
		return nil, fmt.Errorf("money units cannot be negative")
	}
	return &gen.Money{
		Units:        m.Units,
		CurrencyCode: strings.ToUpper(m.CurrencyCode),
	}, nil
}

// New creates a valid gen.Money from minor units and currency
func New(units int64, currency string) (*gen.Money, error) {
	if err := ValidateCurrency(currency); err != nil {
		return nil, err
	}
	if units < 0 {
		return nil, fmt.Errorf("units must be non-negative")
	}
	return &gen.Money{
		Units:        units,
		CurrencyCode: strings.ToUpper(currency),
	}, nil
}

// FromMajorAmount creates Money from a major-unit amount (e.g., 125.50) as string
// Uses decimal for exact parsing.
func FromMajorAmount(amountStr, currency string) (*gen.Money, error) {
	if err := ValidateCurrency(currency); err != nil {
		return nil, err
	}

	d, err := decimal.NewFromString(amountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	if d.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("amount must be non-negative")
	}

	fracDigits := FractionalDigits(currency)
	multiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(fracDigits)))
	minorUnits := d.Mul(multiplier).Round(0)
	units := minorUnits.IntPart()

	return New(units, currency)
}

// ToMajorString converts Money back to major-unit string (e.g., "125.50")
func ToMajorString(m *gen.Money) (string, error) {
	m, err := FromProto(m)
	if err != nil {
		return "", err
	}

	fracDigits := FractionalDigits(m.CurrencyCode)
	d := decimal.NewFromInt(m.Units)

	if fracDigits == 0 {
		return d.String(), nil
	}

	divisor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(fracDigits)))
	major := d.Div(divisor)
	return major.String(), nil
}

// Equals checks if two Money values are equal (same units and currency)
func Equals(a, b *gen.Money) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Units == b.Units && strings.EqualFold(a.CurrencyCode, b.CurrencyCode)
}

// Add adds two Money values (must have same currency)
func Add(a, b *gen.Money) (*gen.Money, error) {
	if !strings.EqualFold(a.CurrencyCode, b.CurrencyCode) {
		return nil, fmt.Errorf("cannot add different currencies: %s vs %s", a.CurrencyCode, b.CurrencyCode)
	}
	return New(a.Units+b.Units, a.CurrencyCode)
}

// MultiplyByInt multiplies money by an integer (e.g., for quantity)
func MultiplyByInt(m *gen.Money, n int64) (*gen.Money, error) {
	if n < 0 {
		return nil, fmt.Errorf("multiplier must be non-negative")
	}
	return New(m.Units*n, m.CurrencyCode)
}
