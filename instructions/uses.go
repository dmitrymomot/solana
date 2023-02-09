package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
	commonx "github.com/solplaydev/solana/common"
	"github.com/solplaydev/solana/token_metadata"
)

// ApproveUseAuthorityParams are the parameters for the ApproveUseAuthority instruction.
type ApproveUseAuthorityParams struct {
	FeePayer        common.PublicKey // required; the account to pay the fees
	Mint            common.PublicKey // required; the token mint to approve use authority for
	MintOwner       common.PublicKey // required; the mint owner
	NewUseAuthority common.PublicKey // required; the new use authority to approve
	NumberOfUses    uint64           // required; the number of uses to approve for the new use authority
}

// Validate checks that the required fields of the params are set.
func (p ApproveUseAuthorityParams) Validate() error {
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.MintOwner == (common.PublicKey{}) {
		return fmt.Errorf("mint owner is required")
	}
	if p.NewUseAuthority == (common.PublicKey{}) {
		return fmt.Errorf("new use authority is required")
	}
	if p.NumberOfUses == 0 {
		return fmt.Errorf("number of uses is required and must be greater than 0")
	}
	return nil
}

// ApproveUseAuthority instructs the mint to approve the new use authority to use the specified number of tokens.
func ApproveUseAuthority(params ApproveUseAuthorityParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		useAuthorityRecord, err := metaplex_token_metadata.GetUseAuthorityRecord(params.Mint, params.NewUseAuthority)
		if err != nil {
			return nil, fmt.Errorf("failed to get use authority record: %w", err)
		}

		ownerAta, _, err := common.FindAssociatedTokenAddress(params.MintOwner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		metadata, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		burner, err := commonx.FindBurnerPubkey()
		if err != nil {
			return nil, fmt.Errorf("failed to find burner pubkey: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.ApproveUseAuthority(metaplex_token_metadata.ApproveUseAuthorityParam{
				UseAuthorityRecord: useAuthorityRecord,
				Owner:              params.MintOwner,
				Payer:              params.FeePayer,
				User:               params.NewUseAuthority,
				OwnerTokenAccount:  ownerAta,
				Metadata:           metadata,
				Mint:               params.Mint,
				Burner:             burner,
				NumberOfUses:       params.NumberOfUses,
			}),
		}, nil
	}
}

// RevokeUseAuthorityParams are the parameters for the RevokeUseAuthority instruction.
type RevokeUseAuthorityParams struct {
	Mint         common.PublicKey // required; the token mint to revoke use authority for
	MintOwner    common.PublicKey // required; the mint owner
	UseAuthority common.PublicKey // required; the use authority to revoke
}

// Validate checks that the required fields of the params are set.
func (p RevokeUseAuthorityParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.MintOwner == (common.PublicKey{}) {
		return fmt.Errorf("mint owner is required")
	}
	if p.UseAuthority == (common.PublicKey{}) {
		return fmt.Errorf("use authority is required")
	}
	return nil
}

// RevokeUseAuthority instructs the mint to revoke the use authority.
func RevokeUseAuthority(params RevokeUseAuthorityParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		useAuthorityRecord, err := metaplex_token_metadata.GetUseAuthorityRecord(params.Mint, params.UseAuthority)
		if err != nil {
			return nil, fmt.Errorf("failed to get use authority record: %w", err)
		}

		ownerAta, _, err := common.FindAssociatedTokenAddress(params.MintOwner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		metadata, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.RevokeUseAuthority(metaplex_token_metadata.RevokeUseAuthorityParam{
				UseAuthorityRecord: useAuthorityRecord,
				Owner:              params.MintOwner,
				User:               params.UseAuthority,
				OwnerTokenAccount:  ownerAta,
				Mint:               params.Mint,
				Metadata:           metadata,
			}),
		}, nil
	}
}

// UseTokenParams are the parameters for the UseToken instruction.
type UseTokenParams struct {
	FeePayer     common.PublicKey // required; the account to pay the fees
	Mint         common.PublicKey // required; the token mint to use
	MintOwner    common.PublicKey // required; the mint owner
	UseAuthority common.PublicKey // required; the use authority to use the token
}

// Validate checks that the required fields of the params are set.
func (p UseTokenParams) Validate() error {
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.MintOwner == (common.PublicKey{}) {
		return fmt.Errorf("mint owner is required")
	}
	if p.UseAuthority == (common.PublicKey{}) {
		return fmt.Errorf("use authority is required")
	}
	return nil
}

// UseToken instructs the mint to use the token.
func UseToken(params UseTokenParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		useAuthorityRecord, err := metaplex_token_metadata.GetUseAuthorityRecord(params.Mint, params.UseAuthority)
		if err != nil {
			return nil, fmt.Errorf("failed to get use authority record: %w", err)
		}

		ata, _, err := common.FindAssociatedTokenAddress(params.MintOwner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		metadata, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		burner, err := commonx.FindBurnerPubkey()
		if err != nil {
			return nil, fmt.Errorf("failed to find burner pubkey: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.Utilize(metaplex_token_metadata.UtilizeParam{
				Metadata:           metadata,
				TokenAccount:       ata,
				Mint:               params.Mint,
				UseAuthority:       params.UseAuthority,
				Owner:              params.MintOwner,
				UseAuthorityRecord: useAuthorityRecord,
				Burner:             burner,
				NumberOfUses:       1,
			}),
		}, nil
	}
}
