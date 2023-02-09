package types

import "github.com/solplaydev/solana/utils"

const (
	// 1 SOL = 1e9 lamports
	SOL uint64 = 1e9

	// SPL token default decimals
	SPLTokenDefaultDecimals uint8 = 9

	// SPL token default multiplier for decimals
	SPLTokenDefaultMultiplier uint64 = 1e9
)

// TokenAmount represents the amount of a token.
type TokenAmount struct {
	Amount         uint64  `json:"amount"`                     // amount in token lamports
	Decimals       uint8   `json:"decimals,omitempty"`         // number of decimals for the token; max 9
	UIAmount       float64 `json:"ui_amount,omitempty"`        // amount in token units, e.g. 1 SOL
	UIAmountString string  `json:"ui_amount_string,omitempty"` // amount in token units as a string, e.g. "1.000000000"
}

// TokenAmountFromLamports converts the given lamports to a token amount.
func NewTokenAmountFromLamports(lamports uint64, decimals uint8) TokenAmount {
	uiAmount := utils.AmountToFloat64(lamports, decimals)
	return TokenAmount{
		Amount:         lamports,
		Decimals:       decimals,
		UIAmount:       uiAmount,
		UIAmountString: utils.Float64ToString(uiAmount),
	}
}

// NewDefaultTokenAmount converts the given lamports to a token amount with the default decimals.
func NewDefaultTokenAmount(lamports uint64) TokenAmount {
	return NewTokenAmountFromLamports(lamports, SPLTokenDefaultDecimals)
}
