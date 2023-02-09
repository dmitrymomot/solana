package transaction

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/client"
	"github.com/solplaydev/solana/instructions"
	"github.com/solplaydev/solana/token_metadata"
)

type (
	// TransactionBuilder is a builder for transactions.
	TransactionBuilder struct {
		client           solanaClient                   // solana client wrapper
		feePayer         *common.PublicKey              // transaction fee payer
		signers          []types.Account                // additional transaction signers
		instructions     []instructions.InstructionFunc // transaction instructions
		isDurrableTx     bool                           // is durable transaction
		durableNonce     *common.PublicKey              // durable nonce account
		durableNonceAuth *common.PublicKey              // durable nonce auth account
	}

	// solanaClient is a wrapper for the solana client.
	solanaClient interface {
		DefaultDecimals() uint8
		GetMinimumBalanceForRentExemption(ctx context.Context, size uint64) (uint64, error)
		GetTokenAccountInfo(ctx context.Context, base58AtaAddr string) (token.TokenAccount, error)
		GetTokenMetadata(ctx context.Context, base58MintAddr string) (*token_metadata.Metadata, error)
		GetMasterEditionSupply(ctx context.Context, masterMint common.PublicKey) (current, max uint64, err error)
		GetEditionInfo(ctx context.Context, base58MintAddr string) (*token_metadata.Edition, error)
		NewTransaction(ctx context.Context, params client.NewTransactionParams) (string, error)
		NewDurableTransaction(ctx context.Context, params client.NewDurableTransactionParams) (string, error)
	}
)

// NewTransactionBuilder creates a new transaction builder.
func NewTransactionBuilder(c solanaClient) *TransactionBuilder {
	return &TransactionBuilder{client: c}
}

// SetDurableNonce sets the transaction as durable via nonce account.
func (tb *TransactionBuilder) SetDurableNonce(nonce, nonceAuth common.PublicKey) *TransactionBuilder {
	tb.isDurrableTx = true
	tb.durableNonce = &nonce
	tb.durableNonceAuth = &nonceAuth
	return tb
}

// AddInstruction adds an instruction to the transaction.
func (tb *TransactionBuilder) AddInstruction(instruction instructions.InstructionFunc) *TransactionBuilder {
	tb.instructions = append(tb.instructions, instruction)
	return tb
}

// SetFeePayer sets the transaction fee payer.
func (tb *TransactionBuilder) SetFeePayer(feePayer common.PublicKey) *TransactionBuilder {
	tb.feePayer = &feePayer
	return tb
}

// AddSigner adds a signer to the transaction.
func (tb *TransactionBuilder) AddSigner(signer types.Account) *TransactionBuilder {
	tb.signers = append(tb.signers, signer)
	return tb
}

// Build builds the transaction.
// Returns the base64 encoded transaction or an error.
func (tb *TransactionBuilder) Build(ctx context.Context) (string, error) {
	instructions := make([]types.Instruction, 0, len(tb.instructions))
	for _, instruction := range tb.instructions {
		subInstructions, err := instruction(ctx, tb.client)
		if err != nil {
			return "", fmt.Errorf("failed to build transaction: %w", err)
		}
		if len(subInstructions) > 0 {
			instructions = append(instructions, subInstructions...)
		}
	}

	if tb.isDurrableTx {
		if tb.durableNonce == nil || *tb.durableNonce == (common.PublicKey{}) {
			return "", fmt.Errorf("failed to build transaction: missing or invalid durable nonce public key")
		}
		if tb.durableNonceAuth == nil || *tb.durableNonceAuth == (common.PublicKey{}) {
			return "", fmt.Errorf("failed to build transaction: missing or invalid durable nonce auth public key")
		}
		if tb.feePayer == nil || *tb.feePayer == (common.PublicKey{}) {
			tb.feePayer = tb.durableNonceAuth
		}

		return tb.client.NewDurableTransaction(ctx, client.NewDurableTransactionParams{
			FeePayer:     tb.feePayer,
			Instructions: instructions,
			Signers:      tb.signers,
			DurableNonce: *tb.durableNonce,
			NonceAuth:    *tb.durableNonceAuth,
		})
	}

	if tb.feePayer == nil || *tb.feePayer == (common.PublicKey{}) {
		return "", fmt.Errorf("failed to build transaction: missing or invalid fee payer public key")
	}

	return tb.client.NewTransaction(ctx, client.NewTransactionParams{
		FeePayer:     *tb.feePayer,
		Instructions: instructions,
		Signers:      tb.signers,
	})
}
