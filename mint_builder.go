package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// MintBuilder is a builder to build a mint
type MintBuilder struct {
	client solanaClient

	feePayer common.PublicKey
	owner    common.PublicKey

	mint *types.Account

	masterEditionMint common.PublicKey
	edition           uint64

	tokenStandard *token_metadata.TokenStandard

	decimals         uint8
	supplyAmount     uint64
	fixedSupply      bool
	maxEditionSupply uint64

	metadataPubkey      common.PublicKey
	metadataInstruction *types.Instruction
}

// Solana client interface
type solanaClient interface {
	GetMinimumBalanceForRentExemption(ctx context.Context, size uint64) (uint64, error)
	GetMasterEditionInfo(ctx context.Context, base58MintAddr string) (*token_metadata.Edition, error)
	NewTransaction(ctx context.Context, params NewTransactionParams) (string, error)
}

// NewMintBuilder creates a new MintBuilder
func NewMintBuilder(client solanaClient) *MintBuilder {
	return &MintBuilder{client: client}
}

// SetMint sets the mint account
func (b *MintBuilder) SetMint(mint types.Account) *MintBuilder {
	b.mint = utils.Pointer(mint)
	return b
}

// SetFeePayer sets the fee payer of the mint
func (b *MintBuilder) SetFeePayer(feePayer common.PublicKey) *MintBuilder {
	b.feePayer = feePayer
	return b
}

// SetFeePayerBase58 sets the fee payer of the mint
func (b *MintBuilder) SetFeePayerBase58(feePayer string) *MintBuilder {
	b.feePayer = common.PublicKeyFromString(feePayer)
	return b
}

// SetOwner sets the owner of the mint
func (b *MintBuilder) SetOwner(owner common.PublicKey) *MintBuilder {
	b.owner = owner
	return b
}

// SetOwnerBase58 sets the owner of the mint
func (b *MintBuilder) SetOwnerBase58(owner string) *MintBuilder {
	b.owner = common.PublicKeyFromString(owner)
	return b
}

// SetMasterEditionMint sets the master edition mint of the mint
func (b *MintBuilder) SetMasterEditionMint(masterEditionMint common.PublicKey) *MintBuilder {
	b.masterEditionMint = masterEditionMint
	return b
}

// SetEditionNumber sets the edition number
func (b *MintBuilder) SetEditionNumber(edition uint64) *MintBuilder {
	b.edition = edition
	return b
}

// SetMasterEditionMintBase58 sets the master edition mint of the mint
func (b *MintBuilder) SetMasterEditionMintBase58(masterEditionMint string) *MintBuilder {
	b.masterEditionMint = common.PublicKeyFromString(masterEditionMint)
	return b
}

// SetTokenStandard sets the token standard of the mint
func (b *MintBuilder) SetTokenStandard(tokenStandard token_metadata.TokenStandard) *MintBuilder {
	b.tokenStandard = &tokenStandard
	return b
}

// SetDecimals sets the decimals of the mint
func (b *MintBuilder) SetDecimals(decimals uint8) *MintBuilder {
	if decimals > 9 {
		decimals = 9
	}
	b.decimals = decimals
	return b
}

// SetSupplyAmount sets the supply amount of the mint
func (b *MintBuilder) SetSupplyAmount(supplyAmount uint64) *MintBuilder {
	b.supplyAmount = supplyAmount
	return b
}

// SetFixedSupply sets the fixed supply of the mint
func (b *MintBuilder) SetFixedSupply(fixedSupply bool) *MintBuilder {
	b.fixedSupply = fixedSupply
	return b
}

// SetMaxEditionSupply sets the max edition supply of the mint
func (b *MintBuilder) SetMaxEditionSupply(maxEditionSupply uint64) *MintBuilder {
	b.maxEditionSupply = maxEditionSupply
	return b
}

// SetMetadataPubkey sets the metadata pubkey of the mint
func (b *MintBuilder) SetMetadataPubkey(metadataPubkey common.PublicKey) *MintBuilder {
	b.metadataPubkey = metadataPubkey
	return b
}

