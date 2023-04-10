package e2e_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/dmitrymomot/solana/client"
	"github.com/dmitrymomot/solana/tests/e2e"
	"github.com/stretchr/testify/require"
)

func TestGetTokensList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new sc
	sc := client.New(client.SetSolanaEndpoint(e2e.SolanaDevnetRPCNode))

	// Get tokens list
	tokensList, err := sc.GetFungibleTokensList(ctx, e2e.Wallet1Pubkey.ToBase58())
	require.NoError(t, err)
	require.NotEmpty(t, tokensList)
	// utils.PrettyPrint(tokensList)
	// fmt.Println("total:", len(tokensList))

	// get tokens metadata
	for _, v := range tokensList {
		metadata, err := sc.GetFungibleTokenMetadata(ctx, v.Mint.ToBase58())
		if err != nil {
			fmt.Println("error:", v.Mint.ToBase58(), "err:", err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, metadata)
		// utils.PrettyPrint(metadata)
	}
}
