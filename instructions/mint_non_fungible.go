package instructions

import (
	"context"
	"fmt"
	"strings"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/metadata"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// MintNonFungibleParam defines the parameters for the MintNonFungible instruction.
type MintNonFungibleParam struct {
	Mint       common.PublicKey  // required; The token mint public key
	Owner      common.PublicKey  // required; The wallet to mint tokens to
	FeePayer   *common.PublicKey // optional; The wallet to pay the fees from; default is Owner
	Collection *common.PublicKey // optional; The collection mint public key
	Creators   *[]Creator        // optional; The creators of the token; FeePayer must be one of the creators; Default is mintTo:100 & FeePayer:0

	MaxEditionSupply     uint64 // optional; The max print edition supply; default is 0
	MetadataURI          string // optional; URI of the token metadata; can be set later
	TokenName            string // optional; Name of the token; used for the token metadata if MetadataURI is not set.
	TokenSymbol          string // optional; Symbol of the token; used for the token metadata if MetadataURI is not set.
	SellerFeeBasisPoints uint16 // optional; The seller fee basis points; default is 0

	UseMethod *token_metadata.TokenUseMethod // optional; The use method; default is nil
	UseLimit  *uint64                        // optional; The use times limit; default is 1; if UseMethod is nil, this field will be ignored
}

// Validate validates the parameters.
func (p MintNonFungibleParam) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("field Mint is required")
	}
	if p.Owner == (common.PublicKey{}) {
		return fmt.Errorf("field Owner is required")
	}
	if p.MetadataURI != "" && !strings.HasPrefix(p.MetadataURI, "http") {
		return fmt.Errorf("field MetadataURI must be a valid URI")
	}
	if p.FeePayer != nil && *p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("invalid fee payer public key")
	}
	if p.MetadataURI == "" && (p.TokenName == "" || p.TokenSymbol == "") {
		return fmt.Errorf("field TokenName and TokenSymbol are required if MetadataURI is not set")
	}
	if p.TokenName != "" && (len(p.TokenName) < 2 || len(p.TokenName) > 32) {
		return fmt.Errorf("token name must be between 2 and 32 characters")
	}
	if p.TokenSymbol != "" && (len(p.TokenSymbol) < 3 || len(p.TokenSymbol) > 10) {
		return fmt.Errorf("token symbol must be between 3 and 10 characters")
	}
	return nil
}

