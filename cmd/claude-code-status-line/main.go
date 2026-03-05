package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/engine"
	"github.com/EvanPluchart/claude-code-status-line/internal/exchange"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/wizard"
)

// version is set by goreleaser at build time.
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		render()

		return
	}

	switch os.Args[1] {
	case "render":
		render()
	case "init":
		initCmd()
	case "config":
		configCmd()
	case "update-rates":
		updateRatesCmd()
	case "update":
		updateCmd()
	case "uninstall":
		uninstall()
	case "version", "--version", "-v":
		fmt.Println("claude-code-status-line " + version)
	case "help", "--help", "-h":
		printHelp()
	default:
		render()
	}
}

func render() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil || len(data) == 0 {
		return
	}

	input, err := parser.Parse(data)
	if err != nil {
		return
	}

	cfg := config.Load()
	output := engine.Render(input, cfg)

	fmt.Print(output)
}

// --- Init command ---

func initCmd() {
	// Check for flags (non-interactive mode)
	hasFlags := false
	locale := ""
	theme := ""
	currency := ""

	for i := 2; i < len(os.Args)-1; i++ {
		switch os.Args[i] {
		case "--locale":
			locale = os.Args[i+1]
			hasFlags = true
		case "--theme":
			theme = os.Args[i+1]
			hasFlags = true
		case "--currency":
			currency = os.Args[i+1]
			hasFlags = true
		}
	}

	if !hasFlags {
		cfg, confirmed, err := wizard.Run(nil, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !confirmed {
			fmt.Fprintln(os.Stderr, "  Setup cancelled.")
			return
		}

		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "  Configuration saved to %s\n", config.ConfigPath())
		refreshExchangeRates()
		registerInClaude()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  Done! Restart Claude Code to see your new statusline.")

		return
	}

	cfg := config.Default()

	if locale != "" {
		cfg.Locale = locale
	}

	if theme != "" {
		cfg.Theme = theme
	}

	if currency != "" {
		cfg.Widgets.Cost.Currency = currency
	}

	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "  Configuration saved to %s\n", config.ConfigPath())

	refreshExchangeRates()
	registerInClaude()

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Done! Restart Claude Code to see your new statusline.")
	fmt.Fprintln(os.Stderr, "  Run 'claude-code-status-line config' to customize further.")
}

// --- Config command ---

func configCmd() {
	if len(os.Args) < 3 {
		existing := config.Load()

		cfg, confirmed, err := wizard.Run(existing, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !confirmed {
			fmt.Fprintln(os.Stderr, "  Changes discarded.")
			return
		}

		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "  Configuration saved to %s\n", config.ConfigPath())

		return
	}

	switch os.Args[2] {
	case "edit":
		openInEditor()
	case "set":
		configSet()
	case "get":
		configGet()
	case "path":
		fmt.Println(config.ConfigPath())
	case "reset":
		configReset()
	default:
		fmt.Fprintf(os.Stderr, "Unknown config command: %s\n", os.Args[2])
		fmt.Fprintln(os.Stderr, "Available: edit, set, get, path, reset")
		os.Exit(1)
	}
}

func openInEditor() {
	path := config.ConfigPath()

	editor := os.Getenv("VISUAL")

	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	if editor == "" {
		for _, e := range []string{"vim", "nano", "vi"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e

				break
			}
		}
	}

	if editor == "" {
		fmt.Fprintf(os.Stderr, "No editor found. Edit manually:\n  %s\n", path)

		return
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
	}
}

func configSet() {
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "Usage: claude-code-status-line config set <key> <value>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Keys:")
		fmt.Fprintln(os.Stderr, "  theme       Theme name (default, minimal, neon, dracula, catppuccin, nord)")
		fmt.Fprintln(os.Stderr, "  locale      Language (en, fr)")
		fmt.Fprintln(os.Stderr, "  currency    Currency code (USD, EUR, GBP, JPY, CAD...)")
		os.Exit(1)
	}

	key := os.Args[3]
	value := os.Args[4]

	cfg := config.Load()

	switch key {
	case "theme":
		cfg.Theme = value
	case "locale":
		cfg.Locale = value
	case "currency":
		cfg.Widgets.Cost.Currency = value
	default:
		fmt.Fprintf(os.Stderr, "Unknown key: %s\nAvailable keys: theme, locale, currency\n", key)
		fmt.Fprintln(os.Stderr, "For advanced options, use: claude-code-status-line config edit")
		os.Exit(1)
	}

	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "%s = %s\n", key, value)
}

