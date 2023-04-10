package e2e_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/dmitrymomot/solana/client"
	"github.com/dmitrymomot/solana/instructions"
	"github.com/dmitrymomot/solana/tests/e2e"
	"github.com/dmitrymomot/solana/transaction"
	typesx "github.com/dmitrymomot/solana/types"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestAirdrop(t *testing.T) {
	t.SkipNow() // uncomment to run this test

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := client.New(client.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	// Request airdrop
	airdropSignature, err := client.RequestAirdrop(ctx, e2e.Wallet1Pubkey.ToBase58(), typesx.SOL)
	require.NoError(t, err)
	require.NotEmpty(t, airdropSignature)
	fmt.Printf("Airdrop signature: %s\n", airdropSignature)
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
	sc := client.New(client.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	minAccountRent, err := sc.GetMinimumBalanceForRentExemption(context.Background(), typesx.AccountSize)
	require.NoError(t, err)
	fmt.Printf("Mint account rent: %d\n", minAccountRent)

	amount := minAccountRent + 100

	// Get sender balance
	startSenderBalance, err := sc.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	assert.Greater(t, startSenderBalance, amount+minAccountRent)

	if startSenderBalance < amount+minAccountRent {
		tx, err := sc.RequestAirdrop(ctx, senderAccount.PublicKey.ToBase58(), typesx.SOL)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		// Wait for transaction to be confirmed
		fmt.Println("Waiting for airdrop transaction to be confirmed")
		status, err := sc.WaitForTransactionConfirmed(ctx, tx, 0)
		require.NoError(t, err)
		require.Equal(t, typesx.TransactionStatusSuccess, status)

		startSenderBalance, err = sc.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
		require.NoError(t, err)
		assert.Greater(t, startSenderBalance, amount+minAccountRent)
		fmt.Printf("Start sender balance: %d\n", startSenderBalance)
	}

	// Get recipient balance
	startRecipientBalance, err := sc.GetSOLBalance(ctx, recipientAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	fmt.Printf("Start recipient balance: %d\n", startRecipientBalance)

	// Create a new transaction
	txb, err := transaction.NewTransactionBuilder(sc).
		SetFeePayer(senderAccount.PublicKey).
		AddInstruction(instructions.TransferSOL(instructions.TransferSOLParams{
			Sender:    senderAccount.PublicKey,
			Recipient: recipientAccount.PublicKey,
			Amount:    amount,
		})).
		AddInstruction(instructions.Memo(
			fmt.Sprintf("Send %d lamports to %s", amount, recipientAccount.PublicKey.ToBase58()),
			senderAccount.PublicKey,
		)).
		Build(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, txb)

	// Sign transaction
	txs, err := sc.SignTransaction(ctx, senderAccount, txb)
	require.NoError(t, err)
	require.NotEmpty(t, txs)

	// Send transaction
	txSignature, err := sc.SendTransaction(ctx, txs)
	require.NoError(t, err)
	require.NotEmpty(t, txSignature)
	fmt.Printf("Transaction signature: %s\n", txSignature)

	// Wait for transaction to be confirmed
	fmt.Println("Waiting for transaction to be confirmed...")
	status, err := sc.WaitForTransactionConfirmed(ctx, txSignature, 0)
	require.NoError(t, err)
	require.Equal(t, typesx.TransactionStatusSuccess, status)

	// Get sender balance
	senderBalance, err := sc.GetSOLBalance(ctx, senderAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Less(t, senderBalance, startSenderBalance)
	fmt.Printf("Sender balance: %d\n", senderBalance)

	// Get recipient balance
	recipientBalance, err := sc.GetSOLBalance(ctx, recipientAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.Greater(t, recipientBalance, startRecipientBalance)
	fmt.Printf("Recipient balance: %d\n", recipientBalance)
}
