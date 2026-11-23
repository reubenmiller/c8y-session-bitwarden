package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/reubenmiller/c8y-session-bitwarden/pkg/bitwarden"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <SESSION_URI>",
	Short: "Get a session from your bitwarden vault by its sessionUri",
	Long: heredoc.Doc(`
		Get a Cumulocity session from your bitwarden vault, without user interaction

		The session is printed in the go-c8y-cli session format, including the password and the
		TOTP code (if the item has a TOTP secret). go-c8y-cli uses this command to log in to a
		session which was selected from "list --json", e.g. in the session picker of "c8y mcp serve".

		Examples
			c8y sessions login --from-extension bitwarden --select bitwarden://3f0b9d4e-1c2a-4f5b-9e8d-7a6b5c4d3e2f
			# Log in to a session using go-c8y-cli (requires the extension to be installed)

			c8y-session-bitwarden get bitwarden://3f0b9d4e-1c2a-4f5b-9e8d-7a6b5c4d3e2f
			# Print the session of the bitwarden item
	`),
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := bitwarden.NewClient("").Get(args[0])
		if err != nil {
			return err
		}
		return printSession(session)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	addSessionInterfaceFlags(getCmd)
}
