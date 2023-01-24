package solana

import (
	"context"
	"errors"

	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// GetTokenMetadata returns the metadata of a token
func (c *Client) GetTokenMetadata(ctx context.Context, base58MintAddr string) (*token_metadata.Metadata, error) {
	if base58MintAddr == "" {
		return nil, utils.StackErrors(
			ErrInvalidPublicKey,
			errors.New("mint address is required"),
		)
	}

	mintPubkey := common.PublicKeyFromString(base58MintAddr)

	metadataAccount, err := token_metadata.DeriveTokenMetadataPubkey(mintPubkey)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	metadataAccountInfo, err := c.solana.GetAccountInfo(ctx, metadataAccount.ToBase58())
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	metadata, err := token_metadata.DeserializeMetadata(metadataAccountInfo.Data)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	if metadata.TokenStandard == token_metadata.TokenStandardNonFungible.String() ||
		metadata.TokenStandard == token_metadata.TokenStandardNonFungibleEdition.String() {

		editionPubkey, err := token_metadata.DeriveEditionPubkey(mintPubkey)
		if err != nil {
			return nil, utils.StackErrors(
				ErrGetTokenMetadata,
				err,
			)
		}

		editionAccountInfo, err := c.solana.GetAccountInfo(ctx, editionPubkey.ToBase58())
		if err != nil {
			return nil, utils.StackErrors(
				ErrGetTokenMetadata,
				err,
			)
		}

		edition, err := token_metadata.DeserializeEdition(editionAccountInfo.Data, c.solana.GetAccountInfo)
		if err != nil {
			return nil, utils.StackErrors(
				ErrGetTokenMetadata,
				err,
			)
		}

		metadata.Edition = edition
	}

	return metadata, nil
}

// GetMasterEditionInfo returns the master edition info of a token
func (c *Client) GetMasterEditionInfo(ctx context.Context, base58MintAddr string) (*token_metadata.Edition, error) {
	mint := common.PublicKeyFromString(base58MintAddr)
	masterEditionPubKey, err := token_metadata.DeriveEditionPubkey(mint)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetMasterEditionInfo,
			err,
		)
	}

	masterEdition, err := c.solana.GetAccountInfo(ctx, masterEditionPubKey.String())
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetMasterEditionInfo,
			err,
		)
	}

	edition, err := token_metadata.DeserializeMasterEdition(masterEdition.Data)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetMasterEditionInfo,
			err,
		)
	}

	return edition, nil
}

// GetEditionInfo returns the edition info of a token
func (c *Client) GetEditionInfo(ctx context.Context, base58MintAddr string) (*token_metadata.Edition, error) {
	mint := common.PublicKeyFromString(base58MintAddr)
	editionPubKey, err := token_metadata.DeriveEditionPubkey(mint)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetEditionInfo,
			err,
		)
	}

	editionData, err := c.solana.GetAccountInfo(ctx, editionPubKey.String())
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetEditionInfo,
			err,
		)
	}

	edition, err := token_metadata.DeserializeEdition(editionData.Data, c.solana.GetAccountInfo)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetEditionInfo,
			err,
		)
	}

	return edition, nil
}

// VerifyCreatorParams is the params for VerifyCreator
type VerifyCreatorParams struct {
	MintAddress    string
	CreatorAddress string
	FeePayer       string
}

// VerifyCreator verifies that the given creator is a valid token creator.
// Returns prepared transaction encoded in base64 string and error if any.
func (c *Client) VerifyCreator(ctx context.Context, params VerifyCreatorParams) (string, error) {
	if params.CreatorAddress == "" {
		return "", utils.StackErrors(
			ErrVerifyCreator,
			errors.New("creator address is required"),
		)
	}
	if params.MintAddress == "" {
		return "", utils.StackErrors(
			ErrVerifyCreator,
			errors.New("mint address is required"),
		)
	}
	if params.FeePayer == "" {
		params.FeePayer = params.CreatorAddress
	}

	// mint address
	nft := common.PublicKeyFromString(params.MintAddress)

	tokenMetadataPubkey, err := token_metadata.DeriveTokenMetadataPubkey(nft)
	if err != nil {
		return "", utils.StackErrors(
			ErrVerifyCreator,
			err,
		)
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		FeePayer: params.FeePayer,
		Instructions: []types.Instruction{
			metaplex_token_metadata.SignMetadata(metaplex_token_metadata.SignMetadataParam{
				Metadata: tokenMetadataPubkey,
				Creator:  common.PublicKeyFromString(params.CreatorAddress),
			}),
		},
	})
	if err != nil {
		return "", utils.StackErrors(
			ErrVerifyCreator,
			err,
		)
	}

	return txb, nil
}

