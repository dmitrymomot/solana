package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/types"
)

// TransferSOLParams defines the parameters for transferring SOL.
type TransferSOLParams struct {
	Sender    common.PublicKey  // required; The wallet to send SOL from
	Recipient common.PublicKey  // required; The wallet to send SOL to
	Amount    uint64            // required; The amount of SOL to send (in lamports)
	Reference *common.PublicKey // optional; public key to use as a reference for the transaction.
}

// TransferSOL transfers SOL from one wallet to another.
// Note: This function does not check if the sender has enough SOL to send. It is the responsibility
// of the caller to check this.
// Amount must be greater than minimum account rent exemption (0.0025 SOL).
func TransferSOL(params TransferSOLParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if params.Sender.ToBase58() == params.Recipient.ToBase58() {
			return nil, fmt.Errorf("sender and recipient must be different")
		}

		if params.Amount <= 0 {
			return nil, fmt.Errorf("amount must be greater than 0")
		}

		instruction := system.Transfer(system.TransferParam{
			From:   params.Sender,
			To:     params.Recipient,
			Amount: params.Amount,
		})

		if params.Reference != nil && *params.Reference != (common.PublicKey{}) {
			instruction.Accounts = append(instruction.Accounts, types.AccountMeta{
				PubKey:     *params.Reference,
				IsSigner:   false,
				IsWritable: false,
			})
		}

		return []types.Instruction{instruction}, nil
	}
}
