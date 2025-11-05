package converter

import (
	"errors"

	"github.com/elangreza/e-commerce/gen"
	"github.com/shopspring/decimal"
)

func getCurrencyDecimalPlaces(currency string) int32 {
	switch currency {
	case "JPY", "IDR", "KRW":
		return 0
	default:
		return 2 // USD, EUR, etc.
	}
}

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

// From gRPC to domain
func MoneyFromProto(p *gen.Money) (Money, error) {
	if p == nil {
		return Money{}, errors.New("money is nil")
	}

	decimals := getCurrencyDecimalPlaces(p.CurrencyCode)
	// Convert units to decimal: units / (10^decimals)
	amount := decimal.NewFromInt(p.Units).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt32(decimals)))

	return Money{
		Amount:   amount,
		Currency: p.CurrencyCode,
	}, nil
}

// From domain to gRPC
func (m Money) ToProto() *gen.Money {
	decimals := getCurrencyDecimalPlaces(m.Currency)
	// Multiply by 10^decimals and round to int64
	units := m.Amount.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt32(decimals))).IntPart()
	return &gen.Money{
		Units:        units,
		CurrencyCode: m.Currency,
	}
}
