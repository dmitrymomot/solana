package solana

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/solplaydev/solana/utils"
)

// DeriveTokenAccount derives an associated token account from a Solana account and a mint address.
// This is a wrapper around the FindAssociatedTokenAddress function from the solana-go-sdk.
// base58WalletAddr is the base58 encoded address of the Solana account.
// base58MintAddr is the base58 encoded address of the token mint.
// The function returns the base58 encoded address of the token account or an error.
func DeriveTokenAccount(base58WalletAddr, base58MintAddr string) (common.PublicKey, error) {
	ata, _, err := common.FindAssociatedTokenAddress(
		common.PublicKeyFromString(base58WalletAddr),
		common.PublicKeyFromString(base58MintAddr),
	)
	if err != nil {
		return common.PublicKey{}, utils.StackErrors(ErrDeriveTokenAccount, err)
	}

	return ata, nil
}

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
func (c *Client) GetTokenSupply(ctx context.Context, base58MintAddr string) (TokenAmount, error) {
	result, err := c.solana.GetTokenSupply(ctx, base58MintAddr)
	if err != nil {
		return TokenAmount{}, utils.StackErrors(ErrGetTokenSupply, err)
	}

	return NewTokenAmountFromLamports(result.Amount, result.Decimals), nil
}
