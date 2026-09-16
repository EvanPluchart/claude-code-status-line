package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// keychainService is the macOS Keychain entry Claude Code stores its credentials under.
const keychainService = "Claude Code-credentials"

// credentials mirrors the Claude Code credentials file layout.
type credentials struct {
	ClaudeAiOauth struct {
		AccessToken string `json:"accessToken"`
		// ExpiresAt is in milliseconds since epoch; 0 when the field is absent.
		ExpiresAt int64 `json:"expiresAt"`
	} `json:"claudeAiOauth"`
}

// oauthToken returns the Claude Code OAuth access token from the local credential store.
//
// This is the user's own Claude Code session token, read locally and sent only to the
// Anthropic usage endpoint. It is never written to the cache nor logged.
func oauthToken() (string, error) {
	if raw, err := keychainCredentials(); err == nil {
		if token, err := tokenFrom(raw); err == nil {
			return token, nil
		}
	}

	raw, err := os.ReadFile(credentialsPath())
	if err != nil {
		return "", fmt.Errorf("read credentials: %w", err)
	}

	return tokenFrom(raw)
}

// credentialsPath returns the Claude Code credentials file path.
func credentialsPath() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, ".credentials.json")
	}

	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".claude", ".credentials.json")
}

// keychainCredentials reads the credentials blob from the macOS Keychain.
func keychainCredentials() ([]byte, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("keychain: unsupported on %s", runtime.GOOS)
	}

	out, err := exec.Command("security", "find-generic-password", "-s", keychainService, "-w").Output()
	if err != nil {
		return nil, fmt.Errorf("keychain lookup: %w", err)
	}

	return out, nil
}

// tokenFrom extracts a non-expired access token from a credentials blob.
func tokenFrom(raw []byte) (string, error) {
	var creds credentials

	if err := json.Unmarshal(raw, &creds); err != nil {
		return "", fmt.Errorf("parse credentials: %w", err)
	}

	token := creds.ClaudeAiOauth.AccessToken

	if token == "" {
		return "", fmt.Errorf("parse credentials: no access token")
	}

	if expiry := creds.ClaudeAiOauth.ExpiresAt; expiry > 0 && time.UnixMilli(expiry).Before(time.Now()) {
		return "", fmt.Errorf("credentials expired")
	}

	return token, nil
}
