package client

import (
	"context"

	"github.com/solplaydev/solana/common"
	"github.com/solplaydev/solana/utils"
)

// RequestAirdrop sends a request to the solana network to airdrop SOL to the given account.
// Returns the transaction hash or an error.
func (c *Client) RequestAirdrop(ctx context.Context, base58Addr string, amount uint64) (string, error) {
	if amount < 1 || amount > 2*1e9 {
		return "", ErrInvalidAirdropAmount
	}

	if err := common.ValidateSolanaWalletAddr(base58Addr); err != nil {
		return "", utils.StackErrors(ErrRequestAirdrop, err)
	}

	tx, err := c.solana.RequestAirdrop(ctx, base58Addr, amount)
	if err != nil {
		return "", utils.StackErrors(ErrRequestAirdrop, err)
	}

	return tx, nil
}
