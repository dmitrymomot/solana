package solana

import (
	"context"
	"errors"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// Mint fungible token.
type InitMintFungibleTokenParams struct {
	FeePayer string // required; base58 encoded address of the fee payer
	Owner    string // optional; base58 encoded address of the owner of the token; default is the fee payer

	Decimals     uint8  // optional; number of decimals; default is 9
	SupplyAmount uint64 // optional; amount of tokens to mint in lamports, e.g: if you want to mint 10 tokens and decimals=9, amount=10*1e9/amount=10000000000; default is 0, then no tokens will be minted
	FixedSupply  bool   // optional; disable minting of the token; false by default

	// Metadata
	Name        string // required; name of the token; max 32 characters
	Symbol      string // required; symbol of the token; max 10 characters
	MetadataURI string // optional; URI of the token metadata; can be set later
}

// Validate validates the parameters.
func (params InitMintFungibleTokenParams) Validate() error {
	if params.Name == "" || params.Symbol == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("name and symbol are required"),
		)
	}

	if len(params.Name) > 32 || len(params.Name) < 2 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("name must be between 2 and 32 characters"),
		)
	}

	if len(params.Symbol) > 10 || len(params.Symbol) < 3 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("symbol must be between 3 and 10 characters"),
		)
	}

	if params.Decimals > 9 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("decimals must be between 0 and 9"),
		)
	}

	if params.FeePayer == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("fee payer is required"),
		)
	}

	return nil
}

// InitMintFungibleToken mints a fungible token.
// These are simple SPL tokens with limited metadata and supply >= 0. Examples are USDC, GBTC and RAY.
// The token is minted to the owner's associated token account.
// Returns the mint address and base64 encoded transaction or an error.
func (c *Client) InitMintFungibleToken(ctx context.Context, params InitMintFungibleTokenParams) (mintAddr, tx string, err error) {
	if err := params.Validate(); err != nil {
		return "", "", utils.StackErrors(ErrMintFungibleToken, err)
	}

	if params.Decimals < 1 || params.Decimals > 9 {
		params.Decimals = c.defaultDecimals
	}

	result, err := c.prepareInitMintTransaction(ctx, initMintTransactionParams{
		FeePayer:      params.FeePayer,
		Owner:         params.Owner,
		TokenStandard: utils.Pointer(token_metadata.TokenStandardFungible),
		Decimals:      params.Decimals,
		SupplyAmount:  params.SupplyAmount,
		FixedSupply:   params.FixedSupply,
		Name:          params.Name,
		Symbol:        params.Symbol,
		MetadataURI:   params.MetadataURI,
	})
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintFungibleToken,
			err,
		)
	}

	return result.Mint.PublicKey.ToBase58(), result.Tx, nil
}

// InitMintFungibleAssetParams contains the parameters for minting a semi-fungible token (asset).
type InitMintFungibleAssetParams struct {
	FeePayer string // required; base58 encoded address of the fee payer
	Owner    string // optional; base58 encoded address of the owner of the token; default is the fee payer

	SupplyAmount uint64 // optional; amount of assets to mint; default is 0, then no tokens will be minted
	FixedSupply  bool   // optional; disable minting of the new tokens; false by default

	// Metadata
	Name        string // required; name of the token; max 32 characters
	Symbol      string // required; symbol of the token; max 10 characters
	MetadataURI string // optional; URI of the token metadata; can be set later
	Collection  string // optional; base58 encoded address of the collection; can be set later
}

// Validate validates the parameters.
func (params InitMintFungibleAssetParams) Validate() error {
	if params.Name == "" || params.Symbol == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("name and symbol are required"),
		)
	}

	if len(params.Name) > 32 || len(params.Name) < 2 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("name must be between 2 and 32 characters"),
		)
	}

	if len(params.Symbol) > 10 || len(params.Symbol) < 3 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("symbol must be between 3 and 10 characters"),
		)
	}

	if params.FeePayer == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("fee payer is required"),
		)
	}

	return nil
}

