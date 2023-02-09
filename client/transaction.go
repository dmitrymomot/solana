package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	sdktypes "github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/types"
	"github.com/solplaydev/solana/utils"
)

// NewTransactionParams is the params for NewTransaction function.
type NewTransactionParams struct {
	FeePayer     common.PublicKey       // transaction fee payer
	Instructions []sdktypes.Instruction // transaction instructions
	Signers      []sdktypes.Account     // transaction signers
}

// NewTransaction creates a new transaction.
// Returns the transaction or an error.
func (c *Client) NewTransaction(ctx context.Context, params NewTransactionParams) (string, error) {
	latestBlockhash, err := c.solana.GetLatestBlockhash(ctx)
	if err != nil {
		return "", utils.StackErrors(
			ErrNewTransaction,
			ErrGetLatestBlockhash,
			err,
		)
	}

	tx, err := sdktypes.NewTransaction(sdktypes.NewTransactionParam{
		Message: sdktypes.NewMessage(sdktypes.NewMessageParam{
			FeePayer:        params.FeePayer,
			RecentBlockhash: latestBlockhash.Blockhash,
			Instructions:    params.Instructions,
		}),
		Signers: params.Signers,
	})
	if err != nil {
		return "", utils.StackErrors(ErrNewTransaction, err)
	}

	txb, err := tx.Serialize()
	if err != nil {
		return "", utils.StackErrors(
			ErrNewTransaction,
			ErrSerializeTransaction,
			err,
		)
	}

	return utils.BytesToBase64(txb), nil
}

// NewDurableTransactionParams are the parameters for NewDurableTransaction function.
type NewDurableTransactionParams struct {
	FeePayer     *common.PublicKey      // optional; if not provided, the fee payer will be the durable nonce
	NonceAuth    common.PublicKey       // required; the nonce authority
	DurableNonce common.PublicKey       // required; the durable nonce
	Instructions []sdktypes.Instruction // required; the transaction instructions
	Signers      []sdktypes.Account     // transaction signers
}

// NewDurableTransaction creates a new durable transaction.
// Returns the serialized transaction or an error.
func (c *Client) NewDurableTransaction(ctx context.Context, params NewDurableTransactionParams) (string, error) {
	nonce, err := c.solana.GetNonceFromNonceAccount(ctx, params.DurableNonce.ToBase58())
	if err != nil {
		return "", utils.StackErrors(
			ErrNewDurableTransaction,
			ErrGetNonceFromNonceAccount,
			err,
		)
	}

	if params.FeePayer != nil && *params.FeePayer == (common.PublicKey{}) {
		return "", utils.StackErrors(
			ErrNewDurableTransaction,
			fmt.Errorf("invalid fee payer public key: %s", params.FeePayer.ToBase58()),
		)
	}

	if params.FeePayer == nil {
		params.FeePayer = &params.NonceAuth
	}

	instr := []sdktypes.Instruction{
		system.AdvanceNonceAccount(system.AdvanceNonceAccountParam{
			Nonce: params.DurableNonce,
			Auth:  params.NonceAuth,
		}),
	}
	instr = append(instr, params.Instructions...)

	tx, err := sdktypes.NewTransaction(sdktypes.NewTransactionParam{
		Message: sdktypes.NewMessage(sdktypes.NewMessageParam{
			FeePayer:        *params.FeePayer,
			RecentBlockhash: nonce,
			Instructions:    instr,
		}),
		Signers: params.Signers,
	})
	if err != nil {
		return "", utils.StackErrors(
			ErrNewDurableTransaction,
			ErrNewTransaction,
			err,
		)
	}

	txb, err := tx.Serialize()
	if err != nil {
		return "", utils.StackErrors(
			ErrNewDurableTransaction,
			ErrSerializeTransaction,
			err,
		)
	}

	return utils.BytesToBase64(txb), nil
}

