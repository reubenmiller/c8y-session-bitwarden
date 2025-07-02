/*
Copyright © 2024 Reuben Miller
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/reubenmiller/c8y-session-bitwarden/pkg/bitwarden"
	"github.com/reubenmiller/c8y-session-bitwarden/pkg/core/picker"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [SEARCH_TERMS]",
	Short: "List sessions stored in your bitwarden vault",
	Long: heredoc.Doc(`
		List Cumulocity IoT sessions from your bitwarden vault

		Notes:
			* the first SEARCH_TERM is used sent to the bw cli command as the '--search <term>'
			  argument, and then a client-side filter is used to apply the additional terms. This is
			  due to the limitation that the bw commands only supports one search term

			* If only 1 match if found, then the session will be selected automatically

		Examples
			c8y-session-bitwarden list --folder c8y
			# Select items from the c8y folder

			c8y-session-bitwarden list --folder c8y example.com
			# Select items from the c8y folder, and match the search term, "example.com"

			c8y-session-bitwarden list --folder c8y example.com dev
			# Select items from the c8y folder, and match the search term, "example.com" AND "dev"
	`),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		folder, err := cmd.Flags().GetString("folder")
		if err != nil {
			return err
		}
		client := bitwarden.NewClient(folder)
		sessions, err := client.List(args...)
		if err != nil {
			return err
		}

		session, err := picker.Pick(sessions, picker.PickerOptions{
			AutoSelectIfOnlyOne: true,
		})
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
	listCmd.Flags().String("loginType", "", "Not used. Kept to satisfy the go-c8y-cli session interface")
	listCmd.Flags().Bool("clear", false, "Not used. Kept to satisfy the go-c8y-cli session interface")
}
