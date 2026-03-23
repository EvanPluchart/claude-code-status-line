package widgets

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

// SessionUsageWidget displays the 5-hour session rate limit: label + bar + percent + (reset time).
type SessionUsageWidget struct{}

func (w *SessionUsageWidget) ID() string { return "session-usage" }

func (w *SessionUsageWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.FiveHour == nil {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)
	rl := ctx.Input.RateLimits.FiveHour

	return renderRateLimitWidget(t.SessionLabel, rl.UsedPercentage, rl.ResetsAt, ctx)
}

// WeeklyUsageWidget displays the 7-day rate limit: label + bar + percent + (reset time).
type WeeklyUsageWidget struct{}

func (w *WeeklyUsageWidget) ID() string { return "weekly-usage" }

func (w *WeeklyUsageWidget) Render(ctx *Context) string {
	if ctx.Input.RateLimits == nil || ctx.Input.RateLimits.SevenDay == nil {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)
	rl := ctx.Input.RateLimits.SevenDay

	return renderRateLimitWidget(t.WeeklyLabel, rl.UsedPercentage, rl.ResetsAt, ctx)
}

// renderRateLimitWidget renders: label bar percent (reset time)
func renderRateLimitWidget(label string, pct float64, resetsAt int64, ctx *Context) string {
	t := i18n.Get(ctx.Config.Locale)

	// Label
	labelStr := ansi.Colorize(label, ctx.Theme.Muted)

	// Bar
	barStr := renderRateLimitBar(pct, ctx)

	// Percent
	pctStr := renderRateLimitPercent(pct, ctx)

	// Reset time
	resetStr := ""

	if resetsAt > 0 {
		remaining := time.Until(time.Unix(resetsAt, 0))

		if remaining > 0 {
			resetStr = " " + ansi.Colorize(fmt.Sprintf("(%s %s)", t.ResetsIn, formatDuration(remaining, t)), ctx.Theme.Muted)
		}
	}

	return labelStr + " " + barStr + " " + pctStr + resetStr
}

// formatDuration formats a duration using i18n units.
func formatDuration(d time.Duration, t i18n.Translations) string {
	totalSec := int64(d.Seconds())

	days := totalSec / 86400
	hours := (totalSec % 86400) / 3600
	minutes := (totalSec % 3600) / 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d%s", days, t.DurationDays))
	}

	if hours > 0 || days > 0 {
		parts = append(parts, fmt.Sprintf("%d%s", hours, t.DurationHours))
	}

	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d%s", minutes, t.DurationMinutes))
	} else {
		parts = append(parts, fmt.Sprintf("%02d%s", minutes, t.DurationMinutes))
	}

	return strings.Join(parts, "")
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
