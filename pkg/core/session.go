package core

import (
	"fmt"
	"strings"
)

var TypeDev = "dev"
var TypeQual = "qual"
var TypeProduction = "prod"

func MarshalSessionType(v string) (string, error) {
	switch v {
	case "dev":
		return TypeDev, nil
	case "qual", "test":
		return TypeQual, nil
	case "prod", "production":
		return TypeProduction, nil
	default:
		return TypeProduction, fmt.Errorf("unknown session type: %s", v)
	}
}

type CumulocitySession struct {
	SessionURI string `json:"sessionUri,omitempty"`
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Tenant     string `json:"tenant,omitempty"`
	TOTP       string `json:"totp,omitempty"`
	TOTPSecret string `json:"totpSecret,omitempty"`
	Mode       string `json:"mode,omitempty"`

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

	if i.Mode != "" {
		fields = append(fields, ", mode=%s")
		args = append(args, i.Mode)
	}

	fields = append(fields, " | uri=%s")
	args = append(args, i.SessionURI)

	return fmt.Sprintf(strings.Join(fields, ""), args...)
}
