package solana

import (
	"context"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memo"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/utils"
)

// TransferSOLParams are the parameters for TransferSOL function.
type TransferSOLParams struct {
	From   string // base58 encoded source account address and fee payer
	To     string // base58 encoded destination account address
	Amount uint64 // amount of SOL to transfer in lamports
	Memo   string // optional
}

// TransferSOL transfers SOL from the given base58 encoded source account to the given base58 encoded destination account.
// Returns the base64 encoded transaction blob or an error.
func (c *Client) TransferSOL(ctx context.Context, params TransferSOLParams) (string, error) {
	if err := ValidateSolanaWalletAddr(params.From); err != nil {
		return "", utils.StackErrors(ErrTransferSOL, err)
	}
	if err := ValidateSolanaWalletAddr(params.To); err != nil {
		return "", utils.StackErrors(ErrTransferSOL, err)
	}
	if params.Amount == 0 {
		return "", utils.StackErrors(ErrTransferSOL, ErrInvalidTransferAmount)
	}

	fromPubKey := common.PublicKeyFromString(params.From)
	toPubKey := common.PublicKeyFromString(params.To)

	instr := []types.Instruction{
		system.Transfer(system.TransferParam{
			From:   fromPubKey,
			To:     toPubKey,
			Amount: params.Amount,
		}),
	}

	if params.Memo != "" {
		instr = append(instr, memo.BuildMemo(memo.BuildMemoParam{
			SignerPubkeys: []common.PublicKey{fromPubKey},
			Memo:          []byte(params.Memo),
		}))
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     params.From,
		Instructions: instr,
	})
	if err != nil {
		return "", utils.StackErrors(ErrTransferSOL, err)
	}

	return txb, nil
}

// TransferSPLTokenParams are the parameters for TransferSPLToken function.
type TransferSPLTokenParams struct {
	FeePayer string // base58 encoded fee payer address
	From     string // base58 encoded source account address and fee payer
	To       string // base58 encoded destination account address
	Mint     string // base58 encoded mint address
	Amount   uint64 // amount of SPL tokens to transfer
	Memo     string // optional
}

// TransferSPLToken transfers SPL tokens from the given base58 encoded source account to the given base58 encoded destination account.
// Returns the base64 encoded transaction blob or an error.
func (c *Client) TransferSPLToken(ctx context.Context, params TransferSPLTokenParams) (string, error) {
	if err := ValidateSolanaWalletAddr(params.From); err != nil {
		return "", utils.StackErrors(ErrTransferSPLToken, err)
	}
	if err := ValidateSolanaWalletAddr(params.To); err != nil {
		return "", utils.StackErrors(ErrTransferSPLToken, err)
	}
	if params.Amount == 0 {
		return "", utils.StackErrors(ErrTransferSPLToken, ErrInvalidTransferAmount)
	}

	if params.FeePayer == "" {
		params.FeePayer = params.From
	} else {
		if err := ValidateSolanaWalletAddr(params.FeePayer); err != nil {
			return "", utils.StackErrors(ErrTransferSPLToken, err)
		}
	}

	feePayer := common.PublicKeyFromString(params.FeePayer)

	fromAtaPubKey, err := DeriveTokenAccount(params.From, params.Mint)
	if err != nil {
		return "", utils.StackErrors(ErrTransferSPLToken, err)
	}

	toAtaPubKey, err := DeriveTokenAccount(params.To, params.Mint)
	if err != nil {
		return "", utils.StackErrors(ErrTransferSPLToken, err)
	}

	instr := []types.Instruction{
		token.Transfer(token.TransferParam{
			From:   fromAtaPubKey,
			To:     toAtaPubKey,
			Auth:   common.PublicKeyFromString(params.FeePayer),
			Amount: params.Amount,
		}),
	}

	if params.Memo != "" {
		instr = append(instr, memo.BuildMemo(memo.BuildMemoParam{
			SignerPubkeys: []common.PublicKey{feePayer},
			Memo:          []byte(params.Memo),
		}))
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     params.FeePayer,
		Instructions: instr,
	})
	if err != nil {
		return "", utils.StackErrors(ErrTransferSPLToken, err)
	}

	return txb, nil
}
