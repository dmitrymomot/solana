package solana

import (
	"context"
	"errors"

	"github.com/portto/solana-go-sdk/common"
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
