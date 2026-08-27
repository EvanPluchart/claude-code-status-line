package wizard

// StepID identifies a wizard step.
type StepID int

const (
	StepLocale StepID = iota
	StepTheme
	StepCurrency
	StepLine1
	StepLine2
	StepLine3
	StepConfirm
)

// Choice is a selectable option.
type Choice struct {
	Value string
	Label string
}

// StepDef describes a wizard step.
type StepDef struct {
	ID      StepID
	Title   string
	IsMulti bool // multi-select (widget lines) vs single-select
	Choices []Choice
}

// maxWidgetsPerLine is the maximum number of widgets per line.
const maxWidgetsPerLine = 8

// widgetOrder is the deterministic display order for widgets.
// separator and spacer are excluded: separators are auto-inserted between widgets.
var widgetOrder = []string{
	"model",
	"effort",
	"thinking",
	"fast-mode",
	"directory",
	"repo",
	"git-branch",
	"git-status",
	"git-changes",
	"git-ahead-behind",
	"pr",
	"worktree",
	"cost",
	"burn-rate",
	"duration",
	"api-time",
	"token-bar",
	"context-percent",
	"token-count",
	"context-remaining",
	"total-tokens",
	"cache-ratio",
	"exceeds-200k",
	"session-usage",
	"weekly-usage",
	"lines-changed",
	"agent",
	"output-style",
	"session-name",
	"vim-mode",
	"timestamp",
	"claude-version",
	"os-info",
	"hostname",
	"nested-repos",
}

// widgetLabels maps widget IDs to human-readable descriptions.
var widgetLabels = map[string]string{
	"model":             "Model name (Opus 5, Sonnet 5...)",
	"effort":            "Reasoning effort level",
	"thinking":          "Extended thinking indicator",
	"fast-mode":         "Fast mode indicator",
	"directory":         "Project directory",
	"repo":              "Remote repository (owner/name)",
	"git-branch":        "Git branch name",
	"git-status":        "Git clean/dirty indicator",
	"git-changes":       "Staged/modified/untracked counts",
	"git-ahead-behind":  "Commits ahead/behind upstream",
	"pr":                "Open pull request + review state",
	"worktree":          "Active worktree name",
	"cost":              "Session cost",
	"burn-rate":         "Spending rate per hour",
	"duration":          "Session duration",
	"api-time":          "Time waiting for the API",
	"token-bar":         "Context usage progress bar",
	"context-percent":   "Context usage percentage",
	"token-count":       "Token count (used/max)",
	"context-remaining": "Tokens left in context",
	"total-tokens":      "Total input/output tokens",
	"cache-ratio":       "Cache hit ratio",
	"exceeds-200k":      "Marker when context exceeds 200k",
	"session-usage":     "Session rate limit (5h bar + reset)",
	"weekly-usage":      "Weekly rate limit (7d bar + reset)",
	"lines-changed":     "Lines added/removed",
	"agent":             "Active agent name",
	"output-style":      "Output style (when not default)",
	"session-name":      "Session name/title",
	"vim-mode":          "Vim mode indicator",
	"timestamp":         "Current time",
	"claude-version":    "Claude Code version",
	"os-info":           "OS and architecture",
	"hostname":          "Machine hostname",
	"nested-repos":      "Nested git repositories count",
}

func buildSteps() []StepDef {
	widgetChoices := make([]Choice, len(widgetOrder))
	for i, id := range widgetOrder {
		widgetChoices[i] = Choice{Value: id, Label: widgetLabels[id]}
	}

	return []StepDef{
		{
			ID: StepLocale, Title: "Language",
			Choices: []Choice{
				{"en", "English"},
				{"fr", "Français"},
			},
		},
		{
			ID: StepTheme, Title: "Theme",
			Choices: []Choice{
				{"default", "Adaptive colors"},
				{"minimal", "Monochrome, subtle"},
				{"neon", "Bright, vibrant"},
				{"dracula", "Dracula palette"},
				{"catppuccin", "Catppuccin Mocha"},
				{"nord", "Nord palette"},
			},
		},
		{
			ID: StepCurrency, Title: "Currency",
			Choices: []Choice{
				{"USD", "US Dollar ($)"},
				{"EUR", "Euro (€)"},
				{"GBP", "British Pound (£)"},
				{"JPY", "Japanese Yen (¥)"},
				{"CAD", "Canadian Dollar (C$)"},
				{"AUD", "Australian Dollar (A$)"},
				{"CHF", "Swiss Franc (CHF)"},
			},
		},
		{
			ID: StepLine1, Title: "Line 1 widgets",
			IsMulti: true, Choices: widgetChoices,
		},
		{
			ID: StepLine2, Title: "Line 2 widgets",
			IsMulti: true, Choices: widgetChoices,
		},
		{
			ID: StepLine3, Title: "Line 3 widgets (optional)",
			IsMulti: true, Choices: widgetChoices,
		},
		{
			ID: StepConfirm, Title: "Save configuration?",
			Choices: []Choice{
				{"yes", "Save and apply"},
				{"no", "Discard changes"},
			},
		},
	}
}
