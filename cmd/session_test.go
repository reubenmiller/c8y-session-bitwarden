package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/reubenmiller/c8y-session-bitwarden/pkg/core"
)

func TestPrintSessionSummariesHasNoSecrets(t *testing.T) {
	buf := &bytes.Buffer{}
	err := printSessionSummaries(buf, []*core.CumulocitySession{{
		SessionURI: "bitwarden://item-1",
		Name:       "Demo",
		Host:       "https://demo.example.com",
		Username:   "demo",
		Password:   "secret-password",
		TOTPSecret: "secret-totp",
		FolderName: "c8y",
		Mode:       "dev",
	}})
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	want := `[{"sessionUri":"bitwarden://item-1","name":"Demo","host":"https://demo.example.com","username":"demo","mode":"dev","group":"c8y"}]`
	if got != want {
		t.Errorf("got %s\nwant %s", got, want)
	}

	buf.Reset()
	if err := printSessionSummaries(buf, nil); err != nil || strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("expected an empty array, got %q (%v)", buf.String(), err)
	}
}
