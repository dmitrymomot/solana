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

// MintFungibleParam defines the parameters for the MintFungible instruction.
type MintFungibleParam struct {
	Mint     common.PublicKey  // required; The token mint public key
	MintTo   common.PublicKey  // required; The wallet to mint tokens to
	FeePayer *common.PublicKey // optional; The wallet to pay the fees from; default is MintTo

	Decimals      uint8  // required; The number of decimals the token has
	SupplyAmount  uint64 // required; The init supply of the token (in token minimal units), e.g: if you want to mint 10 tokens and decimals=9, amount=10*1e9/amount=10000000000; default is 0, then no tokens will be minted
	IsFixedSupply bool   // required; Whether the token has a fixed supply or not. If true, you cannot mint more tokens.
	MetadataURI   string // optional; URI of the token metadata; can be set later
	TokenName     string // optional; Name of the token; used for the token metadata if MetadataURI is not set.
	TokenSymbol   string // optional; Symbol of the token; used for the token metadata if MetadataURI is not set.
}

// Validate checks that the required fields of the params are set.
func (p MintFungibleParam) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("field Mint is required")
	}
	if p.MintTo == (common.PublicKey{}) {
		return fmt.Errorf("field MintTo is required")
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

// MintFungible creates instructions for minting fungible tokens or assets.
// The token mint account must be created before calling this function.
// To mint common fungible tokens, decimals must be greater than 0.
// If decimals is 0, the token is fungible asset.
func MintFungible(params MintFungibleParam) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("invalid params: %w", err)
		}

		if params.FeePayer == nil {
			params.FeePayer = &params.MintTo
		}

		metaPubkey, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to get token metadata pubkey: %w", err)
		}

		var metadataV2 metaplex_token_metadata.DataV2
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

			metadataV2 = metaplex_token_metadata.DataV2{
				Name:   md.Name,
				Symbol: md.Symbol,
				Uri:    params.MetadataURI,
			}
		} else {
			metadataV2 = metaplex_token_metadata.DataV2{
				Name:   params.TokenName,
				Symbol: params.TokenSymbol,
			}
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
				Decimals:   params.Decimals,
				Mint:       params.Mint,
				MintAuth:   params.MintTo,
				FreezeAuth: utils.Pointer(params.MintTo),
			}),
			metaplex_token_metadata.CreateMetadataAccountV2(metaplex_token_metadata.CreateMetadataAccountV2Param{
				Metadata:                metaPubkey,
				Mint:                    params.Mint,
				MintAuthority:           params.MintTo,
				Payer:                   *params.FeePayer,
				UpdateAuthority:         params.MintTo,
				UpdateAuthorityIsSigner: true,
				IsMutable:               true,
				Data:                    metadataV2,
			}),
		}

		if params.SupplyAmount > 0 {
			ownerAta, _, err := common.FindAssociatedTokenAddress(params.MintTo, params.Mint)
			if err != nil {
				return nil, fmt.Errorf("failed to find associated token address: %w", err)
			}

			instructions = append(
				instructions,
				associated_token_account.CreateAssociatedTokenAccount(
					associated_token_account.CreateAssociatedTokenAccountParam{
						Funder:                 *params.FeePayer,
						Owner:                  params.MintTo,
						Mint:                   params.Mint,
						AssociatedTokenAccount: ownerAta,
					},
				),
				token.MintToChecked(token.MintToCheckedParam{
					Mint:     params.Mint,
					Auth:     params.MintTo,
					Signers:  []common.PublicKey{},
					To:       ownerAta,
					Amount:   params.SupplyAmount,
					Decimals: params.Decimals,
				}),
			)
		}

		if params.IsFixedSupply && params.SupplyAmount > 0 {
			instructions = append(instructions, token.SetAuthority(token.SetAuthorityParam{
				Account:  params.Mint,
				AuthType: token.AuthorityTypeMintTokens,
				Auth:     params.MintTo,
				NewAuth:  nil,
				Signers:  []common.PublicKey{},
			}))
		}

		return instructions, nil
	}
}
