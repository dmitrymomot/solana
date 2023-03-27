package metadata

// @deprecated
// This is a temporary solution to support the deprecated metadata format.
type (
	TokenList struct {
		Name     string                  `json:"name"`
		LogoURI  string                  `json:"logoURI"`
		Keywords []string                `json:"keywords"`
		Tags     map[string]TokenListTag `json:"tags"`
		Tokens   []TokenListToken        `json:"tokens"`
	}

	TokenListTag struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	TokenListToken struct {
		ChainID    int                    `json:"chainId"`
		Address    string                 `json:"address"`
		Symbol     string                 `json:"symbol"`
		Name       string                 `json:"name"`
		Decimals   int                    `json:"decimals"`
		LogoURI    string                 `json:"logoURI"`
		Tags       []string               `json:"tags,omitempty"`
		Extensions map[string]interface{} `json:"extensions,omitempty"`
	}
)

// Token list chain IDs
const (
	ChainIdMainnet = 101 // Mainnet-beta
	ChainIdTestnet = 102 // Testnet
)
