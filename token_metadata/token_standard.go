package token_metadata

import "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"

// TokenStandard represents the standard of a token.
type TokenStandard string

// String returns the string representation of the token standard.
func (s TokenStandard) String() string {
	if !s.Valid() {
		return string(TokenStandardUndefined)
	}
	return string(s)
}

// Valid returns true if the token standard is valid.
func (s TokenStandard) Valid() bool {
	return s == TokenStandardNonFungible ||
		s == TokenStandardNonFungibleEdition ||
		s == TokenStandardFungibleAsset ||
		s == TokenStandardFungible
}

// ToSystemTokenStandard returns the system token standard.
func (s TokenStandard) ToSystemTokenStandard() token_metadata.TokenStandard {
	switch s {
	case TokenStandardNonFungible:
		return token_metadata.NonFungible
	case TokenStandardNonFungibleEdition:
		return token_metadata.NonFungibleEdition
	case TokenStandardFungibleAsset:
		return token_metadata.FungibleAsset
	case TokenStandardFungible:
		return token_metadata.Fungible
	default:
		return token_metadata.NonFungible
	}
}

// Token standards enum
const (
	TokenStandardUndefined          TokenStandard = "undefined"
	TokenStandardNonFungible        TokenStandard = "non_fungible"
	TokenStandardNonFungibleEdition TokenStandard = "non_fungible_edition"
	TokenStandardFungibleAsset      TokenStandard = "fungible_asset"
	TokenStandardFungible           TokenStandard = "fungible"
)

// TokenStandardMap is a map of token_metadata.TokenStandard to TokenStandard.
var tokenStandardsMap = map[token_metadata.TokenStandard]TokenStandard{
	token_metadata.NonFungible:        TokenStandardNonFungible,
	token_metadata.NonFungibleEdition: TokenStandardNonFungibleEdition,
	token_metadata.FungibleAsset:      TokenStandardFungibleAsset,
	token_metadata.Fungible:           TokenStandardFungible,
}

// CastToTokenStandard casts token_metadata.TokenStandard to TokenStandard.
func CastToTokenStandard(s token_metadata.TokenStandard) TokenStandard {
	if v, ok := tokenStandardsMap[s]; ok {
		return v
	}
	return TokenStandardUndefined
}