// InitMintFungibleAsset mints a semi-fungible token (asset).
// These are fungible tokens with more extensive metadata and supply >= 0. An example of this kind of token is something the community has been calling "semi-fungible tokens" often used to represent a fungible but attribute-heavy in-game item such as a sword or a piece of wood.
// The token is minted to the owner's associated token account.
// The owner can set the metadata and collection later.
// The owner can also disable minting of the new tokens.
// Returns the mint address and base64 encoded transaction or an error.
func (c *Client) InitMintFungibleAsset(ctx context.Context, params InitMintFungibleAssetParams) (mintAddr, tx string, err error) {
	if err := params.Validate(); err != nil {
		return "", "", utils.StackErrors(ErrMintFungibleAsset, err)
	}

	result, err := c.prepareInitMintTransaction(ctx, initMintTransactionParams{
		FeePayer:      params.FeePayer,
		Owner:         params.Owner,
		TokenStandard: utils.Pointer(token_metadata.TokenStandardFungibleAsset),
		Decimals:      0,
		SupplyAmount:  params.SupplyAmount,
		FixedSupply:   params.FixedSupply,
		Name:          params.Name,
		Symbol:        params.Symbol,
		MetadataURI:   params.MetadataURI,
		Collection:    &token_metadata.Collection{Key: params.Collection},
	})
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintFungibleAsset,
			err,
		)
	}

	return result.Mint.PublicKey.ToBase58(), result.Tx, nil
}

// MintNonFungibleTokenParams contains the parameters for minting a non-fungible token (NFT).
type MintNonFungibleTokenParams struct {
	FeePayer string // required; base58 encoded address of the fee payer
	Owner    string // optional; base58 encoded address of the owner of the token; default is the fee payer

	// Metadata
	Name        string // required; name of the token; max 32 characters
	Symbol      string // required; symbol of the token; max 10 characters
	MetadataURI string // optional; URI of the token metadata; can be set later
	Collection  string // optional; base58 encoded address of the collection; can be set later

	// Minting
	MaxSupply            uint64                   // optional; maximum amount of edition tokens can be minted from master edition; default is 0, then only one token will be minted.
	SellerFeeBasisPoints uint16                   // optional; fee that will be paid to the owner of the master edition when the token is sold; default is 0
	Creators             []token_metadata.Creator // optional; creators of the token; default is fee payer with 100% share; fee payer must be in a creators list; total share must be 100.
	Uses                 *token_metadata.Uses     // optional; uses of the token; default is unlimited
}

// Validate validates the parameters.
func (params MintNonFungibleTokenParams) Validate() error {
	if params.Name == "" || params.Symbol == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("name and symbol are required"),
		)
	}

	if len(params.Name) > 32 || len(params.Name) < 2 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("name must be between 2 and 32 characters"),
		)
	}

	if len(params.Symbol) > 10 || len(params.Symbol) < 3 {
		return utils.StackErrors(
			ErrInvalidParameter,
			errors.New("symbol must be between 3 and 10 characters"),
		)
	}

	if params.FeePayer == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("fee payer is required"),
		)
	}

	if params.Creators != nil || len(params.Creators) > 0 {
		feePayerInCreators := false
		for _, creator := range params.Creators {
			if creator.Address == "" {
				return utils.StackErrors(
					ErrMissedRequiredParameters,
					errors.New("creator address is required"),
				)
			}

			if creator.Address == params.FeePayer {
				feePayerInCreators = true
			}
		}

		if !feePayerInCreators {
			return utils.StackErrors(
				ErrInvalidParameter,
				errors.New("fee payer must be in a creators list"),
			)
		}
	}

	if params.Uses != nil {
		if params.Uses.UseMethod == token_metadata.TokenUseMethodUnknown.String() {
			return utils.StackErrors(
				ErrInvalidParameter,
				errors.New("unknown token use method"),
			)
		}
	}

	return nil
}

