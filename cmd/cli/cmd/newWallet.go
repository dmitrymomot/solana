/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/dmitrymomot/solana/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newWalletCmd represents the newWallet command
var newWalletCmd = &cobra.Command{
	Use:   "new-wallet",
	Short: "Generate a new wallet",
	Long:  "Generate a new wallet with a mnemonic and a private key.",

	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.BgGreen, color.FgHiWhite).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		mnemonic, err := common.NewMnemonic(common.MnemonicLength12)
		if err != nil {
			panic(err)
		}

		wallet, err := common.DeriveAccountFromMnemonicBip44(mnemonic)
		if err != nil {
			panic(err)
		}

		fmt.Println("\n" + green("*** Client created successfully! ***"))
		fmt.Println(cyan("Mnemonic: ") + mnemonic)
		fmt.Println(cyan("Wallet address: ") + wallet.PublicKey.ToBase58())
		fmt.Println(cyan("Wallet private key: ") + common.AccountToBase58(wallet))
	},
}

func init() {
	rootCmd.AddCommand(newWalletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newWalletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newWalletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
