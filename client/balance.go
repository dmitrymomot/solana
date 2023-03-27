package client

import (
	"context"

	"github.com/dmitrymomot/solana/common"
	"github.com/dmitrymomot/solana/types"
	"github.com/dmitrymomot/solana/utils"
)

// GetSOLBalance returns the SOL balance of the given base58 encoded account address.
// Returns the balance or an error.
func (c *Client) GetSOLBalance(ctx context.Context, base58Addr string) (uint64, error) {
	if err := common.ValidateSolanaWalletAddr(base58Addr); err != nil {
		return 0, utils.StackErrors(ErrGetSolBalance, err)
	}

	balance, err := c.rpcClient.GetBalance(ctx, base58Addr)
	if err != nil {
		return 0, utils.StackErrors(ErrGetSolBalance, err)
	}

	return balance, nil
}

// GetTokenBalance returns the SPL token balance of the given base58 encoded account address and SPL token mint address.
// base58Addr is the base58 encoded account address.
// base58MintAddr is the base58 encoded SPL token mint address.
// Returns the balance in lamports and token decimals, or an error.
func (c *Client) GetTokenBalance(ctx context.Context, base58Addr, base58MintAddr string) (types.TokenAmount, error) {
	if err := common.ValidateSolanaWalletAddr(base58Addr); err != nil {
		return types.TokenAmount{}, utils.StackErrors(ErrGetSplTokenBalance, err)
	}
	if err := common.ValidateSolanaWalletAddr(base58MintAddr); err != nil {
		return types.TokenAmount{}, utils.StackErrors(ErrGetSplTokenBalance, err)
	}

	ata, err := common.DeriveTokenAccount(base58Addr, base58MintAddr)
	if err != nil {
		return types.TokenAmount{}, utils.StackErrors(ErrGetSplTokenBalance, ErrFindAssociatedTokenAddress, err)
	}

	return c.GetAtaBalance(ctx, ata.String())
}

// GetAtaBalance returns the SPL token balance of the given base58 encoded associated token account address.
// base58Addr is the base58 encoded associated token account address.
// Returns the balance in lamports and token decimals, or an error.
func (c *Client) GetAtaBalance(ctx context.Context, base58Addr string) (types.TokenAmount, error) {
	balance, err := c.rpcClient.GetTokenAccountBalance(ctx, base58Addr)
	if err != nil {
		return types.TokenAmount{}, utils.StackErrors(ErrGetAtaBalance, ErrGetSplTokenBalance, err)
	}

	return types.NewTokenAmountFromLamports(balance.Amount, balance.Decimals), nil
}
