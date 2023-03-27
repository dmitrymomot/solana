package client

import (
	"net/http"

	"github.com/dmitrymomot/solana/types"
	"github.com/portto/solana-go-sdk/client"
)

type (
	// Solana client wrapper
	Client struct {
		rpcClient       *client.Client
		http            *http.Client
		defaultDecimals uint8
		tokenListPath   string
	}

	ClientOption func(*Client)
)

// WithCustomSolanaClient sets a custom solana client
func WithCustomSolanaClient(solana *client.Client) ClientOption {
	return func(c *Client) {
		if c.rpcClient != nil {
			panic("solana client is already set")
		}
		c.rpcClient = solana
	}
}

// WithCustomDecimals sets the custom default decimals
func WithCustomDecimals(decimals uint8) ClientOption {
	return func(c *Client) {
		c.defaultDecimals = decimals
	}
}

// SetSolanaEndpoint sets the solana endpoint
func SetSolanaEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		if c.rpcClient != nil {
			panic("solana client is already set")
		}
		c.rpcClient = client.NewClient(endpoint)
	}
}

// SetHTTPClient sets the http client
func SetHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if c.http != nil {
			panic("http client is already set")
		}
		c.http = httpClient
	}
}

// SetTokenListPath sets the token list path
func SetTokenListPath(path string) ClientOption {
	return func(c *Client) {
		if c.tokenListPath != "" {
			panic("token list path is already set")
		}
		c.tokenListPath = path
	}
}

// NewClient creates a new client
// endpoint is the endpoint of the solana RPC node
// cnf is the configuration for the client
func New(opts ...ClientOption) *Client {
	c := &Client{defaultDecimals: types.SPLTokenDefaultDecimals}

	for _, opt := range opts {
		opt(c)
	}

	if c.rpcClient == nil {
		panic("missing solana client")
	}

	if c.http == nil {
		c.http = http.DefaultClient
	}

	return c
}

// Solana returns the solana client
func (c *Client) Solana() *client.Client {
	return c.rpcClient
}

// DefaultDecimals returns the default decimals
func (c *Client) DefaultDecimals() uint8 {
	return c.defaultDecimals
}
