package metadata

import "encoding/json"

// PropertyCategory represents the category of a non-fungible token
// E.g. image, video, audio, vr, html
type PropertyCategory string

// Supported property categories
const (
	PropertyCategoryImage PropertyCategory = "image" // PNG, GIF, JPG
	PropertyCategoryVideo PropertyCategory = "video" // MP4, MOV
	PropertyCategoryAudio PropertyCategory = "audio" // MP3, FLAC, WAV
	PropertyCategoryVr    PropertyCategory = "vr"    // 3D models; GLB, GLTF
	PropertyCategoryHtml  PropertyCategory = "html"  // HTML pages; scripts and relative paths within the HTML page are also supported
)

// Attribute represents a display type of attribute of a non-fungible token
type AttributeDisplayType string

// Predefined nft attribute display types
const (
	AttributeDisplayString   AttributeDisplayType = "string"
	AttributeDisplayNumber   AttributeDisplayType = "number"
	AttributeDisplayBoolean  AttributeDisplayType = "boolean"
	AttributeDisplayDate     AttributeDisplayType = "date"
	AttributeDisplayTime     AttributeDisplayType = "time"
	AttributeDisplayDateTime AttributeDisplayType = "datetime"
)

type (
	// Metadata represents the metadata of a fungible/semi-fungible/non-fungible token
	Metadata struct {
		// The name of the asset.
		Name string `json:"name"`

		// The symbol of the asset.
		Symbol string `json:"symbol"`

		// Human readable description of the asset.
		Description string `json:"description,omitempty"`

		// URL to the image of the asset. PNG, GIF and JPG file formats are supported.
		// You may use the ?ext={file_extension} query to provide information on the file type.
		Image string `json:"image,omitempty"`

		// URL to a multi-media attachment of the asset.
		// The supported file formats are
		// MP4 and MOV for video,
		// MP3, FLAC and WAV for audio,
		// GLB for AR/3D assets,
		// You may use the ?ext={file_extension} query to provide information on the file type.
		AnimationURL string `json:"animation_url,omitempty"`

		// URL to an external application or website where users can also view the asset.
		ExternalURL string `json:"external_url,omitempty"`

		// Attribute represents the attribute of a token
		Attributes []Attribute `json:"attributes,omitempty"`

		// Properties represents the properties of a non-fungible token
		Properties PropertiesMap `json:"properties,omitempty"`
	}

	// Attribute represents the attribute of a non-fungible token
	Attribute struct {
		TraitType string      `json:"trait_type"`
		Value     interface{} `json:"value"`

		// Optional
		DisplayType AttributeDisplayType `json:"display_type,omitempty"` // string, number, boolean, date
		MaxValue    int64                `json:"max_value,omitempty"`
		TraitCount  int64                `json:"trait_count,omitempty"`
	}

	// PropertiesMap ...
	PropertiesMap map[string]interface{}

	// Properties represents the properties of a non-fungible token
	Properties struct {
		Files        []File                 `json:"files"`
		Category     PropertyCategory       `json:"category,omitempty"`
		CustomFields map[string]interface{} `json:"ext,omitempty"`
		Collection   *Collection            `json:"collection,omitempty"`
	}

	// @deprecated
	// Collection represents the collection of a non-fungible token
	// Do not use - may be removed in a future release.
	// Use on-chain data instead.
	Collection struct {
		Name   string `json:"name"`
		Family string `json:"family,omitempty"` // Optional
	}

	// File represents the file of a non-fungible token
	File struct {
		// Mandatory
		URI  string `json:"uri"`
		Type string `json:"type,omitempty"`

		// Optional
		CDN bool `json:"cdn,omitempty"`
	}
)

// ToJSON returns the metadata as a JSON string
func (m Metadata) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// String returns the metadata as a JSON string
func (m Metadata) String() string {
	b, _ := m.ToJSON()
	return string(b)
}
