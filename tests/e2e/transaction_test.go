package e2e_test

import (
	"context"
	"testing"

	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana"
	"github.com/stretchr/testify/require"
)

var (
	senderAddr       = "FuQhSmAT6kAmmzCMiiYbzFcTQJFuu6raXAdCFibz4YPR"
	senderPrivateKey = "4xgyc4d8SkRMK4BrdDnhk1Cb3fJBfevZP4Fueiga7wt3aDaDDtYSLJV8V4pY9rci9Qqyo9zwV6dBmV2G7nVYk9sV"

	// recipientAddr       = "RjpQLUttBMdoQ4HKMygScEjkd6S69dZZC9T4W3Z3DKD"
	recipientPrivateKey = "5f34yVBKf7VfcpgW3pD91UcYuMiU7MnAgtMNUooXECcMc2kEGwM2p4LsHFwqK61X2o9TjA5iUpSRYkyYUojmbCrj"
)

func TestRequestAirdrop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(solana.SolanaDevnetRPCURL))

	// Request airdrop
	airdropSignature, err := client.RequestAirdrop(ctx, senderAddr, solana.SOL)
	require.NoError(t, err)
	require.NotEmpty(t, airdropSignature)
	t.Logf("Airdrop signature: %s", airdropSignature)
}

func TestTransaction(t *testing.T) {
	senderAccount, err := types.AccountFromBase58(senderPrivateKey)
	require.NoError(t, err)
	require.NotEmpty(t, senderAccount)

	recipientAccount, err := types.AccountFromBase58(recipientPrivateKey)
	require.NoError(t, err)
	require.NotEmpty(t, recipientAccount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(solana.SolanaDevnetRPCURL))

	// Get sender balance
	senderBalance, err := client.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Greater(t, senderBalance, uint64(0))
	t.Logf("Sender balance: %d", senderBalance)

	// Get recipient balance
	recipientBalance, err := client.GetSOLBalance(ctx, recipientAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Equal(t, uint64(0), recipientBalance)

	amount := uint64(100)

	// Create a new transaction
	txb, err := client.TransferSOL(ctx, solana.TransferSOLParams{
		Base58SourceAddr: senderAccount.PublicKey.ToBase58(),
		Base58DestAddr:   recipientAccount.PublicKey.ToBase58(),
		Lamports:         amount,
		Memo:             "Test transaction",
	})
	require.NoError(t, err)
	require.NotEmpty(t, txb)

	// Sign transaction
	txs, err := client.SignTransaction(ctx, senderAccount, txb)
	require.NoError(t, err)
	require.NotEmpty(t, txs)

	// Send transaction
	txSignature, err := client.SendTransaction(ctx, txs)
	require.NoError(t, err)
	require.NotEmpty(t, txSignature)
	t.Logf("Transaction signature: %s", txSignature)

	// Wait for transaction to be confirmed
	for {
		txStatus, err := client.GetTransactionStatus(ctx, txSignature)
		require.NoError(t, err)
		require.NotEmpty(t, txStatus)

		if txStatus == solana.TransactionStatusFailure || txStatus == solana.TransactionStatusSuccess {
			break
		}
	}

	// Get sender balance
	senderBalance, err = client.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Greater(t, senderBalance, uint64(0))
	t.Logf("Sender balance: %d", senderBalance)

	// Get recipient balance
	recipientBalance, err = client.GetSOLBalance(ctx, recipientAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Equal(t, amount, recipientBalance)
}
