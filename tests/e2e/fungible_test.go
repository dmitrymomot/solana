package e2e_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/portto/solana-go-sdk/program/token"
	"github.com/solplaydev/solana/client"
	"github.com/solplaydev/solana/common"
	"github.com/solplaydev/solana/instructions"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/transaction"
	"github.com/solplaydev/solana/types"
	"github.com/stretchr/testify/require"
)

func TestFungibleToken(t *testing.T) {
	var (
		tokenNameInit   = "Test Token Init"
		tokenSymbolInit = "TSTi"
		tokenName       = "Test Token"
		tokenSymbol     = "TSTt"
		metadataUri     = "https://www.arweave.net/QR1PsBgIbiYoKgGff5Jq2U8QavHChRjBki8XRJ-06mI?ext=json"
		supplyAmount    = 1000000 * types.SPLTokenDefaultMultiplier
		mint            = common.NewAccount()
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new sc
	sc := client.New(client.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	t.Run("mint fungible token", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddSigner(mint).
			AddInstruction(instructions.MintFungible(instructions.MintFungibleParam{
				Mint:          mint.PublicKey,
				MintTo:        e2e.Wallet1Pubkey,
				FeePayer:      &e2e.FeePayerPubkey,
				Decimals:      types.SPLTokenDefaultDecimals,
				SupplyAmount:  supplyAmount,
				IsFixedSupply: true,
				TokenName:     tokenNameInit,
				TokenSymbol:   tokenSymbolInit,
				// MetadataURI:   metadataUri,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		t.Logf("Transaction hash: %s", txHash)
		require.NotEmpty(t, txHash)
		t.Logf("Transaction status: %+v", txStatus)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance
		t.Run("check token balance", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet1Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
			require.EqualValues(t, supplyAmount, balance.Amount)
			require.EqualValues(t, types.SPLTokenDefaultDecimals, balance.Decimals)
		})

		// Check token metadata
		t.Run("check token metadata", func(t *testing.T) {
			metadata, err := sc.GetTokenMetadata(ctx, mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token metadata: %+v", metadata)
			require.EqualValues(t, tokenNameInit, metadata.Data.Name)
			require.EqualValues(t, tokenSymbolInit, metadata.Data.Symbol)
			require.EqualValues(t, token_metadata.TokenStandardFungible, metadata.TokenStandard)
		})
	})

	t.Run("update token metadata", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.UpdateMetadata(instructions.UpdateMetadataParams{
				Mint:            mint.PublicKey,
				UpdateAuthority: e2e.Wallet1Pubkey,
				MetadataUri:     &metadataUri,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		t.Logf("Transaction hash: %s", txHash)
		require.NotEmpty(t, txHash)
		t.Logf("Transaction status: %+v", txStatus)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token metadata
		t.Run("check token metadata", func(t *testing.T) {
			metadata, err := sc.GetTokenMetadata(ctx, mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token metadata: %+v", metadata)
			require.EqualValues(t, tokenName, metadata.Data.Name)
			require.EqualValues(t, tokenSymbol, metadata.Data.Symbol)
			require.EqualValues(t, token_metadata.TokenStandardFungible, metadata.TokenStandard)
		})
	})

	// Transfer token
	t.Run("transfer fungible token: wallet => wallet", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.CreateAssociatedTokenAccountIfNotExists(
				instructions.CreateAssociatedTokenAccountParam{
					Funder: e2e.FeePayerPubkey,
					Owner:  e2e.Wallet2Pubkey,
					Mint:   mint.PublicKey,
				},
			)).
			AddInstruction(instructions.TransferToken(instructions.TransferTokenParam{
				Mint:      mint.PublicKey,
				Amount:    supplyAmount / 2,
				Sender:    e2e.Wallet1Pubkey,
				Recipient: e2e.Wallet2Pubkey,
			})).
			AddInstruction(instructions.Memo(
				fmt.Sprintf("Send %d %s to %s", supplyAmount/2, tokenSymbol, e2e.Wallet2Pubkey.ToBase58()),
			)).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		t.Logf("Transaction hash: %s", txHash)
		require.NotEmpty(t, txHash)
		t.Logf("Transaction status: %+v", txStatus)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance of wallet 1
		t.Run("check token balance of wallet 1", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet1Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
			require.EqualValues(t, supplyAmount/2, balance.Amount)
			require.EqualValues(t, types.SPLTokenDefaultDecimals, balance.Decimals)
		})

		// Check token balance of wallet 2
		t.Run("check token balance of wallet 2", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
			require.EqualValues(t, supplyAmount/2, balance.Amount)
			require.EqualValues(t, types.SPLTokenDefaultDecimals, balance.Decimals)
		})
	})

	// Transfer token
	t.Run("transfer rest of fungible token: wallet => wallet", func(t *testing.T) {
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.CreateAssociatedTokenAccountIfNotExists(
				instructions.CreateAssociatedTokenAccountParam{
					Funder: e2e.FeePayerPubkey,
					Owner:  e2e.Wallet2Pubkey,
					Mint:   mint.PublicKey,
				},
			)).
			AddInstruction(instructions.TransferToken(instructions.TransferTokenParam{
				Mint:      mint.PublicKey,
				Amount:    supplyAmount / 2,
				Sender:    e2e.Wallet1Pubkey,
				Recipient: e2e.Wallet2Pubkey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		t.Logf("Transaction hash: %s", txHash)
		require.NotEmpty(t, txHash)
		t.Logf("Transaction status: %+v", txStatus)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance of wallet 1
		t.Run("check token balance of wallet 1", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet1Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
			require.EqualValues(t, uint64(0), balance.Amount)
			require.EqualValues(t, types.SPLTokenDefaultDecimals, balance.Decimals)
		})

		// Check token balance of wallet 2
		t.Run("check token balance of wallet 2", func(t *testing.T) {
			balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			t.Logf("Token balance: %d, decimals: %d", balance.Amount, balance.Decimals)
			require.EqualValues(t, supplyAmount, balance.Amount)
			require.EqualValues(t, types.SPLTokenDefaultDecimals, balance.Decimals)
		})
	})

	// Close empty token account
	t.Run("close empty token account", func(t *testing.T) {
		// close token account
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.CloseTokenAccount(instructions.CloseTokenAccountParams{
				Owner: e2e.Wallet1Pubkey,
				Mint:  &mint.PublicKey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet1PrivateKey)
		require.NoError(t, err)
		t.Logf("Transaction hash: %s", txHash)
		require.NotEmpty(t, txHash)
		t.Logf("Transaction status: %+v", txStatus)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token account info
		t.Run("check token account info", func(t *testing.T) {
			ata, err := common.DeriveTokenAccount(e2e.Wallet1Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			ataInfo, err := sc.GetTokenAccountInfo(ctx, ata.ToBase58())
			require.Error(t, err)
			require.EqualValues(t, ataInfo, token.TokenAccount{})
		})
	})

	// Burn token
	t.Run("burn fungible token and close token account", func(t *testing.T) {
		balance, err := sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
		require.NoError(t, err)
		t.Logf("Token balance: %d", balance.Amount)

		// Burn token and close token account
		tx, err := transaction.NewTransactionBuilder(sc).
			SetFeePayer(e2e.FeePayerPubkey).
			AddInstruction(instructions.BurnToken(instructions.BurnTokenParams{
				Mint:              mint.PublicKey,
				Amount:            balance.Amount,
				TokenAccountOwner: e2e.Wallet2Pubkey,
			})).
			AddInstruction(instructions.CloseTokenAccount(instructions.CloseTokenAccountParams{
				Owner: e2e.Wallet2Pubkey,
				Mint:  &mint.PublicKey,
			})).
			Build(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		txHash, txStatus, err := e2e.SignAndSendTransaction(ctx, sc, tx, e2e.FeePayerPrivateKey, e2e.Wallet2PrivateKey)
		require.NoError(t, err)
		t.Logf("Transaction hash: %s", txHash)
		require.NotEmpty(t, txHash)
		t.Logf("Transaction status: %+v", txStatus)
		require.EqualValues(t, txStatus, types.TransactionStatusSuccess)

		// Check token balance of wallet 2
		t.Run("check token balance of wallet 2", func(t *testing.T) {
			balance, err = sc.GetTokenBalance(ctx, e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.Error(t, err)
		})

		t.Run("check token account info", func(t *testing.T) {
			ata, err := common.DeriveTokenAccount(e2e.Wallet2Pubkey.ToBase58(), mint.PublicKey.ToBase58())
			require.NoError(t, err)
			ataInfo, err := sc.GetTokenAccountInfo(ctx, ata.ToBase58())
			require.Error(t, err)
			require.EqualValues(t, ataInfo, token.TokenAccount{})
		})
	})
}