// MintNonFungibleToken mints a non-fungible token (NFT).
// These are the "standard" non-fungible tokens with unique metadata and supply = 1 the community is already familiar with and have both a Metadata PDA and a Master Edition (or Edition) PDA. Examples of these are Solana Monkey Business, Stylish Studs and Thugbirdz.
// The token is minted to the owner's associated token account.
func (c *Client) MintNonFungibleToken(ctx context.Context, params MintNonFungibleTokenParams) (mintAddr, tx string, err error) {
	if err := params.Validate(); err != nil {
		return "", "", utils.StackErrors(ErrMintNonFungibleToken, err)
	}

	if params.Creators == nil || len(params.Creators) == 0 {
		params.Creators = []token_metadata.Creator{
			{
				Address:  params.FeePayer,
				Share:    100,
				Verified: true,
			},
		}
	}

	result, err := c.prepareInitMintTransaction(ctx, initMintTransactionParams{
		FeePayer:             params.FeePayer,
		Owner:                params.Owner,
		TokenStandard:        utils.Pointer(token_metadata.TokenStandardNonFungible),
		Decimals:             0,
		SupplyAmount:         1,
		MaxEditionSupply:     params.MaxSupply,
		Name:                 params.Name,
		Symbol:               params.Symbol,
		MetadataURI:          params.MetadataURI,
		Collection:           &token_metadata.Collection{Key: params.Collection},
		SellerFeeBasisPoints: params.SellerFeeBasisPoints,
		Creators:             &params.Creators,
		Uses:                 params.Uses,
	})
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintFungibleAsset,
			err,
		)
	}

	return result.Mint.PublicKey.ToBase58(), result.Tx, nil
}

// MintNonFungibleTokenEditionParams are the parameters for MintNonFungibleTokenEdition.
type MintNonFungibleTokenEditionParams struct {
	FeePayer string // required; base58 encoded address of the fee payer
	Mint     string // required; base58 encoded address of the master edition
	Owner    string // optional; base58 encoded address of the owner; default is fee payer
	MintTo   string // optional; base58 encoded address to mint the token to; default is owner
	SendTo   string // optional; base58 encoded address to send the token to; default is owner
}

// Validate validates the parameters.
func (params MintNonFungibleTokenEditionParams) Validate() error {
	if params.FeePayer == "" || params.Mint == "" {
		return utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("fee payer public key and mint address are required"),
		)
	}

	return nil
}

