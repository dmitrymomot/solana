package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dmitrymomot/solana/utils"
	"github.com/portto/solana-go-sdk/common"
	// "github.com/portto/solana-go-sdk/program/token"
)

const (
	// 1 SOL = 1e9 lamports
	SOL uint64 = 1e9

	// SPL token default decimals
	SPLTokenDefaultDecimals uint8 = 9

	// SPL token default multiplier for decimals
	SPLTokenDefaultMultiplier uint64 = 1e9

	// Wrapped SOL mint address
	WrappedSOLMint = "So11111111111111111111111111111111111111112"

	// DeprecatedTokenListPath is the default token list path
	DeprecatedTokenListPath = "https://raw.githubusercontent.com/solana-labs/token-list/main/src/tokens/solana.tokenlist.json"
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

type (
	// The account that holds the token
	TokenAccount struct {
		Pubkey           common.PublicKey  `json:"pubkey"`
		Mint             common.PublicKey  `json:"mint"`
		Owner            common.PublicKey  `json:"owner"`
		State            TokenAccountState `json:"state"`
		IsNative         bool              `json:"is_native"` // if is wrapped SOL, IsNative is the rent-exempt value
		Balance          TokenAmount       `json:"balance"`
		Delegate         *common.PublicKey `json:"delegate,omitempty"`
		DelegatedBalance *TokenAmount      `json:"delegated_balance,omitempty"`
	}

	TokenAccountState string

	/**
	RPCTokenAccount represents the following json data:
	{
		"pubkey": "DUNMHHh3qLwd7zVfckWHK7DoAk7jaeHiJgouVEQGraEe",
		"account": {
			"lamports": 2039280,
			"owner": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			"rentEpoch": 0,
			"data": {
				"parsed": {
					"info": {
						"isNative": false,
						"mint": "DCQqijZDH6a14o3wQcF8KzrBtsCjTBgfmagYYvQS8ihB",
						"owner": "FuQhSmAT6kAmmzCMiiYbzFcTQJFuu6raXAdCFibz4YPR",
						"state": "initialized",
						"tokenAmount": {
							"amount": "1",
							"decimals": 0,
							"uiAmount": 1,
							"uiAmountString": "1"
						},
						"delegate": "GKv5PeCxKBCDezo4FMVjjRbkUfoou9PRvPKdzaFEwjXi",
						"delegatedAmount": {
							"amount": "1",
							"decimals": 0,
							"uiAmount": 1,
							"uiAmountString": "1"
						},
					},
					"type": "account"
				},
				"program": "spl-token",
				"space": 165
			},
			"executable": false
		}
	}
	*/
	RPCTokenAccount struct {
		Pubkey  string `json:"pubkey"`
		Account struct {
			Lamports  uint64 `json:"lamports"`
			Owner     string `json:"owner"`
			RentEpoch uint64 `json:"rentEpoch"`
			Data      struct {
				Parsed struct {
					Info struct {
						IsNative    bool   `json:"isNative"`
						Mint        string `json:"mint"`
						Owner       string `json:"owner"`
						State       string `json:"state"`
						TokenAmount struct {
							Amount         string  `json:"amount"`
							Decimals       uint8   `json:"decimals"`
							UIAmount       float64 `json:"uiAmount"`
							UIAmountString string  `json:"uiAmountString"`
						} `json:"tokenAmount"`
						Delegate        *string `json:"delegate"`
						DelegatedAmount *struct {
							Amount         string  `json:"amount"`
							Decimals       uint8   `json:"decimals"`
							UIAmount       float64 `json:"uiAmount"`
							UIAmountString string  `json:"uiAmountString"`
						} `json:"delegatedAmount"`
					} `json:"info"`
					Type string `json:"type"`
				} `json:"parsed"`
				Program string `json:"program"`
				Space   uint64 `json:"space"`
			} `json:"data"`
			Executable bool `json:"executable"`
		} `json:"account"`
	}
)

// IsEmpty returns true if the token account is empty.
func (a TokenAccount) IsEmpty() bool {
	return a.Balance.Amount == 0
}

// IsNFT returns true if the token account is an NFT.
func (a TokenAccount) IsNFT() bool {
	return a.Balance.Decimals == 0 && a.Balance.Amount == 1
}

// IsFungibleToken returns true if the token account is a fungible token.
func (a TokenAccount) IsFungibleToken() bool {
	return a.Balance.Decimals > 0 && a.Balance.Amount > 0
}

// IsFungibleAsset returns true if the token account is a fungible asset.
func (a TokenAccount) IsFungibleAsset() bool {
	return a.Balance.Decimals == 0 && a.Balance.Amount > 1
}

// String returns the string representation of the token account state.
func (s TokenAccountState) String() string {
	return string(s)
}

// Predefined token account states.
const (
	TokenAccountStateUninitialized TokenAccountState = "uninitialized"
	TokenAccountStateInitialized   TokenAccountState = "initialized"
	TokenAccountFrozen             TokenAccountState = "frozen"
)

// TokenAccountStateMap is a map of token account states.
// var tokenAccountStates = map[token.TokenAccountState]TokenAccountState{
// 	token.TokenAccountStateUninitialized: TokenAccountStateUninitialized,
// 	token.TokenAccountStateInitialized:   TokenAccountStateInitialized,
// 	token.TokenAccountFrozen:             TokenAccountFrozen,
// }

// NewTokenAccount converts the given json encoded data to a token account.
func NewTokenAccount(data []byte) (TokenAccount, error) {
	rpcResponse := &RPCTokenAccount{}
	if err := json.Unmarshal(data, rpcResponse); err != nil {
		return TokenAccount{}, fmt.Errorf("could not parse token account rpc response data: %w", err)
	}

	amount, err := strconv.ParseUint(rpcResponse.Account.Data.Parsed.Info.TokenAmount.Amount, 10, 64)
	if err != nil {
		return TokenAccount{}, fmt.Errorf("could not parse token balance amount: %w", err)
	}
	balance := TokenAmount{
		Amount:         amount,
		Decimals:       rpcResponse.Account.Data.Parsed.Info.TokenAmount.Decimals,
		UIAmount:       rpcResponse.Account.Data.Parsed.Info.TokenAmount.UIAmount,
		UIAmountString: rpcResponse.Account.Data.Parsed.Info.TokenAmount.UIAmountString,
	}

	var (
		delegate        *common.PublicKey
		delegateBalance *TokenAmount
	)

	if rpcResponse.Account.Data.Parsed.Info.Delegate != nil &&
		*rpcResponse.Account.Data.Parsed.Info.Delegate != "" {
		delegate = utils.Pointer(common.PublicKeyFromString(*rpcResponse.Account.Data.Parsed.Info.Delegate))
	}
	if rpcResponse.Account.Data.Parsed.Info.DelegatedAmount != nil &&
		rpcResponse.Account.Data.Parsed.Info.DelegatedAmount.UIAmount > 0 {
		dAmount, err := strconv.ParseUint(rpcResponse.Account.Data.Parsed.Info.DelegatedAmount.Amount, 10, 64)
		if err != nil {
			return TokenAccount{}, fmt.Errorf("could not parse delegated balance amount: %w", err)
		}
		delegateBalance = &TokenAmount{
			Amount:         dAmount,
			Decimals:       rpcResponse.Account.Data.Parsed.Info.DelegatedAmount.Decimals,
			UIAmount:       rpcResponse.Account.Data.Parsed.Info.DelegatedAmount.UIAmount,
			UIAmountString: rpcResponse.Account.Data.Parsed.Info.DelegatedAmount.UIAmountString,
		}
	}

	return TokenAccount{
		Pubkey:           common.PublicKeyFromString(rpcResponse.Pubkey),
		Mint:             common.PublicKeyFromString(rpcResponse.Account.Data.Parsed.Info.Mint),
		Owner:            common.PublicKeyFromString(rpcResponse.Account.Data.Parsed.Info.Owner),
		State:            TokenAccountState(rpcResponse.Account.Data.Parsed.Info.State),
		IsNative:         rpcResponse.Account.Data.Parsed.Info.IsNative,
		Balance:          balance,
		Delegate:         delegate,
		DelegatedBalance: delegateBalance,
	}, nil
}
