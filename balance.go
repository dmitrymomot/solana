package solana

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/solplaydev/solana/utils"
)

// GetSOLBalance returns the SOL balance of the given base58 encoded account address.
// Returns the balance or an error.
func (c *Client) GetSOLBalance(ctx context.Context, base58Addr string) (uint64, error) {
	if err := ValidateSolanaWalletAddr(base58Addr); err != nil {
		return 0, utils.StackErrors(ErrGetSolBalance, err)
	}

	balance, err := c.solana.GetBalance(ctx, base58Addr)
	if err != nil {
		return 0, utils.StackErrors(ErrGetSolBalance, err)
	}

	return balance, nil
}

// GetTokenBalance returns the SPL token balance of the given base58 encoded account address and SPL token mint address.
// base58Addr is the base58 encoded account address.
// base58MintAddr is the base58 encoded SPL token mint address.
// Returns the balance in lamports and token decimals, or an error.
func (c *Client) GetTokenBalance(ctx context.Context, base58Addr, base58MintAddr string) (TokenAmount, error) {
	if err := ValidateSolanaWalletAddr(base58Addr); err != nil {
		return TokenAmount{}, utils.StackErrors(ErrGetSplTokenBalance, err)
	}
	if err := ValidateSolanaWalletAddr(base58MintAddr); err != nil {
		return TokenAmount{}, utils.StackErrors(ErrGetSplTokenBalance, err)
	}

	ata, _, err := common.FindAssociatedTokenAddress(
		common.PublicKeyFromString(base58Addr),
		common.PublicKeyFromString(base58MintAddr),
	)
	if err != nil {
		return TokenAmount{}, utils.StackErrors(ErrGetSplTokenBalance, ErrFindAssociatedTokenAddress, err)
	}

	return c.GetAtaBalance(ctx, ata.String())
}

// GetAtaBalance returns the SPL token balance of the given base58 encoded associated token account address.
// base58Addr is the base58 encoded associated token account address.
// Returns the balance in lamports and token decimals, or an error.
func (c *Client) GetAtaBalance(ctx context.Context, base58Addr string) (TokenAmount, error) {
	balance, err := c.solana.GetTokenAccountBalance(ctx, base58Addr)
	if err != nil {
		return TokenAmount{}, utils.StackErrors(ErrGetAtaBalance, ErrGetSplTokenBalance, err)
	}

	return NewTokenAmountFromLamports(balance.Amount, balance.Decimals), nil
}
