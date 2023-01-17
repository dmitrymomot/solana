package e2e_test

import (
	"context"
	"testing"

	"github.com/solplaydev/solana"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/stretchr/testify/require"
)

func TestMintNftEdition(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	mintAddr, tx, err := client.MintNonFungibleTokenEdition(ctx, solana.MintNonFungibleTokenEditionParams{
		FeePayer: e2e.FeePayerAddr,
		Owner:    e2e.Wallet1Addr,
		Mint:     e2e.MasterEditionMintAddr,
		// Edition:  2,
	})
	require.NoError(t, err)
	require.NotEmpty(t, tx)
	t.Logf("NFT edition mint address: %s", mintAddr)

	// Sign the transaction by the fee payer
	feePayer, err := solana.AccountFromBase58(e2e.FeePayerPrivateKey)
	require.NoError(t, err)
	tx, err = client.SignTransaction(ctx, feePayer, tx)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	// Sign the transaction by the token owner
	owner, err := solana.AccountFromBase58(e2e.Wallet1PrivateKey)
	require.NoError(t, err)
	tx, err = client.SignTransaction(ctx, owner, tx)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	// Send the transaction
	txHash, err := client.SendTransaction(ctx, tx)
	require.NoError(t, err)
	require.NotEmpty(t, txHash)
	t.Logf("Transaction hash: %s", txHash)

	// Wait for the transaction to be confirmed
	txInfo, err := client.WaitForTransactionConfirmed(ctx, txHash, 0)
	require.NoError(t, err)
	t.Logf("Transaction status: %+v", txInfo)
	require.EqualValues(t, txInfo, solana.TransactionStatusSuccess)

	// Check token balance
	balance, deciamls, err := client.GetTokenBalance(ctx, e2e.Wallet1Addr, mintAddr)
	require.NoError(t, err)
	t.Logf("Token balance: %d, decimals: %d", balance, deciamls)
	require.EqualValues(t, 1, balance)
	require.EqualValues(t, uint8(0), deciamls)
}
