package instructions

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/types"
)

// CreateNonceAccountParams is the params for creating a nonce account.
type CreateNonceAccountParams struct {
	FeePayer               common.PublicKey // required; The fee payer public key.
	Nonce                  common.PublicKey // required; The nonce account public key.
	NonceAuth              common.PublicKey // optional; The nonce account authority public key; default is fee payer.
	NonceAccountMinBalance uint64           // required; The nonce account minimum balance.
}

// Validate validates the params.
func (p CreateNonceAccountParams) Validate() error {
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	if p.Nonce == (common.PublicKey{}) {
		return fmt.Errorf("nonce is required")
	}
	if p.NonceAccountMinBalance == 0 {
		return fmt.Errorf("nonce account minimum balance is required")
	}
	return nil
}

// CreateNonceAccount creates a nonce account.
func CreateNonceAccount(params CreateNonceAccountParams) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("create nonce account: %w", err)
		}

		if params.NonceAuth == (common.PublicKey{}) {
			params.NonceAuth = params.FeePayer
		}

		instructions := []types.Instruction{
			system.CreateAccount(system.CreateAccountParam{
				From:     params.FeePayer,
				New:      params.Nonce,
				Owner:    common.SystemProgramID,
				Lamports: params.NonceAccountMinBalance,
				Space:    system.NonceAccountSize,
			}),
			system.InitializeNonceAccount(system.InitializeNonceAccountParam{
				Nonce: params.Nonce,
				Auth:  params.NonceAuth,
			}),
		}

		return instructions, nil
	}
}
