package token_metadata

import (
	"fmt"

	"github.com/dmitrymomot/solana/metadata"
	"github.com/dmitrymomot/solana/utils"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
)

// TokenMetadataInstructionBuilder is a builder for token metadata instructions
type TokenMetadataInstructionBuilder struct {
	metadata                common.PublicKey
	mint                    common.PublicKey
	mintAuthority           common.PublicKey
	payer                   common.PublicKey
	updateAuthority         *common.PublicKey
	updateAuthorityIsSigner *bool
	isMutable               *bool
	data                    token_metadata.DataV2
}

// NewTokenMetadataInstructionBuilder creates a new TokenMetadataInstructionBuilder
func NewTokenMetadataInstructionBuilder() *TokenMetadataInstructionBuilder {
	return &TokenMetadataInstructionBuilder{}
}

// SetMetadata sets the metadata public key
func (b *TokenMetadataInstructionBuilder) SetMetadata(metadata common.PublicKey) *TokenMetadataInstructionBuilder {
	b.metadata = metadata
	return b
}

// SetMetadataBase58 sets the metadata public key from base58 string
func (b *TokenMetadataInstructionBuilder) SetMetadataBase58(metadata string) *TokenMetadataInstructionBuilder {
	return b.SetMetadata(common.PublicKeyFromString(metadata))
}

// SetMint sets the mint public key
func (b *TokenMetadataInstructionBuilder) SetMint(mint common.PublicKey) *TokenMetadataInstructionBuilder {
	b.mint = mint
	return b
}

// SetMintBase58 sets the mint public key from base58 string
func (b *TokenMetadataInstructionBuilder) SetMintBase58(mint string) *TokenMetadataInstructionBuilder {
	return b.SetMint(common.PublicKeyFromString(mint))
}

// SetMintAuthority sets the mint authority public key
func (b *TokenMetadataInstructionBuilder) SetMintAuthority(mintAuthority common.PublicKey) *TokenMetadataInstructionBuilder {
	b.mintAuthority = mintAuthority
	return b
}

// SetMintAuthorityBase58 sets the mint authority public key from base58 string
func (b *TokenMetadataInstructionBuilder) SetMintAuthorityBase58(mintAuthority string) *TokenMetadataInstructionBuilder {
	return b.SetMintAuthority(common.PublicKeyFromString(mintAuthority))
}

// SetPayer sets the payer public key
func (b *TokenMetadataInstructionBuilder) SetPayer(payer common.PublicKey) *TokenMetadataInstructionBuilder {
	b.payer = payer
	return b
}

// SetPayerBase58 sets the payer public key from base58 string
func (b *TokenMetadataInstructionBuilder) SetPayerBase58(payer string) *TokenMetadataInstructionBuilder {
	return b.SetPayer(common.PublicKeyFromString(payer))
}

// SetUpdateAuthority sets the update authority public key
func (b *TokenMetadataInstructionBuilder) SetUpdateAuthority(updateAuthority common.PublicKey) *TokenMetadataInstructionBuilder {
	b.updateAuthority = utils.Pointer(updateAuthority)
	b.updateAuthorityIsSigner = utils.Pointer(true)
	b.isMutable = utils.Pointer(true)
	return b
}

// SetUpdateAuthorityBase58 sets the update authority public key from base58 string
func (b *TokenMetadataInstructionBuilder) SetUpdateAuthorityBase58(updateAuthority string) *TokenMetadataInstructionBuilder {
	return b.SetUpdateAuthority(common.PublicKeyFromString(updateAuthority))
}

// UpdateAuthorityIsSigner sets the update authority is signer
func (b *TokenMetadataInstructionBuilder) UpdateAuthorityIsSigner(updateAuthorityIsSigner bool) *TokenMetadataInstructionBuilder {
	b.updateAuthorityIsSigner = utils.Pointer(updateAuthorityIsSigner)
	return b
}

// IsMutable sets the is mutable
func (b *TokenMetadataInstructionBuilder) IsMutable(isMutable bool) *TokenMetadataInstructionBuilder {
	b.isMutable = utils.Pointer(isMutable)
	return b
}

// SetData sets the data
func (b *TokenMetadataInstructionBuilder) SetData(data token_metadata.DataV2) *TokenMetadataInstructionBuilder {
	b.data = data
	return b
}

// SetName sets the name
func (b *TokenMetadataInstructionBuilder) SetName(name string) *TokenMetadataInstructionBuilder {
	b.data.Name = name
	return b
}

