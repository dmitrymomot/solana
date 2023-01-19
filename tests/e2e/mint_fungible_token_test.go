package e2e_test

import (
	"context"
	"testing"

	"github.com/portto/solana-go-sdk/common"
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

	// Build token metadata
	mb := token_metadata.NewTokenMetadataInstructionBuilder()
	mb.SetName(tokenName)
	mb.SetSymbol(tokenSymbol)
	mb.SetUri(metadataUri)
	metaPubkey, metadataInstruction, err := mb.Build()
	require.NoError(t, err)
	require.NotNil(t, metadataInstruction)
	require.True(t, metaPubkey != (common.PublicKey{}))

	// Build mint transaction
	b := solana.NewMintBuilder(client)
	b.SetTokenStandard(token_metadata.TokenStandardFungible)
	b.SetFeePayerBase58(e2e.FeePayerAddr)
	b.SetSupplyAmount(supplyAmount)
	b.SetDecimals(solana.SPLTokenDefaultDecimals)
	b.SetFixedSupply(true)

	mintAddr, tx, err := b.Build()
	require.NoError(t, err)
	require.NotNil(t, tx)
	require.NotEmpty(t, mintAddr)
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
	balance, deciamls, err := client.GetTokenBalance(ctx, e2e.Wallet1Addr, mintAddr)
	require.NoError(t, err)
	t.Logf("Token balance: %d, decimals: %d", balance, deciamls)
	require.EqualValues(t, supplyAmount, balance)
	require.EqualValues(t, solana.SPLTokenDefaultDecimals, deciamls)

	// Check token metadata
	metadata, err := client.GetTokenMetadata(ctx, mintAddr)
	require.NoError(t, err)
	t.Logf("Token metadata: %+v", metadata)
	require.EqualValues(t, tokenName, metadata.Data.Name)
	require.EqualValues(t, tokenSymbol, metadata.Data.Symbol)
	require.EqualValues(t, token_metadata.TokenStandardFungible, metadata.TokenStandard)
}
