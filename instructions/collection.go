package instructions

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	metaplex_token_metadata "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/types"
	"github.com/solplaydev/solana/token_metadata"
)

// ApproveCollectionAuthorityParams is the params for ApproveCollectionAuthority
type ApproveCollectionAuthorityParams struct {
	CollectionMint            common.PublicKey  // required; The mint of the collection
	CollectionUpdateAuthority common.PublicKey  // required; The current update authority of the collection metadata
	NewCollectionAuthority    common.PublicKey  // required; The new authority of the collection
	FeePayer                  *common.PublicKey // optional; The fee payer of the transaction; default is the collection update authority
}

// Validate validates the params.
func (p ApproveCollectionAuthorityParams) Validate() error {
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionUpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection update authority is required")
	}
	if p.NewCollectionAuthority == (common.PublicKey{}) {
		return fmt.Errorf("new collection authority is required")
	}
	if p.FeePayer != nil && *p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("invalid fee payer public key")
	}
	return nil
}

// ApproveCollectionAuthority approves the collection authority.
func ApproveCollectionAuthority(params ApproveCollectionAuthorityParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate approve collection authority: %w", err)
		}

		if params.FeePayer == nil {
			params.FeePayer = &params.CollectionUpdateAuthority
		}

		collectionMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection metadata pubkey: %w", err)
		}

		collectionAuthorityRecord, err := token_metadata.DeriveCollectionAuthorityRecord(params.CollectionMint, params.NewCollectionAuthority)
		if err != nil {
			return nil, fmt.Errorf("derive new collection authority record: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.ApproveCollectionAuthority(metaplex_token_metadata.ApproveCollectionAuthorityParam{
				CollectionAuthorityRecord: collectionAuthorityRecord,
				NewCollectionAuthority:    params.NewCollectionAuthority,
				UpdateAuthority:           params.CollectionUpdateAuthority,
				Payer:                     *params.FeePayer,
				CollectionMetadata:        collectionMetadata,
				CollectionMint:            params.CollectionMint,
			}),
		}, nil
	}
}

// RevokeCollectionAuthorityParams is the params for RevokeCollectionAuthority
type RevokeCollectionAuthorityParams struct {
	CollectionMint            common.PublicKey // required; The mint of the collection
	CollectionUpdateAuthority common.PublicKey // required; The current update authority of the collection metadata
	RevokeAuthority           common.PublicKey // required; The authority to revoke
}

// Validate validates the params.
func (p RevokeCollectionAuthorityParams) Validate() error {
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionUpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection update authority is required")
	}
	if p.RevokeAuthority == (common.PublicKey{}) {
		return fmt.Errorf("revoke authority is required")
	}
	return nil
}

// RevokeCollectionAuthority revokes the collection authority.
func RevokeCollectionAuthority(params RevokeCollectionAuthorityParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate revoke collection authority: %w", err)
		}

		collectionMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection metadata pubkey: %w", err)
		}

		collectionAuthorityRecord, err := token_metadata.DeriveCollectionAuthorityRecord(params.CollectionMint, params.RevokeAuthority)
		if err != nil {
			return nil, fmt.Errorf("derive revoke authority record: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.RevokeCollectionAuthority(metaplex_token_metadata.RevokeCollectionAuthorityParam{
				CollectionAuthorityRecord: collectionAuthorityRecord,
				DelegateAuthority:         params.CollectionUpdateAuthority,
				RevokeAuthority:           params.RevokeAuthority,
				CollectionMetadata:        collectionMetadata,
				CollectionMint:            params.CollectionMint,
			}),
		}, nil
	}
}

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
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
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
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
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

// SetAndVerifyCollectionParams is the params for SetAndVerifyCollection
type SetAndVerifyCollectionParams struct {
	Mint                common.PublicKey // required; The mint of the token
	UpdateAuthority     common.PublicKey // required; The update authority of the token
	CollectionMint      common.PublicKey // required; The mint of the collection
	CollectionAuthority common.PublicKey // required; The authority of the collection
	FeePayer            common.PublicKey // required; The fee payer of the transaction
}

