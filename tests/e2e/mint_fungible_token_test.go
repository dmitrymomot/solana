package e2e_test

import (
	"context"
	"testing"

	"github.com/solplaydev/solana"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/stretchr/testify/require"
)

func TestMintFungibleToken_MintFixedSupply(t *testing.T) {
	var (
		tokenName    = "Test Token"
		tokenSymbol  = "TSTt"
		metadataUri  = "https://www.arweave.net/QR1PsBgIbiYoKgGff5Jq2U8QavHChRjBki8XRJ-06mI?ext=json"
		supplyAmount = 1000000 * solana.SPLTokenDefaultMultiplier
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	// Mint a fungible token
	mintAddr, tx, err := client.InitMintFungibleToken(ctx, solana.InitMintFungibleTokenParams{
		FeePayer:     e2e.FeePayerAddr,
		Owner:        e2e.Wallet1Addr,
		SupplyAmount: supplyAmount,
		Decimals:     solana.SPLTokenDefaultDecimals,
		FixedSupply:  true,

		Name:        tokenName,
		Symbol:      tokenSymbol,
		MetadataURI: metadataUri,
	})
	require.NoError(t, err)
	require.NotEmpty(t, tx)
	t.Logf("Mint address: %s", mintAddr)

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
	t.Logf("Transaction hash: %s", txHash)
	require.NotEmpty(t, txHash)

	// Wait for the transaction to be confirmed
	txInfo, err := client.WaitForTransactionConfirmed(ctx, txHash, 0)
	require.NoError(t, err)
	t.Logf("Transaction status: %+v", txInfo)
	require.EqualValues(t, txInfo, solana.TransactionStatusSuccess)

	// Check token balance
	balance, err := client.GetTokenBalance(ctx, e2e.Wallet1Addr, mintAddr)
	require.NoError(t, err)
	t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
	require.EqualValues(t, supplyAmount, balance.Amount)
	require.EqualValues(t, solana.SPLTokenDefaultDecimals, balance.Decimals)

	// Check token metadata
	metadata, err := client.GetTokenMetadata(ctx, mintAddr)
	require.NoError(t, err)
	t.Logf("Token metadata: %+v", metadata)
	require.EqualValues(t, tokenName, metadata.Data.Name)
	require.EqualValues(t, tokenSymbol, metadata.Data.Symbol)
	require.EqualValues(t, token_metadata.TokenStandardFungible, metadata.TokenStandard)
}
