package utils

import (
	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/types"
)

// EncodeTransaction returns a base64 encoded transaction.
func EncodeTransaction(tx types.Transaction) (string, error) {
	txb, err := tx.Serialize()
	if err != nil {
		return "", errors.Wrap(err, "failed to build transaction: serialize")
	}

	return BytesToBase64(txb), nil
}

// DecodeTransaction returns a transaction from a base64 encoded transaction.
func DecodeTransaction(base64Tx string) (types.Transaction, error) {
	txb, err := Base64ToBytes(base64Tx)
	if err != nil {
		return types.Transaction{}, errors.Wrap(err, "failed to deserialize transaction: base64 to bytes")
	}

	tx, err := types.TransactionDeserialize(txb)
	if err != nil {
		return types.Transaction{}, errors.Wrap(err, "failed to deserialize transaction: deserialize")
	}

	return tx, nil
}
