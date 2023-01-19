package token_metadata

import (
	"context"
	"fmt"

	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/solplaydev/solana/metadata"
)

// DeriveTokenMetadataPubkey returns the token metadata program public key.
func DeriveTokenMetadataPubkey(mint common.PublicKey) (common.PublicKey, error) {
	pk, err := token_metadata.GetTokenMetaPubkey(mint)
	if err != nil {
		return common.PublicKey{}, fmt.Errorf("failed to derive token metadata pubkey: %w", err)
	}

	return pk, nil
}

// DeserializeMetadata deserializes the metadata.
func DeserializeMetadata(data []byte) (*Metadata, error) {
	md, err := token_metadata.MetadataDeserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata: %w", err)
	}

	m := &Metadata{
		UpdateAuthority:      md.UpdateAuthority.ToBase58(),
		Mint:                 md.Mint.ToBase58(),
		PrimarySaleHappened:  md.PrimarySaleHappened,
		IsMutable:            md.IsMutable,
		TokenStandard:        CastToTokenStandard(*md.TokenStandard).String(),
		MetadataUri:          md.Data.Uri,
		SellerFeeBasisPoints: md.Data.SellerFeeBasisPoints,
	}

	if md.Data.Uri != "" {
		mdp, err := metadata.MetadataFromURI(md.Data.Uri)
		if err != nil {
			m.Data.Name = md.Data.Name
			m.Data.Symbol = md.Data.Symbol
		}
		m.Data = mdp
	} else {
		m.Data.Name = md.Data.Name
		m.Data.Symbol = md.Data.Symbol
	}

	if md.Collection != nil {
		m.Collection = &Collection{
			Verified: md.Collection.Verified,
			Key:      md.Collection.Key.ToBase58(),
		}
		if md.CollectionDetails != nil {
			m.Collection.Size = md.CollectionDetails.V1.Size
		}
	}

	if md.Uses != nil {
		m.Uses = &Uses{
			UseMethod: CastMetadataUseMethod(md.Uses.UseMethod).String(),
			Remaining: md.Uses.Remaining,
			Total:     md.Uses.Total,
		}
	}

	if md.Data.Creators != nil {
		if m.Creators == nil {
			m.Creators = make([]Creator, 0, len(*md.Data.Creators))
		}
		for _, creator := range *md.Data.Creators {
			m.Creators = append(m.Creators, Creator{
				Address:  creator.Address.ToBase58(),
				Verified: creator.Verified,
				Share:    creator.Share,
			})
		}
	}

	return m, nil
}

// DeriveEditionPubkey returns the edition public key.
func DeriveEditionPubkey(mint common.PublicKey) (common.PublicKey, error) {
	pk, err := token_metadata.GetMasterEdition(mint)
	if err != nil {
		return common.PublicKey{}, fmt.Errorf("failed to derive edition pubkey: %w", err)
	}

	return pk, nil
}

// DeserializeEdition deserializes the edition.
func DeserializeEdition(data []byte, getAccountInfo getAccountInfoFunc) (*Edition, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var edition token_metadata.MasterEditionV2
	if err := borsh.Deserialize(&edition, data); err != nil {
		return nil, fmt.Errorf("failed to deserialize edition key: %w", err)
	}

	e := &Edition{Type: CastToKey(edition.Key).String()}

	if edition.Key == token_metadata.KeyMasterEditionV1 || edition.Key == token_metadata.KeyMasterEditionV2 {
		masterEdition, err := DeserializeMasterEdition(data)
		if err != nil {
			return nil, err
		}

		e.MaxSupply = masterEdition.MaxSupply
		e.Supply = masterEdition.Supply

	} else if edition.Key == token_metadata.KeyEditionV1 {
		var editionData EditionData
		if err := borsh.Deserialize(&editionData, data); err != nil {
			return nil, fmt.Errorf("failed to deserialize edition data: %w", err)
		}

		e.Edition = editionData.Edition

		if editionData.Parent != PubNil && getAccountInfo != nil {
			parent, err := getAccountInfo(ctx, editionData.Parent.ToBase58())
			if err != nil {
				return nil, fmt.Errorf("failed to get parent account info: %w", err)
			}

			masterEdition, err := DeserializeMasterEdition(parent.Data)
			if err != nil {
				return nil, err
			}

			e.MaxSupply = masterEdition.MaxSupply
			e.Supply = masterEdition.Supply
		}
	}

	return e, nil
}

// DeserializeMasterEdition deserializes the master edition data.
func DeserializeMasterEdition(data []byte) (*Edition, error) {
	masterEdition := &token_metadata.MasterEditionV2{}
	if err := borsh.Deserialize(masterEdition, data); err != nil {
		return nil, fmt.Errorf("failed to deserialize master edition: %w", err)
	}
	return &Edition{
		Type:      CastToKey(masterEdition.Key).String(),
		MaxSupply: *masterEdition.MaxSupply,
		Supply:    masterEdition.Supply,
	}, nil
}

// DeriveEditionMarkerPubkey returns the edition marker public key.
func DeriveEditionMarkerPubkey(mint common.PublicKey, edition uint64) (common.PublicKey, error) {
	pk, err := token_metadata.GetEditionMark(mint, edition)
	if err != nil {
		return common.PublicKey{}, fmt.Errorf("failed to derive edition marker pubkey: %w", err)
	}

	return pk, nil
}
