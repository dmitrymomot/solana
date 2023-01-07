package solana

import "errors"

// Predefined package errors
var (
	ErrCreateBip39Entropy          = errors.New("failed to create bip39 entropy")
	ErrCreateBip39Mnemonic         = errors.New("failed to create bip39 mnemonic")
	ErrCreateBip39SeedFromMnemonic = errors.New("failed to create bip39 seed from mnemonic")
	ErrCreateAccountFromSeed       = errors.New("failed to create account from seed")
	ErrDecodeBase58ToAccount       = errors.New("failed to decode base58 to account")
	ErrDeriveKeyFromSeed           = errors.New("failed to derive key from seed")
	ErrInvalidPublicKey            = errors.New("invalid base58 public key")
	ErrDeserializeTransaction      = errors.New("failed to deserialize transaction")
	ErrGetTransactionFee           = errors.New("failed to get transaction fee")
	ErrSendTransaction             = errors.New("failed to send transaction")
	ErrSerializeMessage            = errors.New("failed to serialize message")
	ErrAddSignature                = errors.New("failed to add signature")
	ErrSerializeTransaction        = errors.New("failed to serialize transaction")
	ErrInvalidAirdropAmount        = errors.New("invalid airdrop amount; must be greater than 0 and less or equal 2000000000")
	ErrRequestAirdrop              = errors.New("failed to request airdrop")
)
