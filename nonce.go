package solana

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/types"
)

// CreateNonceAccountParams is the params for creating a nonce account.
type CreateNonceAccountParams struct {
	Base58FeePayerAddr  string
	Base58NonceAddr     string
	Base58NonceAuthAddr string
}

// CreateNonceAccount creates a nonce account.
// base58FeePayerAddr is the base58 encoded fee payer address.
// Returns the base64 encoded transaction, or an error.
// !!! This transaction must be signed by both: the fee payer and the nonce account.
func (c *Client) CreateNonceAccount(ctx context.Context, params CreateNonceAccountParams) ([]byte, error) {
	nonceAccountMinimumBalance, err := c.solana.GetMinimumBalanceForRentExemption(ctx, system.NonceAccountSize)
	if err != nil {
		return nil, ErrGetMinimumBalanceForRentExemption
	}

	feePayerPublicKey := common.PublicKeyFromString(params.Base58FeePayerAddr)
	noncePublicKey := common.PublicKeyFromString(params.Base58NonceAddr)
	nonceAuthPublicKey := common.PublicKeyFromString(params.Base58NonceAuthAddr)

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		Base58FeePayerAddr: params.Base58FeePayerAddr,
		Instructions: []types.Instruction{
			system.CreateAccount(system.CreateAccountParam{
				From:     feePayerPublicKey,
				New:      noncePublicKey,
				Owner:    common.SystemProgramID,
				Lamports: nonceAccountMinimumBalance,
				Space:    system.NonceAccountSize,
			}),
			system.InitializeNonceAccount(system.InitializeNonceAccountParam{
				Nonce: noncePublicKey,
				Auth:  nonceAuthPublicKey,
			}),
		},
	})
	if err != nil {
		return nil, err
	}

	return txb, nil
}
