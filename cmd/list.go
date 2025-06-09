/*
Copyright © 2024 Reuben Miller
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/reubenmiller/c8y-session-bitwarden/pkg/bitwarden"
	"github.com/reubenmiller/c8y-session-bitwarden/pkg/core/picker"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List sessions stored in your bitwarden vault",
	Long:         `List Cumulocity IoT sessions from your bitwarden vault`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		folder, err := cmd.Flags().GetString("folder")
		if err != nil {
			return err
		}
		client := bitwarden.NewClient(folder)
		sessions, err := client.List()
		if err != nil {
			return err
		}

		session, err := picker.Pick(sessions)
		if err != nil {
			return err
		}

		// Check if TOTP secret is present and calc next code
		for _, s := range sessions {
			if session.SessionURI == s.SessionURI {
				session.Password = s.Password
				if s.TOTPSecret != "" {
					totp, toptErr := bitwarden.GetTOTPCodeFromSecret(s.TOTPSecret)
					if toptErr == nil {
						session.TOTP = totp
					}
					break
				}
			}
		}

		out, err := json.MarshalIndent(session, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out)
		return err
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("folder", "c8y", "Folder")
}
