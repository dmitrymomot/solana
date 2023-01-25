package instructions

import "github.com/portto/solana-go-sdk/types"

// InstructionFunc is a function that returns a list of prepared instructions.
type InstructionFunc func() ([]types.Instruction, error)
