package instructions

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
)

// TransferTokenParam defines the parameters for transferring tokens.
type TransferTokenParam struct {
	Mint   common.PublicKey // required; The token mint to send
	Amount uint64           // required; The amount of tokens to send (in token minimal units)

	Sender    *common.PublicKey // required if SenderAta is empty; The wallet to send tokens from
	Recipient *common.PublicKey // required if RecipientAta is empty; The wallet to send tokens to

	FeePayer     *common.PublicKey // optional; The account to pay for the transaction; defaults to the sender
	SenderAta    *common.PublicKey // optional; The associated token account of the sender; if nil, it will be fetched
	RecipientAta *common.PublicKey // optional; The associated token account of the recipient; if nil, it will be fetched
}

// TransferToken transfers tokens from one wallet to another.
// Note: This function does not check if the sender has enough tokens to send. It is the responsibility
// of the caller to check this.
// FeePayer must be provided if Sender is not set.
func TransferToken(params TransferTokenParam) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if params.Sender == nil && params.SenderAta == nil {
			return nil, fmt.Errorf("sender or senderAta must be provided")
		}

		if params.Recipient == nil && params.RecipientAta == nil {
			return nil, fmt.Errorf("recipient or recipientAta must be provided")
		}

		if params.SenderAta == nil {
			senderAta, _, err := common.FindAssociatedTokenAddress(*params.Sender, params.Mint)
			if err != nil {
				return nil, fmt.Errorf("failed to find associated token address for sender wallet: %w", err)
			}
			params.SenderAta = &senderAta
		}

		if params.RecipientAta == nil {
			recipientAta, _, err := common.FindAssociatedTokenAddress(*params.Recipient, params.Mint)
			if err != nil {
				return nil, fmt.Errorf("failed to find associated token address for recipient wallet: %w", err)
			}
			params.RecipientAta = &recipientAta
		}

		if params.FeePayer == nil && params.Sender == nil {
			return nil, fmt.Errorf("feePayer must be provided if sender is not set")
		}

		if params.FeePayer == nil && params.Sender != nil {
			params.FeePayer = params.Sender
		}

		return []types.Instruction{
			token.Transfer(token.TransferParam{
				From:   *params.SenderAta,
				To:     *params.RecipientAta,
				Auth:   *params.FeePayer,
				Amount: params.Amount,
			}),
		}, nil
	}
}
