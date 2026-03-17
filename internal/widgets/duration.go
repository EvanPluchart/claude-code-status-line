package widgets

import (
	"fmt"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

// DurationWidget displays the session duration.
type DurationWidget struct{}

func (w *DurationWidget) ID() string { return "duration" }

func (w *DurationWidget) Render(ctx *Context) string {
	t := i18n.Get(ctx.Config.Locale)
	ms := ctx.Input.Cost.TotalDurationMS
	totalSec := ms / 1000

	months := totalSec / 2592000  // 30 days
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

	if len(parts) > 0 {
		parts = append(parts, fmt.Sprintf("%02d%s", minutes, t.DurationMinutes))
	} else if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d%s%02d%s", minutes, t.DurationMinutes, seconds, t.DurationSeconds))
	} else {
		parts = append(parts, fmt.Sprintf("%d%s", seconds, t.DurationSeconds))
	}

	text := strings.Join(parts, " ")

	color := ctx.Theme.Success
	secFloat := float64(totalSec)

	if secFloat >= ctx.Config.Thresholds.Duration.Red {
		color = ctx.Theme.Danger
	} else if secFloat >= ctx.Config.Thresholds.Duration.Yellow {
		color = ctx.Theme.Warning
	}

	return color + text + ansi.RST
}