// MintNonFungible creates instructions for minting fungible tokens.
func MintNonFungible(params MintNonFungibleParam) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate parameters: %w", err)
		}

		if params.FeePayer == nil {
			params.FeePayer = &params.Owner
		}

		metadataV2 := metaplex_token_metadata.DataV2{
			Name:                 params.TokenName,
			Symbol:               params.TokenSymbol,
			Uri:                  params.MetadataURI,
			SellerFeeBasisPoints: params.SellerFeeBasisPoints,
			Collection: func() *metaplex_token_metadata.Collection {
				if params.Collection != nil {
					return &metaplex_token_metadata.Collection{
						Key: *params.Collection,
					}
				}
				return nil
			}(),
			Uses: func() *metaplex_token_metadata.Uses {
				if params.UseMethod != nil {
					if params.UseLimit == nil {
						params.UseLimit = utils.Pointer[uint64](1)
					}
					return &metaplex_token_metadata.Uses{
						UseMethod: params.UseMethod.ToMetadataUseMethod(),
						Remaining: *params.UseLimit,
						Total:     *params.UseLimit,
					}
				}
				return nil
			}(),
		}

		if params.MetadataURI != "" {
			md, err := metadata.MetadataFromURI(params.MetadataURI)
			if err != nil {
				return nil, fmt.Errorf("failed to get metadata from URI: %w", err)
			}

			if md.Name == "" || len(md.Name) < 2 || len(md.Name) > 32 {
				return nil, fmt.Errorf("metadata name must be between 2 and 32 characters")
			}
			if md.Symbol == "" || len(md.Symbol) < 2 || len(md.Symbol) > 10 {
				return nil, fmt.Errorf("metadata symbol must be between 2 and 10 characters")
			}

			metadataV2.Name = md.Name
			metadataV2.Symbol = md.Symbol
		}

		metaPubkey, err := metaplex_token_metadata.GetTokenMetaPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to get token metadata pubkey: %w", err)
		}

		tokenMasterEditionPubkey, err := metaplex_token_metadata.GetMasterEdition(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to get master edition pubkey: %w", err)
		}

		ownerAta, _, err := common.FindAssociatedTokenAddress(params.Owner, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address: %w", err)
		}

		// Preparing of NFT creators list
		if params.Creators != nil {
			creators := make([]metaplex_token_metadata.Creator, 0, len(*params.Creators))
			totalShare := uint8(0)
			feePayerInCreators := false
			for _, creator := range *params.Creators {
				if creator.Address.ToBase58() == params.FeePayer.ToBase58() {
					feePayerInCreators = true
				}
				totalShare += creator.Share
				creators = append(creators, metaplex_token_metadata.Creator{
					Address: creator.Address,
					Share:   creator.Share,
					Verified: func() bool {
						return creator.Address.ToBase58() == params.Owner.ToBase58()
					}(),
				})
			}

			if !feePayerInCreators {
				creators = append(creators, metaplex_token_metadata.Creator{
					Address:  *params.FeePayer,
					Share:    0,
					Verified: false,
				})
			}

			if totalShare != 100 {
				return nil, fmt.Errorf("creators share must be 100, got %d", totalShare)
			}

			metadataV2.Creators = &creators
		} else {
			creators := []metaplex_token_metadata.Creator{{
				Address:  params.Owner,
				Share:    100,
				Verified: true,
			}}

			if params.FeePayer.ToBase58() != params.Owner.ToBase58() {
				creators = append(creators, metaplex_token_metadata.Creator{
					Address:  *params.FeePayer,
					Share:    0,
					Verified: false,
				})
			}

			metadataV2.Creators = &creators
		}

		rentExemption, err := c.GetMinimumBalanceForRentExemption(ctx, token.MintAccountSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get minimum balance for rent exemption: %w", err)
		}

		instructions := []types.Instruction{
			system.CreateAccount(system.CreateAccountParam{
				From:     *params.FeePayer,
				New:      params.Mint,
				Owner:    common.TokenProgramID,
				Lamports: rentExemption,
				Space:    token.MintAccountSize,
			}),
			token.InitializeMint2(token.InitializeMint2Param{
				Decimals:   0,
				Mint:       params.Mint,
				MintAuth:   params.Owner,
				FreezeAuth: utils.Pointer(params.Owner),
			}),
			metaplex_token_metadata.CreateMetadataAccountV2(metaplex_token_metadata.CreateMetadataAccountV2Param{
				Metadata:                metaPubkey,
				Mint:                    params.Mint,
				MintAuthority:           params.Owner,
				Payer:                   *params.FeePayer,
				UpdateAuthority:         params.Owner,
				UpdateAuthorityIsSigner: true,
				IsMutable:               true,
				Data:                    metadataV2,
			}),
			associated_token_account.CreateAssociatedTokenAccount(
				associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 *params.FeePayer,
					Owner:                  params.Owner,
					Mint:                   params.Mint,
					AssociatedTokenAccount: ownerAta,
				},
			),
			token.MintToChecked(token.MintToCheckedParam{
				Mint:     params.Mint,
				Auth:     params.Owner,
				Signers:  []common.PublicKey{},
				To:       ownerAta,
				Amount:   1,
				Decimals: 0,
			}),
			metaplex_token_metadata.CreateMasterEditionV3(
				metaplex_token_metadata.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            params.Mint,
					UpdateAuthority: params.Owner,
					MintAuthority:   params.Owner,
					Metadata:        metaPubkey,
					Payer:           *params.FeePayer,
					MaxSupply:       utils.Pointer(params.MaxEditionSupply),
				},
			),
		}

		// Verify fee payer if it's not the owner
		if params.Owner.ToBase58() != params.FeePayer.ToBase58() {
			instructions = append(instructions, metaplex_token_metadata.SignMetadata(
				metaplex_token_metadata.SignMetadataParam{
					Metadata: metaPubkey,
					Creator:  *params.FeePayer,
				},
			))
		}

		return instructions, nil
	}
}
