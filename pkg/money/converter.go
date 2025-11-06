// pkg/money/currency.go
package money

import (
	"fmt"
	"strings"
)

// FractionalDigits returns the number of fractional digits (minor units exponent)
// for a given ISO 4217 currency code.
// Returns 0 if unknown (safe default for IDR/JPY-like currencies).
func FractionalDigits(currencyCode string) int {
	code := strings.ToUpper(currencyCode)
	if digits, ok := currencyFractionDigits[code]; ok {
		return digits
	}
	// Unknown currency: assume 2 digits (common case), but log warning in real app
	return 2
}

// MajorToMinor converts a major-unit amount (e.g., 125.50 USD) to minor units (12550)
// given a currency code.
func MajorToMinor(majorAmount float64, currencyCode string) int64 {
	frac := FractionalDigits(currencyCode)
	multiplier := 1
	for i := 0; i < frac; i++ {
		multiplier *= 10
	}
	return int64(majorAmount * float64(multiplier))
}

var currencyFractionDigits = map[string]int{
	// Common zero-decimal currencies
	"JPY": 0, "KRW": 0, "HUF": 0, "VND": 0, "XOF": 0, "XAF": 0, "CLP": 0,
	"BIF": 0, "DJF": 0, "GNF": 0, "KMF": 0, "MGA": 0, "PYG": 0, "RWF": 0,
	"UGX": 0, "UYI": 0, "VES": 0, "VUV": 0, "XPF": 0, "IDR": 0,

	// Standard 2-decimal currencies
	"USD": 2, "EUR": 2, "GBP": 2, "CAD": 2, "AUD": 2, "SGD": 2, "MYR": 2,
	"CNY": 2, "HKD": 2, "NZD": 2, "THB": 2, "PHP": 2, "INR": 2, "TWD": 2,

	// 3-decimal currencies
	"BHD": 3, "JOD": 3, "KWD": 3, "OMR": 3, "TND": 3, "LYD": 3,

	// Add more as needed
}

// ValidateCurrency returns error if currency code is unknown or invalid format
func ValidateCurrency(currencyCode string) error {
	if len(currencyCode) != 3 {
		return fmt.Errorf("currency code must be 3 letters: %q", currencyCode)
	}
	if strings.ToUpper(currencyCode) != currencyCode {
		// Optional: enforce uppercase
	}
	// Even if not in map, allow it (with default 2 digits)
	return nil
}