// Validate validates the params.
func (p SetAndVerifyCollectionParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("token mint is required")
	}
	if p.UpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("update authority is required")
	}
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection authority is required")
	}
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("fee payer is required")
	}
	return nil
}

// SetAndVerifyCollection sets and verifies the collection.
func SetAndVerifyCollection(params SetAndVerifyCollectionParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate set and verify collection: %w", err)
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
			metaplex_token_metadata.SetAndVerifyCollection(metaplex_token_metadata.SetAndVerifyCollectionParam{
				Metadata:                       tokenMetadata,
				CollectionAuthority:            params.CollectionAuthority,
				Payer:                          params.FeePayer,
				UpdateAuthority:                params.UpdateAuthority,
				CollectionMint:                 params.CollectionMint,
				CollectionMetadata:             collectionMetadata,
				CollectionMasterEditionAccount: collectionMaster,
				CollectionAuthorityRecord:      collectionAuthorityRecord,
			}),
		}, nil
	}
}

// SetCollectionSize is the params for SetCollectionSize
type SetCollectionSizeParams struct {
	CollectionMint      common.PublicKey // required; The mint of the collection
	CollectionAuthority common.PublicKey // required; The authority of the collection
	Size                uint64           // required; The size of the collection
}

// Validate validates the params.
func (p SetCollectionSizeParams) Validate() error {
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection authority is required")
	}
	if p.Size == 0 {
		return fmt.Errorf("collection size is required")
	}
	return nil
}

// SetCollectionSize sets the size of the collection.
func SetCollectionSize(params SetCollectionSizeParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate set collection size: %w", err)
		}

		collectionMetadata, err := token_metadata.DeriveTokenMetadataPubkey(params.CollectionMint)
		if err != nil {
			return nil, fmt.Errorf("derive collection metadata pubkey: %w", err)
		}

		collectionAuthorityRecord, err := token_metadata.DeriveCollectionAuthorityRecord(params.CollectionMint, params.CollectionAuthority)
		if err != nil {
			return nil, fmt.Errorf("derive collection authority record: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.SetCollectionSize(metaplex_token_metadata.SetCollectionSizeParam{
				CollectionMetadata:        collectionMetadata,
				CollectionAuthority:       params.CollectionAuthority,
				CollectionMint:            params.CollectionMint,
				CollectionAuthorityRecord: collectionAuthorityRecord,
				Size:                      params.Size,
			}),
		}, nil
	}
}

// VerifySizedCollectionItemParams is the params for VerifySizedCollectionItem
type VerifySizedCollectionItemParams struct {
	Mint                common.PublicKey  // required; The mint of the token
	CollectionMint      common.PublicKey  // required; The mint of the collection
	CollectionAuthority common.PublicKey  // required; The authority of the collection
	FeePayer            *common.PublicKey // optional; The fee payer of the transaction; default is the collection authority
}

// Validate validates the params.
func (p VerifySizedCollectionItemParams) Validate() error {
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
		return fmt.Errorf("invalid fee payer public key")
	}
	return nil
}

// VerifySizedCollectionItem verifies the collection.
func VerifySizedCollectionItem(params VerifySizedCollectionItemParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate verify sized collection item: %w", err)
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

		collectionAuthorityRecord, err := token_metadata.DeriveCollectionAuthorityRecord(params.CollectionMint, params.CollectionAuthority)
		if err != nil {
			return nil, fmt.Errorf("derive collection authority record: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.VerifySizedCollectionItem(metaplex_token_metadata.VerifySizedCollectionItemParam{
				Metadata:                       tokenMetadata,
				CollectionAuthority:            params.CollectionAuthority,
				Payer:                          *params.FeePayer,
				CollectionMint:                 params.CollectionMint,
				CollectionMetadata:             collectionMetadata,
				CollectionMasterEditionAccount: collectionMaster,
				CollectionAuthorityRecord:      collectionAuthorityRecord,
			}),
		}, nil
	}
}

