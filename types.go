package solana

import "github.com/portto/solana-go-sdk/rpc"

const (
	// 1 SOL = 1e9 lamports
	SOL uint64 = 1e9

	// SPL token default decimals
	SPLTokenDefaultDecimals uint8 = 9

	// Solana Devnet RPC URL
	SolanaDevnetRPCURL = "https://api.devnet.solana.com"

	// Solana Mainnet RPC URL
	SolanaMainnetRPCURL = "https://api.mainnet-beta.solana.com"

	// Solana Testnet RPC URL
	SolanaTestnetRPCURL = "https://api.testnet.solana.com"
)

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
