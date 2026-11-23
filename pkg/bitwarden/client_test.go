package bitwarden

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMapToSessionTenant(t *testing.T) {
	cases := []struct {
		name         string
		username     string
		fields       []BWField
		wantTenant   string
		wantUsername string
	}{
		{name: "tenant from username prefix", username: "t100/demo-user", wantTenant: "t100", wantUsername: "demo-user"},
		{name: "tenant field takes precedence", username: "t100/demo-user", fields: []BWField{{Name: "tenant", Value: "t200"}}, wantTenant: "t200", wantUsername: "demo-user"},
		{name: "tenant field without prefix", username: "demo-user", fields: []BWField{{Name: "tenant", Value: "t200"}}, wantTenant: "t200", wantUsername: "demo-user"},
		{name: "no tenant", username: "demo-user", wantTenant: "", wantUsername: "demo-user"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := &BWItem{
				ID:     "item-1",
				Login:  BWLogin{Username: tc.username, Uris: []BWUri{{URI: "https://demo.example.com"}}},
				Fields: tc.fields,
			}
			got := mapToSession(item, map[string]string{})
			if got.Tenant != tc.wantTenant {
				t.Errorf("tenant: got %q, want %q", got.Tenant, tc.wantTenant)
			}
			if got.Username != tc.wantUsername {
				t.Errorf("username: got %q, want %q", got.Username, tc.wantUsername)
			}
		})
	}
}

func TestItemID(t *testing.T) {
	cases := map[string]string{
		"bitwarden://3f0b9d4e-1c2a-4f5b-9e8d-7a6b5c4d3e2f": "3f0b9d4e-1c2a-4f5b-9e8d-7a6b5c4d3e2f",
		"3f0b9d4e-1c2a-4f5b-9e8d-7a6b5c4d3e2f":             "3f0b9d4e-1c2a-4f5b-9e8d-7a6b5c4d3e2f",
	}
	for input, want := range cases {
		got, err := ItemID(input)
		if err != nil || got != want {
			t.Errorf("ItemID(%q) = %q, %v; want %q", input, got, err, want)
		}
	}
	for _, input := range []string{"", "bitwarden://", "onepassword://vault/item"} {
		if _, err := ItemID(input); err == nil {
			t.Errorf("ItemID(%q) should fail", input)
		}
	}
}

// fakeBW puts a fake bw cli on the PATH which prints the given output (or fails with the
// given message), and records its arguments
func fakeBW(t *testing.T, stdout string, failure string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("uses a shell script as the bw cli")
	}
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	script := "#!/bin/sh\necho \"$*\" > " + argsFile + "\n"
	if failure != "" {
		script += "echo '" + failure + "' >&2\nexit 1\n"
	} else {
		script += "cat <<'JSON'\n" + stdout + "\nJSON\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "bw"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BW_SESSION", "test-session")
	return argsFile
}

func TestGet(t *testing.T) {
	argsFile := fakeBW(t, `{"id":"item-1","name":"Demo","login":{"username":"t100/demo","password":"secret","totp":"","uris":[{"uri":"https://demo.example.com"}]},"fields":[{"name":"mode","value":"dev"}]}`, "")
	s, err := NewClient("").Get("bitwarden://item-1")
	if err != nil {
		t.Fatal(err)
	}
	if s.Host != "https://demo.example.com" || s.Tenant != "t100" || s.Username != "demo" || s.Password != "secret" || s.Mode != "dev" || s.SessionURI != "bitwarden://item-1" {
		t.Errorf("unexpected session: %+v", s)
	}
	args, _ := os.ReadFile(argsFile)
	if got := strings.TrimSpace(string(args)); got != "--nointeraction get item item-1" {
		t.Errorf("unexpected bw arguments: %s", got)
	}
}

func TestGetNotASession(t *testing.T) {
	fakeBW(t, `{"id":"item-1","name":"Note","login":{"uris":[]}}`, "")
	if _, err := NewClient("").Get("bitwarden://item-1"); err == nil || !strings.Contains(err.Error(), "not a Cumulocity session") {
		t.Errorf("expected an error, got %v", err)
	}
}

func TestBWErrorIsReturned(t *testing.T) {
	fakeBW(t, "", "Vault is locked.")
	_, err := NewClient("").List()
	if err == nil || !strings.Contains(err.Error(), "Vault is locked.") {
		t.Errorf("expected the bw error, got %v", err)
	}
	_, err = NewClient("").Get("bitwarden://item-1")
	if err == nil || !strings.Contains(err.Error(), "Vault is locked.") {
		t.Errorf("expected the bw error, got %v", err)
	}
}
