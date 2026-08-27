package wizard

import (
	"strings"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/engine"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/themes"
)

// previewPlaceholders provides fallback text for widgets that need
// external resources (git, etc.) unavailable during preview.
var previewPlaceholders = map[string]engine.Placeholder{
	"git-branch": func(t themes.Theme) string {
		return ansi.Colorize("feature/auth", t.Info)
	},
	"git-status": func(t themes.Theme) string {
		return ansi.Colorize("✓", t.Success)
	},
	"git-changes": func(t themes.Theme) string {
		return ansi.Colorize("●2", t.Success) + " " + ansi.Colorize("✚1", t.Warning)
	},
	"git-ahead-behind": func(t themes.Theme) string {
		return ansi.Colorize("↑2", t.Success)
	},
	"nested-repos": func(t themes.Theme) string {
		return ansi.Colorize("2 repos", t.Muted)
	},
	"lines-changed": func(t themes.Theme) string {
		return ansi.Colorize("+42", t.Success) + " " + ansi.Colorize("-7", t.Danger)
	},
}

// SampleInput returns a realistic Input for preview rendering.
func SampleInput() *parser.Input {
	return &parser.Input{
		CWD:         "/home/user/projects/my-project",
		SessionID:   "preview-session",
		SessionName: "Add OAuth login flow",
		Version:     "2.1.90",
		Model: parser.Model{
			ID:          "claude-opus-5",
			DisplayName: "Opus",
		},
		Workspace: parser.Workspace{
			CurrentDir: "/home/user/projects/my-project",
			ProjectDir: "/home/user/projects/my-project",
			Repo:       &parser.Repo{Host: "github.com", Owner: "acme", Name: "my-project"},
		},
		OutputStyle: &parser.OutputStyle{Name: "default"},
		Cost: parser.Cost{
			TotalCostUSD:       0.42,
			TotalDurationMS:    754000, // 12m34s
			TotalAPIDurationMS: 312000,
			TotalLinesAdded:    42,
			TotalLinesRemoved:  7,
		},
		ContextWindow: parser.ContextWindow{
			ContextWindowSize: 200000,
			UsedPercentage:    45.0,
			RemainingPct:      55.0,
			TotalInputTokens:  75000,
			TotalOutputTokens: 15000,
			CurrentUsage: &parser.CurrentUsage{
				InputTokens:              55000,
				OutputTokens:             2000,
				CacheCreationInputTokens: 5000,
				CacheReadInputTokens:     30000,
			},
		},
		Effort:   &parser.Effort{Level: "high"},
		Thinking: &parser.Thinking{Enabled: true},
		RateLimits: &parser.RateLimits{
			FiveHour: &parser.RateLimit{
				UsedPercentage: 23.5,
				ResetsAt:       time.Now().Add(3*time.Hour + 42*time.Minute).Unix(),
			},
			SevenDay: &parser.RateLimit{
				UsedPercentage: 41.2,
				ResetsAt:       time.Now().Add(4*24*time.Hour + 8*time.Hour).Unix(),
			},
		},
		Vim: &parser.Vim{Mode: "NORMAL"},
		PR:  &parser.PullRequest{Number: 1234, ReviewState: "pending"},
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

	// Apply widget lines from toggle order, auto-inserting separators.
	lines := []config.LineConfig{{}, {}, {}}

	for i, step := range []StepID{StepLine1, StepLine2, StepLine3} {
		if order, ok := m.toggleOrder[step]; ok {
			lines[i].Widgets = withSeparators(order)
		}
	}

	cfg.Lines = lines

	return cfg
}

// RenderPreview renders the statusline with placeholder fallbacks
// for widgets that need a real repository (git widgets).
func RenderPreview(cfg *config.Config, input *parser.Input) string {
	return engine.RenderWith(input, cfg, previewPlaceholders)
}

// renderPreview renders the statusline preview inside a box.
func renderPreview(m *Model) string {
	cfg := buildPreviewConfig(m)
	raw := RenderPreview(cfg, SampleInput())

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