// SetSymbol sets the symbol
func (b *TokenMetadataInstructionBuilder) SetSymbol(symbol string) *TokenMetadataInstructionBuilder {
	b.data.Symbol = symbol
	return b
}

// SetUri sets the uri
func (b *TokenMetadataInstructionBuilder) SetUri(uri string) *TokenMetadataInstructionBuilder {
	b.data.Uri = uri
	if uri != "" {
		m, err := metadata.MetadataFromURI(uri)
		if err != nil || m == nil {
			return b
		}
		b.SetName(m.Name)
		b.SetSymbol(m.Symbol)
	}
	return b
}

// SetSellerFeeBasisPoints sets the seller fee basis points
func (b *TokenMetadataInstructionBuilder) SetSellerFeeBasisPoints(sellerFeeBasisPoints uint16) *TokenMetadataInstructionBuilder {
	b.data.SellerFeeBasisPoints = sellerFeeBasisPoints
	return b
}

// SetCreator sets
func (b *TokenMetadataInstructionBuilder) SetCreator(creator common.PublicKey, share uint8) *TokenMetadataInstructionBuilder {
	if b.data.Creators == nil {
		b.data.Creators = &[]token_metadata.Creator{}
	}
	creators := append(*b.data.Creators, token_metadata.Creator{
		Address:  creator,
		Verified: false,
		Share:    share,
	})
	b.data.Creators = &creators

	return b
}

// SetCreatorBase58 sets the creator public key encoded in base58
func (b *TokenMetadataInstructionBuilder) SetCreatorBase58(creator string, share uint8) *TokenMetadataInstructionBuilder {
	return b.SetCreator(common.PublicKeyFromString(creator), share)
}

// SetCollection sets the on-chain collection
// This is used to group tokens together.
// Collection should be verified by the owner of the collection.
func (b *TokenMetadataInstructionBuilder) SetCollection(collection common.PublicKey) *TokenMetadataInstructionBuilder {
	b.data.Collection = &token_metadata.Collection{
		Key: collection,
	}
	return b
}

// SetCollectionBase58 sets the on-chain collection public key encoded in base58
func (b *TokenMetadataInstructionBuilder) SetCollectionBase58(collection string) *TokenMetadataInstructionBuilder {
	return b.SetCollection(common.PublicKeyFromString(collection))
}

// SetUses provides a way to set the number of times a token can be used.
// This is useful for NFTs that can be used multiple times.
// For example, a ticket NFT that can be used once.
func (b *TokenMetadataInstructionBuilder) SetUses(method TokenUseMethod, remaining, total uint64) *TokenMetadataInstructionBuilder {
	if method == TokenUseMethodUnknown {
		return b
	}
	if remaining > total {
		remaining = total
	}
	if total == 0 {
		total = 1
	}

	b.data.Uses = &token_metadata.Uses{
		UseMethod: method.ToMetadataUseMethod(),
		Remaining: remaining,
		Total:     total,
	}

	return b
}

// Build builds the instruction
func (b *TokenMetadataInstructionBuilder) Build() (meta common.PublicKey, instruction types.Instruction, err error) {
	if b.mint == PubNil {
		return PubNil, types.Instruction{}, fmt.Errorf("mint public key is required")
	}

	if b.payer == PubNil {
		return PubNil, types.Instruction{}, fmt.Errorf("payer public key is required")
	}

	if b.metadata == PubNil {
		mintMeta, err := token_metadata.GetTokenMetaPubkey(b.mint)
		if err != nil {
			return PubNil, types.Instruction{}, fmt.Errorf("failed to get metadata account pubkey: %w", err)
		}
		b.metadata = mintMeta
	}

	if b.mintAuthority == PubNil {
		b.mintAuthority = b.payer
	}

	if b.updateAuthority == nil {
		b.SetUpdateAuthority(b.payer)
	}

	params := token_metadata.CreateMetadataAccountV2Param{
		Metadata:                b.metadata,
		Mint:                    b.mint,
		MintAuthority:           b.mintAuthority,
		Payer:                   b.payer,
		UpdateAuthority:         *b.updateAuthority,
		UpdateAuthorityIsSigner: *b.updateAuthorityIsSigner,
		IsMutable:               *b.isMutable,
		Data:                    b.data,
	}

	return b.metadata, token_metadata.CreateMetadataAccountV2(params), nil
}
