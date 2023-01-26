package instructions

import (
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

type (
	// InstructionFunc is a function that returns a list of prepared instructions.
	InstructionFunc func() ([]types.Instruction, error)

	// Creator is the creator of the token metadata.
	Creator struct {
		Address common.PublicKey // required; The creator public key
		Share   uint8            // required; The share of the creator
	}
)
