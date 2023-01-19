package metadata

import (
	"errors"
	"fmt"
	"strings"

	"github.com/solplaydev/solana/utils"
)

// FungibleAssetMetadataBuilder is a builder to build fungible asset metadata
type FungibleAssetMetadataBuilder struct {
	name         string      // required
	symbol       string      // required
	description  string      // required
	image        string      // required
	animationURL string      // optional
	externalURL  string      // optional
	attributes   []Attribute // optional
}

// NewFungibleAssetMetadataBuilder creates a new FungibleAssetMetadataBuilder
func NewFungibleAssetMetadataBuilder() *FungibleAssetMetadataBuilder {
	return &FungibleAssetMetadataBuilder{}
}

// SetName sets the name of the asset
func (b *FungibleAssetMetadataBuilder) SetName(name string) *FungibleAssetMetadataBuilder {
	b.name = name
	return b
}

// SetSymbol sets the symbol of the asset
func (b *FungibleAssetMetadataBuilder) SetSymbol(symbol string) *FungibleAssetMetadataBuilder {
	b.symbol = strings.ToUpper(symbol)
	return b
}

// SetDescription sets the description of the asset
func (b *FungibleAssetMetadataBuilder) SetDescription(description string) *FungibleAssetMetadataBuilder {
	b.description = description
	return b
}

// SetImage sets the image of the asset
func (b *FungibleAssetMetadataBuilder) SetImage(image string) *FungibleAssetMetadataBuilder {
	b.image = image
	return b
}

// SetAnimationURL sets the animation URL of the asset
func (b *FungibleAssetMetadataBuilder) SetAnimationURL(animationURL string) *FungibleAssetMetadataBuilder {
	b.animationURL = animationURL
	return b
}

// SetExternalURL sets the external URL of the asset
func (b *FungibleAssetMetadataBuilder) SetExternalURL(externalURL string) *FungibleAssetMetadataBuilder {
	b.externalURL = externalURL
	return b
}

// SetAttributes sets the attributes of the asset
func (b *FungibleAssetMetadataBuilder) SetAttributes(attributes []Attribute) *FungibleAssetMetadataBuilder {
	b.attributes = attributes
	return b
}

// SetAttribute adds an attribute to the asset
func (b *FungibleAssetMetadataBuilder) SetAttribute(key string, value any) *FungibleAssetMetadataBuilder {
	if b.attributes == nil {
		b.attributes = make([]Attribute, 0)
	}

	var displayType AttributeDisplayType
	valueType := utils.GetVarType(value)
	switch valueType {
	case "string", "byte", "rune":
		displayType = AttributeDisplayString
	case "int", "float", "int64", "float64", "uint", "uint64", "int32", "float32", "uint32", "int16", "uint16", "int8", "uint8":
		displayType = AttributeDisplayNumber
	case "bool", "boolean":
		displayType = AttributeDisplayBoolean
	default:
		displayType = AttributeDisplayString
	}

	b.attributes = append(b.attributes, Attribute{
		TraitType:   key,
		Value:       fmt.Sprintf("%v", value),
		DisplayType: displayType,
	})

	return b
}

// SetAttributeStruct adds an attribute to the asset
func (b *FungibleAssetMetadataBuilder) SetAttributeStruct(attr Attribute) *FungibleAssetMetadataBuilder {
	if b.attributes == nil {
		b.attributes = make([]Attribute, 0)
	}

	b.attributes = append(b.attributes, attr)

	return b
}

// Build builds the fungible asset metadata
func (b *FungibleAssetMetadataBuilder) Build() (*Metadata, error) {
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
		Name:         b.name,
		Symbol:       b.symbol,
		Description:  b.description,
		Image:        b.image,
		AnimationURL: b.animationURL,
		ExternalURL:  b.externalURL,
		Attributes:   b.attributes,
	}, nil
}
