package e2e

import (
	"context"
	"fmt"

	"github.com/solplaydev/solana/client"
	"github.com/solplaydev/solana/common"
	"github.com/solplaydev/solana/types"
)

// SignAndSendTransaction signs a transaction by the fee payer and the wallet1 and sends it.
// Returns the transaction hash and status or an error.
func SignAndSendTransaction(ctx context.Context, client *client.Client, tx string, signers ...string) (string, types.TransactionStatus, error) {
	if tx == "" {
		return "", types.TransactionStatusUnknown, fmt.Errorf("empty transaction")
	}
	if len(signers) > 0 {
		for _, signer := range signers {
			if signer == "" {
				return "", types.TransactionStatusUnknown, fmt.Errorf("empty signer")
			}

			signerAcc, err := common.AccountFromBase58(signer)
			if err != nil {
				return "", types.TransactionStatusUnknown, fmt.Errorf("failed to create signer account: %w", err)
			}
			tx, err = client.SignTransaction(ctx, signerAcc, tx)
			if err != nil {
				return "", types.TransactionStatusUnknown, fmt.Errorf("failed to sign transaction by signer: %w", err)
			}
		}
	}

	// Send the transaction
	txHash, err := client.SendTransaction(ctx, tx)
	if err != nil {
		return "", types.TransactionStatusFailure, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for the transaction to be confirmed
	txInfo, err := client.WaitForTransactionConfirmed(ctx, txHash, 0)
	if err != nil {
		return "", types.TransactionStatusUnknown, fmt.Errorf("failed to wait for transaction to be confirmed: %w", err)
	}

	return txHash, txInfo, nil
}
