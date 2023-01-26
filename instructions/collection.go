package instructions

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
)

// CreateCollectionParams is the params for CreateCollection
type VerifyCollectionItemParams struct {
	Mint                common.PublicKey  // required; The mint of the token
	CollectionMint      common.PublicKey  // required; The mint of the collection
	CollectionAuthority common.PublicKey  // required; The authority of the collection
	FeePayer            *common.PublicKey // optional; The fee payer of the collection; default is collection authority
}

// Validate validates the params.
func (p VerifyCollectionItemParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("token mint is required")
	}
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection authority is required")
	}
	if p.FeePayer != nil && *p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is invalid")
	}
	return nil
}

// VerifyCollectionItem verifies the collection.
func VerifyCollectionItem(params VerifyCollectionItemParams) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate verify collection: %w", err)
		}

		if params.FeePayer == nil {
			params.FeePayer = &params.CollectionAuthority
		}

		tokenMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("derive token metadata pubkey: %w", err)
		}

		collectionMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection metadata pubkey: %w", err)
		}

		collectionMaster, err := token_metadata.DeriveEditionPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection master edition pubkey: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.VerifyCollection(metaplex_token_metadata.VerifyCollectionParam{
				Metadata:                       tokenMetadata,
				Payer:                          *params.FeePayer,
				CollectionAuthority:            params.CollectionAuthority,
				CollectionMint:                 params.CollectionMint,
				CollectionMetadata:             collectionMetadata,
				CollectionMasterEditionAccount: collectionMaster,
			}),
		}, nil
	}
}

// UnverifyCollectionItemParams is the params for UnverifyCollectionItem
type UnverifyCollectionItemParams struct {
	Mint                common.PublicKey // required; The mint of the token
	CollectionMint      common.PublicKey // required; The mint of the collection
	CollectionAuthority common.PublicKey // required; The authority of the collection
}

// Validate validates the params.
func (p UnverifyCollectionItemParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("token mint is required")
	}
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection authority is required")
	}
	return nil
}

// UnverifyCollectionItem unverifies the collection.
func UnverifyCollectionItem(params UnverifyCollectionItemParams) InstructionFunc {
	return func() ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate unverify collection: %w", err)
		}

		tokenMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.Mint)
		if err != nil {
			return nil, fmt.Errorf("derive token metadata pubkey: %w", err)
		}

		collectionMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection metadata pubkey: %w", err)
		}

		collectionMaster, err := token_metadata.DeriveEditionPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection master edition pubkey: %w", err)
		}

		collectionAuthorityRecord, err := token_metadata.DeriveCollectionAuthorityRecord(params.CollectionMint, params.CollectionAuthority)
		if err != nil {
			return nil, fmt.Errorf("derive collection authority record: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.UnverifyCollection(metaplex_token_metadata.UnverifyCollectionParam{
				Metadata:                       tokenMetadata,
				CollectionAuthority:            params.CollectionAuthority,
				CollectionMint:                 params.CollectionMint,
				CollectionMetadata:             collectionMetadata,
				CollectionMasterEditionAccount: collectionMaster,
				CollectionAuthorityRecord:      collectionAuthorityRecord,
			}),
		}, nil
	}
}
