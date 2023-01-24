package utils

import (
	"fmt"
	"math"
	"strings"
)

// AmountToFloat64 converts amount lamports to float64 with given decimals.
func AmountToFloat64(amount uint64, decimals uint8) float64 {
	return float64(amount) / math.Pow10(int(decimals))
}

// AmountToUint64 converts amount from float64 to uint64 with given decimals.
func AmountToUint64(amount float64, decimals uint8) uint64 {
	return uint64(amount * math.Pow10(int(decimals)))
}

// AmountToString converts amount lamports to string with given decimals.
func AmountToString(amount uint64, decimals uint8) string {
	f := AmountToFloat64(amount, decimals)
	return Float64ToString(f)
}

// IntAmountToFloat64 converts int64 amount lamports to float64 with given decimals.
func IntAmountToFloat64(amount int64, decimals uint8) float64 {
	return float64(amount) / math.Pow10(int(decimals))
}

// Float64ToString converts float64 to string with minimum number of decimals.
// For example, 1.000000000 will be converted to "1", 1.100000000 will be converted to "1.1".
func Float64ToString(amount float64) string {
	// Convert to string with 9 decimals.
	s := fmt.Sprintf("%.9f", amount)

	// Trim trailing zeros.
	s = TrimRightZeros(s)

	// Trim trailing dot.
	s = strings.TrimRight(s, ".")

	return s
}
