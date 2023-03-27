package client

import (
	"context"
	"encoding/json"
	"fmt"

	commonx "github.com/dmitrymomot/solana/common"
	"github.com/dmitrymomot/solana/types"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/rpc"
)

// GetFungibleTokensList gets the list of fungible tokens for the given wallet address.
func (c *Client) GetFungibleTokensList(ctx context.Context, walletAddr string) ([]types.TokenAccount, error) {
	return c.getTokensList(ctx, walletAddr, true)
}

// GetNonFungibleTokensList gets the list of non-fungible tokens for the given wallet address.
// Result includes the list of NFTs and the list of semi-fungible tokens (assets).
func (c *Client) GetNonFungibleTokensList(ctx context.Context, walletAddr string) ([]types.TokenAccount, error) {
	return c.getTokensList(ctx, walletAddr, false)
}

// GetFungibleTokensList gets the list of fungible tokens for the given wallet address.
func (c *Client) getTokensList(ctx context.Context, walletAddr string, fungible bool) ([]types.TokenAccount, error) {
	if err := commonx.ValidateSolanaWalletAddr(walletAddr); err != nil {
		return nil, err
	}

	getTokenAccountsByOwnerResponse, err := c.rpcClient.RpcClient.GetTokenAccountsByOwnerWithConfig(
		ctx,
		walletAddr,
		rpc.GetTokenAccountsByOwnerConfigFilter{
			ProgramId: common.TokenProgramID.ToBase58(),
		},
		rpc.GetTokenAccountsByOwnerConfig{
			Encoding: rpc.AccountEncodingJsonParsed,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not get fungible tokens list: %w", err)
	}

	if getTokenAccountsByOwnerResponse.Error != nil {
		return nil, fmt.Errorf("could not get fungible tokens list: %s", getTokenAccountsByOwnerResponse.Error.Message)
	}

	var tokenAccounts []types.TokenAccount
	for _, v := range getTokenAccountsByOwnerResponse.Result.Value {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("could not marshal account data: %w", err)
		}

		acc, err := types.NewTokenAccount(b)
		if err != nil {
			return nil, fmt.Errorf("NewTokenAccount: %w", err)
		}

		if !acc.IsEmpty() && acc.IsFungibleToken() == fungible {
			tokenAccounts = append(tokenAccounts, acc)
		}
	}

	return tokenAccounts, nil
}
