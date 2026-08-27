package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// MaxLines is the maximum number of statusline lines rendered.
const MaxLines = 3

// LineConfig defines which widgets appear on a line.
type LineConfig struct {
	Widgets []string `yaml:"widgets"`
}

// CostOptions configures the cost widget.
type CostOptions struct {
	Currency string `yaml:"currency"`
	Decimals int    `yaml:"decimals"`
}

// TokenBarOptions configures the token bar widget.
type TokenBarOptions struct {
	Width      int    `yaml:"width"`
	FilledChar string `yaml:"filled_char"`
	EmptyChar  string `yaml:"empty_char"`
}

// SeparatorOptions configures the separator widget.
type SeparatorOptions struct {
	Char string `yaml:"char"`
}

// ModelOptions configures the model widget.
type ModelOptions struct {
	ShortName bool `yaml:"short_name"`
}

// TimestampOptions configures the timestamp widget.
type TimestampOptions struct {
	ShowSeconds bool `yaml:"show_seconds"`
}

// LinesChangedOptions configures the lines-changed widget.
// Source is "git" (working tree diff) or "session" (lines changed by Claude Code).
type LinesChangedOptions struct {
	Source string `yaml:"source"`
}

// DirectoryOptions configures the directory widget.
// Depth is the number of trailing path segments to display (0 = basename only).
type DirectoryOptions struct {
	Depth int `yaml:"depth"`
}

// GitOptions configures git-related widgets.
type GitOptions struct {
	MaxBranchLength int `yaml:"max_branch_length"`
}

// WidgetOptions holds per-widget configuration.
type WidgetOptions struct {
	Cost         CostOptions         `yaml:"cost"`
	TokenBar     TokenBarOptions     `yaml:"token_bar"`
	Separator    SeparatorOptions    `yaml:"separator"`
	Model        ModelOptions        `yaml:"model"`
	Timestamp    TimestampOptions    `yaml:"timestamp"`
	LinesChanged LinesChangedOptions `yaml:"lines_changed"`
	Directory    DirectoryOptions    `yaml:"directory"`
	Git          GitOptions          `yaml:"git"`
}

// ThresholdGroup holds color threshold values.
type ThresholdGroup struct {
	Green  float64 `yaml:"green"`
	Yellow float64 `yaml:"yellow"`
	Orange float64 `yaml:"orange"`
	Red    float64 `yaml:"red"`
}

// Thresholds holds all threshold configurations.
type Thresholds struct {
	Context   ThresholdGroup `yaml:"context"`
	Cost      ThresholdGroup `yaml:"cost"`
	Duration  ThresholdGroup `yaml:"duration"`
	RateLimit ThresholdGroup `yaml:"rate_limit"`
}

// Config is the user configuration file structure.
type Config struct {
	Version    int           `yaml:"version"`
	Locale     string        `yaml:"locale"`
	Theme      string        `yaml:"theme"`
	Lines      []LineConfig  `yaml:"lines"`
	Widgets    WidgetOptions `yaml:"widgets"`
	Thresholds Thresholds    `yaml:"thresholds"`
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		Version: 1,
		Locale:  "en",
		Theme:   "default",
		Lines: []LineConfig{
			{Widgets: []string{"model", "separator", "directory", "separator", "git-branch", "git-status", "separator", "duration", "separator", "cost"}},
			{Widgets: []string{"token-bar", "context-percent", "token-count", "separator", "effort", "thinking", "fast-mode"}},
			{Widgets: []string{}},
		},
		Widgets: WidgetOptions{
			Cost:         CostOptions{Currency: "USD", Decimals: 2},
			TokenBar:     TokenBarOptions{Width: 16, FilledChar: "━", EmptyChar: "─"},
			Separator:    SeparatorOptions{Char: "│"},
			LinesChanged: LinesChangedOptions{Source: "git"},
			Git:          GitOptions{MaxBranchLength: 40},
		},
		Thresholds: Thresholds{
			Context:   ThresholdGroup{Green: 0, Yellow: 60, Orange: 80, Red: 90},
			Cost:      ThresholdGroup{Green: 0, Yellow: 0.25, Orange: 1.0, Red: 5.0},
			Duration:  ThresholdGroup{Green: 0, Yellow: 60, Orange: 600, Red: 1800},
			RateLimit: ThresholdGroup{Green: 0, Yellow: 60, Orange: 80, Red: 90},
		},
	}
}

