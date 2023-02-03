package instructions

import (
	"context"
	"fmt"
	"strings"

	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/metadata"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// UpdateMetadataParams is the params for UpdateMetadata
type UpdateMetadataParams struct {
	Mint               common.PublicKey  // required; The mint of the token
	UpdateAuthority    common.PublicKey  // required; The update authority of the token
	NewUpdateAuthority *common.PublicKey // optional; The new update authority of the token

	MetadataUri          *string                        // optional; new metadata json uri
	SellerFeeBasisPoints *uint16                        // optional; new seller fee basis points
	Creators             *[]Creator                     // optional; new creators list
	PrimarySaleHappened  *bool                          // optional; new primary sale happened
	IsMutable            *bool                          // optional; new is mutable
	Collection           *common.PublicKey              // optional; new collection public key
	UseMethod            *token_metadata.TokenUseMethod // optional; new use method
	UseLimit             *uint64                        // optional; new use limit; default is 1; if use method is empty, use limit will be ignored
	UseRemaining         *uint64                        // optional; new use remaining; default equals use limit; if use method is empty, use remaining will be ignored
}

// Validate validates the params.
func (p UpdateMetadataParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("mint is required")
	}
	if p.UpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("update authority is required")
	}
	if p.NewUpdateAuthority != nil && *p.NewUpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("new update authority is invalid")
	}
	if p.MetadataUri != nil &&
		(*p.MetadataUri == "" ||
			(!strings.HasPrefix(*p.MetadataUri, "http://") &&
				!strings.HasPrefix(*p.MetadataUri, "https://"))) {
		return fmt.Errorf("metadata uri is invalid")
	}
	if p.SellerFeeBasisPoints != nil && *p.SellerFeeBasisPoints > 10000 {
		return fmt.Errorf("seller fee basis points must be less than or equal to 10000")
	}
	if p.Collection != nil && *p.Collection == (common.PublicKey{}) {
		return fmt.Errorf("collection public key is invalid")
	}
	if p.UseLimit != nil && p.UseMethod != nil && *p.UseLimit == 0 {
		return fmt.Errorf("use limit must be greater than 0")
	}
	if p.UseMethod != nil && !p.UseMethod.Valid() {
		return fmt.Errorf("use method is invalid")
	}
	return nil
}

// UpdateMetadata updates the metadata of the token.
func UpdateMetadata(params UpdateMetadataParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate update metadata: %w", err)
		}

		tokenMetadataPubkey, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
		}

		oldMetadata, err := c.GetTokenMetadata(ctx, params.Mint.ToBase58())
		if err != nil {
			return nil, fmt.Errorf("failed to get current token metadata: %w", err)
		}

		instructions := []types.Instruction{
			metaplex_token_metadata.UpdateMetadataAccountV2(metaplex_token_metadata.UpdateMetadataAccountV2Param{
				MetadataAccount: tokenMetadataPubkey,
				UpdateAuthority: params.UpdateAuthority,
				NewUpdateAuthority: func() *common.PublicKey {
					if params.NewUpdateAuthority != nil {
						return params.NewUpdateAuthority
					}
					return nil
				}(),
				PrimarySaleHappened: func() *bool {
					if params.PrimarySaleHappened != nil {
						return params.PrimarySaleHappened
					}
					return nil
				}(),
				IsMutable: func() *bool {
					if params.IsMutable != nil {
						return params.IsMutable
					}
					return nil
				}(),
				Data: func() *metaplex_token_metadata.DataV2 {
					if params.MetadataUri != nil ||
						params.SellerFeeBasisPoints != nil ||
						params.Creators != nil ||
						params.Collection != nil ||
						params.UseMethod != nil {

						var (
							name   = oldMetadata.Data.Name
							symbol = oldMetadata.Data.Symbol
						)
						if params.MetadataUri != nil {
							metadata, _ := metadata.MetadataFromURI(*params.MetadataUri)
							if metadata != nil {
								name = metadata.Name
								symbol = metadata.Symbol
							}
						}

						return &metaplex_token_metadata.DataV2{
							Name:   name,
							Symbol: symbol,
							Uri: func() string {
								if params.MetadataUri != nil {
									return *params.MetadataUri
								}
								return oldMetadata.MetadataUri
							}(),
							SellerFeeBasisPoints: func() uint16 {
								if params.SellerFeeBasisPoints != nil {
									return *params.SellerFeeBasisPoints
								}
								return oldMetadata.SellerFeeBasisPoints
							}(),
							Creators: func() *[]metaplex_token_metadata.Creator {
								if params.Creators != nil && len(*params.Creators) > 0 {
									creators := make([]metaplex_token_metadata.Creator, 0, len(*params.Creators))
									for _, creator := range *params.Creators {
										creators = append(creators, metaplex_token_metadata.Creator{
											Address: creator.Address,
											Share:   creator.Share,
											Verified: func() bool {
												return creator.Address.ToBase58() == params.UpdateAuthority.ToBase58()
											}(),
										})
									}
									return &creators
								} else if len(oldMetadata.Creators) > 0 {
									creators := make([]metaplex_token_metadata.Creator, 0, len(oldMetadata.Creators))
									for _, creator := range oldMetadata.Creators {
										creators = append(creators, metaplex_token_metadata.Creator{
											Address:  common.PublicKeyFromString(creator.Address),
											Share:    creator.Share,
											Verified: creator.Verified,
										})
									}
									return &creators
								}
								return nil
							}(),
							Collection: func() *metaplex_token_metadata.Collection {
								if params.Collection != nil {
									return &metaplex_token_metadata.Collection{
										Key: *params.Collection,
									}
								} else if oldMetadata.Collection != nil {
									return &metaplex_token_metadata.Collection{
										Key:      common.PublicKeyFromString(oldMetadata.Collection.Key),
										Verified: oldMetadata.Collection.Verified,
									}
								}
								return nil
							}(),
							Uses: func() *metaplex_token_metadata.Uses {
								if params.UseMethod != nil {
									if params.UseLimit == nil || *params.UseLimit == 0 {
										params.UseLimit = utils.Pointer[uint64](1)
									}
									if params.UseRemaining == nil || *params.UseRemaining == 0 {
										params.UseRemaining = params.UseLimit
									}
									return &metaplex_token_metadata.Uses{
										UseMethod: params.UseMethod.ToMetadataUseMethod(),
										Remaining: *params.UseRemaining,
										Total:     *params.UseLimit,
									}
								} else if oldMetadata.Uses != nil {
									useMethod := token_metadata.TokenUseMethod(oldMetadata.Uses.UseMethod)
									if useMethod.Valid() {
										return &metaplex_token_metadata.Uses{
											UseMethod: useMethod.ToMetadataUseMethod(),
											Remaining: oldMetadata.Uses.Remaining,
											Total:     oldMetadata.Uses.Total,
										}
									}
								}
								return nil
							}(),
						}
					}

					return nil
				}(),
			}),
		}

		return instructions, nil
	}
}
