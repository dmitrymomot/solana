package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
)

// TransferTokenParam defines the parameters for transferring tokens.
type TransferTokenParam struct {
	Sender    common.PublicKey // required if SenderAta is empty; The wallet to send tokens from
	Recipient common.PublicKey // required if RecipientAta is empty; The wallet to send tokens to
	Mint      common.PublicKey // required; The token mint to send
	Amount    uint64           // required; The amount of tokens to send (in token minimal units)
}

// Validate validates the parameters.
func (p TransferTokenParam) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("missed or invalid mint public key")
	}
	if p.Amount == 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if p.Sender == (common.PublicKey{}) {
		return fmt.Errorf("missed or invalid sender public key")
	}
	if p.Recipient == (common.PublicKey{}) {
		return fmt.Errorf("missed or invalid recipient public key")
	}
	return nil
}

// TransferToken transfers tokens from one wallet to another.
// Note: This function does not check if the sender has enough tokens to send. It is the responsibility
// of the caller to check this.
// FeePayer must be provided if Sender is not set.
func TransferToken(params TransferTokenParam) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("invalid given data: %w", err)
		}

		senderAta, _, err := common.FindAssociatedTokenAddress(params.Sender, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address for sender wallet: %w", err)
		}

		recipientAta, _, err := common.FindAssociatedTokenAddress(params.Recipient, params.Mint)
		if err != nil {
			return nil, fmt.Errorf("failed to find associated token address for recipient wallet: %w", err)
		}

		return []types.Instruction{
			token.Transfer(token.TransferParam{
				From:   senderAta,
				To:     recipientAta,
				Auth:   params.Sender,
				Amount: params.Amount,
			}),
		}, nil
	}
}
