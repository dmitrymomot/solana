package solana

import (
	"github.com/portto/solana-go-sdk/client"
)

type (
	// Solana client wrapper
	Client struct {
		solana          *client.Client
		defaultDecimals uint8
	}

	ClientOption func(*Client)
)

// WithCustomSolanaClient sets a custom solana client
func WithCustomSolanaClient(solana *client.Client) ClientOption {
	return func(c *Client) {
		if c.solana != nil {
			panic("solana client is already set")
		}
		c.solana = solana
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
		if c.solana != nil {
			panic("solana client is already set")
		}
		c.solana = client.NewClient(endpoint)
	}
}

// NewClient creates a new client
// endpoint is the endpoint of the solana RPC node
// cnf is the configuration for the client
func New(opts ...ClientOption) *Client {
	c := &Client{defaultDecimals: 9}

	for _, opt := range opts {
		opt(c)
	}

	if c.solana == nil {
		panic("missing solana client")
	}

	return c
}

// Solana returns the solana client
func (c *Client) Solana() *client.Client {
	return c.solana
}

// DefaultDecimals returns the default decimals
func (c *Client) DefaultDecimals() uint8 {
	return c.defaultDecimals
}
