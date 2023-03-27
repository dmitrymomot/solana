package instructions

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/solana/token_metadata"
	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
)

// BurnNftParams are the parameters for the BurnNft instruction.
type BurnNftParams struct {
	Mint           common.PublicKey  // required; the token mint to burn
	MintOwner      common.PublicKey  // required; the mint owner
	CollectionMint *common.PublicKey // optional; the collection mint
}

// Validate checks that the required fields of the params are set.
func (p BurnNftParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.MintOwner == (common.PublicKey{}) {
		return fmt.Errorf("mint owner is required")
	}
	if p.CollectionMint != nil && *p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("invalid collection mint public key")
	}
	return nil
}

// BurnNft burns the specified NFT.
func BurnNft(params BurnNftParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		ownerAta, _, err := common.FindAssociatedTokenAddress(params.MintOwner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		metadata, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		masterEdition, err := token_metadata.DeriveEditionPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive master edition pubkey: %w", err)
		}

		var collectionMetadata *common.PublicKey
		if params.CollectionMint != nil {
			cmd, err := token_metadata.DeriveTokenMetadataPubkey(*params.CollectionMint)
			if err != nil {
				return nil, fmt.Errorf("failed to derive collection metadata pubkey: %w", err)
			}
			collectionMetadata = &cmd
		}

		return []types.Instruction{
			metaplex_token_metadata.BurnNft(metaplex_token_metadata.BurnNftParam{
				Metadata:             metadata,
				Owner:                params.MintOwner,
				Mint:                 params.Mint,
				TokenAccount:         ownerAta,
				MasterEditionAccount: masterEdition,
				CollectionMetadata:   collectionMetadata,
			}),
		}, nil
	}
}

// BurnNftEditionParams are the parameters for the BurnNftEdition instruction.
type BurnNftEditionParams struct {
	MasterMint       common.PublicKey // required; the master mint of the edition to burn
	MasterMintOwner  common.PublicKey // required; the master mint owner
	EditionMint      common.PublicKey // required; the edition mint to burn
	EditionMintOwner common.PublicKey // required; the edition mint owner
}

// Validate checks that the required fields of the params are set.
func (p BurnNftEditionParams) Validate() error {
	if p.MasterMint == (common.PublicKey{}) {
		return fmt.Errorf("master mint is required")
	}
	if p.MasterMintOwner == (common.PublicKey{}) {
		return fmt.Errorf("master mint owner is required")
	}
	if p.EditionMint == (common.PublicKey{}) {
		return fmt.Errorf("edition mint is required")
	}
	if p.EditionMintOwner == (common.PublicKey{}) {
		return fmt.Errorf("edition mint owner is required")
	}
	return nil
}

// BurnNftEdition burns the specified NFT edition.
func BurnNftEdition(params BurnNftEditionParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		editionAta, _, err := common.FindAssociatedTokenAddress(params.EditionMintOwner, params.EditionMint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		masterAta, _, err := common.FindAssociatedTokenAddress(params.MasterMintOwner, params.MasterMint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		metadata, err := token_metadata.DeriveTokenMetadataPubkey(params.EditionMint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		masterEdition, err := token_metadata.DeriveEditionPubkey(params.MasterMint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive master edition pubkey: %w", err)
		}

		printEdition, err := token_metadata.DeriveEditionPubkey(params.EditionMint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive print edition pubkey: %w", err)
		}

		editionInfo, err := c.GetEditionInfo(ctx, params.EditionMint.ToBase58())
		if err != nil {
			return nil, fmt.Errorf("failed to get print edition info: %w", err)
		}

		editionMarker, err := token_metadata.DeriveEditionMarkerPubkey(params.MasterMint, editionInfo.Edition)
		if err != nil {
			return nil, fmt.Errorf("derive new edition marker pubkey: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.BurnEditionNft(metaplex_token_metadata.BurnEditionNftParam{
				Metadata:                  metadata,
				Owner:                     params.EditionMintOwner,
				PrintEditionMint:          params.EditionMint,
				MasterEditionMint:         params.MasterMint,
				PrintEditionTokenAccount:  editionAta,
				MasterEditionTokenAccount: masterAta,
				MasterEditionAccount:      masterEdition,
				PrintEditionAccount:       printEdition,
				EditionMarkerAccount:      editionMarker,
			}),
		}, nil
	}
}

// BurnTokenParams are the parameters for the BurnToken instruction.
type BurnTokenParams struct {
	Mint              common.PublicKey // optional; the mint to burn
	TokenAccountOwner common.PublicKey // optional; the token account owner
	Amount            uint64           // optional; the amount to burn in token units
}

// Validate checks that the required fields of the params are set.
func (p BurnTokenParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.TokenAccountOwner == (common.PublicKey{}) {
		return fmt.Errorf("token account owner is required")
	}
	return nil
}

// BurnToken burns the specified token.
func BurnToken(params BurnTokenParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate params: %w", err)
		}

		ata, _, err := common.FindAssociatedTokenAddress(params.TokenAccountOwner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		return []types.Instruction{
			token.Burn(token.BurnParam{
				Account: ata,
				Mint:    params.Mint,
				Auth:    params.TokenAccountOwner,
				Amount:  params.Amount,
			}),
		}, nil
	}
}
