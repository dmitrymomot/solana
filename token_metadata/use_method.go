package token_metadata

import "github.com/portto/solana-go-sdk/program/metaplex/token_metadata"

// TokenUseMethod represents the use method of a token.
type TokenUseMethod string

// TokenUseMethod enum.
const (
	TokenUseMethodUnknown TokenUseMethod = "unknown"
	TokenUseMethodBurn    TokenUseMethod = "burn"
	TokenUseMethodSingle  TokenUseMethod = "single"
	TokenUseMethodMulti   TokenUseMethod = "multiple"
)

// String returns the string representation of the token use method.
func (m TokenUseMethod) String() string {
	if !m.Valid() {
		return string(TokenUseMethodUnknown)
	}
	return string(m)
}

// Valid returns true if the token use method is valid.
func (m TokenUseMethod) Valid() bool {
	return m == TokenUseMethodBurn || m == TokenUseMethodSingle || m == TokenUseMethodMulti
}

// ToMetadataUseMethod converts the token use method to token_metadata.UseMethod.
func (m TokenUseMethod) ToMetadataUseMethod() token_metadata.UseMethod {
	switch m {
	case TokenUseMethodBurn:
		return token_metadata.Burn
	case TokenUseMethodSingle:
		return token_metadata.Single
	case TokenUseMethodMulti:
		return token_metadata.Multiple
	default:
		return token_metadata.Burn
	}
}

// CastMetadataUseMethod converts the token_metadata.UseMethod to TokenUseMethod.
func CastMetadataUseMethod(m token_metadata.UseMethod) TokenUseMethod {
	switch m {
	case token_metadata.Burn:
		return TokenUseMethodBurn
	case token_metadata.Single:
		return TokenUseMethodSingle
	case token_metadata.Multiple:
		return TokenUseMethodMulti
	default:
		return TokenUseMethodUnknown
	}
}