// UnverifySizedCollectionItemParams is the params for UnverifySizedCollectionItem
type UnverifySizedCollectionItemParams struct {
	Mint                common.PublicKey  // required; The mint of the token
	CollectionMint      common.PublicKey  // required; The mint of the collection
	CollectionAuthority common.PublicKey  // required; The authority of the collection
	FeePayer            *common.PublicKey // optional; The fee payer of the transaction; default is the collection authority
}

// Validate validates the params.
func (p UnverifySizedCollectionItemParams) Validate() error {
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
		return fmt.Errorf("invalid fee payer public key")
	}
	return nil
}

// UnverifySizedCollectionItem unverifies the collection.
func UnverifySizedCollectionItem(params UnverifySizedCollectionItemParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate unverify sized collection item: %w", err)
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

		collectionAuthorityRecord, err := token_metadata.DeriveCollectionAuthorityRecord(params.CollectionMint, params.CollectionAuthority)
		if err != nil {
			return nil, fmt.Errorf("derive collection authority record: %w", err)
		}

		return []types.Instruction{
			metaplex_token_metadata.UnverifySizedCollectionItem(metaplex_token_metadata.UnverifySizedCollectionItemParam{
				Metadata:                       tokenMetadata,
				CollectionAuthority:            params.CollectionAuthority,
				Payer:                          *params.FeePayer,
				CollectionMint:                 params.CollectionMint,
				CollectionMetadata:             collectionMetadata,
				CollectionMasterEditionAccount: collectionMaster,
				CollectionAuthorityRecord:      collectionAuthorityRecord,
			}),
		}, nil
	}
}

// SetAndVerifySizedCollectionItemParams is the params for SetAndVerifySizedCollectionItem
type SetAndVerifySizedCollectionItemParams struct {
	Mint                common.PublicKey // required; The mint of the token
	MintUpdateAuthority common.PublicKey // required; The mint update authority of the token
	CollectionMint      common.PublicKey // required; The mint of the collection
	CollectionAuthority common.PublicKey // required; The authority of the collection
	FeePayer            common.PublicKey // required; The fee payer of the transaction
}

// Validate validates the params.
func (p SetAndVerifySizedCollectionItemParams) Validate() error {
	if p.Mint == (common.PublicKey{}) {
		return fmt.Errorf("token mint is required")
	}
	if p.MintUpdateAuthority == (common.PublicKey{}) {
		return fmt.Errorf("token mint update authority is required")
	}
	if p.CollectionMint == (common.PublicKey{}) {
		return fmt.Errorf("collection mint is required")
	}
	if p.CollectionAuthority == (common.PublicKey{}) {
		return fmt.Errorf("collection authority is required")
	}
	if p.FeePayer == (common.PublicKey{}) {
		return fmt.Errorf("invalid fee payer public key")
	}
	return nil
}

// SetAndVerifySizedCollectionItem sets and verifies the collection.
func SetAndVerifySizedCollectionItem(params SetAndVerifySizedCollectionItemParams) InstructionFunc {
	return func(ctx context.Context, c Client) ([]types.Instruction, error) {
		if err := params.Validate(); err != nil {
			return nil, fmt.Errorf("validate set and verify sized collection item: %w", err)
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
			metaplex_token_metadata.SetAndVerifySizedCollectionItem(metaplex_token_metadata.SetAndVerifySizedCollectionItemParam{
				Metadata:                       tokenMetadata,
				CollectionAuthority:            params.CollectionAuthority,
				Payer:                          params.FeePayer,
				UpdateAuthority:                params.MintUpdateAuthority,
				CollectionMint:                 params.CollectionMint,
				CollectionMetadata:             collectionMetadata,
				CollectionMasterEditionAccount: collectionMaster,
				CollectionAuthorityRecord:      collectionAuthorityRecord,
			}),
		}, nil
	}
}
