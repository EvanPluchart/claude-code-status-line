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
const maxWidgetsPerLine = 6

// widgetOrder is the deterministic display order for widgets.
// separator and spacer are excluded: separators are auto-inserted between widgets.
var widgetOrder = []string{
	"model",
	"directory",
	"git-branch",
	"git-status",
	"cost",
	"duration",
	"token-bar",
	"context-percent",
	"token-count",
	"total-tokens",
	"cache-ratio",
	"lines-changed",
	"vim-mode",
	"timestamp",
	"os-info",
	"nested-repos",
	"session-bar",
	"session-percent",
	"weekly-bar",
	"weekly-percent",
}

// widgetLabels maps widget IDs to human-readable descriptions.
var widgetLabels = map[string]string{
	"model":           "Model name (Opus 4.6, Sonnet...)",
	"directory":       "Project directory",
	"git-branch":      "Git branch name",
	"git-status":      "Git clean/dirty indicator",
	"cost":            "Session cost",
	"duration":        "Session duration",
	"token-bar":       "Token usage progress bar",
	"context-percent": "Context usage percentage",
	"token-count":     "Token count (used/max)",
	"total-tokens":    "Total input/output tokens",
	"cache-ratio":     "Cache hit ratio",
	"lines-changed":   "Lines added/removed",
	"vim-mode":        "Vim mode indicator",
	"timestamp":       "Current time",
	"os-info":         "OS and architecture",
	"nested-repos":    "Nested git repositories count",
	"session-bar":     "Session rate limit bar (5h)",
	"session-percent": "Session rate limit % (5h)",
	"weekly-bar":      "Weekly rate limit bar (7d)",
	"weekly-percent":  "Weekly rate limit % (7d)",
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
				{"fr", "Francais"},
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
