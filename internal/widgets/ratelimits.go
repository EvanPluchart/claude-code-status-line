package widgets

import (
	"fmt"
	"math"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
)

// SessionBarWidget displays a progress bar for the 5-hour session rate limit.
type SessionBarWidget struct{}

func (w *SessionBarWidget) ID() string { return "session-bar" }

func (w *SessionBarWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.FiveHour == nil {
		return ""
	}

	pct := ctx.Input.RateLimits.FiveHour.UsedPercentage

	return renderRateLimitBar(pct, ctx)
}

// SessionPercentWidget displays the 5-hour session rate limit percentage.
type SessionPercentWidget struct{}

func (w *SessionPercentWidget) ID() string { return "session-percent" }

func (w *SessionPercentWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.FiveHour == nil {
		return ""
	}

	pct := ctx.Input.RateLimits.FiveHour.UsedPercentage

	return renderRateLimitPercent(pct, ctx)
}

// WeeklyBarWidget displays a progress bar for the 7-day rate limit.
type WeeklyBarWidget struct{}

func (w *WeeklyBarWidget) ID() string { return "weekly-bar" }

func (w *WeeklyBarWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.SevenDay == nil {
		return ""
	}

	pct := ctx.Input.RateLimits.SevenDay.UsedPercentage

	return renderRateLimitBar(pct, ctx)
}

// WeeklyPercentWidget displays the 7-day rate limit percentage.
type WeeklyPercentWidget struct{}

func (w *WeeklyPercentWidget) ID() string { return "weekly-percent" }

func (w *WeeklyPercentWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.SevenDay == nil {
		return ""
	}

	pct := ctx.Input.RateLimits.SevenDay.UsedPercentage

	return renderRateLimitPercent(pct, ctx)
}

// renderRateLimitBar renders a progress bar for a rate limit percentage.
func renderRateLimitBar(pct float64, ctx *Context) string {
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

// renderRateLimitPercent renders a colored percentage for a rate limit.
func renderRateLimitPercent(pct float64, ctx *Context) string {
	rounded := int(math.Round(pct))
	text := fmt.Sprintf("%d%%", rounded)

	color := ctx.Theme.Success

	if rounded >= 90 {
		color = ctx.Theme.Danger + ansi.Bold
	} else if rounded >= 70 {
		color = ctx.Theme.Warning
	}

	return color + text + ansi.RST
}
