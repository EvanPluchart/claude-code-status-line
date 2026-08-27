package widgets

import (
	"fmt"
	"math"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

// formatTokens renders a token count in a compact form (1.2k, 3.4M).
func formatTokens(count int) string {
	if count >= 1_000_000 {
		return trimZero(fmt.Sprintf("%.1f", float64(count)/1_000_000)) + "M"
	}

	if count >= 1_000 {
		return trimZero(fmt.Sprintf("%.1f", float64(count)/1_000)) + "k"
	}

	return fmt.Sprintf("%d", count)
}

// trimZero strips a trailing ".0" so that 200.0k becomes 200k.
func trimZero(s string) string {
	return strings.TrimSuffix(s, ".0")
}

// thresholdColor maps a value to a theme color using a threshold group.
// Values at or above red are rendered bold.
func thresholdColor(value float64, t config.ThresholdGroup, ctx *Context) string {
	if value >= t.Red {
		return ctx.Theme.Danger + ansi.Bold
	}

	if value >= t.Orange {
		return ctx.Theme.Danger
	}

	if value >= t.Yellow {
		return ctx.Theme.Warning
	}

	return ctx.Theme.Success
}

// contextColor returns the color for a context usage percentage.
func contextColor(pct float64, ctx *Context) string {
	return thresholdColor(pct, ctx.Config.Thresholds.Context, ctx)
}

// renderBar renders a progress bar colored by the given threshold group.
func renderBar(pct float64, t config.ThresholdGroup, ctx *Context) string {
	width := ctx.Config.Widgets.TokenBar.Width
	filledChar := ctx.Config.Widgets.TokenBar.FilledChar
	emptyChar := ctx.Config.Widgets.TokenBar.EmptyChar

	if width <= 0 {
		width = 16
	}

	if filledChar == "" {
		filledChar = "━"
	}

	if emptyChar == "" {
		emptyChar = "─"
	}

	clamped := math.Max(0, math.Min(100, pct))
	filled := int(math.Round(clamped * float64(width) / 100))
	empty := width - filled

	filledStr := strings.Repeat(filledChar, filled)
	emptyStr := strings.Repeat(emptyChar, empty)

	return ansi.Colorize(filledStr, thresholdColor(pct, t, ctx)) + ansi.Colorize(emptyStr, ctx.Theme.BarEmpty)
}

// renderPercent renders a percentage colored by the given threshold group.
func renderPercent(pct float64, t config.ThresholdGroup, ctx *Context) string {
	rounded := int(math.Round(pct))

	return ansi.Colorize(fmt.Sprintf("%d%%", rounded), thresholdColor(float64(rounded), t, ctx))
}

// TokenBarWidget displays a progress bar for context usage.
type TokenBarWidget struct{}

func (w *TokenBarWidget) ID() string { return "token-bar" }

func (w *TokenBarWidget) Render(ctx *Context) string {
	return renderBar(ctx.Input.ContextWindow.UsedPercentage, ctx.Config.Thresholds.Context, ctx)
}

// ContextPercentWidget displays the context usage percentage.
type ContextPercentWidget struct{}

func (w *ContextPercentWidget) ID() string { return "context-percent" }

func (w *ContextPercentWidget) Render(ctx *Context) string {
	return renderPercent(ctx.Input.ContextWindow.UsedPercentage, ctx.Config.Thresholds.Context, ctx)
}

// TokenCountWidget displays used/max tokens.
type TokenCountWidget struct{}

func (w *TokenCountWidget) ID() string { return "token-count" }

func (w *TokenCountWidget) Render(ctx *Context) string {
	cw := ctx.Input.ContextWindow
	text := fmt.Sprintf("(%s/%s)", formatTokens(cw.UsedTokens()), formatTokens(cw.ContextWindowSize))

	return ansi.Colorize(text, ctx.Theme.Muted)
}

// ContextRemainingWidget displays the number of tokens left before the context is full.
type ContextRemainingWidget struct{}

func (w *ContextRemainingWidget) ID() string { return "context-remaining" }

func (w *ContextRemainingWidget) Render(ctx *Context) string {
	cw := ctx.Input.ContextWindow
	remaining := cw.ContextWindowSize - cw.UsedTokens()

	if remaining < 0 {
		remaining = 0
	}

	t := i18n.Get(ctx.Config.Locale)
	text := fmt.Sprintf("%s %s", formatTokens(remaining), t.RemainingLabel)

	return ansi.Colorize(text, contextColor(cw.UsedPercentage, ctx))
}

// TotalTokensWidget displays total input/output tokens.
type TotalTokensWidget struct{}

func (w *TotalTokensWidget) ID() string { return "total-tokens" }

func (w *TotalTokensWidget) Render(ctx *Context) string {
	input := ctx.Input.ContextWindow.TotalInputTokens
	output := ctx.Input.ContextWindow.TotalOutputTokens

	text := fmt.Sprintf("↑%s ↓%s", formatTokens(input), formatTokens(output))

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

	ratio := int(math.Round(float64(usage.CacheReadInputTokens) / float64(total) * 100))
	t := i18n.Get(ctx.Config.Locale)

	return ansi.Colorize(fmt.Sprintf("%s: %d%%", t.CacheLabel, ratio), ctx.Theme.Muted)
}

// Exceeds200KWidget shows a marker when the conversation exceeds 200k tokens.
type Exceeds200KWidget struct{}

func (w *Exceeds200KWidget) ID() string { return "exceeds-200k" }

func (w *Exceeds200KWidget) Render(ctx *Context) string {
	if !ctx.Input.Exceeds200K {
		return ""
	}

	return ansi.ColorBold(">200k", ctx.Theme.Warning)
}
