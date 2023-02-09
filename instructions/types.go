package instructions

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
)

type (
	// InstructionFunc is a function that returns a list of prepared instructions.
	InstructionFunc func(ctx context.Context, c Client) ([]types.Instruction, error)

	// Creator is the creator of the token metadata.
	Creator struct {
		Address common.PublicKey // required; The creator public key
		Share   uint8            // required; The share of the creator
	}

	// Client is the interface that wraps the basic methods of the client.
	Client interface {
		DefaultDecimals() uint8
		GetMinimumBalanceForRentExemption(ctx context.Context, size uint64) (uint64, error)
		GetTokenAccountInfo(ctx context.Context, base58AtaAddr string) (token.TokenAccount, error)
		GetTokenMetadata(ctx context.Context, base58MintAddr string) (*token_metadata.Metadata, error)
		GetMasterEditionSupply(ctx context.Context, masterMint common.PublicKey) (current, max uint64, err error)
		GetEditionInfo(ctx context.Context, base58MintAddr string) (*token_metadata.Edition, error)
	}
)
