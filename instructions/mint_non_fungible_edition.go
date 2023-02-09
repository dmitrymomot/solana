package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// MintNonFungibleEditionParam defines the parameters for the MintNonFungibleEdition instruction.
type MintNonFungibleEditionParam struct {
	FeePayer           common.PublicKey // required; The wallet to pay the fees from
	MasterEditionMint  common.PublicKey // required; The master edition mint public key
	MasterEditionOwner common.PublicKey // required; The master edition owner public key
	EditionMint        common.PublicKey // required; The new edition mint public key
	EditionOwner       common.PublicKey // optional; The new edition owner public key; defaults to the master edition owner
}

// Validate checks the parameters for the MintNonFungibleEdition instruction.
func (p MintNonFungibleEditionParam) Validate() error {
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	if p.MasterEditionMint == (common.PublicKey{}) {
		return fmt.Errorf("master edition mint is required")
	}
	if p.MasterEditionOwner == (common.PublicKey{}) {
		return fmt.Errorf("master edition owner is required")
	}
	if p.EditionMint == (common.PublicKey{}) {
		return fmt.Errorf("print edition mint is required")
	}
	return nil
}

// MintNonFungibleEdition creates instructions for minting fungible tokens.
func MintNonFungibleEdition(params MintNonFungibleEditionParam) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}
		if params.EditionOwner == (common.PublicKey{}) {
			params.EditionOwner = params.MasterEditionOwner
		}

		masterOwnerAta, _, err := common.FindAssociatedTokenAddress(params.MasterEditionOwner, params.MasterEditionMint)
		if err != nil {
			return nil, fmt.Errorf("find associated token address for master edition mint: %w", err)
		}

		masterEditionPublicKey, err := token_metadata.DeriveEditionPubkey(params.MasterEditionMint)
		if err != nil {
			return nil, fmt.Errorf("derive master edition pubkey: %w", err)
		}

		masterMetaPublicKey, err := token_metadata.DeriveTokenMetadataPubkey(params.MasterEditionMint)
		if err != nil {
			return nil, fmt.Errorf("derive master metadata pubkey: %w", err)
		}

		newMintOwnerAta, _, err := common.FindAssociatedTokenAddress(params.EditionOwner, params.EditionMint)
		if err != nil {
			return nil, fmt.Errorf("find associated token address for new edition mint: %w", err)
		}

		newMintMetaPublicKey, err := token_metadata.DeriveTokenMetadataPubkey(params.EditionMint)
		if err != nil {
			return nil, fmt.Errorf("derive new edition metadata pubkey: %w", err)
		}

		newMintEditionPublicKey, err := token_metadata.DeriveEditionPubkey(params.EditionMint)
		if err != nil {
			return nil, fmt.Errorf("derive new edition pubkey: %w", err)
		}

		current, max, err := c.GetMasterEditionSupply(ctx, params.MasterEditionMint)
		if err != nil {
			return nil, fmt.Errorf("failed to get master edition supply: %w", err)
		}
		if current >= max {
			return nil, fmt.Errorf("master edition supply is already at max")
		}

		editionNumber := current + 1
		newMintEditionMark, err := token_metadata.DeriveEditionMarkerPubkey(params.MasterEditionMint, editionNumber)
		if err != nil {
			return nil, fmt.Errorf("derive new edition marker pubkey: %w", err)
		}

		rentExemption, err := c.GetMinimumBalanceForRentExemption(ctx, token.MintAccountSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get minimum balance for rent exemption: %w", err)
		}

		instructions := []types.Instruction{
			system.CreateAccount(system.CreateAccountParam{
				From:     params.FeePayer,
				New:      params.EditionMint,
				Owner:    common.TokenProgramID,
				Lamports: rentExemption,
				Space:    token.MintAccountSize,
			}),
			token.InitializeMint2(token.InitializeMint2Param{
				Decimals:   0,
				Mint:       params.EditionMint,
				MintAuth:   params.MasterEditionOwner,
				FreezeAuth: utils.Pointer(params.MasterEditionOwner),
			}),
			associated_token_account.CreateAssociatedTokenAccount(
				associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 params.FeePayer,
					Owner:                  params.EditionOwner,
					Mint:                   params.EditionMint,
					AssociatedTokenAccount: newMintOwnerAta,
				},
			),
			token.MintTo(token.MintToParam{
				Mint:   params.EditionMint,
				Auth:   params.MasterEditionOwner,
				To:     newMintOwnerAta,
				Amount: 1,
			}),
			metaplex_token_metadata.MintNewEditionFromMasterEditionViaToken(
				metaplex_token_metadata.MintNewEditionFromMasterEditionViaTokeParam{
					NewMetaData:                newMintMetaPublicKey,
					NewEdition:                 newMintEditionPublicKey,
					MasterEdition:              masterEditionPublicKey,
					NewMint:                    params.EditionMint,
					NewMintAuthority:           params.MasterEditionOwner,
					Payer:                      params.FeePayer,
					TokenAccountOwner:          params.MasterEditionOwner,
					TokenAccount:               masterOwnerAta,
					NewMetadataUpdateAuthority: params.MasterEditionOwner,
					MasterMetadata:             masterMetaPublicKey,

					EditionMark: newMintEditionMark,
					Edition:     editionNumber,
				},
			),
		}

		return instructions, nil
	}
}
