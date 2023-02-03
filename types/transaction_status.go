package types

import "github.com/portto/solana-go-sdk/rpc"

// TransactionStatus represents the status of a transaction.
type TransactionStatus uint8

// TransactionStatus enum.
const (
	TransactionStatusUnknown TransactionStatus = iota
	TransactionStatusSuccess
	TransactionStatusInProgress
	TransactionStatusFailure
)

// TransactionStatusStrings is a map of TransactionStatus to string.
var transactionStatusStrings = map[TransactionStatus]string{
	TransactionStatusUnknown:    "unknown",
	TransactionStatusSuccess:    "success",
	TransactionStatusInProgress: "in_progress",
	TransactionStatusFailure:    "failure",
}

// String returns the string representation of the transaction status.
func (s TransactionStatus) String() string {
	return transactionStatusStrings[s]
}

// ParseTransactionStatus parses the transaction status from the given string.
func ParseTransactionStatus(s rpc.Commitment) TransactionStatus {
	switch s {
	case rpc.CommitmentFinalized:
		return TransactionStatusSuccess
	case rpc.CommitmentConfirmed, rpc.CommitmentProcessed:
		return TransactionStatusInProgress
	default:
		return TransactionStatusUnknown
	}
}
