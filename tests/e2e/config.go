package e2e

import (
	"github.com/dmitrymomot/go-env"
	"github.com/dmitrymomot/solana/common"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

var (
	SolanaDevnetRPCNode = env.GetString("SOLANA_RPC_ENDPOINT", "https://api.devnet.solana.com")

	FeePayerPubkey     = common.PublicKeyFromString("71fjb18P3CCaCNgRrUGbMsVQx2vB9XbdhL1BfmRahEPq")
	FeePayerPrivateKey = env.MustString("FEE_PAYER_PRIVATE_KEY")

	Wallet1Pubkey     = common.PublicKeyFromString("FuQhSmAT6kAmmzCMiiYbzFcTQJFuu6raXAdCFibz4YPR")
	Wallet1PrivateKey = env.MustString("WALLET_1_PRIVATE_KEY")

	Wallet2Pubkey     = common.PublicKeyFromString("RjpQLUttBMdoQ4HKMygScEjkd6S69dZZC9T4W3Z3DKD")
	Wallet2PrivateKey = env.MustString("WALLET_2_PRIVATE_KEY")

	TokenMintPubkey         = common.PublicKeyFromString("3GYtjt6Qi93no13nQED5siMMU4fR8zRDPi6V55Vg2mez")
	AssetMintPubkey         = common.PublicKeyFromString("7HyvGUEjxsGJLFX5foWTCWZVpEtabAF5LV63wP2Ei41d")
	MasterEditionMintPubkey = common.PublicKeyFromString("6JPksXSsPNoFqiiyttozoawrFSD7K4pZX31xb6Qw322w")
	EditionMintPubkey2      = common.PublicKeyFromString("FUA9uHMomQUnYRf6BEBsxaEZMwvmiGSXMi3bDnPp1pfu")
	EditionMintPubkey4      = common.PublicKeyFromString("FphT9gWeUQr6gYdD3i8PMneXGHDQyHjyedZY75uAa6Gv")

	CollectionPubkey     = common.PublicKeyFromString("5kyJBiH1ybSnhMniwr2CyaL7LitMKaLGA4HpGZtRcD6e")
	CollectionPrivateKey = env.MustString("COLLECTION_PRIVATE_KEY")

	ArweaveWalletPath = env.GetString("ARWEAVE_WALLET_PATH", "./wallet.json")
	ArweaveClientURL  = "https://arweave.net"
)
