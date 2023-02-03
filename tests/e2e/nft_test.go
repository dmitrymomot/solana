package e2e_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/solplaydev/solana/client"
	"github.com/solplaydev/solana/common"
	"github.com/solplaydev/solana/instructions"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/transaction"
	"github.com/solplaydev/solana/types"
	"github.com/solplaydev/solana/utils"
	"github.com/stretchr/testify/require"
)

func TestNFT(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new sc
	sc := client.New(client.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	var (
		tokenName   = "Test NFT"
		tokenSymbol = "TSTn"
		metadataUri = "https://www.arweave.net/jQ6ecVJtPZwaC-tsSYftEqaKsC8R3winHH2Z2hLxiBk?ext=json"
		collection  = common.NewAccount()
		mint        = common.NewAccount()
		editionMint = common.NewAccount()
	)

	// Display account public keys
	{
		fmt.Println("Collection:", collection.PublicKey.ToBase58())
		fmt.Println("Master NFT:", mint.PublicKey.ToBase58())
		fmt.Println("Edition NFT:", editionMint.PublicKey.ToBase58())
	}
	// Mint collection
	t.Run("Mint collection", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddSigner(collection).
			AddInstruction(instructions.MintNonFungible(instructions.MintNonFungibleParam{
				Mint:        collection.PublicKey,
				Owner:       e2e.Wallet1Pubkey,
				FeePayer:    &e2e.FeePayerPubkey,
				TokenName:   "Test collection",
				TokenSymbol: "TSTc",
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)
	})

	// Mint NFT
	t.Run("Mint NFT", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddSigner(mint).
			AddInstruction(instructions.MintNonFungible(instructions.MintNonFungibleParam{
				Mint:                 mint.PublicKey,
				Owner:                e2e.Wallet1Pubkey,
				FeePayer:             &e2e.FeePayerPubkey,
				Collection:           &collection.PublicKey,
				MetadataURI:          metadataUri,
				SellerFeeBasisPoints: 1000,
				MaxEditionSupply:     10,
				UseMethod:            utils.Pointer(token_metadata.TokenUseMethodBurn),
				UseLimit:             utils.Pointer[uint64](1),
			})).
			AddInstruction(instructions.VerifyCollectionItem(instructions.VerifyCollectionItemParams{
				Mint:                mint.PublicKey,
				CollectionMint:      collection.PublicKey,
				CollectionAuthority: e2e.Wallet1Pubkey,
				FeePayer:            &e2e.FeePayerPubkey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token metadata
		t.Run("check token metadata", func(t *testing.T) {
			metadata, err := sc.GetTokenMetadata(ctx, mint.PublicKey.ToBase58())
			require.NoError(t, err)
			require.EqualValues(t, tokenName, metadata.Data.Name)
			require.EqualValues(t, tokenSymbol, metadata.Data.Symbol)
			require.EqualValues(t, token_metadata.TokenStandardNonFungible, metadata.TokenStandard)
			require.EqualValues(t, collection.PublicKey.ToBase58(), metadata.Collection.Key)
			require.True(t, metadata.Collection.Verified)
		})

		// Check token balance
		t.Run("check token balance", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet1Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			require.EqualValues(t, uint64(1), balance.Amount)
		})

		// Check token supply
		t.Run("check token supply", func(t *testing.T) {
			current, max, err := sc.GetMasterEditionSupply(ctx, mint.PublicKey)
			require.NoError(t, err)
			require.EqualValues(t, uint64(0), current)
			require.EqualValues(t, uint64(10), max)
		})
	})

	// Mint NFT edition
	t.Run("Mint NFT edition", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddSigner(editionMint).
			AddInstruction(instructions.MintNonFungibleEdition(instructions.MintNonFungibleEditionParam{
				FeePayer:           e2e.FeePayerPubkey,
				MasterEditionMint:  mint.PublicKey,
				MasterEditionOwner: e2e.Wallet1Pubkey,
				EditionMint:        editionMint.PublicKey,
				EditionOwner:       e2e.Wallet2Pubkey,
			})).
			// AddInstruction(instructions.ApproveUseAuthority(instructions.ApproveUseAuthorityParams{
			// 	FeePayer:        e2e.FeePayerPubkey,
			// 	Mint:            editionMint.PublicKey,
			// 	MintOwner:       e2e.Wallet2Pubkey,
			// 	NewUseAuthority: e2e.Wallet1Pubkey,
			// 	NumberOfUses:    1,
			// })).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(
			ctx, sc, tx,
			e2e.FeePayerPrivateKey,
			e2e.Wallet1PrivateKey,
			// e2e.Wallet2PrivateKey,
		)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance
		t.Run("check token balance: wallet 2", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), editionMint.PublicKey.ToBase58())
			require.NoError(t, err)
			require.EqualValues(t, uint64(1), balance.Amount)
		})

		// Check token metadata
		t.Run("check token edition metadata", func(t *testing.T) {
			metadata, err := sc.GetTokenMetadata(ctx, editionMint.PublicKey.ToBase58())
			require.NoError(t, err)
			require.EqualValues(t, tokenName, metadata.Data.Name)
			require.EqualValues(t, tokenSymbol, metadata.Data.Symbol)
			require.EqualValues(t, token_metadata.TokenStandardNonFungibleEdition, metadata.TokenStandard)
			require.EqualValues(t, collection.PublicKey.ToBase58(), metadata.Collection.Key)
			require.True(t, metadata.Collection.Verified)
		})
	})

	// Transfer NFT
	t.Run("Transfer NFT", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.CreateAssociatedTokenAccount(
				instructions.CreateAssociatedTokenAccountParam{
					Funder: e2e.FeePayerPubkey,
					Owner:  e2e.Wallet2Pubkey,
					Mint:   mint.PublicKey,
				},
			)).
			AddInstruction(instructions.TransferToken(instructions.TransferTokenParam{
				Mint:      mint.PublicKey,
				Amount:    1,
				Sender:    e2e.Wallet1Pubkey,
				Recipient: e2e.Wallet2Pubkey,
			})).
			AddInstruction(instructions.CloseTokenAccount(instructions.CloseTokenAccountParams{
				Owner: e2e.Wallet1Pubkey,
				Mint:  &mint.PublicKey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance
		t.Run("check token balance: wallet 1", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet1Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.Error(t, err)
			require.EqualValues(t, uint64(0), balance.Amount)
		})

		// Check token balance
		t.Run("check token balance: wallet 2", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			require.EqualValues(t, uint64(1), balance.Amount)
		})
	})

	// Burn NFT edition
	t.Run("Burn NFT edition", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.BurnNftEdition(instructions.BurnNftEditionParams{
				MasterMint:       mint.PublicKey,
				MasterMintOwner:  e2e.Wallet2Pubkey,
				EditionMint:      editionMint.PublicKey,
				EditionMintOwner: e2e.Wallet2Pubkey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet2PrivateKey)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance
		t.Run("check token balance: wallet 2", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), editionMint.PublicKey.ToBase58())
			require.Error(t, err)
			require.EqualValues(t, uint64(0), balance.Amount)
		})

		// Check token metadata
		t.Run("check nft edition metadata", func(t *testing.T) {
			_, err = sc.GetTokenMetadata(ctx, editionMint.PublicKey.ToBase58())
			require.Error(t, err)
		})
	})

	// Burn NFT
	t.Run("Burn master edition NFT", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.BurnNft(instructions.BurnNftParams{
				Mint:           mint.PublicKey,
				MintOwner:      e2e.Wallet2Pubkey,
				CollectionMint: &collection.PublicKey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet2PrivateKey)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token metadata
		t.Run("check master metadata", func(t *testing.T) {
			_, err = sc.GetTokenMetadata(ctx, mint.PublicKey.ToBase58())
			require.Error(t, err)
		})

		// Check token balance
		t.Run("check token balance: wallet 2", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.Error(t, err)
			require.EqualValues(t, uint64(0), balance.Amount)
		})
	})

	// Burn NFT collection
	t.Run("Burn NFT collection", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.BurnNft(instructions.BurnNftParams{
				Mint:      collection.PublicKey,
				MintOwner: e2e.Wallet1Pubkey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		fmt.Println("tx:", txHash, "status:", txStatus)
		require.NotEmpty(t, txHash)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token metadata
		t.Run("check master metadata", func(t *testing.T) {
			_, err = sc.GetTokenMetadata(ctx, collection.PublicKey.ToBase58())
			require.Error(t, err)
		})

		// Check token balance
		t.Run("check token balance: wallet 1", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet1Pubkey.ToBase58(), collection.PublicKey.ToBase58())
			require.Error(t, err)
			require.EqualValues(t, uint64(0), balance.Amount)
		})
	})
}
