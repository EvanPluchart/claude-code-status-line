package widgets

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
)

// TimestampWidget displays the current time.
type TimestampWidget struct{}

func (w *TimestampWidget) ID() string { return "timestamp" }

func (w *TimestampWidget) Render(ctx *Context) string {
	now := time.Now()
	format := "15:04"

	if ctx.Config.Widgets.Timestamp.ShowSeconds {
		format = "15:04:05"
	}

	return ansi.Colorize(now.Format(format), ctx.Theme.Muted)
}

// OSInfoWidget displays the OS and architecture.
type OSInfoWidget struct{}

func (w *OSInfoWidget) ID() string { return "os-info" }

func (w *OSInfoWidget) Render(ctx *Context) string {
	osNames := map[string]string{
		"darwin":  "macOS",
		"linux":   "Linux",
		"windows": "Windows",
		"freebsd": "FreeBSD",
	}

	name := osNames[runtime.GOOS]
	if name == "" {
		name = runtime.GOOS
	}

	return ansi.Colorize(name+" "+runtime.GOARCH, ctx.Theme.Muted)
}

// SeparatorWidget displays a visual separator.
type SeparatorWidget struct{}

func (w *SeparatorWidget) ID() string { return "separator" }

func (w *SeparatorWidget) Render(ctx *Context) string {
	char := ctx.Config.Widgets.Separator.Char
	if char == "" {
		char = "\u2502"
	}

	return ansi.Colorize(" "+char+" ", ctx.Theme.Separator)
}

// SpacerWidget adds flexible space.
type SpacerWidget struct{}

func (w *SpacerWidget) ID() string { return "spacer" }

func (w *SpacerWidget) Render(_ *Context) string { return " " }

// VimModeWidget displays the current vim mode.
type VimModeWidget struct{}

func (w *VimModeWidget) ID() string { return "vim-mode" }

func (w *VimModeWidget) Render(ctx *Context) string {
	if ctx.Input.Vim == nil {
		return ""
	}

	mode := strings.ToUpper(ctx.Input.Vim.Mode)
	color := ctx.Theme.Info

	if mode == "INSERT" {
		color = ctx.Theme.Success
	}

	return ansi.ColorBold(mode, color)
}

// LinesChangedWidget displays lines added/removed.
type LinesChangedWidget struct{}

func (w *LinesChangedWidget) ID() string { return "lines-changed" }

func (w *LinesChangedWidget) Render(ctx *Context) string {
	added := ctx.Input.Cost.TotalLinesAdded
	removed := ctx.Input.Cost.TotalLinesRemoved

	if added == 0 && removed == 0 {
		return ""
	}

	return ansi.Colorize(fmt.Sprintf("+%d", added), ctx.Theme.Success) +
		" " +
		ansi.Colorize(fmt.Sprintf("-%d", removed), ctx.Theme.Danger)
}
