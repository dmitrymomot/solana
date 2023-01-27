package instructions

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
)

// VerifyCreatorParams is the params for VerifyCreator
type VerifyCreatorParams struct {
	Mint    common.PublicKey // required; The mint of the token
	Creator common.PublicKey // required; The creator of the token
}

// Validate validates the params.
func (p VerifyCreatorParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.Creator == (common.PublicKey{}) {
		return fmt.Errorf("creator is required")
	}
	return nil
}

// VerifyCreator verifies the creator of the token metadata.
func VerifyCreator(params VerifyCreatorParams) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("verify creator: %w", err)
		}

		tokenMetadataPubkey, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		instructions := []types.Instruction{
			metaplex_token_metadata.SignMetadata(
				metaplex_token_metadata.SignMetadataParam{
					Metadata: tokenMetadataPubkey,
					Creator:  params.Creator,
				},
			),
		}

		return instructions, nil
	}
}

// RemoveCreatorVerificationParams is the params for RemoveCreatorVerification
type RemoveCreatorVerificationParams struct {
	Mint    common.PublicKey // required; The mint of the token
	Creator common.PublicKey // required; The creator of the token
}

// Validate validates the params.
func (p RemoveCreatorVerificationParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.Creator == (common.PublicKey{}) {
		return fmt.Errorf("creator is required")
	}
	return nil
}

// RemoveCreatorVerification removes the creator verification of the token metadata.
func RemoveCreatorVerification(params RemoveCreatorVerificationParams) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("remove creator verification: %w", err)
		}

		tokenMetadataPubkey, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.RemoveCreatorVerification(
				metaplex_token_metadata.RemoveCreatorVerificationParam{
					Metadata: tokenMetadataPubkey,
					Creator:  params.Creator,
				},
			),
		}, nil
	}
}
