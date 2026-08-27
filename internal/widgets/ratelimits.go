package widgets

import (
	"fmt"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
)

// SessionUsageWidget displays the 5-hour session rate limit: label + bar + percent + (reset time).
type SessionUsageWidget struct{}

func (w *SessionUsageWidget) ID() string { return "session-usage" }

func (w *SessionUsageWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.FiveHour == nil {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)

	return renderRateLimit(t.SessionLabel, ctx.Input.RateLimits.FiveHour, ctx)
}

// WeeklyUsageWidget displays the 7-day rate limit: label + bar + percent + (reset time).
type WeeklyUsageWidget struct{}

func (w *WeeklyUsageWidget) ID() string { return "weekly-usage" }

func (w *WeeklyUsageWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.SevenDay == nil {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)

	return renderRateLimit(t.WeeklyLabel, ctx.Input.RateLimits.SevenDay, ctx)
}

// renderRateLimit renders: label bar percent (reset time).
func renderRateLimit(label string, rl *parser.RateLimit, ctx *Context) string {
	t := i18n.Get(ctx.Config.Locale)
	thresholds := ctx.Config.Thresholds.RateLimit

	out := ansi.Colorize(label, ctx.Theme.Muted) + " " +
		renderBar(rl.UsedPercentage, thresholds, ctx) + " " +
		renderPercent(rl.UsedPercentage, thresholds, ctx)

	if rl.ResetsAt > 0 {
		remaining := time.Until(time.Unix(rl.ResetsAt, 0))

		if remaining > 0 {
			reset := fmt.Sprintf("(%s %s)", t.ResetsIn, formatCompactDuration(remaining, t))
			out += " " + ansi.Colorize(reset, ctx.Theme.Muted)
		}
	}

	return out
}
