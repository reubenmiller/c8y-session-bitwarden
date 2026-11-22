/*
Copyright © 2024 Reuben Miller
*/
package cmd

import (
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
			set-session --from-extension bitwarden -- --folder c8y
			# Select a session using go-c8y-cli (requires the extension to be installed)

			c8y-session-bitwarden list --folder c8y
			# Select items from the c8y folder

			c8y-session-bitwarden list --folder c8y example.com
			# Select items from the c8y folder, and match the search term, "example.com"

			c8y-session-bitwarden list --folder c8y example.com dev
			# Select items from the c8y folder, and match the search term, "example.com" AND "dev"

			c8y-session-bitwarden list --folder c8y --json
			# Print the sessions as json (without secrets), e.g. for the session picker of "c8y mcp serve"
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

		if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
			return printSessionSummaries(cmd.OutOrStdout(), sessions)
		}

		session, err := picker.Pick(sessions, picker.PickerOptions{
			AutoSelectIfOnlyOne: true,
		})
		if err != nil {
			return err
		}

		// The picker only returns a copy of the session without secrets
		for _, s := range sessions {
			if session.SessionURI == s.SessionURI {
				return printSession(s)
			}
		}
		return fmt.Errorf("the selected session was not found. sessionUri=%s", session.SessionURI)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("folder", "c8y", "Folder")
	listCmd.Flags().Bool("json", false, "Print the sessions as a json array (without secrets) instead of selecting one")
	addSessionInterfaceFlags(listCmd)
}
