package client

import (
	"fmt"
	"strconv"

	"github.com/dmitrymomot/solana/utils"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/types"
)

// SignTransaction signs a transaction and returns a base64 encoded transaction.
func SignTransaction(txSource string, signer types.Account) (string, error) {
	tx, err := utils.DecodeTransaction(txSource)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: base64 to bytes: %w", err)
	}

	msg, err := tx.Message.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: serialize message: %w", err)
	}

	if err := tx.AddSignature(signer.Sign(msg)); err != nil {
		return "", fmt.Errorf("failed to sign transaction: add signature: %w", err)
	}

	result, err := utils.EncodeTransaction(tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: encode transaction: %w", err)
	}

	return result, nil
}

// CheckSolTransferTransaction checks if a transaction is a SOL transfer transaction.
// Verifies that destination account has been credited with the correct amount.
func CheckSolTransferTransaction(meta *client.TransactionMeta, tx types.Transaction, destination string, amount uint64) error {
	var destIdx int
	for i, acc := range tx.Message.Accounts {
		if acc.ToBase58() == destination {
			destIdx = i
			break
		}
	}

	txAmount := meta.PostBalances[destIdx] - meta.PreBalances[destIdx]
	if txAmount != int64(amount) {
		return fmt.Errorf("amount is not equal to the amount in the transaction: %d != %d", amount, txAmount)
	}

	return nil
}

// CheckTokenTransferTransaction checks if a transaction is a token transfer transaction.
// Verifies that destination account has been credited with the correct amount of the token.
func CheckTokenTransferTransaction(meta *client.TransactionMeta, tx types.Transaction, mint, destination string, amount uint64) error {
	var preBalance uint64
	var postBalance uint64

	for _, balance := range meta.PreTokenBalances {
		if balance.Mint == mint && balance.Owner == destination {
			amount, err := strconv.ParseUint(balance.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse pre balance: %w", err)
			}
			preBalance = amount
			break
		}
	}

	for _, balance := range meta.PostTokenBalances {
		if balance.Mint == mint && balance.Owner == destination {
			amount, err := strconv.ParseUint(balance.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse post balance: %w", err)
			}
			postBalance = amount
			break
		}
	}

	if postBalance-preBalance != amount {
		return fmt.Errorf("amount is not equal to the amount in the transaction: %d != %d", amount, postBalance-preBalance)
	}

	return nil
}
