package bitwarden

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/cli/safeexec"
	"github.com/pquerna/otp/totp"
	session "github.com/reubenmiller/c8y-session-bitwarden/pkg/core"
)

type Client struct {
	Folder string
}

func NewClient(folder string) *Client {
	return &Client{
		Folder: folder,
	}
}

// BWItem bitwarden item containing the login information
type BWItem struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Login    BWLogin   `json:"login"`
	Fields   []BWField `json:"fields"`
	FolderID string    `json:"folderId"`
}

func (bwi *BWItem) HasTenantField() bool {
	for _, field := range bwi.Fields {
		name := strings.ToLower(field.Name)
		if strings.Contains(name, "tenant") && strings.TrimSpace(field.Value) != "" {
			return true
		}
	}
	return false
}

func (bwi *BWItem) Skip() bool {
	return len(bwi.Login.Uris) == 0
}

// BWLogin bitwarden login credentials
type BWLogin struct {
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	TOTPSecret string  `json:"totp"`
	Uris       []BWUri `json:"uris"`
}

// BWField bitwarden custom fields
type BWField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  int32  `json:"type"`
}

func (b *BWLogin) MatchesUri(search string) bool {
	for _, uri := range b.Uris {
		if strings.Contains(strings.ToLower(uri.URI), search) {
			return true
		}
	}
	return false
}

// BWUri bitwarden URI associated with the login credentials
type BWUri struct {
	URI string `json:"uri"`
}

func mapToSession(item *BWItem, folders map[string]string) *session.CumulocitySession {

	out := &session.CumulocitySession{
		SessionURI: fmt.Sprintf("bitwarden://%s", item.ID),
		Name:       item.Name,
		Username:   item.Login.Username,
		Password:   item.Login.Password,
		FolderID:   item.FolderID,
		TOTPSecret: item.Login.TOTPSecret,
	}

	// Include folder name (for humans)
	if folderName, found := folders[item.FolderID]; found {
		out.FolderName = folderName
	}

	if len(item.Login.Uris) > 0 {
		out.Host = item.Login.Uris[0].URI
	}

	if len(item.Fields) > 0 {
		for _, field := range item.Fields {
			if strings.HasPrefix(strings.ToLower(field.Name), "tenant") {
				out.Tenant = field.Value
			}
			if strings.EqualFold(field.Name, "type") {
				v, typeErr := session.MarshalSessionType(field.Value)
				if typeErr != nil {
					slog.Error("Unknown session type, so using default instead.", "got", field.Value, "default", v)
				}
				out.Mode = v
			}
		}
	}

	if strings.Contains(item.Login.Username, "/") {
		parts := strings.SplitN(item.Login.Username, "/", 2)
		if len(parts) == 2 {
			if out.Tenant != "" {
				out.Tenant = parts[0]
			}
			out.Username = parts[1]
		}
	}

	return out
}

func isUID(v string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(v)
}

type Folder struct {
	Object string `json:"object"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

func (c *Client) ListFolders(name ...string) (map[string]string, error) {
	folders := make([]Folder, 0)

	args := []string{
		"list",
		"folders",
	}
	if len(name) > 0 {
		args = append(args, "--search", name[0])
	}

	err := c.exec(args, &folders)

	folderMap := make(map[string]string)
	for _, folder := range folders {
		folderMap[folder.ID] = folder.Name
	}

	return folderMap, err
}

func (c *Client) exec(args []string, data any) error {
	if _, err := safeexec.LookPath("bw"); err != nil {
		return err
	}

	if v := os.Getenv("BW_SESSION"); v == "" {
		return fmt.Errorf("bitwarden cli 'BW_SESSION' env variable is not set")
	}

	bw := exec.Command("bw", args...)
	stdout, err := bw.StdoutPipe()
	if err != nil {
		return err
	}
	err = bw.Start()
	if err != nil {
		return err
	}
	parseErr := json.NewDecoder(stdout).Decode(data)
	if parseErr != nil {
		parseErr = fmt.Errorf("failed to parse json output. %w", parseErr)
	}

	// wait for command to finish in background
	go bw.Wait()

	return parseErr
}

func (c *Client) List(name ...string) ([]*session.CumulocitySession, error) {
	cmdArgs := []string{
		"list", "items",
	}

	var folders map[string]string
	var folderErr error

	if c.Folder != "" {
		if isUID(c.Folder) {
			// Filter by folder id (no additional lookup required)
			cmdArgs = append(cmdArgs, "--folderid", c.Folder)
		} else {
			// Filter by folder name/pattern (additional lookup required)
			folders, folderErr = c.ListFolders(c.Folder)
			if folderErr != nil {
				return nil, folderErr
			}
		}
	}

	slog.Debug("Starting", "time", time.Now().Format(time.RFC3339Nano))

	items := make([]BWItem, 0)
	c.exec(cmdArgs, &items)

	sessions := make([]*session.CumulocitySession, 0)
	for _, item := range items {
		if item.Skip() {
			continue
		}

		if len(folders) > 0 {
			if _, found := folders[item.FolderID]; !found {
				continue
			}
		}

		sessions = append(sessions, mapToSession(&item, folders))
	}
	return sessions, nil
}

func GetTOTPCode(secret string, t time.Time) (string, error) {
	if t.Year() == 0 {
		t = time.Now()
	}
	return totp.GenerateCode(secret, t)
}

func GetTOTPCodeFromSecret(secret string) (string, error) {
	now := time.Now()
	totpTime := now
	totpPeriod := 30
	totpNextTransition := totpPeriod - now.Second()%30
	if totpNextTransition < 5 {
		totpTime = now.Add(30 * time.Second)
	}
	return GetTOTPCode(secret, totpTime)
}
