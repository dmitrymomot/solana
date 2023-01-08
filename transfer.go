package solana

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memo"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/utils"
)

// TransferSOLParams are the parameters for TransferSOL function.
type TransferSOLParams struct {
	Base58SourceAddr string // base58 encoded source account address and fee payer
	Base58DestAddr   string // base58 encoded destination account address
	Lamports         uint64 // amount of SOL to transfer in lamports
	Memo             string // optional
}

// TransferSOL transfers SOL from the given base58 encoded source account to the given base58 encoded destination account.
// Returns the base64 encoded transaction blob or an error.
func (c *Client) TransferSOL(ctx context.Context, params TransferSOLParams) ([]byte, error) {
	if err := ValidateSolanaWalletAddr(params.Base58SourceAddr); err != nil {
		return nil, utils.StackErrors(ErrTransferSOL, err)
	}
	if err := ValidateSolanaWalletAddr(params.Base58DestAddr); err != nil {
		return nil, utils.StackErrors(ErrTransferSOL, err)
	}
	if params.Lamports == 0 {
		return nil, utils.StackErrors(ErrTransferSOL, ErrInvalidTransferAmount)
	}

	fromPubKey := common.PublicKeyFromString(params.Base58SourceAddr)
	toPubKey := common.PublicKeyFromString(params.Base58DestAddr)

	instr := []types.Instruction{
		system.Transfer(system.TransferParam{
			From:   fromPubKey,
			To:     toPubKey,
			Amount: params.Lamports,
		}),
	}

	if params.Memo != "" {
		instr = append(instr, memo.BuildMemo(memo.BuildMemoParam{
			SignerPubkeys: []common.PublicKey{fromPubKey},
			Memo:          []byte(params.Memo),
		}))
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		Base58FeePayerAddr: params.Base58SourceAddr,
		Instructions:       instr,
	})
	if err != nil {
		return nil, utils.StackErrors(ErrTransferSOL, err)
	}

	return txb, nil
}
