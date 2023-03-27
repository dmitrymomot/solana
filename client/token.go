package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dmitrymomot/solana/metadata"
	"github.com/dmitrymomot/solana/token_metadata"
	"github.com/dmitrymomot/solana/types"
	"github.com/dmitrymomot/solana/utils"
	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/token"
)

// GetTokenAccountInfo returns the token account information for a given token account address.
// This is a wrapper around the GetTokenAccount function from the solana-go-sdk.
// base58AtaAddr is the base58 encoded address of the associated token account.
// The function returns the token account information or an error.
func (c *Client) GetTokenAccountInfo(ctx context.Context, base58AtaAddr string) (token.TokenAccount, error) {
	ta, err := c.rpcClient.GetTokenAccount(ctx, base58AtaAddr)
	if err != nil {
		return token.TokenAccount{}, utils.StackErrors(ErrGetTokenAccount, err)
	}

	return ta, nil
}

// GetMintInfo returns the token mint information for a given mint address.
func (c *Client) GetMintInfo(ctx context.Context, base58MintAddr string) (token.MintAccount, error) {
	accInfo, err := c.rpcClient.GetAccountInfo(ctx, base58MintAddr)
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
	result, err := c.rpcClient.GetTokenSupply(ctx, base58MintAddr)
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

	metadataAccountInfo, err := c.rpcClient.GetAccountInfo(ctx, metadataAccount.ToBase58())
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

		editionAccountInfo, err := c.rpcClient.GetAccountInfo(ctx, editionPubkey.ToBase58())
		if err != nil {
			return nil, utils.StackErrors(ErrGetTokenMetadata, err)
		}

		edition, err := token_metadata.DeserializeEdition(editionAccountInfo.Data, c.rpcClient.GetAccountInfo)
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

	masterEdition, err := c.rpcClient.GetAccountInfo(ctx, masterEditionPubKey.String())
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

	editionData, err := c.rpcClient.GetAccountInfo(ctx, editionPubKey.String())
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetEditionInfo,
			err,
		)
	}

	edition, err := token_metadata.DeserializeEdition(editionData.Data, c.rpcClient.GetAccountInfo)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetEditionInfo,
			err,
		)
	}

	return edition, nil
}

// GetFungibleTokenMetadata returns the on-chain SPL token metadata by the given base58 encoded SPL token mint address.
// Returns the token metadata or an error.
func (c *Client) GetFungibleTokenMetadata(ctx context.Context, base58MintAddr string) (result *metadata.Metadata, err error) {
	// fallback to the deprecated metadata account if the given mint address has no on-chain metadata
	defer func() {
		if result == nil || err != nil || result.Name == "" || result.Symbol == "" || result.Image == "" {
			if depr, err := c.getDeprecatedTokenMetadata(ctx, base58MintAddr); err == nil {
				if result == nil {
					result = depr
				} else {
					if result.Name == "" {
						result.Name = depr.Name
					}
					if result.Symbol == "" {
						result.Symbol = depr.Symbol
					}
					if result.Image == "" {
						result.Image = depr.Image
					}
					if result.ExternalURL == "" {
						result.ExternalURL = depr.ExternalURL
					}
					if result.Description == "" {
						result.Description = depr.Description
					}
				}
			}
		}
	}()

	metadataAccount, err := metaplex_token_metadata.GetTokenMetaPubkey(common.PublicKeyFromString(base58MintAddr))
	if err != nil {
		return result, fmt.Errorf("failed to get token metadata account: %w", err)
	}

	accountInfo, err := c.rpcClient.GetAccountInfo(ctx, metadataAccount.ToBase58())
	if err != nil {
		return result, fmt.Errorf("failed to get account info: %w", err)
	}

	md, err := metaplex_token_metadata.MetadataDeserialize(accountInfo.Data)
	if err != nil {
		return result, fmt.Errorf("failed to deserialize metadata: %w", err)
	}

	result = &metadata.Metadata{
		Name:   md.Data.Name,
		Symbol: md.Data.Symbol,
	}

	if md.Data.Uri != "" && strings.HasPrefix(md.Data.Uri, "http") {
		mde, err := metadata.MetadataFromURI(md.Data.Uri)
		if err != nil {
			return result, fmt.Errorf("failed to get additional metadata from uri: %w", err)
		}

		result.Description = mde.Description
		result.Image = mde.Image
		result.ExternalURL = mde.ExternalURL
	}

	return result, nil
}

// @deprecated
// getDeprecatedTokenMetadata returns the deprecated SPL token metadata by the given base58 encoded SPL token mint address.
// This is a temporary solution to support the deprecated metadata format.
// Returns the token metadata or an error.
// Works only with mainnet.
func (c *Client) getDeprecatedTokenMetadata(_ context.Context, base58MintAddr string) (*metadata.Metadata, error) {
	if c.tokenListPath == "" || base58MintAddr == "" {
		return nil, fmt.Errorf("failed to get token metadata: token list path or mint address is empty")
	}

	resp, err := http.Get(c.tokenListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to download token list from uri: %w", err)
	}
	defer resp.Body.Close()

	var tokenList metadata.TokenList
	if err := json.NewDecoder(resp.Body).Decode(&tokenList); err != nil {
		return nil, fmt.Errorf("failed to decode token list from uri: %w", err)
	}

	// Find token metadata.
	var tokenMeta metadata.TokenListToken
	for _, token := range tokenList.Tokens {
		if token.Address == base58MintAddr && token.ChainID == metadata.ChainIdMainnet {
			tokenMeta = token
			break
		}
	}

	result := metadata.Metadata{
		Name:   tokenMeta.Name,
		Symbol: tokenMeta.Symbol,
		Image:  tokenMeta.LogoURI,
	}

	if tokenMeta.Extensions != nil {
		if tokenMeta.Extensions["description"] != nil {
			result.Description = tokenMeta.Extensions["description"].(string)
		}
		if tokenMeta.Extensions["website"] != nil {
			result.ExternalURL = tokenMeta.Extensions["website"].(string)
		} else if tokenMeta.Extensions["twitter"] != nil {
			result.ExternalURL = tokenMeta.Extensions["twitter"].(string)
		} else if tokenMeta.Extensions["discord"] != nil {
			result.ExternalURL = tokenMeta.Extensions["discord"].(string)
		}
	}

	return &result, nil
}
