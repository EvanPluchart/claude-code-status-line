package widgets

import "github.com/EvanPluchart/claude-code-status-line/internal/ansi"

var modelNames = map[string]string{
	"claude-opus-4-6":            "Opus 4.6",
	"claude-opus-4-20250514":     "Opus 4",
	"claude-sonnet-4-6":          "Sonnet 4.6",
	"claude-sonnet-4-5-20250929": "Sonnet 4.5",
	"claude-sonnet-4-20250514":   "Sonnet 4",
	"claude-haiku-4-5-20251001":  "Haiku 4.5",
	"claude-haiku-3-5-20241022":  "Haiku 3.5",
}

var shortNames = map[string]string{
	"claude-opus-4-6":            "O4.6",
	"claude-opus-4-20250514":     "O4",
	"claude-sonnet-4-6":          "S4.6",
	"claude-sonnet-4-5-20250929": "S4.5",
	"claude-sonnet-4-20250514":   "S4",
	"claude-haiku-4-5-20251001":  "H4.5",
	"claude-haiku-3-5-20241022":  "H3.5",
}

// ModelWidget displays the current Claude model name.
type ModelWidget struct{}

func (w *ModelWidget) ID() string { return "model" }

func (w *ModelWidget) Render(ctx *Context) string {
	lookup := modelNames

	if ctx.Config.Widgets.Model.ShortName {
		lookup = shortNames
	}

	name, ok := lookup[ctx.Input.Model.ID]
	if !ok {
		name = ctx.Input.Model.DisplayName
	}

	pct := ctx.Input.ContextWindow.UsedPercentage

	color := ctx.Theme.Primary

	if pct >= 90 {
		color = ctx.Theme.Danger
	} else if pct >= 70 {
		color = ctx.Theme.Warning
	}

	return ansi.ColorBold(name, color)
}
