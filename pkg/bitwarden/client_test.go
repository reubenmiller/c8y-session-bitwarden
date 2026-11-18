package bitwarden

import "testing"

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
