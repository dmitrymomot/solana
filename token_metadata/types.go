package token_metadata

import (
	"context"

	"github.com/dmitrymomot/solana/metadata"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
)

// PubNil is a nil public key
var PubNil = common.PublicKey{}

type (
	// Metadata represents the metadata of a token.
	Metadata struct {
		UpdateAuthority      string             `json:"update_authority"`
		Mint                 string             `json:"mint"`
		PrimarySaleHappened  bool               `json:"primary_sale_happened,omitempty"`
		IsMutable            bool               `json:"is_mutable"`
		EditionNonce         *uint8             `json:"edition_nonce,omitempty"`
		TokenStandard        string             `json:"token_standard"`
		Collection           *Collection        `json:"collection,omitempty"`
		Uses                 *Uses              `json:"uses,omitempty"`
		Edition              *Edition           `json:"edition,omitempty"`
		MetadataUri          string             `json:"metadata_uri,omitempty"`
		SellerFeeBasisPoints uint16             `json:"seller_fee_basis_points,omitempty"`
		Creators             []Creator          `json:"creators,omitempty"`
		Data                 *metadata.Metadata `json:"data,omitempty"`
	}

	Edition struct {
		Type      string `json:"type,omitempty"`
		Supply    uint64 `json:"supply,omitempty"`
		MaxSupply uint64 `json:"max_supply,omitempty"`
		Edition   uint64 `json:"edition,omitempty"`
	}

	EditionKey struct {
		Key token_metadata.Key
	}

	EditionData struct {
		Key     token_metadata.Key
		Parent  common.PublicKey
		Edition uint64
	}

	Collection struct {
		Verified bool   `json:"verified"`
		Key      string `json:"key"`
		Size     uint64 `json:"size,omitempty"`
	}

	Uses struct {
		UseMethod string `json:"use_method"`
		Total     uint64 `json:"total"`
		Remaining uint64 `json:"remaining"`
	}

	Creator struct {
		Address  string `json:"address"`
		Verified bool   `json:"verified"`
		Share    uint8  `json:"share"`
	}

	// GetAccountInfoFunc is a function that returns the account info of a given address.
	getAccountInfoFunc func(ctx context.Context, base58Addr string) (client.AccountInfo, error)
)
