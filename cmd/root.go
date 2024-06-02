/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "c8y-session-bitwarden",
	Short: "go-c8y-cli bitwarden session selector",
	Long: `Select a session from your bitwarden password manager

Pre-requisites:

 * bitwarden-cli (bw) - https://github.com/bitwarden/clients

Login to your bitwarden account from the command line

	$ bw login

Then export the bitwarden session variable (as suggested in the command's output)
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
