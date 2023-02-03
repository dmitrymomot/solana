package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/token"
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
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
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

// CreateAssociatedTokenAccountIfNotExists creates an associated token account for
// the given owner and mint if it does not exist.
func CreateAssociatedTokenAccountIfNotExists(params CreateAssociatedTokenAccountParam) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		ata, _, err := common.FindAssociatedTokenAddress(params.Owner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		if info, err := c.GetTokenAccountInfo(ctx, ata.ToBase58()); err == nil {
			if info.Mint.ToBase58() == params.Mint.ToBase58() {
				return nil, nil
			}
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

// FreezeDelegatedAccountParams are the parameters for the FreezeDelegatedAccount instruction.
type FreezeDelegatedAccountParams struct {
	FeePayer           common.PublicKey // required; the account to pay the fees
	FreezeTokenAccount common.PublicKey // required; the public key of account to freeze
}

// Validate checks that the required fields of the params are set.
func (p FreezeDelegatedAccountParams) Validate() error {
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	if p.FreezeTokenAccount == (common.PublicKey{}) {
		return fmt.Errorf("freeze token account is required")
	}
	return nil
}

// FreezeDelegatedAccount freezes the specified delegated account.
func FreezeDelegatedAccount(params FreezeDelegatedAccountParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.FreezeDelegatedAccount(metaplex_token_metadata.FreezeDelegatedAccountParam{
				Delegate:     common.PublicKey{},
				TokenAccount: params.FreezeTokenAccount,
				Edition:      common.PublicKey{},
				Mint:         common.PublicKey{},
			}),
		}, nil
	}
}

// UnfreezeDelegatedAccountParams are the parameters for the UnfreezeDelegatedAccount instruction.
type UnfreezeDelegatedAccountParams struct {
	FeePayer             common.PublicKey // required; the account to pay the fees
	UnfreezeTokenAccount common.PublicKey // required; the public key of account to unfreeze
}

// Validate checks that the required fields of the params are set.
func (p UnfreezeDelegatedAccountParams) Validate() error {
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	if p.UnfreezeTokenAccount == (common.PublicKey{}) {
		return fmt.Errorf("unfreeze token account is required")
	}
	return nil
}

// UnfreezeDelegatedAccount unfreezes the specified delegated account.
func UnfreezeDelegatedAccount(params UnfreezeDelegatedAccountParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.ThawDelegatedAccount(metaplex_token_metadata.ThawDelegatedAccountParam{
				Delegate:     common.PublicKey{},
				TokenAccount: params.UnfreezeTokenAccount,
				Edition:      common.PublicKey{},
				Mint:         common.PublicKey{},
			}),
		}, nil
	}
}

// CloseTokenAccountParams are the parameters for the CloseTokenAccount instruction.
type CloseTokenAccountParams struct {
	Owner             common.PublicKey  // required; the owner of the token account
	CloseTokenAccount *common.PublicKey // required if Mint is empty; the public key of account to close
	Mint              *common.PublicKey // required if CloseTokenAccount is empty; the mint of the token account
}

// Validate checks that the required fields of the params are set.
func (p CloseTokenAccountParams) Validate() error {
	if p.Owner == (common.PublicKey{}) {
		return fmt.Errorf("owner is required")
	}
	if p.CloseTokenAccount != nil && *p.CloseTokenAccount == (common.PublicKey{}) {
		return fmt.Errorf("invalid close token account public key")
	}
	if p.Mint != nil && *p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("invalid mint public key")
	}
	if p.CloseTokenAccount == nil && p.Mint == nil {
		return fmt.Errorf("one of close token account or mint must be set")
	}
	return nil
}

// CloseTokenAccount closes the specified token account.
func CloseTokenAccount(params CloseTokenAccountParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		if params.CloseTokenAccount == nil && params.Mint != nil {
			ata, _, err := common.FindAssociatedTokenAddress(params.Owner, *params.Mint)
			if err != nil {
				return nil, fmt.Errorf("failed to find associated token address: %w", err)
			}
			params.CloseTokenAccount = &ata
		}

		return []types.Instruction{
			token.CloseAccount(token.CloseAccountParam{
				Account: *params.CloseTokenAccount,
				Auth:    params.Owner,
				To:      params.Owner,
			}),
		}, nil
	}
}

// CloseMintAccountParams are the parameters for the CloseMintAccount instruction.
type CloseMintAccountParams struct {
	UpdateAuthority common.PublicKey // required; the owner of the mint account
	Mint            common.PublicKey // required; the public key of mint account to close
}

// Validate checks that the required fields of the params are set.
func (p CloseMintAccountParams) Validate() error {
	if p.UpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("update authority is required")
	}
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	return nil
}
