package wizard

import (
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/themes"
	"github.com/EvanPluchart/claude-code-status-line/internal/widgets"
)

// previewPlaceholders provides fallback text for widgets that need
// external resources (git, etc.) unavailable during preview.
var previewPlaceholders = map[string]func(themes.Theme) string{
	"git-branch": func(t themes.Theme) string {
		return ansi.Colorize("main", t.Info)
	},
	"git-status": func(t themes.Theme) string {
		return ansi.Colorize("\u2713", t.Success)
	},
	"nested-repos": func(t themes.Theme) string {
		return ansi.Colorize("2 repos", t.Muted)
	},
}

// sampleInput returns a realistic Input for preview rendering.
func sampleInput() *parser.Input {
	return &parser.Input{
		CWD:       "/home/user/my-project",
		SessionID: "preview-session",
		Model: parser.Model{
			ID:          "claude-opus-4-6",
			DisplayName: "Opus 4.6",
		},
		Workspace: parser.Workspace{
			CurrentDir: "/home/user/my-project",
			ProjectDir: "/home/user/my-project",
		},
		Cost: parser.Cost{
			TotalCostUSD:    0.42,
			TotalDurationMS: 754000, // 12m34s
			TotalLinesAdded: 42, TotalLinesRemoved: 7,
		},
		ContextWindow: parser.ContextWindow{
			ContextWindowSize: 200000,
			UsedPercentage:    45.0,
			RemainingPct:      55.0,
			TotalInputTokens:  75000,
			TotalOutputTokens: 15000,
			CurrentUsage: &parser.CurrentUsage{
				InputTokens:              70000,
				OutputTokens:             20000,
				CacheCreationInputTokens: 5000,
				CacheReadInputTokens:     30000,
			},
		},
		Vim: &parser.Vim{Mode: "normal"},
	}
}

// withSeparators inserts "separator" between each widget in the list.
func withSeparators(wids []string) []string {
	if len(wids) == 0 {
		return nil
	}
	result := make([]string, 0, len(wids)*2-1)
	for i, w := range wids {
		if i > 0 {
			result = append(result, "separator")
		}
		result = append(result, w)
	}
	return result
}

// buildPreviewConfig creates a Config from the current wizard selections,
// using the hovered value for the current single-select step.
func buildPreviewConfig(m *Model) *config.Config {
	cfg := config.Default()

	currentStep := m.steps[m.stepIndex]

	applySelection := func(id StepID, apply func(string)) {
		if !currentStep.IsMulti && currentStep.ID == id {
			apply(currentStep.Choices[m.cursor].Value)
		} else if v, ok := m.selections[id]; ok {
			apply(v)
		}
	}

	applySelection(StepLocale, func(v string) { cfg.Locale = v })
	applySelection(StepTheme, func(v string) { cfg.Theme = v })
	applySelection(StepCurrency, func(v string) { cfg.Widgets.Cost.Currency = v })

	// Apply widget lines from toggle order, auto-inserting separators
	lines := []config.LineConfig{{}, {}, {}}
	for i, step := range []StepID{StepLine1, StepLine2, StepLine3} {
		if order, ok := m.toggleOrder[step]; ok {
			lines[i].Widgets = withSeparators(order)
		}
	}
	cfg.Lines = lines

	return cfg
}

// renderPreviewLines renders the statusline with placeholder fallbacks
// for widgets that return empty (e.g. git widgets without a real repo).
func renderPreviewLines(cfg *config.Config, input *parser.Input) string {
	theme := themes.Get(cfg.Theme)
	ctx := &widgets.Context{
		Input:  input,
		Config: cfg,
		Theme:  theme,
	}

	var lines []string
	for _, line := range cfg.Lines {
		if len(line.Widgets) == 0 {
			continue
		}
		var parts []string
		for _, wid := range line.Widgets {
			w := widgets.Get(wid)
			if w == nil {
				continue
			}
			rendered := w.Render(ctx)
			if rendered == "" {
				if ph, ok := previewPlaceholders[wid]; ok {
					rendered = ph(theme)
				}
			}
			if rendered != "" {
				parts = append(parts, rendered)
			}
		}
		if len(parts) > 0 {
			lines = append(lines, strings.Join(parts, ""))
		}
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// renderPreview renders the statusline preview inside a box.
func renderPreview(m *Model) string {
	cfg := buildPreviewConfig(m)
	input := sampleInput()
	raw := renderPreviewLines(cfg, input)

	if raw == "" {
		raw = "  (no widgets selected)\n"
	}

	var b strings.Builder
	b.WriteString("  \x1b[2m╭─ Preview ─────────────────────────────────────────╮\x1b[0m\n")
	for _, line := range strings.Split(strings.TrimRight(raw, "\n"), "\n") {
		b.WriteString("  \x1b[2m│\x1b[0m ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString("  \x1b[2m╰───────────────────────────────────────────────────╯\x1b[0m")

	return b.String()
}
