package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestAirdrop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	// Request airdrop
	airdropSignature, err := client.RequestAirdrop(ctx, e2e.Wallet1Addr, solana.SOL)
	require.NoError(t, err)
	require.NotEmpty(t, airdropSignature)
	t.Logf("Airdrop signature: %s", airdropSignature)
}

func TestTransaction(t *testing.T) {
	senderAccount, err := types.AccountFromBase58(e2e.Wallet1PrivateKey)
	require.NoError(t, err)
	require.NotEmpty(t, senderAccount)

	recipientAccount, err := types.AccountFromBase58(e2e.Wallet2PrivateKey)
	require.NoError(t, err)
	require.NotEmpty(t, recipientAccount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	minAccountRent, err := client.GetMinimumBalanceForRentExemption(context.Background(), solana.AccountSize)
	require.NoError(t, err)
	t.Logf("Mint account rent: %d", minAccountRent)

	amount := minAccountRent + 100

	// Get sender balance
	startSenderBalance, err := client.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	assert.Greater(t, startSenderBalance, amount+minAccountRent)
	t.Logf("Start sender balance: %d", startSenderBalance)

	if startSenderBalance < amount+minAccountRent {
		t.Log("Requesting airdrop...")
		tx, err := client.RequestAirdrop(ctx, senderAccount.PublicKey.ToBase58(), solana.SOL)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		// Wait for transaction to be confirmed
		t.Log("Waiting for airdrop transaction to be confirmed")
		status, err := client.WaitForTransactionConfirmed(ctx, tx, 0)
		require.NoError(t, err)
		require.Equal(t, solana.TransactionStatusSuccess, status)

		startSenderBalance, err := client.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
		require.NoError(t, err)
		assert.Greater(t, startSenderBalance, amount+minAccountRent)
		t.Logf("Start sender balance: %d", startSenderBalance)
	}

	// Get recipient balance
	startRecipientBalance, err := client.GetSOLBalance(ctx, recipientAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	t.Logf("Start recipient balance: %d", startRecipientBalance)

	// Create a new transaction
	txb, err := client.TransferSOL(ctx, solana.TransferSOLParams{
		From:   senderAccount.PublicKey.ToBase58(),
		To:     recipientAccount.PublicKey.ToBase58(),
		Amount: amount,
		Memo:   "Test transaction " + time.Now().String(),
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
	t.Log("Waiting for transaction to be confirmed...")
	status, err := client.WaitForTransactionConfirmed(ctx, txSignature, 0)
	require.NoError(t, err)
	require.Equal(t, solana.TransactionStatusSuccess, status)

	// Get sender balance
	senderBalance, err := client.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Less(t, senderBalance, startSenderBalance)
	t.Logf("Sender balance: %d", senderBalance)

	// Get recipient balance
	recipientBalance, err := client.GetSOLBalance(ctx, recipientAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Greater(t, recipientBalance, startRecipientBalance)
	t.Logf("Recipient balance: %d", recipientBalance)
}
