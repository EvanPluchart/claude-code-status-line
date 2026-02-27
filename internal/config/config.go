package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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

// WidgetOptions holds per-widget configuration.
type WidgetOptions struct {
	Cost      CostOptions      `yaml:"cost"`
	TokenBar  TokenBarOptions  `yaml:"token_bar"`
	Separator SeparatorOptions `yaml:"separator"`
	Model     ModelOptions     `yaml:"model"`
	Timestamp TimestampOptions `yaml:"timestamp"`
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
	Context  ThresholdGroup `yaml:"context"`
	Cost     ThresholdGroup `yaml:"cost"`
	Duration ThresholdGroup `yaml:"duration"`
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
			{Widgets: []string{"model", "separator", "directory", "separator", "git-branch", "separator", "duration", "separator", "cost"}},
			{Widgets: []string{"token-bar", "context-percent", "token-count"}},
			{Widgets: []string{}},
		},
		Widgets: WidgetOptions{
			Cost:      CostOptions{Currency: "USD", Decimals: 2},
			TokenBar:  TokenBarOptions{Width: 16, FilledChar: "\u2501", EmptyChar: "\u2500"},
			Separator: SeparatorOptions{Char: "\u2502"},
		},
		Thresholds: Thresholds{
			Context:  ThresholdGroup{Green: 0, Yellow: 50, Orange: 70, Red: 90},
			Cost:     ThresholdGroup{Green: 0, Yellow: 0.25, Orange: 1.0, Red: 5.0},
			Duration: ThresholdGroup{Green: 0, Yellow: 60, Red: 1800},
		},
	}
}

// ConfigDir returns the configuration directory path.
func ConfigDir() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".claude-statusline")
}

// ConfigPath returns the configuration file path.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yml")
}

// ClaudeSettingsPath returns the Claude Code settings file path.
func ClaudeSettingsPath() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".claude", "settings.json")
}

// BackupPath returns the statusline backup file path.
func BackupPath() string {
	return filepath.Join(ConfigDir(), "statusline.backup.json")
}

// Load reads the config from disk. Returns defaults if not found.
func Load() *Config {
	cfg := Default()
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	var userCfg Config

	if err := yaml.Unmarshal(data, &userCfg); err != nil {
		return cfg
	}

	// Merge: user values override defaults
	if userCfg.Locale != "" {
		cfg.Locale = userCfg.Locale
	}

	if userCfg.Theme != "" {
		cfg.Theme = userCfg.Theme
	}

	if len(userCfg.Lines) > 0 {
		cfg.Lines = userCfg.Lines
	}

	if userCfg.Widgets.Cost.Currency != "" {
		cfg.Widgets.Cost.Currency = userCfg.Widgets.Cost.Currency
	}

	if userCfg.Widgets.Cost.Decimals > 0 {
		cfg.Widgets.Cost.Decimals = userCfg.Widgets.Cost.Decimals
	}

	if userCfg.Widgets.TokenBar.Width > 0 {
		cfg.Widgets.TokenBar.Width = userCfg.Widgets.TokenBar.Width
	}

	if userCfg.Widgets.TokenBar.FilledChar != "" {
		cfg.Widgets.TokenBar.FilledChar = userCfg.Widgets.TokenBar.FilledChar
	}

	if userCfg.Widgets.TokenBar.EmptyChar != "" {
		cfg.Widgets.TokenBar.EmptyChar = userCfg.Widgets.TokenBar.EmptyChar
	}

	if userCfg.Widgets.Separator.Char != "" {
		cfg.Widgets.Separator.Char = userCfg.Widgets.Separator.Char
	}

	cfg.Widgets.Model.ShortName = userCfg.Widgets.Model.ShortName
	cfg.Widgets.Timestamp.ShowSeconds = userCfg.Widgets.Timestamp.ShowSeconds

	return cfg
}

// Save writes the config to disk.
func Save(cfg *Config) error {
	dir := ConfigDir()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), data, 0o644)
}
