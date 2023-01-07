package solana

import (
	"context"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/types"
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
		return nil, errors.Wrap(ErrGetLatestBlockhash, err.Error())
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        common.PublicKeyFromString(params.Base58FeePayerAddr),
			RecentBlockhash: latestBlockhash.Blockhash,
			Instructions:    params.Instructions,
		}),
	})
	if err != nil {
		return nil, errors.Wrap(ErrNewTransaction, err.Error())
	}

	txb, err := tx.Serialize()
	if err != nil {
		return nil, errors.Wrap(ErrSerializeTransaction, err.Error())
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
		return nil, errors.Wrap(ErrGetNonceFromNonceAccount, err.Error())
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
		return nil, errors.Wrap(ErrNewTransaction, err.Error())
	}

	txb, err := tx.Serialize()
	if err != nil {
		return nil, errors.Wrap(ErrSerializeTransaction, err.Error())
	}

	return txb, nil
}

// GetTransactionFee gets the fee for a transaction.
// Returns the fee or error.
func (c *Client) GetTransactionFee(ctx context.Context, txSource []byte) (uint64, error) {
	tx, err := types.TransactionDeserialize(txSource)
	if err != nil {
		return 0, errors.Wrap(ErrDeserializeTransaction, err.Error())
	}

	fee, err := c.solana.GetFeeForMessage(ctx, tx.Message)
	if err != nil {
		return 0, errors.Wrap(ErrGetTransactionFee, err.Error())
	}

	return *fee, nil
}

// Sign transaction
// returns the signed transaction or an error
func (c *Client) SignTransaction(ctx context.Context, wallet types.Account, txSource []byte) ([]byte, error) {
	tx, err := types.TransactionDeserialize(txSource)
	if err != nil {
		return nil, errors.Wrap(ErrDeserializeTransaction, err.Error())
	}

	msg, err := tx.Message.Serialize()
	if err != nil {
		return nil, errors.Wrap(ErrSerializeMessage, err.Error())
	}

	if err := tx.AddSignature(wallet.Sign(msg)); err != nil {
		return nil, errors.Wrap(ErrAddSignature, err.Error())
	}

	result, err := tx.Serialize()
	if err != nil {
		return nil, errors.Wrap(ErrSerializeTransaction, err.Error())
	}

	return result, nil
}

// Send transaction
// returns the transaction hash or an error
func (c *Client) SendTransaction(ctx context.Context, txSource []byte) (string, error) {
	tx, err := types.TransactionDeserialize(txSource)
	if err != nil {
		return "", errors.Wrap(ErrDeserializeTransaction, err.Error())
	}

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		return "", errors.Wrap(ErrSendTransaction, err.Error())
	}

	return txhash, nil
}

// GetTransactionStatus gets the transaction status.
// Returns the transaction status or an error.
func (c *Client) GetTransactionStatus(ctx context.Context, txhash string) (TransactionStatus, error) {
	status, err := c.solana.GetSignatureStatus(ctx, txhash)
	if err != nil {
		return TransactionStatusUnknown, errors.Wrap(ErrGetTransactionStatus, err.Error())
	}

	if status.Err != nil {
		return TransactionStatusFailure, nil
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
