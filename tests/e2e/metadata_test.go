package e2e_test

import (
	"context"
	"testing"

	"github.com/solplaydev/solana"
	"github.com/solplaydev/solana/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokenMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new client
	client := solana.New(solana.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	t.Run("master edition metadata", func(t *testing.T) {
		md, err := client.GetTokenMetadata(ctx, e2e.MasterEditionMintAddr)
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
		assert.NotNil(t, md.Edition)
		assert.EqualValues(t, md.Edition.Type, solana.EditionMasterEdition)
		assert.EqualValues(t, md.Edition.MaxSupply, uint64(1000))
		assert.GreaterOrEqual(t, md.Edition.Supply, uint64(4))
		assert.EqualValues(t, md.Edition.Edition, uint64(0))
	})

	t.Run("edition metadata #2", func(t *testing.T) {
		md, err := client.GetTokenMetadata(ctx, e2e.EditionMintAddr2)
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
		assert.NotNil(t, md.Edition)
		assert.EqualValues(t, md.Edition.Type, solana.EditionPrintedEdition)
		assert.EqualValues(t, md.Edition.MaxSupply, uint64(1000))
		assert.GreaterOrEqual(t, md.Edition.Supply, uint64(4))
		assert.EqualValues(t, md.Edition.Edition, uint64(2))
	})
}
