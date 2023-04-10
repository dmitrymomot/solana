package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dmitrymomot/solana/types"
	"github.com/dmitrymomot/solana/utils"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/rpc"
	sdktypes "github.com/portto/solana-go-sdk/types"
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
	latestBlockhash, err := c.rpcClient.GetLatestBlockhash(ctx)
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

	txb, err := utils.EncodeTransaction(tx)
	if err != nil {
		return "", utils.StackErrors(
			ErrNewTransaction,
			ErrSerializeTransaction,
			err,
		)
	}

	return txb, nil
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
	nonce, err := c.rpcClient.GetNonceFromNonceAccount(ctx, params.DurableNonce.ToBase58())
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

	txb, err := utils.EncodeTransaction(tx)
	if err != nil {
		return "", utils.StackErrors(
			ErrNewDurableTransaction,
			ErrSerializeTransaction,
			err,
		)
	}

	return txb, nil
}

// GetTransactionFee gets the fee for a transaction.
// Returns the fee or error.
func (c *Client) GetTransactionFee(ctx context.Context, txSource string) (uint64, error) {
	tx, err := utils.DecodeTransaction(txSource)
	if err != nil {
		return 0, utils.StackErrors(ErrGetTransactionFee, ErrDeserializeTransaction, err)
	}

	fee, err := c.rpcClient.GetFeeForMessage(ctx, tx.Message)
	if err != nil {
		return 0, utils.StackErrors(ErrGetTransactionFee, err)
	}

	return *fee, nil
}

// Sign transaction
// returns the signed transaction or an error
func (c *Client) SignTransaction(ctx context.Context, wallet sdktypes.Account, txSource string) (string, error) {
	tx, err := utils.DecodeTransaction(txSource)
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

	result, err := utils.EncodeTransaction(tx)
	if err != nil {
		return "", utils.StackErrors(ErrSignTransaction, ErrSerializeTransaction, err)
	}

	return result, nil
}

// Send transaction
// returns the transaction hash or an error
func (c *Client) SendTransaction(ctx context.Context, txSource string, i ...uint8) (string, error) {
	var tryN uint8 = 1
	if len(i) > 0 {
		tryN = i[0]
	}

	tx, err := utils.DecodeTransaction(txSource)
	if err != nil {
		return "", utils.StackErrors(ErrSendTransaction, ErrDeserializeTransaction, err)
	}

	txhash, err := c.rpcClient.SendTransaction(ctx, tx)
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
	status, err := c.rpcClient.GetSignatureStatus(ctx, txhash)
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
	mintAccountRent, err := c.rpcClient.GetMinimumBalanceForRentExemption(ctx, size)
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

// GetOldestTransactionForWallet returns the oldest transaction by the given base58 encoded public key.
// Returns the transaction or an error.
func (c *Client) GetOldestTransactionForWallet(
	ctx context.Context,
	base58Addr string,
	offsetTxSignature string,
) (string, *client.Transaction, error) {
	limit := 1000
	result, err := c.rpcClient.GetSignaturesForAddressWithConfig(ctx, base58Addr, client.GetSignaturesForAddressConfig{
		Limit:      limit,
		Before:     offsetTxSignature,
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to get signatures for address: %s: %w", base58Addr, err)
	}

	if l := len(result); l == 0 {
		return "", nil, ErrNoTransactionsFound
	} else if l < limit {
		tx := result[l-1]
		if tx.Err != nil {
			return "", nil, fmt.Errorf("transaction failed: %v", tx.Err)
		}
		if tx.Signature == "" {
			return "", nil, ErrNoTransactionsFound
		}
		if tx.BlockTime == nil || *tx.BlockTime == 0 || *tx.BlockTime > time.Now().Unix() {
			return "", nil, ErrTransactionNotConfirmed
		}

		resp, err := c.GetTransaction(ctx, tx.Signature)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get oldest transaction for wallet: %s: %w", base58Addr, err)
		}

		return tx.Signature, resp, nil
	}

	return c.GetOldestTransactionForWallet(ctx, base58Addr, result[limit-1].Signature)
}

// GetTransaction returns the transaction by the given base58 encoded transaction signature.
// Returns the transaction or an error.
func (c *Client) GetTransaction(ctx context.Context, txSignature string) (*client.Transaction, error) {
	tx, err := c.rpcClient.GetTransaction(ctx, txSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	if tx == nil || tx.Meta == nil {
		return nil, ErrTransactionNotFound
	}
	if tx.Meta.Err != nil {
		return nil, fmt.Errorf("transaction failed: %v", tx.Meta.Err)
	}

	return tx, nil
}

// ValidateTransactionByReference returns the transaction by the given reference.
// Returns transaction signature or an error if the transaction is not found or the transaction failed.
func (c *Client) ValidateTransactionByReference(ctx context.Context, reference, destination string, amount uint64, mint string) (string, error) {
	txSign, tx, err := c.GetOldestTransactionForWallet(ctx, reference, "")
	if err != nil {
		return "", fmt.Errorf("failed to validate transaction for reference %s: %w", reference, err)
	}

	if mint == "" || mint == "SOL" || mint == "So11111111111111111111111111111111111111112" {
		if err := CheckSolTransferTransaction(tx.Meta, tx.Transaction, destination, amount); err != nil {
			return "", fmt.Errorf("failed to validate transaction for reference %s: %w", reference, err)
		}
		return txSign, nil
	}

	if err := CheckTokenTransferTransaction(tx.Meta, tx.Transaction, mint, destination, amount); err != nil {
		return "", fmt.Errorf("failed to validate transaction for reference %s: %w", reference, err)
	}

	return txSign, nil
}
