package widgets

import (
	"fmt"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
)

// DurationWidget displays the session duration.
type DurationWidget struct{}

func (w *DurationWidget) ID() string { return "duration" }

func (w *DurationWidget) Render(ctx *Context) string {
	ms := ctx.Input.Cost.TotalDurationMS
	totalSec := ms / 1000
	hours := totalSec / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60

	var text string

	if hours > 0 {
		text = fmt.Sprintf("%dh%02dm", hours, minutes)
	} else if minutes > 0 {
		text = fmt.Sprintf("%dm%02ds", minutes, seconds)
	} else {
		text = fmt.Sprintf("%ds", seconds)
	}

	color := ctx.Theme.Success
	secFloat := float64(totalSec)

	if secFloat >= ctx.Config.Thresholds.Duration.Red {
		color = ctx.Theme.Danger
	} else if secFloat >= ctx.Config.Thresholds.Duration.Yellow {
		color = ctx.Theme.Warning
	}

	return color + text + ansi.RST
}
