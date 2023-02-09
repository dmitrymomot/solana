package client

import (
	"context"
	"errors"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/types"
	"github.com/solplaydev/solana/utils"
)

// GetTokenAccountInfo returns the token account information for a given token account address.
// This is a wrapper around the GetTokenAccount function from the solana-go-sdk.
// base58AtaAddr is the base58 encoded address of the associated token account.
// The function returns the token account information or an error.
func (c *Client) GetTokenAccountInfo(ctx context.Context, base58AtaAddr string) (token.TokenAccount, error) {
	ta, err := c.solana.GetTokenAccount(ctx, base58AtaAddr)
	if err != nil {
		return token.TokenAccount{}, utils.StackErrors(ErrGetTokenAccount, err)
	}

	return ta, nil
}

// GetMintInfo returns the token mint information for a given mint address.
func (c *Client) GetMintInfo(ctx context.Context, base58MintAddr string) (token.MintAccount, error) {
	accInfo, err := c.solana.GetAccountInfo(ctx, base58MintAddr)
	if err != nil {
		return token.MintAccount{}, utils.StackErrors(ErrGetMintInfo, err)
	}

	mintInfo, err := token.MintAccountFromData(accInfo.Data)
	if err != nil {
		return token.MintAccount{}, utils.StackErrors(ErrGetMintInfo, err)
	}

	return mintInfo, nil
}

// GetTokenSupply returns the token supply for a given mint address.
// This is a wrapper around the GetTokenSupply function from the solana-go-sdk.
// base58MintAddr is the base58 encoded address of the token mint.
// The function returns the token supply and decimals or an error.
func (c *Client) GetTokenSupply(ctx context.Context, base58MintAddr string) (types.TokenAmount, error) {
	result, err := c.solana.GetTokenSupply(ctx, base58MintAddr)
	if err != nil {
		return types.TokenAmount{}, utils.StackErrors(ErrGetTokenSupply, err)
	}

	return types.NewTokenAmountFromLamports(result.Amount, result.Decimals), nil
}

// GetMasterEditionSupply returns the current supply and max supply of a master edition
func (c *Client) GetMasterEditionSupply(ctx context.Context, masterMint common.PublicKey) (current, max uint64, err error) {
	editionInfo, err := c.GetMasterEditionInfo(ctx, masterMint.ToBase58())
	if err != nil || editionInfo == nil {
		return 0, 0, utils.StackErrors(
			ErrGetMasterEditionCurrentSupply,
			err,
		)
	}
	if editionInfo.Type == "" || editionInfo.Type != token_metadata.KeyMasterEdition.String() {
		return 0, 0, utils.StackErrors(
			ErrGetMasterEditionCurrentSupply,
			ErrTokenIsNotMasterEdition,
		)
	}
	if editionInfo.MaxSupply == 0 || editionInfo.Supply == editionInfo.MaxSupply {
		return 0, 0, utils.StackErrors(
			ErrGetMasterEditionCurrentSupply,
			ErrMaxSupplyReached,
		)
	}

	return editionInfo.Supply, editionInfo.MaxSupply, nil
}

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
		return nil, utils.StackErrors(ErrGetTokenMetadata, err)
	}

	metadataAccountInfo, err := c.solana.GetAccountInfo(ctx, metadataAccount.ToBase58())
	if err != nil {
		return nil, utils.StackErrors(ErrGetTokenMetadata, err)
	}

	// log.Println(utils.PrettyPrint(metadataAccountInfo))

	metadata, err := token_metadata.DeserializeMetadata(metadataAccountInfo.Data)
	if err != nil {
		return nil, utils.StackErrors(ErrGetTokenMetadata, err)
	}

	if metadata.TokenStandard == token_metadata.TokenStandardNonFungible.String() ||
		metadata.TokenStandard == token_metadata.TokenStandardNonFungibleEdition.String() {

		editionPubkey, err := token_metadata.DeriveEditionPubkey(mintPubkey)
		if err != nil {
			return nil, utils.StackErrors(ErrGetTokenMetadata, err)
		}

		editionAccountInfo, err := c.solana.GetAccountInfo(ctx, editionPubkey.ToBase58())
		if err != nil {
			return nil, utils.StackErrors(ErrGetTokenMetadata, err)
		}

		edition, err := token_metadata.DeserializeEdition(editionAccountInfo.Data, c.solana.GetAccountInfo)
		if err != nil {
			return nil, utils.StackErrors(ErrGetTokenMetadata, err)
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
