package token_metadata

import "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"

// Token metadata edition key
type Key string

// String returns the string representation of the key.
func (k Key) String() string {
	return string(k)
}

// Predefined token editions
const (
	KeyUndefined                 Key = "undefined"
	KeyMasterEdition             Key = "master_edition"
	KeyPrintedEdition            Key = "edition"
	KeyReservationList           Key = "reservation_list"
	KeyMetadata                  Key = "metadata"
	KeyEditionMarker             Key = "edition_marker"
	KeyUseAuthorityRecord        Key = "use_authority_record"
	KeyCollectionAuthorityRecord Key = "collection_authority_record"
)

// Map token_metadata.Key to string
var keysMap = map[token_metadata.Key]Key{
	token_metadata.KeyUninitialized:             KeyUndefined,
	token_metadata.KeyEditionV1:                 KeyPrintedEdition,
	token_metadata.KeyMasterEditionV1:           KeyMasterEdition,
	token_metadata.KeyReservationListV1:         KeyReservationList,
	token_metadata.KeyMetadataV1:                KeyMetadata,
	token_metadata.KeyReservationListV2:         KeyReservationList,
	token_metadata.KeyMasterEditionV2:           KeyMasterEdition,
	token_metadata.KeyEditionMarker:             KeyEditionMarker,
	token_metadata.KeyUseAuthorityRecord:        KeyUseAuthorityRecord,
	token_metadata.KeyCollectionAuthorityRecord: KeyCollectionAuthorityRecord,
}

// Cast token_metadata.Key to Key
func CastToKey(k token_metadata.Key) Key {
	if val, ok := keysMap[k]; ok {
		return val
	}
	return KeyUndefined
}
