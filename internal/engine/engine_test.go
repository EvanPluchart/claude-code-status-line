package engine

import (
	"regexp"
	"strings"
	"testing"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func plain(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func sampleInput() *parser.Input {
	input, _ := parser.Parse([]byte(`{
		"cwd": "/tmp/app",
		"model": {"id": "claude-sonnet-5", "display_name": "Sonnet"},
		"cost": {"total_cost_usd": 0.5, "total_duration_ms": 90000},
		"context_window": {"context_window_size": 200000, "used_percentage": 25}
	}`))

	return input
}

func TestRenderJoinsWidgetsWithSpaces(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Widgets: []string{"model", "context-percent", "separator", "cost"}},
	}

	got := plain(Render(sampleInput(), cfg))

	if got != "Sonnet 5 25% │ $0.50\n" {
		t.Errorf("Render = %q", got)
	}
}

func TestRenderStripsDanglingSeparators(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		// effort/vim-mode render empty: their separators must collapse.
		{Widgets: []string{"separator", "model", "separator", "effort", "separator", "vim-mode", "separator", "cost", "separator"}},
	}

	got := plain(Render(sampleInput(), cfg))

	if got != "Sonnet 5 │ $0.50\n" {
		t.Errorf("Render = %q", got)
	}
}

func TestRenderSkipsEmptyLinesAndUnknownWidgets(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Widgets: []string{"model"}},
		{Widgets: []string{}},
		{Widgets: []string{"does-not-exist"}},
	}

	got := plain(Render(sampleInput(), cfg))

	if got != "Sonnet 5\n" {
		t.Errorf("Render = %q", got)
	}
}

func TestRenderCapsLineCount(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Widgets: []string{"model"}},
		{Widgets: []string{"model"}},
		{Widgets: []string{"model"}},
		{Widgets: []string{"model"}},
	}

	got := Render(sampleInput(), cfg)

	if lines := strings.Count(got, "\n"); lines != config.MaxLines {
		t.Errorf("rendered %d lines, want %d", lines, config.MaxLines)
	}
}

func TestRenderEmptyConfig(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{{Widgets: []string{"vim-mode"}}}

	if got := Render(sampleInput(), cfg); got != "" {
		t.Errorf("expected empty output, got %q", got)
	}
}
