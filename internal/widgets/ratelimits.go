package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/usage"
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

// ModelUsageWidget displays the per-model weekly rate limits (e.g. Fable): name + bar + percent + (reset time).
//
// Claude Code does not send these windows in the statusline payload, so they come from the
// on-disk cache refreshed in the background by the usage package.
type ModelUsageWidget struct{}

func (w *ModelUsageWidget) ID() string { return "model-usage" }

func (w *ModelUsageWidget) Render(ctx *Context) string {
	models := usage.Models(ctx.Input.Version)

	if len(models) == 0 {
		return ""
	}

	parts := make([]string, 0, len(models))

	for _, model := range models {
		limit := &parser.RateLimit{UsedPercentage: model.Percent, ResetsAt: model.ResetsAt}
		parts = append(parts, renderRateLimit(model.DisplayName, limit, ctx))
	}

	return strings.Join(parts, " ")
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
