package core

import (
	"fmt"
	"strings"
)

type CumulocitySession struct {
	SessionURI string `json:"sessionUri,omitempty"`
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Tenant     string `json:"tenant,omitempty"`
	TOTP       string `json:"totp,omitempty"`
	TOTPSecret string `json:"totpSecret,omitempty"`

	// Bitwarden specific
	FolderID   string `json:"folderId,omitempty"`
	FolderName string `json:"folderName,omitempty"`
}

func (i CumulocitySession) FilterValue() string {
	return strings.Join([]string{i.SessionURI, i.Host, i.Username}, " ")
}
func (i CumulocitySession) Title() string { return i.Host }
func (i CumulocitySession) Description() string {

	fields := []string{
		"Username=%s",
	}
	args := []any{
		i.Username,
	}

	if i.Tenant != "" {
		fields = append(fields, ", Tenant=%s")
		args = append(args, i.Tenant)
	}

	if i.FolderName != "" {
		fields = append(fields, ", Folder=%s")
		args = append(args, i.FolderName)
	}

	fields = append(fields, " | uri=%s")
	args = append(args, i.SessionURI)

	return fmt.Sprintf(strings.Join(fields, ""), args...)
}