// SetMetadataPubkeyBase58 sets the metadata pubkey of the mint
func (b *MintBuilder) SetMetadataPubkeyBase58(metadataPubkey string) *MintBuilder {
	return b.SetMetadataPubkey(common.PublicKeyFromString(metadataPubkey))
}

// SetMetadataInstruction sets the metadata instruction of the mint
func (b *MintBuilder) SetMetadataInstruction(metadataInstruction *types.Instruction) *MintBuilder {
	b.metadataInstruction = metadataInstruction
	return b
}

// Build builds a mint transaction ready to be signed and sent to the network.
//
// If the owner is not set, the fee payer will be used as the owner. If the fee payer
// is not set, the owner will be used as the fee payer.
//
// Returns:
//
//	mintPubkey 	- is the encoded in base58 public key of the mint account that will be created.
//	tx 			- is the encoded in base64 transaction that will be signed and sent to the network.
//	err 		- is the error if any.
func (b *MintBuilder) Build() (mintPubkey, tx string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if b.mint == nil {
		mint := NewAccount()
		b.mint = &mint
	}

	if b.feePayer == (common.PublicKey{}) && b.owner == (common.PublicKey{}) {
		return mintPubkey, tx, errors.New("one of fee payer or owner is required")
	}

	if b.feePayer == (common.PublicKey{}) && b.owner != (common.PublicKey{}) {
		b.feePayer = b.owner
	}

	if b.owner == (common.PublicKey{}) && b.feePayer != (common.PublicKey{}) {
		b.owner = b.feePayer
	}

	if b.tokenStandard == nil {
		return mintPubkey, tx, errors.New("token standard is required")
	}

	switch *b.tokenStandard {
	case token_metadata.TokenStandardFungible, token_metadata.TokenStandardFungibleAsset:
		return b.buildMintFungible(ctx, *b.mint)
	case token_metadata.TokenStandardNonFungible:
		return b.buildMintNonFungible(ctx, *b.mint)
	case token_metadata.TokenStandardNonFungibleEdition:
		return b.buildMintNonFungibleEdition(ctx, *b.mint)
	default:
		return mintPubkey, tx, errors.New("invalid token standard")
	}
}

func (b *MintBuilder) buildMintFungible(ctx context.Context, mint types.Account) (mintPubkey, tx string, err error) {
	if *b.tokenStandard != token_metadata.TokenStandardFungible &&
		*b.tokenStandard != token_metadata.TokenStandardFungibleAsset {
		return mintPubkey, tx, errors.New("token standard is not fungible")
	}

	rentExemptionBalance, err := b.client.GetMinimumBalanceForRentExemption(ctx, MintAccountSize)
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to get minimum balance for rent exemption: %w", err)
	}

	if b.metadataInstruction == nil {
		return mintPubkey, tx, errors.New("metadata instruction is required for fungible and non-fungible tokens")
	}

	if *b.tokenStandard != token_metadata.TokenStandardFungibleAsset {
		b.decimals = 0 // decimals are not supported for fungible tokens
	}

	if b.decimals > 9 {
		b.decimals = 9 // decimals are limited to 9
	}

	instructions := []types.Instruction{
		system.CreateAccount(system.CreateAccountParam{
			From:     b.feePayer,
			New:      mint.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalance,
			Space:    token.MintAccountSize,
		}),
		token.InitializeMint(token.InitializeMintParam{
			Decimals: b.decimals,
			Mint:     mint.PublicKey,
			MintAuth: b.owner,
		}),
		*b.metadataInstruction,
	}

	if b.supplyAmount > 0 {
		ownerAta, _, err := common.FindAssociatedTokenAddress(b.owner, mint.PublicKey)
		if err != nil {
			return mintPubkey, tx, fmt.Errorf("failed to find associated token address for mint %s: %w", mint.PublicKey.ToBase58(), err)
		}

		instructions = append(
			instructions,
			associated_token_account.CreateAssociatedTokenAccount(
				associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 b.feePayer,
					Owner:                  b.owner,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ownerAta,
				},
			),
			token.MintTo(token.MintToParam{
				Mint:    mint.PublicKey,
				To:      ownerAta,
				Auth:    b.owner,
				Signers: []common.PublicKey{},
				Amount:  b.supplyAmount,
			}),
		)
	}

	if b.fixedSupply && b.supplyAmount > 0 {
		instructions = append(instructions, token.SetAuthority(token.SetAuthorityParam{
			Account:  mint.PublicKey,
			AuthType: token.AuthorityTypeMintTokens,
			Auth:     b.owner,
			NewAuth:  nil,
			Signers:  []common.PublicKey{},
		}))
	}

	txb, err := b.client.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     b.feePayer.ToBase58(),
		Instructions: instructions,
		Signers:      []types.Account{mint},
	})
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to create transaction: %w", err)
	}

	return mint.PublicKey.ToBase58(), txb, nil
}

