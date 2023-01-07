package solana

import (
	"errors"

	"github.com/portto/solana-go-sdk/rpc"
)

// Predefined package errors
var (
	ErrCreateBip39Entropy                = errors.New("failed to create bip39 entropy")
	ErrCreateBip39Mnemonic               = errors.New("failed to create bip39 mnemonic")
	ErrCreateBip39SeedFromMnemonic       = errors.New("failed to create bip39 seed from mnemonic")
	ErrCreateAccountFromSeed             = errors.New("failed to create account from seed")
	ErrDecodeBase58ToAccount             = errors.New("failed to decode base58 to account")
	ErrDeriveKeyFromSeed                 = errors.New("failed to derive key from seed")
	ErrInvalidPublicKey                  = errors.New("invalid base58 public key")
	ErrDeserializeTransaction            = errors.New("failed to deserialize transaction")
	ErrGetTransactionFee                 = errors.New("failed to get transaction fee")
	ErrSendTransaction                   = errors.New("failed to send transaction")
	ErrSerializeMessage                  = errors.New("failed to serialize message")
	ErrAddSignature                      = errors.New("failed to add signature")
	ErrSerializeTransaction              = errors.New("failed to serialize transaction")
	ErrInvalidAirdropAmount              = errors.New("invalid airdrop amount; must be greater than 0 and less or equal 2000000000")
	ErrRequestAirdrop                    = errors.New("failed to request airdrop")
	ErrGetSolBalance                     = errors.New("failed to get SOL balance")
	ErrFindAssociatedTokenAddress        = errors.New("failed to find associated token address")
	ErrGetSplTokenBalance                = errors.New("failed to get SPL token balance")
	ErrInvalidTransferAmount             = errors.New("invalid transfer amount; must be greater than 0")
	ErrGetLatestBlockhash                = errors.New("failed to get latest blockhash")
	ErrNewTransaction                    = errors.New("failed to create new transaction")
	ErrGetNonceFromNonceAccount          = errors.New("failed to get nonce from nonce account")
	ErrGetMinimumBalanceForRentExemption = errors.New("failed to get minimum balance for rent exemption")
	ErrGetTransactionStatus              = errors.New("failed to get transaction status")
)

// UnwrapJsonRpError unwraps the error from rpc.JsonRpcError to a standard error.
func UnwrapJsonRpError(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*rpc.JsonRpcError); ok {
		return errors.New(e.Message)
	}

	return err
}
