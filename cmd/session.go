package cmd

import (
	"encoding/json"
	"io"

	"github.com/reubenmiller/c8y-session-bitwarden/pkg/bitwarden"
	"github.com/reubenmiller/c8y-session-bitwarden/pkg/core"
	"github.com/reubenmiller/go-c8y-cli/pkg/extsdk"
	"github.com/spf13/cobra"
)

// sessionSummary is a session listed for go-c8y-cli, without secrets (the same json format as
// extsdk.SessionSummary in newer versions of the go-c8y-cli extsdk)
type sessionSummary struct {
	SessionUri string `json:"sessionUri"`
	Name       string `json:"name,omitempty"`
	Host       string `json:"host"`
	Tenant     string `json:"tenant,omitempty"`
	Username   string `json:"username,omitempty"`
	Mode       string `json:"mode,omitempty"`
	LoginType  string `json:"loginType,omitempty"`
	Group      string `json:"group,omitempty"`
}

// printSessionSummaries writes the sessions as a json array, without secrets
func printSessionSummaries(w io.Writer, sessions []*core.CumulocitySession) error {
	out := make([]sessionSummary, 0, len(sessions))
	for _, s := range sessions {
		out = append(out, sessionSummary{
			SessionUri: s.SessionURI,
			Name:       s.Name,
			Host:       s.Host,
			Tenant:     s.Tenant,
			Username:   s.Username,
			Mode:       s.Mode,
			LoginType:  s.LoginType,
			Group:      s.FolderName,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(out)
}

// printSession writes the session with its password, and the current TOTP code (if the item
// has a TOTP secret), in the go-c8y-cli session format
func printSession(s *core.CumulocitySession) error {
	if s.TOTPSecret != "" {
		if totp, err := bitwarden.GetTOTPCodeFromSecret(s.TOTPSecret); err == nil {
			s.TOTP = totp
		}
	}
	return extsdk.PrintSession(extsdk.Session{
		Host:        s.Host,
		Tenant:      s.Tenant,
		Username:    s.Username,
		Password:    s.Password,
		TOTP:        s.TOTP,
		Mode:        s.Mode,
		LoginType:   s.LoginType,
		Description: s.Name,
		SessionUri:  s.SessionURI,
	})
}

// addSessionInterfaceFlags adds hidden flags which are only there to satisfy the go-c8y-cli
// session provider interface. go-c8y-cli passes these flags to the provider when they are set
// by the user
func addSessionInterfaceFlags(cmd *cobra.Command) {
	cmd.Flags().String("loginType", "", "Not used. Kept to satisfy the go-c8y-cli session interface")
	cmd.Flags().Bool("clear", false, "Not used. Kept to satisfy the go-c8y-cli session interface")
	cmd.Flags().String("browserCallback", "", "Not used. Kept to satisfy the go-c8y-cli session interface")
	cmd.Flags().String("sso-scopes", "", "Not used. Kept to satisfy the go-c8y-cli session interface")
	cmd.Flags().String("sso-audience", "", "Not used. Kept to satisfy the go-c8y-cli session interface")
	for _, name := range []string{"loginType", "clear", "browserCallback", "sso-scopes", "sso-audience"} {
		_ = cmd.Flags().MarkHidden(name)
	}
}
