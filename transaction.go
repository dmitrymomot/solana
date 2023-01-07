package solana

import (
	"context"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/types"
)

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