// UpdateMetadataParams is the params for UpdateMetadata
type UpdateMetadataParams struct {
	Mint                 string                   // required; token mint address, encoded in base58
	Owner                string                   // required; token owner address encoded in base58
	FeePayer             string                   // optional; defaults to owner
	NewUpdateAuthority   string                   // optional; new update authority address encoded in base58
	Name                 string                   // optional; new name
	Symbol               string                   // optional; new symbol
	Uri                  string                   // optional; new uri
	SellerFeeBasisPoints uint16                   // optional; new seller fee basis points
	Creators             []token_metadata.Creator // optional; new creators addresses encoded in base58
	PrimarySaleHappened  bool                     // optional; new primary sale happened
	IsMutable            bool                     // optional; new is mutable
	Collection           string                   // optional; new collection address encoded in base58
	Uses                 *token_metadata.Uses     // optional; new uses
}

// UpdateMetadata updates the metadata of a token.
// Returns prepared transaction encoded in base64 string and error if any.
func (c *Client) UpdateMetadata(ctx context.Context, params UpdateMetadataParams) (string, error) {
	if params.Owner == "" {
		return "", utils.StackErrors(
			ErrUpdateMetadata,
			errors.New("token owner address is required"),
		)
	}
	if params.Mint == "" {
		return "", utils.StackErrors(
			ErrUpdateMetadata,
			errors.New("token mint address is required"),
		)
	}
	if params.FeePayer == "" {
		params.FeePayer = params.Owner
	}

	// mint address
	nft := common.PublicKeyFromString(params.Mint)

	tokenMetadataPubkey, err := token_metadata.DeriveTokenMetadataPubkey(nft)
	if err != nil {
		return "", utils.StackErrors(
			ErrUpdateMetadata,
			err,
		)
	}

	metadataAccountInfo, err := c.solana.GetAccountInfo(ctx, tokenMetadataPubkey.ToBase58())
	if err != nil {
		return "", utils.StackErrors(
			ErrUpdateMetadata,
			ErrGetTokenMetadata,
			err,
		)
	}

	metadata, err := metaplex_token_metadata.MetadataDeserialize(metadataAccountInfo.Data)
	if err != nil {
		return "", utils.StackErrors(
			ErrUpdateMetadata,
			ErrGetTokenMetadata,
			err,
		)
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		FeePayer: params.FeePayer,
		Instructions: []types.Instruction{
			metaplex_token_metadata.UpdateMetadataAccountV2(metaplex_token_metadata.UpdateMetadataAccountV2Param{
				MetadataAccount: tokenMetadataPubkey,
				UpdateAuthority: common.PublicKeyFromString(params.Owner),
				NewUpdateAuthority: func() *common.PublicKey {
					if params.NewUpdateAuthority != "" {
						return utils.Pointer(common.PublicKeyFromString(params.NewUpdateAuthority))
					}
					return nil
				}(),
				PrimarySaleHappened: func() *bool {
					if params.PrimarySaleHappened {
						return utils.Pointer(params.PrimarySaleHappened)
					}
					return nil
				}(),
				IsMutable: func() *bool {
					if params.IsMutable {
						return utils.Pointer(params.IsMutable)
					}
					return nil
				}(),
				Data: func() *metaplex_token_metadata.DataV2 {
					if params.Name != "" ||
						params.Symbol != "" ||
						params.Uri != "" ||
						params.SellerFeeBasisPoints != 0 ||
						len(params.Creators) != 0 ||
						params.Collection != "" ||
						params.Uses != nil {
						return &metaplex_token_metadata.DataV2{
							Name: func() string {
								if params.Name != "" {
									return params.Name
								}
								return metadata.Data.Name
							}(),
							Symbol: func() string {
								if params.Symbol != "" {
									return params.Symbol
								}
								return metadata.Data.Symbol
							}(),
							Uri: func() string {
								if params.Uri != "" {
									return params.Uri
								}
								return metadata.Data.Uri
							}(),
							SellerFeeBasisPoints: func() uint16 {
								if params.SellerFeeBasisPoints != 0 {
									return params.SellerFeeBasisPoints
								}
								return metadata.Data.SellerFeeBasisPoints
							}(),
							Creators: func() *[]metaplex_token_metadata.Creator {
								if len(params.Creators) != 0 {
									creators := make([]metaplex_token_metadata.Creator, 0, len(params.Creators))
									for _, creator := range params.Creators {
										creators = append(creators, metaplex_token_metadata.Creator{
											Address: common.PublicKeyFromString(creator.Address),
											Share:   creator.Share,
										})
									}
									return &creators
								}
								return metadata.Data.Creators
							}(),
							Collection: func() *metaplex_token_metadata.Collection {
								if params.Collection != "" {
									return &metaplex_token_metadata.Collection{
										Key: common.PublicKeyFromString(params.Collection),
									}
								}
								return metadata.Collection
							}(),
							Uses: func() *metaplex_token_metadata.Uses {
								if params.Uses != nil {
									return &metaplex_token_metadata.Uses{
										UseMethod: token_metadata.TokenUseMethod(params.Uses.UseMethod).ToMetadataUseMethod(),
										Remaining: params.Uses.Remaining,
										Total:     params.Uses.Total,
									}
								}
								return metadata.Uses
							}(),
						}
					}
					return nil
				}(),
			}),
		},
	})
	if err != nil {
		return "", utils.StackErrors(
			ErrUpdateMetadata,
			err,
		)
	}

	return txb, nil
}