// ConfigDir returns the configuration directory path.
func ConfigDir() string {
	if dir := os.Getenv("CLAUDE_STATUSLINE_DIR"); dir != "" {
		return dir
	}

	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".claude-statusline")
}

// ConfigPath returns the configuration file path.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yml")
}

// ClaudeSettingsPath returns the Claude Code settings file path.
func ClaudeSettingsPath() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "settings.json")
	}

	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".claude", "settings.json")
}

// BackupPath returns the statusline backup file path.
func BackupPath() string {
	return filepath.Join(ConfigDir(), "statusline.backup.json")
}

// Load reads the config from disk. Returns defaults if not found or invalid.
func Load() *Config {
	cfg, _ := LoadFrom(ConfigPath())

	return cfg
}

// LoadFrom reads a config file and merges it over the defaults.
// It always returns a usable config; the error reports why user values were ignored.
func LoadFrom(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := Merge(cfg, data); err != nil {
		return Default(), err
	}

	return cfg, nil
}

// Merge unmarshals YAML data over an existing config: only keys present
// in the data override the current values.
func Merge(cfg *Config, data []byte) error {
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	Normalize(cfg)

	return nil
}

// Normalize fills invalid or missing values with defaults and caps the line count.
func Normalize(cfg *Config) {
	def := Default()

	if cfg.Locale == "" {
		cfg.Locale = def.Locale
	}

	if cfg.Theme == "" {
		cfg.Theme = def.Theme
	}

	if len(cfg.Lines) == 0 {
		cfg.Lines = def.Lines
	}

	if len(cfg.Lines) > MaxLines {
		cfg.Lines = cfg.Lines[:MaxLines]
	}

	if cfg.Widgets.Cost.Currency == "" {
		cfg.Widgets.Cost.Currency = def.Widgets.Cost.Currency
	}

	if cfg.Widgets.Cost.Decimals <= 0 {
		cfg.Widgets.Cost.Decimals = def.Widgets.Cost.Decimals
	}

	if cfg.Widgets.TokenBar.Width <= 0 {
		cfg.Widgets.TokenBar.Width = def.Widgets.TokenBar.Width
	}

	if cfg.Widgets.TokenBar.FilledChar == "" {
		cfg.Widgets.TokenBar.FilledChar = def.Widgets.TokenBar.FilledChar
	}

	if cfg.Widgets.TokenBar.EmptyChar == "" {
		cfg.Widgets.TokenBar.EmptyChar = def.Widgets.TokenBar.EmptyChar
	}

	if cfg.Widgets.Separator.Char == "" {
		cfg.Widgets.Separator.Char = def.Widgets.Separator.Char
	}

	if cfg.Widgets.LinesChanged.Source != "session" {
		cfg.Widgets.LinesChanged.Source = "git"
	}

	if cfg.Widgets.Git.MaxBranchLength <= 0 {
		cfg.Widgets.Git.MaxBranchLength = def.Widgets.Git.MaxBranchLength
	}

	normalizeThresholds(&cfg.Thresholds.Context, def.Thresholds.Context)
	normalizeThresholds(&cfg.Thresholds.Cost, def.Thresholds.Cost)
	normalizeThresholds(&cfg.Thresholds.Duration, def.Thresholds.Duration)
	normalizeThresholds(&cfg.Thresholds.RateLimit, def.Thresholds.RateLimit)
}

// normalizeThresholds fills missing levels so that yellow <= orange <= red always holds.
func normalizeThresholds(t *ThresholdGroup, def ThresholdGroup) {
	if t.Yellow <= 0 && t.Orange <= 0 && t.Red <= 0 {
		*t = def

		return
	}

	if t.Red <= 0 {
		t.Red = def.Red
	}

	if t.Orange <= 0 {
		t.Orange = t.Red
	}

	if t.Yellow <= 0 {
		t.Yellow = t.Orange
	}
}

// Save writes the config to disk.
func Save(cfg *Config) error {
	dir := ConfigDir()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	if err := os.WriteFile(ConfigPath(), data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// UsesWidget reports whether any line contains the given widget ID.
func (c *Config) UsesWidget(id string) bool {
	for _, line := range c.Lines {
		for _, w := range line.Widgets {
			if w == id {
				return true
			}
		}
	}

	return false
}