// MintNonFungibleTokenEdition mints a non-fungible token (NFT) edition.
func (c *Client) MintNonFungibleTokenEdition(ctx context.Context, params MintNonFungibleTokenEditionParams) (mintAddr, tx string, err error) {
	if err := params.Validate(); err != nil {
		return "", "", utils.StackErrors(ErrMintNonFungibleTokenEdition, err)
	}

	// If owner is not set, use fee payer as owner.
	if params.Owner == "" {
		params.Owner = params.FeePayer
	}

	// If mint to is not set, use owner as mint to.
	if params.MintTo == "" {
		params.MintTo = params.Owner
	}

	// Function scoped variables.
	var (
		feePayerPublicKey common.PublicKey = common.PublicKeyFromString(params.FeePayer)
		ownerPublicKey    common.PublicKey = common.PublicKeyFromString(params.Owner)
		mintToPublicKey   common.PublicKey = common.PublicKeyFromString(params.MintTo)
		sendToPublicKey   *common.PublicKey

		masterOwner            common.PublicKey = ownerPublicKey
		masterOwnerAta         common.PublicKey
		masterMintPublicKey    common.PublicKey = common.PublicKeyFromString(params.Mint)
		masterMetaPublicKey    common.PublicKey
		masterEditionPublicKey common.PublicKey

		newMintOwner            common.PublicKey = mintToPublicKey
		newMintOwnerAta         common.PublicKey
		newMint                 types.Account = NewAccount()
		newMintMetaPublicKey    common.PublicKey
		newMintEditionPublicKey common.PublicKey
		newMintEditionMark      common.PublicKey
		newMintEdition          uint64
	)

	// Get next edition number.
	{
		editionInfo, err := c.GetMasterEditionInfo(ctx, masterMintPublicKey.ToBase58())
		if err != nil || editionInfo == nil {
			return "", "", utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				err,
			)
		}
		if editionInfo.Type == "" || editionInfo.Type != token_metadata.KeyMasterEdition.String() {
			return "", "", utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				ErrTokenIsNotMasterEdition,
			)
		}
		if editionInfo.MaxSupply == 0 || editionInfo.Supply == editionInfo.MaxSupply {
			return "", "", utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				ErrMaxSupplyReached,
			)
		}

		newMintEdition = editionInfo.Supply + 1
	}

	masterOwnerAta, _, err = common.FindAssociatedTokenAddress(masterOwner, masterMintPublicKey)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrFindAssociatedTokenAddress,
			err,
		)
	}

	masterMetaPublicKey, err = token_metadata.DeriveTokenMetadataPubkey(masterMintPublicKey)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetTokenMetaPubkey,
			err,
		)
	}

	masterEditionPublicKey, err = token_metadata.DeriveEditionPubkey(masterMintPublicKey)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetMasterEdition,
			err,
		)
	}

	newMintOwnerAta, _, err = common.FindAssociatedTokenAddress(newMintOwner, newMint.PublicKey)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrFindAssociatedTokenAddress,
			err,
		)
	}

	newMintMetaPublicKey, err = token_metadata.DeriveTokenMetadataPubkey(newMint.PublicKey)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetTokenMetaPubkey,
			err,
		)
	}

	newMintEditionPublicKey, err = token_metadata.DeriveEditionPubkey(newMint.PublicKey)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetMasterEdition,
			err,
		)
	}

	newMintEditionMark, err = token_metadata.DeriveEditionMarkerPubkey(masterMintPublicKey, newMintEdition)
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			ErrGetEditionMark,
			err,
		)
	}

	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(ctx, MintAccountSize)
	if err != nil {
		return "", "", utils.StackErrors(ErrMintNonFungibleTokenEdition, err)
	}

	instructions := []types.Instruction{
		system.CreateAccount(system.CreateAccountParam{
			From:     feePayerPublicKey,
			New:      newMint.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalance,
			Space:    token.MintAccountSize,
		}),
		token.InitializeMint(token.InitializeMintParam{
			Decimals:   0,
			Mint:       newMint.PublicKey,
			MintAuth:   newMintOwner,
			FreezeAuth: utils.Pointer(newMintOwner),
		}),
		associated_token_account.CreateAssociatedTokenAccount(
			associated_token_account.CreateAssociatedTokenAccountParam{
				Funder:                 feePayerPublicKey,
				Owner:                  newMintOwner,
				Mint:                   newMint.PublicKey,
				AssociatedTokenAccount: newMintOwnerAta,
			},
		),
		token.MintTo(token.MintToParam{
			Mint:   newMint.PublicKey,
			Auth:   ownerPublicKey,
			To:     newMintOwnerAta,
			Amount: 1,
		}),
		metaplex_token_metadata.MintNewEditionFromMasterEditionViaToken(
			metaplex_token_metadata.MintNewEditionFromMasterEditionViaTokeParam{
				NewMetaData:                newMintMetaPublicKey,
				NewEdition:                 newMintEditionPublicKey,
				MasterEdition:              masterEditionPublicKey,
				NewMint:                    newMint.PublicKey,
				NewMintAuthority:           newMintOwner,
				Payer:                      feePayerPublicKey,
				TokenAccountOwner:          masterOwner,
				TokenAccount:               masterOwnerAta,
				NewMetadataUpdateAuthority: newMintOwner,
				MasterMetadata:             masterMetaPublicKey,

				EditionMark: newMintEditionMark,
				Edition:     newMintEdition,
			},
		),
	}

	if sendToPublicKey == nil {
		recipientAta, _, err := common.FindAssociatedTokenAddress(*sendToPublicKey, newMint.PublicKey)
		if err != nil {
			return "", "", utils.StackErrors(
				ErrMintNonFungibleTokenEdition,
				ErrFindAssociatedTokenAddress,
				err,
			)
		}
		instructions = append(instructions,
			associated_token_account.CreateAssociatedTokenAccount(
				associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 feePayerPublicKey,
					Owner:                  *sendToPublicKey,
					Mint:                   newMint.PublicKey,
					AssociatedTokenAccount: recipientAta,
				},
			),
			token.Transfer(
				token.TransferParam{
					From:   newMintOwnerAta,
					To:     recipientAta,
					Auth:   newMintOwner,
					Amount: 1,
				},
			),
		)
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     params.FeePayer,
		Instructions: instructions,
		Signers:      []types.Account{newMint},
	})
	if err != nil {
		return "", "", utils.StackErrors(
			ErrMintNonFungibleTokenEdition,
			err,
		)
	}

	return newMint.PublicKey.String(), txb, nil
}

