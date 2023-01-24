package e2e_test

import (
	"context"
	"testing"

	"github.com/solplaydev/solana"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokenMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	t.Run("master edition metadata", func(t *testing.T) {
		md, err := client.GetTokenMetadata(ctx, e2e.MasterEditionMintAddr)
		require.NoError(t, err)
		require.NotEmpty(t, md)
		assert.NotEmpty(t, md.MetadataUri)
		assert.NotEmpty(t, md.SellerFeeBasisPoints)
		assert.GreaterOrEqual(t, len(md.Creators), 1)
		if assert.NotNil(t, md.Data) {
			assert.NotEmpty(t, md.Data.Name)
			assert.NotEmpty(t, md.Data.Symbol)
			assert.NotEmpty(t, md.Data.Description)
			assert.NotEmpty(t, md.Data.Image)
		}
		assert.NotNil(t, md.Edition)
		assert.EqualValues(t, md.Edition.Type, token_metadata.KeyMasterEdition.String())
		assert.EqualValues(t, md.Edition.MaxSupply, uint64(1000))
		assert.GreaterOrEqual(t, md.Edition.Supply, uint64(4))
		assert.EqualValues(t, md.Edition.Edition, uint64(0))
	})

	t.Run("printed edition metadata #2", func(t *testing.T) {
		md, err := client.GetTokenMetadata(ctx, e2e.EditionMintAddr2)
		require.NoError(t, err)
		require.NotEmpty(t, md)
		assert.NotEmpty(t, md.MetadataUri)
		assert.NotEmpty(t, md.SellerFeeBasisPoints)
		assert.GreaterOrEqual(t, len(md.Creators), 1)
		if assert.NotNil(t, md.Data) {
			assert.NotEmpty(t, md.Data.Name)
			assert.NotEmpty(t, md.Data.Symbol)
			assert.NotEmpty(t, md.Data.Description)
			assert.NotEmpty(t, md.Data.Image)
		}
		assert.NotNil(t, md.Edition)
		assert.EqualValues(t, md.Edition.Type, token_metadata.KeyPrintedEdition.String())
		assert.EqualValues(t, md.Edition.MaxSupply, uint64(1000))
		assert.GreaterOrEqual(t, md.Edition.Supply, uint64(4))
		assert.EqualValues(t, md.Edition.Edition, uint64(2))
	})
}

func TestUpdateNftMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	// Get the metadata for the master edition
	md, err := client.GetTokenMetadata(ctx, e2e.MasterEditionMintAddr)
	require.NoError(t, err)
	require.NotEmpty(t, md)

	t.Run("update metadata", func(t *testing.T) {
		tx, err := client.UpdateMetadata(ctx, solana.UpdateMetadataParams{
			Mint:                 e2e.MasterEditionMintAddr,
			Owner:                e2e.Wallet1Addr,
			FeePayer:             e2e.FeePayerAddr,
			SellerFeeBasisPoints: md.SellerFeeBasisPoints + 100,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tx)

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
	})

	t.Run("check updated metadata", func(t *testing.T) {
		md2, err := client.GetTokenMetadata(ctx, e2e.MasterEditionMintAddr)
		require.NoError(t, err)
		require.NotEmpty(t, md2)
		assert.EqualValues(t, md2.SellerFeeBasisPoints, md.SellerFeeBasisPoints+100)
	})
}
