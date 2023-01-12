package solana

import (
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/rpc"
)

// Predefined Solana account sizes
const (
	AccountSize       uint64 = 165 // 165 bytes
	FeeCalculatorSize uint64 = 8   // 8 bytes
	NonceAccountSize  uint64 = 80  // 80 bytes
	StakeAccountSize  uint64 = 200 // 200 bytes
	TokenAccountSize  uint64 = 165 // 165 bytes
	MintAccountSize   uint64 = 82  // 82 bytes
)

// Lookup table sizes
const (
	LookupTableMetaSize     uint64 = 56  // 56 bytes
	LookupTableMaxAddresses uint   = 256 // 256 addresses
)

const (
	// 1 SOL = 1e9 lamports
	SOL uint64 = 1e9

	// SPL token default decimals
	SPLTokenDefaultDecimals uint8 = 9

	// SPL token default multiplier for decimals
	SPLTokenDefaultMultiplier uint64 = 1e9

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

// TokenUseMethod represents the use method of a token.
type TokenUseMethod string

// TokenUseMethod enum.
const (
	TokenUseMethodUnknown TokenUseMethod = "unknown"
	TokenUseMethodBurn    TokenUseMethod = "burn"
	TokenUseMethodSingle  TokenUseMethod = "single"
	TokenUseMethodMulti   TokenUseMethod = "multiple"
)

// String returns the string representation of the token use method.
func (m TokenUseMethod) String() string {
	if !m.Valid() {
		return TokenUseMethodUnknown.String()
	}
	return string(m)
}

// Valid returns true if the token use method is valid.
func (m TokenUseMethod) Valid() bool {
	return m == TokenUseMethodBurn || m == TokenUseMethodSingle || m == TokenUseMethodMulti
}

// ToMetadataUseMethod converts the token use method to token_metadata.UseMethod.
func (m TokenUseMethod) ToMetadataUseMethod() token_metadata.UseMethod {
	return StringToUseMethod(m.String())
}

// Cast token_metadata.UseMethod to string
func UseMethodToString(m token_metadata.UseMethod) string {
	switch m {
	case token_metadata.Burn:
		return TokenUseMethodBurn.String()
	case token_metadata.Single:
		return TokenUseMethodSingle.String()
	case token_metadata.Multiple:
		return TokenUseMethodMulti.String()
	default:
		return TokenUseMethodUnknown.String()
	}
}

// Cast string to token_metadata.UseMethod
func StringToUseMethod(m string) token_metadata.UseMethod {
	switch m {
	case TokenUseMethodBurn.String():
		return token_metadata.Burn
	case TokenUseMethodSingle.String():
		return token_metadata.Single
	case TokenUseMethodMulti.String():
		return token_metadata.Multiple
	default:
		return token_metadata.Burn
	}
}

// Predefined token editions
const (
	EditionUndefined                 string = "undefined"
	EditionMasterEdition             string = "master_edition"
	EditionPrintedEdition            string = "edition"
	EditionReservationList           string = "reservation_list"
	EditionMetadata                  string = "metadata"
	EditionEditionMarker             string = "edition_marker"
	EditionUseAuthorityRecord        string = "use_authority_record"
	EditionCollectionAuthorityRecord string = "collection_authority_record"
)

// token_metadata.EditionKey to string
func EditionKeyToString(k token_metadata.Key) string {
	switch k {
	case token_metadata.KeyUninitialized:
		return EditionUndefined
	case token_metadata.KeyEditionV1:
		return EditionPrintedEdition
	case token_metadata.KeyMasterEditionV1, token_metadata.KeyMasterEditionV2:
		return EditionMasterEdition
	case token_metadata.KeyReservationListV1:
		return EditionReservationList
	case token_metadata.KeyMetadataV1:
		return EditionMetadata
	case token_metadata.KeyReservationListV2:
		return EditionReservationList
	case token_metadata.KeyEditionMarker:
		return EditionEditionMarker
	case token_metadata.KeyUseAuthorityRecord:
		return EditionUseAuthorityRecord
	case token_metadata.KeyCollectionAuthorityRecord:
		return EditionCollectionAuthorityRecord
	default:
		return EditionUndefined
	}
}

// Token standards enum
const (
	TokenStandardUndefined          string = "undefined"
	TokenStandardNonFungible        string = "non_fungible"
	TokenStandardNonFungibleEdition string = "non_fungible_edition"
	TokenStandardFungibleAsset      string = "fungible_asset"
	TokenStandardFungible           string = "fungible"
)

// Cast token_metadata.TokenStandard to string
func TokenStandardToString(ts token_metadata.TokenStandard) string {
	switch ts {
	case token_metadata.NonFungible:
		return TokenStandardNonFungible
	case token_metadata.NonFungibleEdition:
		return TokenStandardNonFungibleEdition
	case token_metadata.FungibleAsset:
		return TokenStandardFungibleAsset
	case token_metadata.Fungible:
		return TokenStandardFungible
	default:
		return TokenStandardUndefined
	}
}
