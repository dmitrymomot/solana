package common

import "errors"

// Predefiend errors
var (
	ErrCreateBip39Entropy                  = errors.New("failed to create bip39 entropy")
	ErrCreateBip39Mnemonic                 = errors.New("failed to create bip39 mnemonic")
	ErrCreateBip39SeedFromMnemonic         = errors.New("failed to create bip39 seed from mnemonic")
	ErrCreateAccountFromSeed               = errors.New("failed to create account from seed")
	ErrDecodeBase58ToAccount               = errors.New("failed to decode base58 to account")
	ErrDeriveKeyFromSeed                   = errors.New("failed to derive key from seed")
	ErrInvalidPublicKey                    = errors.New("invalid base58 public key")
	ErrInvalidPublicKeyLength              = errors.New("invalid public key length")
	ErrNewMnemonic                         = errors.New("failed to create new mnemonic")
	ErrDeriveAccountFromMnemonicBip44      = errors.New("failed to derive account from mnemonic bip44")
	ErrDeriveAccountsListFromMnemonicBip44 = errors.New("failed to derive accounts list from mnemonic bip44")
	ErrDeriveAccountFromMnemonicBip39      = errors.New("failed to derive account from mnemonic bip39")
	ErrDeriveTokenAccount                  = errors.New("failed to derive associated token account")
	ErrInvalidWalletAddress                = errors.New("invalid wallet address: must be a base58 encoded public key")
)