func (b *MintBuilder) buildMintNonFungible(ctx context.Context, mint types.Account) (mintPubkey, tx string, err error) {
	if *b.tokenStandard != token_metadata.TokenStandardNonFungible {
		return mintPubkey, tx, errors.New("token standard is not NFT")
	}

	rentExemptionBalance, err := b.client.GetMinimumBalanceForRentExemption(ctx, MintAccountSize)
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to get minimum balance for rent exemption: %w", err)
	}

	if b.metadataInstruction == nil {
		return mintPubkey, tx, errors.New("metadata instruction is required for fungible and non-fungible tokens")
	}

	b.decimals = 0     // decimals must be 0 for nft
	b.supplyAmount = 1 // supply amount must be 1 for nft

	ownerAta, _, err := common.FindAssociatedTokenAddress(b.owner, mint.PublicKey)
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to find associated token address for mint %s: %w", mint.PublicKey.ToBase58(), err)
	}

	tokenMasterEditionPubkey, err := token_metadata.DeriveEditionPubkey(mint.PublicKey)
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to derive master edition pubkey: %w", err)
	}

	instructions := []types.Instruction{
		system.CreateAccount(system.CreateAccountParam{
			From:     b.feePayer,
			New:      mint.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalance,
			Space:    token.MintAccountSize,
		}),
		token.InitializeMint(token.InitializeMintParam{
			Decimals:   b.decimals,
			Mint:       mint.PublicKey,
			MintAuth:   b.owner,
			FreezeAuth: utils.Pointer(b.owner),
		}),
		*b.metadataInstruction,
		associated_token_account.CreateAssociatedTokenAccount(
			associated_token_account.CreateAssociatedTokenAccountParam{
				Funder:                 b.feePayer,
				Owner:                  b.owner,
				Mint:                   mint.PublicKey,
				AssociatedTokenAccount: ownerAta,
			},
		),
		token.MintTo(token.MintToParam{
			Mint:    mint.PublicKey,
			To:      ownerAta,
			Auth:    b.owner,
			Signers: []common.PublicKey{},
			Amount:  b.supplyAmount,
		}),
		metaplex_token_metadata.CreateMasterEditionV3(
			metaplex_token_metadata.CreateMasterEditionParam{
				Edition:         tokenMasterEditionPubkey,
				Mint:            mint.PublicKey,
				UpdateAuthority: b.owner,
				MintAuthority:   b.owner,
				Metadata:        b.metadataPubkey,
				Payer:           b.feePayer,
				MaxSupply:       utils.Pointer(b.maxEditionSupply),
			},
		),
	}

	txb, err := b.client.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     b.feePayer.ToBase58(),
		Instructions: instructions,
		Signers:      []types.Account{mint},
	})
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to create transaction to mint NFT: %w", err)
	}

	return mint.PublicKey.ToBase58(), txb, nil
}