type (
	initMintTransactionParams struct {
		FeePayer string
		Owner    string

		TokenStandard *token_metadata.TokenStandard

		Decimals     uint8
		SupplyAmount uint64
		FixedSupply  bool

		Name        string
		Symbol      string
		MetadataURI string
		Collection  *token_metadata.Collection
		Uses        *token_metadata.Uses

		MaxEditionSupply     uint64
		SellerFeeBasisPoints uint16
		Creators             *[]token_metadata.Creator
	}

	initMintTransactionResult struct {
		Tx   string
		Mint types.Account
	}
)

// Prepare initialize mint transaction.
func (c *Client) prepareInitMintTransaction(ctx context.Context, params initMintTransactionParams) (*initMintTransactionResult, error) {
	feePayerPubKey := common.PublicKeyFromString(params.FeePayer)
	if params.Owner == "" {
		params.Owner = params.FeePayer
	}
	ownerPubKey := common.PublicKeyFromString(params.Owner)
	mint := NewAccount()

	ownerAta, _, err := common.FindAssociatedTokenAddress(ownerPubKey, mint.PublicKey)
	if err != nil {
		return nil, utils.StackErrors(
			ErrMintFungibleToken,
			ErrFindAssociatedTokenAddress,
			err,
		)
	}

	metaPubkey, err := metaplex_token_metadata.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		return nil, utils.StackErrors(
			ErrMintFungibleToken,
			ErrGetTokenMetaPubkey,
			err,
		)
	}

	if params.Decimals > 9 {
		params.Decimals = c.defaultDecimals
	}

	var (
		collection        *metaplex_token_metadata.Collection
		collectionDetails *metaplex_token_metadata.CollectionDetails
	)
	if params.Collection != nil {
		collection = &metaplex_token_metadata.Collection{
			Key:      common.PublicKeyFromString(params.Collection.Key),
			Verified: false,
		}
		if params.Collection.Size > 0 {
			collectionDetails = &metaplex_token_metadata.CollectionDetails{
				V1: metaplex_token_metadata.CollectionDetailsV1{
					Size: params.Collection.Size,
				},
			}
		}
	}
	_ = collectionDetails // TODO: add instruction to have ability create sized collections

	var uses *metaplex_token_metadata.Uses
	if params.Uses != nil {
		if useMethod := token_metadata.TokenUseMethod(params.Uses.UseMethod); useMethod.Valid() {
			uses = &metaplex_token_metadata.Uses{
				UseMethod: useMethod.ToMetadataUseMethod(),
				Remaining: params.Uses.Remaining,
				Total:     params.Uses.Total,
			}
		}
	}

	var creators *[]metaplex_token_metadata.Creator
	if params.Creators != nil && len(*params.Creators) > 0 {
		creators = &[]metaplex_token_metadata.Creator{}
		totalShare := uint8(0)
		feePayerInCreators := false
		for _, c := range *params.Creators {
			*creators = append(*creators, metaplex_token_metadata.Creator{
				Address:  common.PublicKeyFromString(c.Address),
				Share:    c.Share,
				Verified: c.Address == params.Owner,
			})
			if c.Address == params.FeePayer {
				feePayerInCreators = true
			}
			totalShare += c.Share
		}
		if totalShare != 100 {
			return nil, utils.StackErrors(
				ErrInvalidParameter,
				errors.New("total creators share must be 100"),
			)
		}
		if !feePayerInCreators {
			return nil, utils.StackErrors(
				ErrInvalidParameter,
				errors.New("fee payer must be in creators list"),
			)
		}
	}

	if params.TokenStandard == nil {
		return nil, utils.StackErrors(
			ErrMissedRequiredParameters,
			errors.New("token standard is required field"),
		)
	}

	var freezeAuth *common.PublicKey
	if *params.TokenStandard == token_metadata.TokenStandardNonFungible || *params.TokenStandard == token_metadata.TokenStandardNonFungibleEdition {
		freezeAuth = utils.Pointer(ownerPubKey)
	}

	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(ctx, MintAccountSize)
	if err != nil {
		return nil, utils.StackErrors(ErrMintFungibleToken, err)
	}

	instructions := []types.Instruction{
		system.CreateAccount(system.CreateAccountParam{
			From:     feePayerPubKey,
			New:      mint.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalance,
			Space:    token.MintAccountSize,
		}),
		token.InitializeMint(token.InitializeMintParam{
			Decimals:   params.Decimals,
			Mint:       mint.PublicKey,
			MintAuth:   ownerPubKey,
			FreezeAuth: freezeAuth,
		}),
		metaplex_token_metadata.CreateMetadataAccountV2(metaplex_token_metadata.CreateMetadataAccountV2Param{
			Metadata:                metaPubkey,
			Mint:                    mint.PublicKey,
			MintAuthority:           ownerPubKey,
			Payer:                   feePayerPubKey,
			UpdateAuthority:         ownerPubKey,
			UpdateAuthorityIsSigner: true,
			IsMutable:               true,
			Data: metaplex_token_metadata.DataV2{
				Name:                 params.Name,
				Symbol:               params.Symbol,
				Uri:                  params.MetadataURI,
				SellerFeeBasisPoints: params.SellerFeeBasisPoints,
				Creators:             creators,
				Collection:           collection,
				Uses:                 uses,
			},
		}),
	}

	if params.SupplyAmount > 0 {
		instructions = append(
			instructions,
			associated_token_account.CreateAssociatedTokenAccount(
				associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 feePayerPubKey,
					Owner:                  ownerPubKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ownerAta,
				},
			),
			token.MintToChecked(token.MintToCheckedParam{
				Mint:     mint.PublicKey,
				Auth:     ownerPubKey,
				Signers:  []common.PublicKey{},
				To:       ownerAta,
				Amount:   params.SupplyAmount,
				Decimals: params.Decimals,
			}),
		)
	}

	if params.FixedSupply && params.SupplyAmount > 0 &&
		*params.TokenStandard != token_metadata.TokenStandardNonFungibleEdition &&
		*params.TokenStandard != token_metadata.TokenStandardNonFungible {

		instructions = append(instructions, token.SetAuthority(token.SetAuthorityParam{
			Account:  mint.PublicKey,
			AuthType: token.AuthorityTypeMintTokens,
			Auth:     ownerPubKey,
			NewAuth:  nil,
			Signers:  []common.PublicKey{},
		}))
	}

	if *params.TokenStandard == token_metadata.TokenStandardNonFungible {
		tokenMasterEditionPubkey, err := metaplex_token_metadata.GetMasterEdition(mint.PublicKey)
		if err != nil {
			return nil, utils.StackErrors(
				ErrGetMasterEditionPubKey,
				err,
			)
		}

		instructions = append(instructions, metaplex_token_metadata.CreateMasterEditionV3(
			metaplex_token_metadata.CreateMasterEditionParam{
				Edition:         tokenMasterEditionPubkey,
				Mint:            mint.PublicKey,
				UpdateAuthority: ownerPubKey,
				MintAuthority:   ownerPubKey,
				Metadata:        metaPubkey,
				Payer:           feePayerPubKey,
				MaxSupply:       utils.Pointer(params.MaxEditionSupply),
			},
		))
	}

	// Verify fee payer if it's not the owner
	if params.Owner != params.FeePayer {
		instructions = append(instructions, metaplex_token_metadata.SignMetadata(
			metaplex_token_metadata.SignMetadataParam{
				Metadata: metaPubkey,
				Creator:  feePayerPubKey,
			},
		))
	}

	txb, err := c.NewTransaction(ctx, NewTransactionParams{
		FeePayer:     params.FeePayer,
		Instructions: instructions,
		Signers:      []types.Account{mint},
	})
	if err != nil {
		return nil, utils.StackErrors(
			ErrMintFungibleToken,
			ErrNewTransaction,
			err,
		)
	}

	return &initMintTransactionResult{
		Mint: mint,
		Tx:   txb,
	}, nil
}
