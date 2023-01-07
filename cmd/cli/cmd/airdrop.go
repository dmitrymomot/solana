/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strconv"

	"github.com/fatih/color"
	"github.com/solplaydev/solana"
	"github.com/spf13/cobra"
)

// airdropCmd represents the airdrop command
var airdropCmd = &cobra.Command{
	Use:   "airdrop",
	Short: "Request airdrop",
	Long:  "Request airdrop to a wallet address.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			color.Red("Missing wallet address.")
			return
		}

		if len(args) > 2 {
			color.Red("Too many arguments.")
			return
		}

		addr := args[0]
		if err := solana.ValidateSolanaWalletAddr(addr); err != nil {
			color.Red("Invalid wallet address.")
			return
		}

		amount := solana.SOL
		if len(args) == 2 {
			var err error
			amount, err = strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				color.Red("Invalid amount.")
				return
			}
		}

		client := solana.New(solana.SetSolanaEndpoint(solana.SolanaDevnetRPCURL))
		tx, err := client.RequestAirdrop(cmd.Context(), addr, amount)
		if err != nil {
			color.Red(err.Error())
			return
		}

		color.Green("Airdrop requested successfully!")
		color.Cyan("Transaction hash: %s", tx)
	},
}

func init() {
	rootCmd.AddCommand(airdropCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// airdropCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// airdropCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