func (b *MintBuilder) buildMintNonFungibleEdition(ctx context.Context, newMint types.Account) (mintPubkey, tx string, err error) {
	if *b.tokenStandard != token_metadata.TokenStandardNonFungibleEdition {
		return mintPubkey, tx, errors.New("token standard is not NFT edition")
	}

	if b.masterEditionMint == (common.PublicKey{}) {
		return mintPubkey, tx, errors.New("master mint pubkey is required for NFT edition")
	}

	// Get next edition number.
	if b.edition == 0 {
		editionInfo, err := b.client.GetMasterEditionInfo(ctx, b.masterEditionMint.ToBase58())
		if err != nil || editionInfo == nil {
			return mintPubkey, tx, utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				err,
			)
		}
		if editionInfo.Type == "" || editionInfo.Type != EditionMasterEdition {
			return mintPubkey, tx, utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				ErrTokenIsNotMasterEdition,
			)
		}
		if editionInfo.MaxSupply == 0 || editionInfo.Supply == editionInfo.MaxSupply {
			return mintPubkey, tx, utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				ErrMaxSupplyReached,
			)
		}

		b.edition = editionInfo.Supply + 1
	}

	masterOwner := b.owner
	newMintOwner := b.owner

	masterOwnerAta, _, err := common.FindAssociatedTokenAddress(masterOwner, b.masterEditionMint)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrFindAssociatedTokenAddress,
			err,
		)
	}

	masterMetaPublicKey, err := token_metadata.DeriveTokenMetadataPubkey(b.masterEditionMint)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetTokenMetaPubkey,
			err,
		)
	}

	masterEditionPublicKey, err := token_metadata.DeriveEditionPubkey(b.masterEditionMint)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetMasterEdition,
			err,
		)
	}

	newMintOwnerAta, _, err := common.FindAssociatedTokenAddress(newMintOwner, newMint.PublicKey)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrFindAssociatedTokenAddress,
			err,
		)
	}

	newMintMetaPublicKey, err := token_metadata.DeriveTokenMetadataPubkey(newMint.PublicKey)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetTokenMetaPubkey,
			err,
		)
	}

	newMintEditionPublicKey, err := token_metadata.DeriveEditionPubkey(newMint.PublicKey)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetMasterEdition,
			err,
		)
	}

	newMintEditionMark, err := token_metadata.DeriveEditionMarkerPubkey(b.masterEditionMint, b.edition)
	if err != nil {
		return mintPubkey, tx, utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetEditionMark,
			err,
		)
	}

	rentExemptionBalance, err := b.client.GetMinimumBalanceForRentExemption(ctx, MintAccountSize)
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to get minimum balance for rent exemption: %w", err)
	}

	instructions := []types.Instruction{
		system.CreateAccount(system.CreateAccountParam{
			From:     b.feePayer,
			New:      newMint.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalance,
			Space:    token.MintAccountSize,
		}),
		token.InitializeMint(token.InitializeMintParam{
			Decimals:   b.decimals,
			Mint:       newMint.PublicKey,
			MintAuth:   b.owner,
			FreezeAuth: utils.Pointer(b.owner),
		}),
		associated_token_account.CreateAssociatedTokenAccount(
			associated_token_account.CreateAssociatedTokenAccountParam{
				Funder:                 b.feePayer,
				Owner:                  b.owner,
				Mint:                   newMint.PublicKey,
				AssociatedTokenAccount: newMintOwnerAta,
			},
		),
		token.MintTo(token.MintToParam{
			Mint:   newMint.PublicKey,
			Auth:   b.owner,
			To:     newMintOwnerAta,
			Amount: 1,
		}),
		metaplex_token_metadata.MintNewEditionFromMasterEditionViaToken(
			metaplex_token_metadata.MintNewEditionFromMasterEditionViaTokeParam{
				NewMetaData:                newMintMetaPublicKey,
				NewEdition:                 newMintEditionPublicKey,
				MasterEdition:              masterEditionPublicKey,
				NewMint:                    newMint.PublicKey,
				NewMintAuthority:           b.owner,
				Payer:                      b.feePayer,
				TokenAccountOwner:          masterOwner,
				TokenAccount:               masterOwnerAta,
				NewMetadataUpdateAuthority: newMintOwner,
				MasterMetadata:             masterMetaPublicKey,

				EditionMark: newMintEditionMark,
				Edition:     b.edition,
			},
		),
	}

	txb, err := b.client.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     b.feePayer.ToBase58(),
		Instructions: instructions,
		Signers:      []types.Account{newMint},
	})
	if err != nil {
		return mintPubkey, tx, fmt.Errorf("failed to create transaction: %w", err)
	}

	return newMint.PublicKey.ToBase58(), txb, nil
}
