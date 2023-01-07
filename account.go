package solana

import (
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/pkg/hdwallet"
	"github.com/portto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
)

// Predefined mnemonic lengths
const (
	MnemonicLength12 MnemonicLength = 128 // 128 bits of entropy
	MnemonicLength24 MnemonicLength = 256 // 256 bits of entropy
)

// Mnemonic length type
type MnemonicLength int

// NewMnemonic generates a new mnemonic phrase
func NewMnemonic(len MnemonicLength) (string, error) {
	entropy, err := bip39.NewEntropy(int(len))
	if err != nil {
		return "", ErrCreateBip39Entropy
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", ErrCreateBip39Mnemonic
	}

	return mnemonic, nil
}

// DeriveFromMnemonicBip44 derives an Solana account from a mnemonic phrase
// Compatible with BIP44 (phantom wallet)
func DeriveFromMnemonicBip44(mnemonic string) (types.Account, error) {
	return deriveFromMnemonicBip44(mnemonic, 0)
}

// DeriveListFromMnemonicBip44 derives a list of Solana accounts from a mnemonic phrase
// Compatible with BIP44 (phantom wallet)
func DeriveListFromMnemonicBip44(mnemonic string, count int) ([]types.Account, error) {
	accounts := make([]types.Account, count)

	for i := 0; i < count; i++ {
		account, err := deriveFromMnemonicBip44(mnemonic, i)
		if err != nil {
			return nil, err
		}

		accounts[i] = account
	}

	return accounts, nil
}

// DeriveFromMnemonicBip39 derives an Solana account from a mnemonic phrase
// Compatible with BIP39 (solana cli tool)
func DeriveFromMnemonicBip39(mnemonic string) (types.Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return types.Account{}, ErrCreateBip39SeedFromMnemonic
	}

	account, err := types.AccountFromSeed(seed[:32])
	if err != nil {
		return types.Account{}, ErrCreateAccountFromSeed
	}

	return account, nil
}

// ToBase58 converts an Solana account to a base58 encoded string
func ToBase58(a types.Account) string {
	return base58.Encode(a.PrivateKey)
}

// FromBase58 creates an Solana account from a base58 encoded string
func FromBase58(s string) (types.Account, error) {
	b, err := base58.Decode(s)
	if err != nil {
		return types.Account{}, ErrDecodeBase58ToAccount
	}

	return types.AccountFromBytes(b)
}

// deriveFromMnemonicBip44 derives an Solana account from a mnemonic phrase
// Compatible with BIP44 (phantom wallet)
func deriveFromMnemonicBip44(mnemonic string, path int) (types.Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return types.Account{}, ErrCreateBip39SeedFromMnemonic
	}

	derivedKey, err := hdwallet.Derived(fmt.Sprintf("m/44'/501'/%d'/0'", path), seed)
	if err != nil {
		return types.Account{}, ErrDeriveKeyFromSeed
	}

	account, err := types.AccountFromSeed(derivedKey.PrivateKey)
	if err != nil {
		return types.Account{}, ErrCreateAccountFromSeed
	}

	return account, nil
}
