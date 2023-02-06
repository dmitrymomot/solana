package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
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

// FreezeTokenAccountParams are the parameters for the FreezeTokenAccount instruction.
type FreezeTokenAccountParams struct {
	FreezeAuth        common.PublicKey  // required; the account to authorize the freeze/unfreeze
	Mint              common.PublicKey  // required; the mint of the token account
	TokenAccount      *common.PublicKey // optional; the public key of account to freeze; if not set, the associated token account will be derived from the mint and token account owner.
	TokenAccountOwner *common.PublicKey // optional; the owner of the token account;
}

// Validate checks that the required fields of the params are set.
func (p FreezeTokenAccountParams) Validate() error {
	if p.FreezeAuth == (common.PublicKey{}) {
		return fmt.Errorf("freeze auth is required")
	}
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.TokenAccount != nil && *p.TokenAccount == (common.PublicKey{}) {
		return fmt.Errorf("invalid token account public key")
	}
	if p.TokenAccountOwner != nil && *p.TokenAccountOwner == (common.PublicKey{}) {
		return fmt.Errorf("invalid token account owner public key")
	}
	if p.TokenAccount == nil && p.TokenAccountOwner == nil {
		return fmt.Errorf("must be set at least one of token account or token account owner")
	}
	return nil
}

// FreezeTokenAccount freezes the specified token account.
func FreezeTokenAccount(params FreezeTokenAccountParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		if params.TokenAccount == nil && params.TokenAccountOwner != nil {
			ata, _, err := common.FindAssociatedTokenAddress(*params.TokenAccountOwner, params.Mint)
			if err != nil {
				return nil, fmt.Errorf("failed to find associated token address: %w", err)
			}
			params.TokenAccount = &ata
		}

		return []types.Instruction{
			token.FreezeAccount(token.FreezeAccountParam{
				Account: *params.TokenAccount,
				Mint:    params.Mint,
				Auth:    params.FreezeAuth,
			}),
		}, nil
	}
}

// UnfreezeTokenAccountParams are the parameters for the UnfreezeTokenAccount instruction.
type UnfreezeTokenAccountParams struct {
	FreezeAuth        common.PublicKey  // required; the account to authorize the freeze/unfreeze
	Mint              common.PublicKey  // required; the mint of the token account
	TokenAccount      *common.PublicKey // optional; the public key of account to freeze; if not set, the associated token account will be derived from the mint and token account owner.
	TokenAccountOwner *common.PublicKey // optional; the owner of the token account;
}

// Validate checks that the required fields of the params are set.
func (p UnfreezeTokenAccountParams) Validate() error {
	if p.FreezeAuth == (common.PublicKey{}) {
		return fmt.Errorf("freeze auth is required")
	}
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.TokenAccount != nil && *p.TokenAccount == (common.PublicKey{}) {
		return fmt.Errorf("invalid token account public key")
	}
	if p.TokenAccountOwner != nil && *p.TokenAccountOwner == (common.PublicKey{}) {
		return fmt.Errorf("invalid token account owner public key")
	}
	if p.TokenAccount == nil && p.TokenAccountOwner == nil {
		return fmt.Errorf("must be set at least one of token account or token account owner")
	}
	return nil
}

// UnfreezeTokenAccount unfreezes the specified token account.
func UnfreezeTokenAccount(params UnfreezeTokenAccountParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		if params.TokenAccount == nil && params.TokenAccountOwner != nil {
			ata, _, err := common.FindAssociatedTokenAddress(*params.TokenAccountOwner, params.Mint)
			if err != nil {
				return nil, fmt.Errorf("failed to find associated token address: %w", err)
			}
			params.TokenAccount = &ata
		}

		return []types.Instruction{
			token.ThawAccount(token.ThawAccountParam{
				Account: *params.TokenAccount,
				Mint:    params.Mint,
				Auth:    params.FreezeAuth,
			}),
		}, nil
	}
}
