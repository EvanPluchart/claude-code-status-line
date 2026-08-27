package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

// formatSessionDuration renders a duration as "1h 05min", "12m34s", "3w 2d 4h 00m"...
func formatSessionDuration(d time.Duration, t i18n.Translations) string {
	totalSec := int64(d.Seconds())

	if totalSec < 0 {
		totalSec = 0
	}

	months := totalSec / 2592000 // 30 days
	weeks := (totalSec % 2592000) / 604800
	days := (totalSec % 604800) / 86400
	hours := (totalSec % 86400) / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60

	var parts []string

	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d%s", months, t.DurationMonths))
	}

	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%d%s", weeks, t.DurationWeeks))
	}

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d%s", days, t.DurationDays))
	}

	if hours > 0 || len(parts) > 0 {
		parts = append(parts, fmt.Sprintf("%d%s", hours, t.DurationHours))
	}

	switch {
	case len(parts) > 0:
		parts = append(parts, fmt.Sprintf("%02d%s", minutes, t.DurationMinutes))
	case minutes > 0:
		parts = append(parts, fmt.Sprintf("%d%s%02d%s", minutes, t.DurationMinutes, seconds, t.DurationSeconds))
	default:
		parts = append(parts, fmt.Sprintf("%d%s", seconds, t.DurationSeconds))
	}

	return strings.Join(parts, " ")
}

// formatCompactDuration renders a duration as "3h42m", "4d08h05m" or "12m" (used for reset timers).
func formatCompactDuration(d time.Duration, t i18n.Translations) string {
	totalSec := int64(d.Seconds())

	if totalSec < 0 {
		totalSec = 0
	}

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

// DurationWidget displays the session duration.
type DurationWidget struct{}

func (w *DurationWidget) ID() string { return "duration" }

func (w *DurationWidget) Render(ctx *Context) string {
	t := i18n.Get(ctx.Config.Locale)
	d := time.Duration(ctx.Input.Cost.TotalDurationMS) * time.Millisecond
	text := formatSessionDuration(d, t)

	return ansi.Colorize(text, thresholdColor(d.Seconds(), ctx.Config.Thresholds.Duration, ctx))
}

// APITimeWidget displays the time spent waiting for the API and its share of the session.
type APITimeWidget struct{}

func (w *APITimeWidget) ID() string { return "api-time" }

func (w *APITimeWidget) Render(ctx *Context) string {
	apiMS := ctx.Input.Cost.TotalAPIDurationMS
	totalMS := ctx.Input.Cost.TotalDurationMS

	if apiMS <= 0 {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)
	text := t.APILabel + " " + formatSessionDuration(time.Duration(apiMS)*time.Millisecond, t)

	if totalMS > 0 {
		share := int(float64(apiMS) / float64(totalMS) * 100)

		if share > 100 {
			share = 100
		}

		text += fmt.Sprintf(" (%d%%)", share)
	}

	return ansi.Colorize(text, ctx.Theme.Muted)
}
