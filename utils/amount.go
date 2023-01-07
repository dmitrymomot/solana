package utils

import "math"

// AmountToFloat64 converts amount lamports to float64 with given decimals.
func AmountToFloat64(amount uint64, decimals uint8) float64 {
	return float64(amount) / math.Pow10(int(decimals))
}

// AmountToUint64 converts amount from float64 to uint64 with given decimals.
func AmountToUint64(amount float64, decimals uint8) uint64 {
	return uint64(amount * math.Pow10(int(decimals)))
}

// IntAmountToFloat64 converts int64 amount lamports to float64 with given decimals.
func IntAmountToFloat64(amount int64, decimals uint8) float64 {
	return float64(amount) / math.Pow10(int(decimals))
}
