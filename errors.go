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
)
