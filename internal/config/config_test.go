package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromMissingFileReturnsDefaults(t *testing.T) {
	cfg, err := LoadFrom(filepath.Join(t.TempDir(), "missing.yml"))

	if err == nil {
		t.Error("expected an error for a missing file")
	}

	if cfg.Theme != "default" || cfg.Locale != "en" || len(cfg.Lines) != 3 {
		t.Errorf("defaults not applied: %+v", cfg)
	}
}

func TestMergeOverridesOnlyPresentKeys(t *testing.T) {
	cfg := Default()

	data := []byte(`
theme: dracula
lines:
  - widgets: [model, cost]
widgets:
  cost:
    currency: EUR
  model:
    short_name: true
  lines_changed:
    source: session
thresholds:
  context:
    yellow: 50
    red: 95
`)

	if err := Merge(cfg, data); err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	if cfg.Theme != "dracula" {
		t.Errorf("theme = %q", cfg.Theme)
	}

	if cfg.Locale != "en" {
		t.Errorf("locale should keep default, got %q", cfg.Locale)
	}

	if len(cfg.Lines) != 1 || len(cfg.Lines[0].Widgets) != 2 {
		t.Errorf("lines not overridden: %+v", cfg.Lines)
	}

	if cfg.Widgets.Cost.Currency != "EUR" || cfg.Widgets.Cost.Decimals != 2 {
		t.Errorf("cost options = %+v", cfg.Widgets.Cost)
	}

	if !cfg.Widgets.Model.ShortName {
		t.Error("short_name should be true")
	}

	if cfg.Widgets.LinesChanged.Source != "session" {
		t.Errorf("lines_changed.source = %q", cfg.Widgets.LinesChanged.Source)
	}

	if cfg.Widgets.TokenBar.Width != 16 || cfg.Widgets.TokenBar.FilledChar != "━" {
		t.Errorf("token bar defaults lost: %+v", cfg.Widgets.TokenBar)
	}

	ctxTh := cfg.Thresholds.Context

	if ctxTh.Yellow != 50 || ctxTh.Red != 95 || ctxTh.Orange != 80 {
		t.Errorf("context thresholds = %+v (orange should keep its default)", ctxTh)
	}

	if cfg.Thresholds.Cost.Red != 5.0 {
		t.Errorf("cost thresholds should keep defaults: %+v", cfg.Thresholds.Cost)
	}
}

func TestMergeInvalidYAML(t *testing.T) {
	cfg := Default()

	if err := Merge(cfg, []byte("theme: [unclosed")); err == nil {
		t.Error("expected an error for invalid YAML")
	}
}

func TestNormalizeFixesInvalidValues(t *testing.T) {
	cfg := &Config{
		Lines: []LineConfig{{}, {}, {}, {Widgets: []string{"model"}}},
		Widgets: WidgetOptions{
			Cost:         CostOptions{Decimals: -1},
			TokenBar:     TokenBarOptions{Width: -5},
			LinesChanged: LinesChangedOptions{Source: "bogus"},
		},
	}

	Normalize(cfg)

	if len(cfg.Lines) != MaxLines {
		t.Errorf("lines should be capped to %d, got %d", MaxLines, len(cfg.Lines))
	}

	if cfg.Widgets.Cost.Decimals != 2 || cfg.Widgets.TokenBar.Width != 16 {
		t.Errorf("invalid values not normalized: %+v", cfg.Widgets)
	}

	if cfg.Widgets.LinesChanged.Source != "git" {
		t.Errorf("bogus source should fall back to git, got %q", cfg.Widgets.LinesChanged.Source)
	}

	if cfg.Thresholds.Context.Red != 90 {
		t.Errorf("thresholds should be defaulted: %+v", cfg.Thresholds.Context)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_STATUSLINE_DIR", dir)

	cfg := Default()
	cfg.Theme = "nord"
	cfg.Widgets.Directory.Depth = 2

	if err := Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "config.yml")); err != nil {
		t.Fatalf("config file not written: %v", err)
	}

	loaded := Load()

	if loaded.Theme != "nord" || loaded.Widgets.Directory.Depth != 2 {
		t.Errorf("round trip lost values: %+v", loaded)
	}
}

func TestUsesWidget(t *testing.T) {
	cfg := Default()

	if !cfg.UsesWidget("model") {
		t.Error("default config should use model")
	}

	if cfg.UsesWidget("vim-mode") {
		t.Error("default config should not use vim-mode")
	}
}
