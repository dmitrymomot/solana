package e2e_test

import (
	"context"
	"testing"

	"github.com/solplaydev/solana"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mintAddrDevnet = "2HRNabptCW4eMnCAS5xPiiBQTj5BPDmbENc7JpvwvHJA"

func TestGetTokenMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	md, err := client.GetTokenMetadata(ctx, mintAddrDevnet)
	require.NoError(t, err)
	require.NotEmpty(t, md)
	assert.NotEmpty(t, md.MetadataUri)
	assert.NotEmpty(t, md.SellerFeeBasisPoints)
	assert.GreaterOrEqual(t, len(md.Creators), 1)
	if assert.NotNil(t, md.Data) {
		assert.NotEmpty(t, md.Data.Name)
		assert.NotEmpty(t, md.Data.Symbol)
		assert.NotEmpty(t, md.Data.Description)
		assert.NotEmpty(t, md.Data.Image)
	}
}
