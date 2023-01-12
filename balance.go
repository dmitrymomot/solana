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
func (c *Client) GetTokenBalance(ctx context.Context, base58Addr, base58MintAddr string) (uint64, uint8, error) {
	if err := ValidateSolanaWalletAddr(base58Addr); err != nil {
		return 0, 0, utils.StackErrors(ErrGetSplTokenBalance, err)
	}
	if err := ValidateSolanaWalletAddr(base58MintAddr); err != nil {
		return 0, 0, utils.StackErrors(ErrGetSplTokenBalance, err)
	}

	ata, _, err := common.FindAssociatedTokenAddress(
		common.PublicKeyFromString(base58Addr),
		common.PublicKeyFromString(base58MintAddr),
	)
	if err != nil {
		return 0, 0, utils.StackErrors(ErrGetSplTokenBalance, ErrFindAssociatedTokenAddress, err)
	}

	balance, decimals, err := c.solana.GetTokenAccountBalance(ctx, ata.ToBase58())
	if err != nil {
		return 0, 0, utils.StackErrors(ErrGetSplTokenBalance, err)
	}

	return balance, decimals, nil
}

// GetAtaBalance returns the SPL token balance of the given base58 encoded associated token account address.
// base58Addr is the base58 encoded associated token account address.
// Returns the balance in lamports and token decimals, or an error.
func (c *Client) GetAtaBalance(ctx context.Context, base58Addr string) (uint64, uint8, error) {
	if err := ValidateSolanaWalletAddr(base58Addr); err != nil {
		return 0, 0, utils.StackErrors(ErrGetAtaBalance, err)
	}

	balance, decimals, err := c.solana.GetTokenAccountBalance(ctx, base58Addr)
	if err != nil {
		return 0, 0, utils.StackErrors(ErrGetAtaBalance, ErrGetSplTokenBalance, err)
	}

	return balance, decimals, nil
}
