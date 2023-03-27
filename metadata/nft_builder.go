package metadata

import (
	"errors"
	"fmt"

	"github.com/dmitrymomot/solana/utils"
)

// NFTMetadataBuilder is a builder to build non-fungible asset metadata
type NFTMetadataBuilder struct {
	name             string                 // required
	symbol           string                 // required
	description      string                 // required
	image            string                 // required
	animationURL     string                 // optional
	externalURL      string                 // optional
	attributes       []Attribute            // optional
	files            []File                 // optional
	category         PropertyCategory       // optional
	collection       *Collection            // optional; deprecated, use on-chain program instead
	customProperties map[string]interface{} // optional; any properties that you want to add to the metadata
}

// NewNFTMetadataBuilder creates a new NFTMetadataBuilder
func NewNFTMetadataBuilder() *NFTMetadataBuilder {
	return &NFTMetadataBuilder{}
}

// SetName sets the name of the asset
func (b *NFTMetadataBuilder) SetName(name string) *NFTMetadataBuilder {
	b.name = name
	return b
}

// SetSymbol sets the symbol of the asset
func (b *NFTMetadataBuilder) SetSymbol(symbol string) *NFTMetadataBuilder {
	b.symbol = symbol
	return b
}

// SetDescription sets the description of the asset
func (b *NFTMetadataBuilder) SetDescription(description string) *NFTMetadataBuilder {
	b.description = description
	return b
}

// SetImage sets the image of the asset
func (b *NFTMetadataBuilder) SetImage(image string) *NFTMetadataBuilder {
	b.image = image
	return b
}

// SetAnimationURL sets the animation URL of the asset
func (b *NFTMetadataBuilder) SetAnimationURL(animationURL string) *NFTMetadataBuilder {
	b.animationURL = animationURL
	return b
}

// SetExternalURL sets the external URL of the asset
func (b *NFTMetadataBuilder) SetExternalURL(externalURL string) *NFTMetadataBuilder {
	b.externalURL = externalURL
	return b
}

// Attributes sets the attributes of the asset
func (b *NFTMetadataBuilder) SetAttributes(attributes []Attribute) *NFTMetadataBuilder {
	b.attributes = attributes
	return b
}

// SetAttribute adds an attribute to the asset
func (b *NFTMetadataBuilder) SetAttribute(key string, value any) *NFTMetadataBuilder {
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
func (b *NFTMetadataBuilder) SetAttributeStruct(attr Attribute) *NFTMetadataBuilder {
	if b.attributes == nil {
		b.attributes = make([]Attribute, 0)
	}

	b.attributes = append(b.attributes, attr)

	return b
}

// SetCategory sets the category of the non-fungible token
func (b *NFTMetadataBuilder) SetCategory(category PropertyCategory) *NFTMetadataBuilder {
	b.category = category
	return b
}

// SetCustomProperty adds a custom property to the non-fungible token
func (b *NFTMetadataBuilder) SetCustomProperty(key string, value any) *NFTMetadataBuilder {
	if b.customProperties == nil {
		b.customProperties = make(map[string]interface{})
	}

	b.customProperties[key] = value
	return b
}

// SetCustomProperties adds custom properties to the non-fungible token
func (b *NFTMetadataBuilder) SetCustomProperties(properties map[string]interface{}) *NFTMetadataBuilder {
	b.customProperties = properties
	return b
}

// @deprecated
// SetCollection sets the collection of the non-fungible token
func (b *NFTMetadataBuilder) SetCollection(name, family string) *NFTMetadataBuilder {
	b.collection = &Collection{
		Name:   name,
		Family: family,
	}
	return b
}

// SetFiles sets the files of the non-fungible token
func (b *NFTMetadataBuilder) SetFiles(files []File) *NFTMetadataBuilder {
	b.files = files
	return b
}

// SetFile adds a file to the non-fungible token
func (b *NFTMetadataBuilder) SetFile(fileURI string, fileType string) *NFTMetadataBuilder {
	if b.files == nil {
		b.files = make([]File, 0)
	}

	if fileType == "" {
		fileType = utils.GetFileTypeByURI(fileURI)
	}
	if fileType == "" {
		fileType = "unknown"
	}

	b.files = append(b.files, File{
		URI:  fileURI,
		Type: fileType,
	})

	return b
}

// SetFileWithCDN adds a file to the non-fungible token
func (b *NFTMetadataBuilder) SetFileWithCDN(fileURI string, fileType string) *NFTMetadataBuilder {
	if b.files == nil {
		b.files = make([]File, 0)
	}

	if fileType == "" {
		fileType = utils.GetFileTypeByURI(fileURI)
	}
	if fileType == "" {
		fileType = "unknown"
	}

	b.files = append(b.files, File{
		URI:  fileURI,
		Type: fileType,
		CDN:  true,
	})

	return b
}

// Build builds the fungible asset metadata
func (b *NFTMetadataBuilder) Build() (*Metadata, error) {
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
	if b.files == nil {
		b.SetFile(b.image, "")
	}

	props := PropertiesMap{}
	if b.customProperties != nil {
		for k, v := range b.customProperties {
			props[k] = v
		}
	}
	if b.category != "" {
		props["category"] = b.category
	}
	if b.collection != nil {
		props["collection"] = b.collection
	}
	if b.files != nil {
		props["files"] = b.files
	}

	return &Metadata{
		Name:         b.name,
		Symbol:       b.symbol,
		Description:  b.description,
		Image:        b.image,
		AnimationURL: b.animationURL,
		ExternalURL:  b.externalURL,
		Attributes:   b.attributes,
		Properties:   props,
	}, nil
}
