package metadata

import "errors"

// FungibleTokenMetadataBuilder is a builder to build fungible token metadata
type FungibleTokenMetadataBuilder struct {
	name        string // required
	symbol      string // required
	description string // required
	image       string // required
	externalURL string // optional
}

// NewFungibleTokenMetadataBuilder creates a new FungibleTokenMetadataBuilder
func NewFungibleTokenMetadataBuilder() *FungibleTokenMetadataBuilder {
	return &FungibleTokenMetadataBuilder{}
}

// SetName sets the name of the token
func (b *FungibleTokenMetadataBuilder) SetName(name string) *FungibleTokenMetadataBuilder {
	b.name = name
	return b
}

// SetSymbol sets the symbol of the token
func (b *FungibleTokenMetadataBuilder) SetSymbol(symbol string) *FungibleTokenMetadataBuilder {
	b.symbol = symbol
	return b
}

// SetDescription sets the description of the token
func (b *FungibleTokenMetadataBuilder) SetDescription(description string) *FungibleTokenMetadataBuilder {
	b.description = description
	return b
}

// SetImage sets the image of the token
func (b *FungibleTokenMetadataBuilder) SetImage(image string) *FungibleTokenMetadataBuilder {
	b.image = image
	return b
}

// SetExternalURL sets the external url of the token
func (b *FungibleTokenMetadataBuilder) SetExternalURL(externalURL string) *FungibleTokenMetadataBuilder {
	b.externalURL = externalURL
	return b
}

// Build builds the fungible token metadata
func (b *FungibleTokenMetadataBuilder) Build() (*Metadata, error) {
	if b.name == "" {
		return nil, errors.New("name is required")
	}
	if b.symbol == "" {
		return nil, errors.New("symbol is required")
	}
	if b.description == "" {
		return nil, errors.New("description is required")
	}
	if b.image == "" {
		return nil, errors.New("image is required")
	}
	return &Metadata{
		Name:        b.name,
		Symbol:      b.symbol,
		Description: b.description,
		Image:       b.image,
		ExternalURL: b.externalURL,
	}, nil
}
