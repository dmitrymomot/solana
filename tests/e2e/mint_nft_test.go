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

func TestMintNFT(t *testing.T) {
	var (
		tokenName   = "Test NFT"
		tokenSymbol = "TSTn"
		metadataUri = "https://www.arweave.net/jQ6ecVJtPZwaC-tsSYftEqaKsC8R3winHH2Z2hLxiBk?ext=json"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	// NFT mint public key
	var mint *string

	t.Run("Mint a Master Edition NFT", func(t *testing.T) {
		// Mint a non-fungible token
		mintAddr, tx, err := client.MintNonFungibleToken(ctx, solana.MintNonFungibleTokenParams{
			FeePayer: e2e.FeePayerAddr,
			Owner:    e2e.Wallet1Addr,

			Name:                 tokenName,
			Symbol:               tokenSymbol,
			MetadataURI:          metadataUri,
			Collection:           e2e.CollectionAddr,
			MaxSupply:            1000,
			SellerFeeBasisPoints: 1000,
			Creators: []token_metadata.Creator{
				{
					Address: e2e.FeePayerAddr,
					Share:   10,
				},
				{
					Address: e2e.Wallet1Addr,
					Share:   85,
				},
				{
					Address: e2e.Wallet2Addr,
					Share:   5,
				},
			},
			Uses: &token_metadata.Uses{
				UseMethod: token_metadata.TokenUseMethodBurn.String(),
				Total:     1,
				Remaining: 1,
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, tx)
		t.Logf("Mint address: %s", mintAddr)
		mint = &mintAddr

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
		balance, err := client.GetTokenBalance(ctx, e2e.Wallet1Addr, mintAddr)
		require.NoError(t, err)
		t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
		require.EqualValues(t, 1, balance.Amount)
		require.EqualValues(t, uint8(0), balance.Decimals)

		// Check token metadata
		metadata, err := client.GetTokenMetadata(ctx, mintAddr)
		require.NoError(t, err)
		t.Logf("Token metadata: %+v", metadata)
		require.EqualValues(t, tokenName, metadata.Data.Name)
		require.EqualValues(t, tokenSymbol, metadata.Data.Symbol)
		require.EqualValues(t, token_metadata.TokenStandardNonFungible.String(), metadata.TokenStandard)
	})

	t.Run("Verify NFT creator", func(t *testing.T) {
		if mint == nil {
			t.Skip("Mint address is not set")
		}

		// Get the metadata for the master edition
		md, err := client.GetTokenMetadata(ctx, *mint)
		require.NoError(t, err)
		require.NotEmpty(t, md)

		for _, creator := range md.Creators {
			if creator.Address == e2e.Wallet1Addr || creator.Address == e2e.FeePayerAddr {
				assert.Truef(t, creator.Verified, "creator %s should be verified", creator.Address)
			} else {
				assert.Falsef(t, creator.Verified, "creator %s should not be verified", creator.Address)
			}
		}

		tx, err := client.VerifyCreator(ctx, solana.VerifyCreatorParams{
			MintAddress:    *mint,
			CreatorAddress: e2e.Wallet2Addr,
			FeePayer:       e2e.FeePayerAddr,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		// Sign the transaction by the fee payer
		feePayer, err := solana.AccountFromBase58(e2e.FeePayerPrivateKey)
		require.NoError(t, err)
		tx, err = client.SignTransaction(ctx, feePayer, tx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		// Sign the transaction by the creator
		creator, err := solana.AccountFromBase58(e2e.Wallet2PrivateKey)
		require.NoError(t, err)
		tx, err = client.SignTransaction(ctx, creator, tx)
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

		// Get the metadata for the master edition
		newMd, err := client.GetTokenMetadata(ctx, *mint)
		require.NoError(t, err)
		require.NotEmpty(t, newMd)

		// check creators
		for _, creator := range newMd.Creators {
			assert.Truef(t, creator.Verified, "creator %s should be verified", creator.Address)
		}
	})
}
