package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/solplaydev/solana/utils"
)

// Supported property categories
const (
	PropertyCategoryImage PropertyCategory = "image" // PNG, GIF, JPG
	PropertyCategoryVideo PropertyCategory = "video" // MP4, MOV
	PropertyCategoryAudio PropertyCategory = "audio" // MP3, FLAC, WAV
	PropertyCategoryVr    PropertyCategory = "vr"    // 3D models; GLB, GLTF
	PropertyCategoryHtml  PropertyCategory = "html"  // HTML pages; scripts and relative paths within the HTML page are also supported
)

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
		Properties *PropertiesMap `json:"properties,omitempty"`
	}

	// Attribute represents a display type of attribute of a non-fungible token
	AttributeDisplayType string

	// Attribute represents the attribute of a non-fungible token
	Attribute struct {
		TraitType string      `json:"trait_type"`
		Value     interface{} `json:"value"`

		// Optional
		DisplayType AttributeDisplayType `json:"display_type,omitempty"` // string, number, boolean, date
		MaxValue    int64                `json:"max_value,omitempty"`
		TraitCount  int64                `json:"trait_count,omitempty"`
	}

	// PropertyCategory represents the category of a non-fungible token
	// E.g. image, video, audio, vr, html
	PropertyCategory string

	// PropertiesMap ...
	PropertiesMap map[string]interface{}

	// Properties represents the properties of a non-fungible token
	Properties struct {
		Files        []File                 `json:"files"`
		Category     PropertyCategory       `json:"category,omitempty"`
		CustomFields map[string]interface{} `json:"ext,omitempty"`
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

// FungibleMetadataParams represents the parameters of a fungible token
type FungibleMetadataParams struct {
	Name        string // required
	Symbol      string // required
	Description string // required
	Image       string // required
	ExternalURL string // optional
}

// NewFungibleMetadata returns the metadata of a fungible token
func NewFungibleMetadata(params FungibleMetadataParams) (Metadata, error) {
	if params.Name == "" || params.Symbol == "" || params.Description == "" || params.Image == "" {
		return Metadata{}, utils.StackErrors(
			ErrMetaRequiredFields,
			errors.New("name, symbol, description and image are required"),
		)
	}

	return Metadata{
		Name:        params.Name,
		Symbol:      params.Symbol,
		Description: params.Description,
		Image:       params.Image,
		ExternalURL: params.ExternalURL,
	}, nil
}

// FungibleAssetMetadataParams represents the parameters of a semi-fungible token (fungible asset)
type FungibleAssetMetadataParams struct {
	Name         string      // required
	Symbol       string      // required
	Description  string      // required
	Image        string      // required
	AnimationURL string      // optional
	ExternalURL  string      // optional
	Attributes   []Attribute // optional
}

// NewFungibleAssetMetadata returns the metadata of a semi-fungible token (fungible asset)
func NewFungibleAssetMetadata(params FungibleAssetMetadataParams) (Metadata, error) {
	if params.Name == "" || params.Symbol == "" || params.Description == "" || params.Image == "" {
		return Metadata{}, utils.StackErrors(
			ErrMetaRequiredFields,
			errors.New("name, symbol, description and image are required"),
		)
	}

	return Metadata{
		Name:         params.Name,
		Symbol:       params.Symbol,
		Description:  params.Description,
		Image:        params.Image,
		AnimationURL: params.AnimationURL,
		ExternalURL:  params.ExternalURL,
		Attributes:   params.Attributes,
	}, nil
}

// NonFungibleMetadataParams represents the parameters of a non-fungible token
type NonFungibleMetadataParams struct {
	Name             string                 // required
	Symbol           string                 // required
	Description      string                 // required
	Image            string                 // required
	AnimationURL     string                 // optional
	ExternalURL      string                 // optional
	Category         PropertyCategory       // optional
	Attributes       []Attribute            // optional
	Files            []File                 // optional
	CustomProperties map[string]interface{} // optional; any properties that you want to add to the metadata
}

// NewNonFungibleMetadata returns the metadata of a non-fungible token
func NewNonFungibleMetadata(params NonFungibleMetadataParams) (Metadata, error) {
	if params.Name == "" || params.Symbol == "" || params.Description == "" || params.Image == "" {
		return Metadata{}, utils.StackErrors(
			ErrMetaRequiredFields,
			errors.New("name, symbol, description and image are required"),
		)
	}

	files := params.Files
	if len(files) == 0 {
		files = make([]File, 0)
	}
	if params.Image != "" && len(files) == 0 {
		if !isFileExistsInArray(params.Image, files) {
			files = append(files, NewFile(params.Image, ""))
		}
	}
	if params.AnimationURL != "" && len(files) == 0 {
		if !isFileExistsInArray(params.AnimationURL, files) {
			files = append(files, NewFile(params.AnimationURL, ""))
		}
	}

	if params.Category == "" {
		params.Category = PropertyCategoryImage
	}

	properties := PropertiesMap{}
	if params.CustomProperties != nil {
		properties = params.CustomProperties
	}
	properties["files"] = files
	properties["category"] = params.Category

	return Metadata{
		Name:         params.Name,
		Symbol:       params.Symbol,
		Description:  params.Description,
		Image:        params.Image,
		AnimationURL: params.AnimationURL,
		ExternalURL:  params.ExternalURL,
		Attributes:   params.Attributes,
		Properties:   &properties,
	}, nil
}

// NewAttribute returns a new attribute
func NewAttribute(key string, value interface{}) Attribute {
	attr := Attribute{
		TraitType: key,
		Value:     fmt.Sprintf("%v", value),
	}

	valueType := utils.GetVarType(value)
	switch valueType {
	case "string", "byte", "rune":
		attr.DisplayType = AttributeDisplayString
	case "int", "float", "int64", "float64", "uint", "uint64", "int32", "float32", "uint32", "int16", "uint16", "int8", "uint8":
		attr.DisplayType = AttributeDisplayNumber
	case "bool", "boolean":
		attr.DisplayType = AttributeDisplayBoolean
	default:
		attr.DisplayType = AttributeDisplayString
	}

	return attr
}

// NewFile returns a new file for a non-fungible token
func NewFile(uri string, fileType string) File {
	if fileType == "" {
		fileType = utils.GetFileTypeByURI(uri)
	}
	return File{
		URI:  uri,
		Type: fileType,
	}
}

// NewFileWithCDN returns a new file for a non-fungible token with CDN enabled
func NewFileWithCDN(uri string, fileType string) File {
	if fileType == "" {
		fileType = utils.GetFileTypeByURI(uri)
	}
	return File{
		URI:  uri,
		Type: fileType,
		CDN:  true,
	}
}

// IsFileExistsInArray returns true if the file exists in the metadata files array
func isFileExistsInArray(url string, files []File) bool {
	for _, f := range files {
		if f.URI == url {
			return true
		}
	}

	return false
}

// TokenMetadata represents the metadata of a token
type (
	TokenMetadata struct {
		UpdateAuthority      string      `json:"update_authority"`
		Mint                 string      `json:"mint"`
		PrimarySaleHappened  bool        `json:"primary_sale_happened,omitempty"`
		IsMutable            bool        `json:"is_mutable"`
		EditionNonce         *uint8      `json:"edition_nonce,omitempty"`
		TokenStandard        string      `json:"token_standard"`
		Collection           *Collection `json:"collection,omitempty"`
		Uses                 *Uses       `json:"uses,omitempty"`
		Edition              *Edition    `json:"edition,omitempty"`
		MetadataUri          string      `json:"metadata_uri,omitempty"`
		SellerFeeBasisPoints uint16      `json:"seller_fee_basis_points,omitempty"`
		Creators             []Creator   `json:"creators,omitempty"`
		Data                 *Metadata   `json:"data,omitempty"`
	}

	Edition struct {
		Key       string  `json:"key"`
		Supply    uint64  `json:"supply"`
		MaxSupply *uint64 `json:"max_supply"`
	}

	Collection struct {
		Verified bool   `json:"verified"`
		Key      string `json:"key"`
		Size     uint64 `json:"size,omitempty"`
	}

	Uses struct {
		UseMethod string `json:"use_method"`
		Remaining uint64 `json:"remaining"`
		Total     uint64 `json:"total"`
	}

	Creator struct {
		Address  string `json:"address"`
		Verified bool   `json:"verified"`
		Share    uint8  `json:"share"`
	}
)

// GetTokenMetadata returns the metadata of a token
func (c *Client) GetTokenMetadata(ctx context.Context, base58MintAddr string) (TokenMetadata, error) {
	if base58MintAddr == "" {
		return TokenMetadata{}, utils.StackErrors(
			ErrInvalidPublicKey,
			errors.New("mint address is required"),
		)
	}

	mint := common.PublicKeyFromString(base58MintAddr)
	metadataAccount, err := token_metadata.GetTokenMetaPubkey(mint)
	if err != nil {
		return TokenMetadata{}, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	metadataAccountInfo, err := c.solana.GetAccountInfo(ctx, metadataAccount.ToBase58())
	if err != nil {
		return TokenMetadata{}, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	metadata, err := token_metadata.MetadataDeserialize(metadataAccountInfo.Data)
	if err != nil {
		return TokenMetadata{}, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	m := TokenMetadata{
		UpdateAuthority:      metadata.UpdateAuthority.ToBase58(),
		Mint:                 metadata.Mint.ToBase58(),
		PrimarySaleHappened:  metadata.PrimarySaleHappened,
		IsMutable:            metadata.IsMutable,
		EditionNonce:         metadata.EditionNonce,
		TokenStandard:        TokenStandardToString(*metadata.TokenStandard),
		MetadataUri:          metadata.Data.Uri,
		SellerFeeBasisPoints: metadata.Data.SellerFeeBasisPoints,
	}

	m.Edition, _ = c.GetMasterEditionInfo(ctx, base58MintAddr)

	m.Data, err = c.DownloadMetadata(ctx, metadata.Data.Uri)
	if err != nil {
		return TokenMetadata{}, utils.StackErrors(
			ErrGetTokenMetadata,
			err,
		)
	}

	if metadata.Data.Creators != nil {
		if m.Creators == nil {
			m.Creators = make([]Creator, 0, len(*metadata.Data.Creators))
		}
		for _, creator := range *metadata.Data.Creators {
			m.Creators = append(m.Creators, Creator{
				Address:  creator.Address.ToBase58(),
				Verified: creator.Verified,
				Share:    creator.Share,
			})
		}
	}

	if metadata.Collection != nil {
		m.Collection = &Collection{
			Verified: metadata.Collection.Verified,
			Key:      metadata.Collection.Key.ToBase58(),
		}
		if metadata.CollectionDetails != nil {
			m.Collection.Size = metadata.CollectionDetails.V1.Size
		}
	}

	if metadata.Uses != nil {
		m.Uses = &Uses{
			UseMethod: UseMethodToString(metadata.Uses.UseMethod),
			Remaining: metadata.Uses.Remaining,
			Total:     metadata.Uses.Total,
		}
	}

	return m, nil
}

// NewMetadataFromJSON returns the metadata from a JSON encoded string
func NewMetadataFromJSON(jsonData []byte) (*Metadata, error) {
	metadata := &Metadata{}
	if err := json.Unmarshal(jsonData, metadata); err != nil {
		return nil, utils.StackErrors(ErrMetaUnmarshal, err)
	}

	return metadata, nil
}

// DownloadMetadata downloads the metadata from the given url
func (c *Client) DownloadMetadata(ctx context.Context, url string) (*Metadata, error) {
	if url == "" {
		return nil, nil
	}

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, utils.StackErrors(
			ErrDownloadMetadata,
			err,
		)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.StackErrors(
			ErrDownloadMetadata,
			err,
		)
	}

	m := &Metadata{}
	if err := json.Unmarshal(body, m); err != nil {
		return nil, utils.StackErrors(
			ErrDownloadMetadata,
			ErrMetaUnmarshal,
			err,
		)
	}

	return m, nil
}

func (c *Client) GetMasterEditionInfo(ctx context.Context, base58MintAddr string) (*Edition, error) {
	mint := common.PublicKeyFromString(base58MintAddr)
	masterEditionPubKey, err := token_metadata.GetMasterEdition(mint)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetMasterEditionInfo,
			err,
		)
	}

	masterEdition, err := c.solana.GetAccountInfo(ctx, masterEditionPubKey.String())
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetMasterEditionInfo,
			err,
		)
	}

	data, err := masterEditionDeserialize(masterEdition.Data)
	if err != nil {
		return nil, utils.StackErrors(
			ErrGetMasterEditionInfo,
			err,
		)
	}

	return &Edition{
		Key:       EditionKeyToString(data.Key),
		Supply:    data.Supply,
		MaxSupply: data.MaxSupply,
	}, nil
}

func masterEditionDeserialize(data []byte) (token_metadata.MasterEditionV2, error) {
	var masterEdition token_metadata.MasterEditionV2
	err := borsh.Deserialize(&masterEdition, data)
	if err != nil {
		return token_metadata.MasterEditionV2{}, utils.StackErrors(
			ErrMasterEditionDeserialize,
			err,
		)
	}

	return masterEdition, nil
}
