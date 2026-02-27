package widgets

import (
	"fmt"
	"math"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

func formatTokens(count int) string {
	if count >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(count)/1_000_000)
	}

	if count >= 1_000 {
		return fmt.Sprintf("%.1fk", float64(count)/1_000)
	}

	return fmt.Sprintf("%d", count)
}

func barColor(pct float64, ctx *Context) string {
	if pct >= 90 {
		return ctx.Theme.Danger
	}

	if pct >= 70 {
		return ctx.Theme.Warning
	}

	return ctx.Theme.Success
}

// TokenBarWidget displays a progress bar for context usage.
type TokenBarWidget struct{}

func (w *TokenBarWidget) ID() string { return "token-bar" }

func (w *TokenBarWidget) Render(ctx *Context) string {
	pct := ctx.Input.ContextWindow.UsedPercentage
	width := ctx.Config.Widgets.TokenBar.Width
	filledChar := ctx.Config.Widgets.TokenBar.FilledChar
	emptyChar := ctx.Config.Widgets.TokenBar.EmptyChar

	if width == 0 {
		width = 16
	}

	if filledChar == "" {
		filledChar = "\u2501"
	}

	if emptyChar == "" {
		emptyChar = "\u2500"
	}

	clamped := math.Max(0, math.Min(100, pct))
	filled := int(math.Round(clamped * float64(width) / 100))
	empty := width - filled

	filledStr := strings.Repeat(filledChar, filled)
	emptyStr := strings.Repeat(emptyChar, empty)

	return barColor(pct, ctx) + filledStr + ansi.RST + ansi.Colorize(emptyStr, ctx.Theme.BarEmpty)
}

// ContextPercentWidget displays the context usage percentage.
type ContextPercentWidget struct{}

func (w *ContextPercentWidget) ID() string { return "context-percent" }

func (w *ContextPercentWidget) Render(ctx *Context) string {
	pct := int(math.Round(ctx.Input.ContextWindow.UsedPercentage))
	text := fmt.Sprintf("%d%%", pct)

	color := ctx.Theme.Success

	if pct >= 90 {
		color = ctx.Theme.Danger + ansi.Bold
	} else if pct >= 70 {
		color = ctx.Theme.Warning
	}

	return color + text + ansi.RST
}

// TokenCountWidget displays used/max tokens.
type TokenCountWidget struct{}

func (w *TokenCountWidget) ID() string { return "token-count" }

func (w *TokenCountWidget) Render(ctx *Context) string {
	cw := ctx.Input.ContextWindow
	used := cw.CurrentUsage.InputTokens + cw.CurrentUsage.CacheCreationInputTokens + cw.CurrentUsage.CacheReadInputTokens

	if used == 0 {
		used = int(float64(cw.ContextWindowSize) * cw.UsedPercentage / 100)
	}

	text := fmt.Sprintf("(%s/%s)", formatTokens(used), formatTokens(cw.ContextWindowSize))

	return ansi.Colorize(text, ctx.Theme.Muted)
}

// TotalTokensWidget displays total input/output tokens.
type TotalTokensWidget struct{}

func (w *TotalTokensWidget) ID() string { return "total-tokens" }

func (w *TotalTokensWidget) Render(ctx *Context) string {
	input := ctx.Input.ContextWindow.TotalInputTokens
	output := ctx.Input.ContextWindow.TotalOutputTokens

	text := fmt.Sprintf("\u2191%s \u2193%s", formatTokens(input), formatTokens(output))

	return ansi.Colorize(text, ctx.Theme.Muted)
}

// CacheRatioWidget displays the cache hit ratio.
type CacheRatioWidget struct{}

func (w *CacheRatioWidget) ID() string { return "cache-ratio" }

func (w *CacheRatioWidget) Render(ctx *Context) string {
	usage := ctx.Input.ContextWindow.CurrentUsage
	total := usage.InputTokens + usage.CacheReadInputTokens + usage.CacheCreationInputTokens

	if total == 0 {
		return ""
	}

	ratio := int(float64(usage.CacheReadInputTokens) / float64(total) * 100)
	t := i18n.Get(ctx.Config.Locale)

	return ansi.Colorize(fmt.Sprintf("%s: %d%%", t.CacheLabel, ratio), ctx.Theme.Muted)
}
