package instructions

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memo"
	"github.com/portto/solana-go-sdk/types"
)

// Memo is the memo instruction.
func Memo(str string, signers ...common.PublicKey) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		return []types.Instruction{
			memo.BuildMemo(memo.BuildMemoParam{
				SignerPubkeys: signers,
				Memo:          []byte(str),
			}),
		}, nil
	}
}
