package solana

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/utils"
)

// NewTransactionParams is the params for NewTransaction function.
type NewTransactionParams struct {
	Base58FeePayerAddr string
	Instructions       []types.Instruction
}

// NewTransaction creates a new transaction.
// Returns the transaction or an error.
func (c *Client) NewTransaction(ctx context.Context, params NewTransactionParams) ([]byte, error) {
	latestBlockhash, err := c.solana.GetLatestBlockhash(ctx)
	if err != nil {
		return nil, utils.StackErrors(ErrGetLatestBlockhash, err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        common.PublicKeyFromString(params.Base58FeePayerAddr),
			RecentBlockhash: latestBlockhash.Blockhash,
			Instructions:    params.Instructions,
		}),
	})
	if err != nil {
		return nil, utils.StackErrors(ErrNewTransaction, err)
	}

	txb, err := tx.Serialize()
	if err != nil {
		return nil, utils.StackErrors(ErrSerializeTransaction, err)
	}

	return txb, nil
}

// NewDurableTransactionParams are the parameters for NewDurableTransaction function.
type NewDurableTransactionParams struct {
	Base58FeePayerAddr     string
	Base58NonceAuthAddr    string
	Base58DurableNonceAddr string
	Instructions           []types.Instruction
}

// NewDurableTransaction creates a new durable transaction.
// base58FeePayerAddr is the base58 encoded fee payer address.
// base58DurableNonceAddr is the base58 encoded durable nonce address.
// instructions is the transaction instructions.
// Returns the serialized transaction or an error.
func (c *Client) NewDurableTransaction(ctx context.Context, params NewDurableTransactionParams) ([]byte, error) {
	nonce, err := c.solana.GetNonceFromNonceAccount(ctx, params.Base58DurableNonceAddr)
	if err != nil {
		return nil, utils.StackErrors(ErrGetNonceFromNonceAccount, err)
	}

	feePayerPublicKey := common.PublicKeyFromString(params.Base58FeePayerAddr)

	base58NonceAuthAddr := params.Base58NonceAuthAddr
	if base58NonceAuthAddr == "" {
		base58NonceAuthAddr = params.Base58FeePayerAddr
	}

	instr := []types.Instruction{
		system.AdvanceNonceAccount(system.AdvanceNonceAccountParam{
			Nonce: common.PublicKeyFromString(params.Base58DurableNonceAddr),
			Auth:  common.PublicKeyFromString(base58NonceAuthAddr),
		}),
	}
	instr = append(instr, params.Instructions...)

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayerPublicKey,
			RecentBlockhash: nonce,
			Instructions:    instr,
		}),
	})
	if err != nil {
		return nil, utils.StackErrors(ErrNewTransaction, err)
	}

	txb, err := tx.Serialize()
	if err != nil {
		return nil, utils.StackErrors(ErrSerializeTransaction, err)
	}

	return txb, nil
}

// GetTransactionFee gets the fee for a transaction.
// Returns the fee or error.
func (c *Client) GetTransactionFee(ctx context.Context, txSource []byte) (uint64, error) {
	tx, err := types.TransactionDeserialize(txSource)
	if err != nil {
		return 0, utils.StackErrors(ErrGetTransactionFee, ErrDeserializeTransaction, err)
	}

	fee, err := c.solana.GetFeeForMessage(ctx, tx.Message)
	if err != nil {
		return 0, utils.StackErrors(ErrGetTransactionFee, err)
	}

	return *fee, nil
}

// Sign transaction
// returns the signed transaction or an error
func (c *Client) SignTransaction(ctx context.Context, wallet types.Account, txSource []byte) ([]byte, error) {
	tx, err := types.TransactionDeserialize(txSource)
	if err != nil {
		return nil, utils.StackErrors(ErrSignTransaction, ErrDeserializeTransaction, err)
	}

	msg, err := tx.Message.Serialize()
	if err != nil {
		return nil, utils.StackErrors(ErrSignTransaction, ErrSerializeMessage, err)
	}

	if err := tx.AddSignature(wallet.Sign(msg)); err != nil {
		return nil, utils.StackErrors(ErrSignTransaction, ErrAddSignature, err)
	}

	result, err := tx.Serialize()
	if err != nil {
		return nil, utils.StackErrors(ErrSignTransaction, ErrSerializeTransaction, err)
	}

	return result, nil
}

// Send transaction
// returns the transaction hash or an error
func (c *Client) SendTransaction(ctx context.Context, txSource []byte) (string, error) {
	tx, err := types.TransactionDeserialize(txSource)
	if err != nil {
		return "", utils.StackErrors(ErrSendTransaction, ErrDeserializeTransaction, err)
	}

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		if strings.Contains(err.Error(), "without insufficient funds for rent") {
			return "", utils.StackErrors(ErrSendTransaction, ErrWithoutInsufficientFound, err)
		}
		return "", utils.StackErrors(ErrSendTransaction, err)
	}

	return txhash, nil
}

// GetTransactionStatus gets the transaction status.
// Returns the transaction status or an error.
func (c *Client) GetTransactionStatus(ctx context.Context, txhash string) (TransactionStatus, error) {
	status, err := c.solana.GetSignatureStatus(ctx, txhash)
	if err != nil {
		return TransactionStatusUnknown, utils.StackErrors(ErrGetTransactionStatus, err)
	}

	if status == nil {
		return TransactionStatusUnknown, nil
	}

	if status.Err != nil {
		return TransactionStatusFailure, fmt.Errorf("transaction failed: %v", status.Err)
	}

	result := TransactionStatusUnknown

	if status.Confirmations != nil && *status.Confirmations > 0 {
		result = TransactionStatusInProgress
	}

	if status.ConfirmationStatus != nil {
		result = ParseTransactionStatus(*status.ConfirmationStatus)
	}

	return result, nil
}

// GetMinimumBalanceForRentExemption gets the minimum balance for rent exemption.
// Returns the minimum balance in lamports or an error.
func (c *Client) GetMinimumBalanceForRentExemption(ctx context.Context, size uint64) (uint64, error) {
	mintAccountRent, err := c.solana.GetMinimumBalanceForRentExemption(ctx, size)
	if err != nil {
		return 0, utils.StackErrors(ErrGetMinimumBalanceForRentExemption, err)
	}

	return mintAccountRent, nil
}

// WaitForTransactionConfirmed waits for a transaction to be confirmed.
// Returns the transaction status or an error.
func (c *Client) WaitForTransactionConfirmed(ctx context.Context, txhash string) (TransactionStatus, error) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return TransactionStatusUnknown, utils.StackErrors(ErrWaitForTransaction, ErrContextDone)
		case <-tick.C:
			status, err := c.GetTransactionStatus(ctx, txhash)
			if err != nil {
				return TransactionStatusUnknown, utils.StackErrors(ErrWaitForTransaction, err)
			}

			if status == TransactionStatusInProgress || status == TransactionStatusUnknown {
				continue
			}

			if status == TransactionStatusFailure || status == TransactionStatusSuccess {
				return status, nil
			}
		}
	}
}
