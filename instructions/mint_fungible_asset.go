package instructions

import (
	"fmt"
	"strings"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/metadata"
)

// MintFungibleAssetParam defines the parameters for the MintFungibleAsset instruction.
type MintFungibleAssetParam struct {
	Mint                       common.PublicKey  // required; The token mint public key
	MintTo                     common.PublicKey  // required; The wallet to mint tokens to
	FeePayer                   *common.PublicKey // optional; The wallet to pay the fees from; default is MintTo
	MinBalanceForRentExemption uint64            // required; The minimum balance required to create the token account

	SupplyAmount  uint64 // required; The init supply of the token (in token minimal units), e.g: if you want to mint 10 tokens and decimals=9, amount=10*1e9/amount=10000000000; default is 0, then no tokens will be minted
	IsFixedSupply bool   // required; Whether the token has a fixed supply or not. If true, you cannot mint more tokens.
	MetadataURI   string // optional; URI of the token metadata; can be set later
}

// MintFungibleAsset creates instructions for minting fungible tokens.
func MintFungibleAsset(params MintFungibleAssetParam) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if params.Mint == (common.PublicKey{}) {
			return nil, fmt.Errorf("field Mint is required")
		}
		if params.MintTo == (common.PublicKey{}) {
			return nil, fmt.Errorf("field MintTo is required")
		}
		if params.FeePayer == nil {
			params.FeePayer = &params.MintTo
		}
		if params.MetadataURI != "" && !strings.HasPrefix(params.MetadataURI, "http") {
			return nil, fmt.Errorf("field MetadataURI must be a valid URI")
		} else if params.MetadataURI == "" {
			return nil, fmt.Errorf("field MetadataURI is required")
		}

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

		metaPubkey, err := metaplex_token_metadata.GetTokenMetaPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to get token metadata pubkey: %w", err)
		}

		instructions := []types.Instruction{
			system.CreateAccount(system.CreateAccountParam{
				From:     *params.FeePayer,
				New:      params.Mint,
				Owner:    common.TokenProgramID,
				Lamports: params.MinBalanceForRentExemption,
				Space:    token.MintAccountSize,
			}),
			token.InitializeMint(token.InitializeMintParam{
				Decimals: 0,
				Mint:     params.Mint,
				MintAuth: params.MintTo,
			}),
			metaplex_token_metadata.CreateMetadataAccountV2(metaplex_token_metadata.CreateMetadataAccountV2Param{
				Metadata:                metaPubkey,
				Mint:                    params.Mint,
				MintAuthority:           params.MintTo,
				Payer:                   *params.FeePayer,
				UpdateAuthority:         params.MintTo,
				UpdateAuthorityIsSigner: true,
				IsMutable:               true,
				Data: metaplex_token_metadata.DataV2{
					Name:   md.Name,
					Symbol: md.Symbol,
					Uri:    params.MetadataURI,
				},
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
					Decimals: 0,
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
