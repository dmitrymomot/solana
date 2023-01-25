package instructions

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	"github.com/portto/solana-go-sdk/types"
)

// CreateAssociatedTokenAccountParam defines the parameters for creating an associated token account.
type CreateAssociatedTokenAccountParam struct {
	Funder common.PublicKey
	Owner  common.PublicKey
	Mint   common.PublicKey
}

// CreateAssociatedTokenAccount creates an associated token account for the given owner and mint.
func CreateAssociatedTokenAccount(params CreateAssociatedTokenAccountParam) InstructionFunc {
	return func() ([]types.Instruction, error) {
		ata, _, err := common.FindAssociatedTokenAddress(params.Owner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		return []types.Instruction{
			associated_token_account.CreateAssociatedTokenAccount(
				associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 params.Funder,
					Owner:                  params.Owner,
					Mint:                   params.Mint,
					AssociatedTokenAccount: ata,
				},
			),
		}, nil
	}
}