func configGet() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "Usage: claude-code-status-line config get <key>")
		fmt.Fprintln(os.Stderr, "Keys: theme, locale, currency")
		os.Exit(1)
	}

	key := os.Args[3]
	cfg := config.Load()

	switch key {
	case "theme":
		fmt.Println(cfg.Theme)
	case "locale":
		fmt.Println(cfg.Locale)
	case "currency":
		fmt.Println(cfg.Widgets.Cost.Currency)
	default:
		fmt.Fprintf(os.Stderr, "Unknown key: %s\n", key)
		os.Exit(1)
	}
}

func configReset() {
	cfg := config.Default()

	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Configuration reset to defaults.")
}

// --- Register / Unregister ---

func registerInClaude() {
	settingsPath := config.ClaudeSettingsPath()

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "  Claude Code settings not found. Add manually:")
		fmt.Fprintln(os.Stderr, `    "statusLine": { "type": "command", "command": "claude-code-status-line" }`)

		return
	}

	var settings map[string]interface{}

	if err := json.Unmarshal(data, &settings); err != nil {
		fmt.Fprintf(os.Stderr, "  Error parsing settings: %v\n", err)

		return
	}

	// Backup existing statusline
	if existing, ok := settings["statusLine"]; ok {
		backupData, _ := json.MarshalIndent(existing, "", "  ")
		_ = os.MkdirAll(config.ConfigDir(), 0o755)
		_ = os.WriteFile(config.BackupPath(), backupData, 0o644)
		fmt.Fprintln(os.Stderr, "  Previous statusline configuration backed up.")
	}

	settings["statusLine"] = map[string]interface{}{
		"type":    "command",
		"command": "claude-code-status-line",
	}

	out, _ := json.MarshalIndent(settings, "", "  ")

	if err := os.WriteFile(settingsPath, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  Error writing settings: %v\n", err)

		return
	}

	fmt.Fprintln(os.Stderr, "  Registered in Claude Code settings.")
}

func uninstall() {
	settingsPath := config.ClaudeSettingsPath()

	data, err := os.ReadFile(settingsPath)
	if err == nil {
		var settings map[string]interface{}

		if err := json.Unmarshal(data, &settings); err == nil {
			backupData, backupErr := os.ReadFile(config.BackupPath())

			if backupErr == nil {
				var backup interface{}
				_ = json.Unmarshal(backupData, &backup)
				settings["statusLine"] = backup
				fmt.Fprintln(os.Stderr, "Previous statusline configuration restored.")
			} else {
				delete(settings, "statusLine")
			}

			out, _ := json.MarshalIndent(settings, "", "  ")
			_ = os.WriteFile(settingsPath, out, 0o644)
			fmt.Fprintln(os.Stderr, "Unregistered from Claude Code settings.")
		}
	}

	_ = os.RemoveAll(config.ConfigDir())
	fmt.Fprintln(os.Stderr, "Configuration directory removed.")
	fmt.Fprintln(os.Stderr, "\nClaude Code Status Line has been uninstalled.")
}

// --- Update command ---

func updateCmd() {
	if version == "dev" {
		fmt.Fprintln(os.Stderr, "Development build, skipping update check.")
		return
	}

	latest, err := fetchLatestVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	current := strings.TrimPrefix(version, "v")
	if latest == current {
		fmt.Fprintln(os.Stderr, "Already up to date ("+current+").")
		return
	}

	fmt.Fprintf(os.Stderr, "New version available: %s → %s\n", current, latest)

	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
		os.Exit(1)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving executable path: %v\n", err)
		os.Exit(1)
	}

	lowerPath := strings.ToLower(execPath)
	if strings.Contains(lowerPath, "cellar") ||
		strings.Contains(lowerPath, "homebrew") ||
		strings.Contains(lowerPath, "linuxbrew") {
		fmt.Fprintln(os.Stderr, "Installed via Homebrew. Run:")
		fmt.Fprintln(os.Stderr, "  brew upgrade claude-code-status-line")
		return
	}

	if err := selfUpdate(execPath, latest); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Updated to %s successfully.\n", latest)
}

func ghToken() string {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	out, err := exec.Command("gh", "auth", "token").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	return ""
}

func ghRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token := ghToken(); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	return http.DefaultClient.Do(req)
}