// GetTransactionFee gets the fee for a transaction.
// Returns the fee or error.
func (c *Client) GetTransactionFee(ctx context.Context, txSource string) (uint64, error) {
	txb, err := utils.Base64ToBytes(txSource)
	if err != nil {
		return 0, utils.StackErrors(ErrGetTransactionFee, err)
	}

	tx, err := sdktypes.TransactionDeserialize(txb)
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
func (c *Client) SignTransaction(ctx context.Context, wallet sdktypes.Account, txSource string) (string, error) {
	txb, err := utils.Base64ToBytes(txSource)
	if err != nil {
		return "", utils.StackErrors(ErrGetTransactionFee, err)
	}

	tx, err := sdktypes.TransactionDeserialize(txb)
	if err != nil {
		return "", utils.StackErrors(ErrSignTransaction, ErrDeserializeTransaction, err)
	}

	msg, err := tx.Message.Serialize()
	if err != nil {
		return "", utils.StackErrors(ErrSignTransaction, ErrSerializeMessage, err)
	}

	if err := tx.AddSignature(wallet.Sign(msg)); err != nil {
		return "", utils.StackErrors(ErrSignTransaction, ErrAddSignature, err)
	}

	result, err := tx.Serialize()
	if err != nil {
		return "", utils.StackErrors(ErrSignTransaction, ErrSerializeTransaction, err)
	}

	return utils.BytesToBase64(result), nil
}

// Send transaction
// returns the transaction hash or an error
func (c *Client) SendTransaction(ctx context.Context, txSource string, i ...uint8) (string, error) {
	var tryN uint8 = 1
	if len(i) > 0 {
		tryN = i[0]
	}

	txb, err := utils.Base64ToBytes(txSource)
	if err != nil {
		return "", utils.StackErrors(ErrGetTransactionFee, err)
	}

	tx, err := sdktypes.TransactionDeserialize(txb)
	if err != nil {
		return "", utils.StackErrors(ErrSendTransaction, ErrDeserializeTransaction, err)
	}

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		if strings.Contains(err.Error(), "without insufficient funds for rent") {
			return "", utils.StackErrors(ErrSendTransaction, ErrWithoutInsufficientFound, err)
		}

		// retry if blockhash not found
		if strings.Contains(err.Error(), "BlockhashNotFound") && tryN < 3 {
			return c.SendTransaction(ctx, txSource, tryN+1)
		}

		return "", utils.StackErrors(ErrSendTransaction, err)
	}

	return txhash, nil
}

// GetTransactionStatus gets the transaction status.
// Returns the transaction status or an error.
func (c *Client) GetTransactionStatus(ctx context.Context, txhash string) (types.TransactionStatus, error) {
	status, err := c.solana.GetSignatureStatus(ctx, txhash)
	if err != nil {
		return types.TransactionStatusUnknown, utils.StackErrors(ErrGetTransactionStatus, err)
	}
	if status == nil {
		return types.TransactionStatusUnknown, nil
	}
	if status.Err != nil {
		return types.TransactionStatusFailure, fmt.Errorf("transaction failed: %v", status.Err)
	}

	result := types.TransactionStatusUnknown
	if status.Confirmations != nil && *status.Confirmations > 0 {
		result = types.TransactionStatusInProgress
	}
	if status.ConfirmationStatus != nil {
		result = types.ParseTransactionStatus(*status.ConfirmationStatus)
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
func (c *Client) WaitForTransactionConfirmed(ctx context.Context, txhash string, maxDuration time.Duration) (types.TransactionStatus, error) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	if maxDuration == 0 {
		maxDuration = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, maxDuration)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return types.TransactionStatusUnknown, utils.StackErrors(ErrWaitForTransaction, ErrContextDone)
		case <-tick.C:
			status, err := c.GetTransactionStatus(ctx, txhash)
			if err != nil {
				return types.TransactionStatusUnknown, utils.StackErrors(ErrWaitForTransaction, err)
			}
			if status == types.TransactionStatusInProgress || status == types.TransactionStatusUnknown {
				continue
			}
			if status == types.TransactionStatusFailure || status == types.TransactionStatusSuccess {
				return status, nil
			}
		}
	}
}