func fetchLatestVersion() (string, error) {
	resp, err := ghRequest("https://api.github.com/repos/EvanPluchart/claude-code-status-line/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

func selfUpdate(execPath, latest string) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}

	assetName := fmt.Sprintf("claude-code-status-line_%s_%s_%s.%s", latest, goos, goarch, ext)

	fmt.Fprintf(os.Stderr, "Downloading %s...\n", assetName)

	resp, err := downloadAsset(latest, assetName)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	tmpDir, err := os.MkdirTemp("", "ccsl-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	binaryName := "claude-code-status-line"
	if goos == "windows" {
		binaryName += ".exe"
	}

	var newBinaryPath string
	if ext == "zip" {
		newBinaryPath, err = extractZip(resp.Body, tmpDir, binaryName)
	} else {
		newBinaryPath, err = extractTarGz(resp.Body, tmpDir, binaryName)
	}
	if err != nil {
		return err
	}

	// Replace current binary
	oldPath := execPath + ".old"
	_ = os.Remove(oldPath)

	if err := os.Rename(execPath, oldPath); err != nil {
		return fmt.Errorf("failed to move old binary: %w", err)
	}

	if err := os.Rename(newBinaryPath, execPath); err != nil {
		// Try to restore old binary
		_ = os.Rename(oldPath, execPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	_ = os.Remove(oldPath)

	return nil
}

func downloadAsset(tag, assetName string) (*http.Response, error) {
	// Try direct download first (works for public repos)
	url := fmt.Sprintf(
		"https://github.com/EvanPluchart/claude-code-status-line/releases/download/v%s/%s",
		tag, assetName,
	)

	resp, err := ghRequest(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return resp, nil
	}
	if resp != nil {
		resp.Body.Close()
	}

	// For private repos, find asset via API and download with octet-stream accept header
	apiURL := "https://api.github.com/repos/EvanPluchart/claude-code-status-line/releases/tags/v" + tag
	resp, err = ghRequest(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	var release struct {
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			req, err := http.NewRequest("GET", asset.URL, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Accept", "application/octet-stream")
			if token := ghToken(); token != "" {
				req.Header.Set("Authorization", "token "+token)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("failed to download asset: %w", err)
			}
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				return nil, fmt.Errorf("asset download returned status %d", resp.StatusCode)
			}
			return resp, nil
		}
	}

	return nil, fmt.Errorf("asset %s not found in release v%s", assetName, tag)
}

func extractTarGz(r io.Reader, destDir, binaryName string) (string, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return "", fmt.Errorf("failed to decompress: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read archive: %w", err)
		}

		if filepath.Base(header.Name) == binaryName && header.Typeflag == tar.TypeReg {
			outPath := filepath.Join(destDir, binaryName)
			out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0o755)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}
			out.Close()
			return outPath, nil
		}
	}

	return "", fmt.Errorf("binary %s not found in archive", binaryName)
}

func extractZip(r io.Reader, destDir, binaryName string) (string, error) {
	// zip needs random access, so write to temp file first
	tmpFile, err := os.CreateTemp(destDir, "archive-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, r); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write archive: %w", err)
	}
	tmpFile.Close()

	zr, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open zip: %w", err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("failed to read zip entry: %w", err)
			}
			defer rc.Close()

			outPath := filepath.Join(destDir, binaryName)
			out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0o755)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(out, rc); err != nil {
				out.Close()
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}
			out.Close()
			return outPath, nil
		}
	}

	return "", fmt.Errorf("binary %s not found in archive", binaryName)
}

func refreshExchangeRates() {
	fmt.Fprint(os.Stderr, "  Fetching exchange rates...")

	if err := exchange.Refresh(); err != nil {
		fmt.Fprintln(os.Stderr, " skipped (will use fallback rates).")

		return
	}

	fmt.Fprintln(os.Stderr, " done.")
}

func updateRatesCmd() {
	fmt.Fprintln(os.Stderr, "Refreshing exchange rates...")

	if err := exchange.Refresh(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Fallback rates will be used.")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Exchange rates updated.")
}

func printHelp() {
	help := `claude-code-status-line - A customizable statusline for Claude Code

Usage:
  claude-code-status-line [command]

Commands:
  render        Render the statusline (default, reads stdin)
  init          Interactive setup (config + register in Claude Code)
  config        Edit configuration
  update        Update to the latest version
  update-rates  Refresh exchange rates cache
  uninstall     Remove config and unregister from Claude Code
  version       Show version

Config subcommands:
  config              Interactive wizard with live preview
  config edit         Open config in $EDITOR
  config set <k> <v>  Set a config value (theme, locale, currency)
  config get <k>      Get a config value
  config path         Print config file path
  config reset        Reset to defaults

Init flags (non-interactive):
  --locale      Set locale (en, fr)
  --theme       Set theme name
  --currency    Set currency code

Config file: ~/.claude-statusline/config.yml`

	fmt.Println(help)
}
