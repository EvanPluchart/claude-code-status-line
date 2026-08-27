package widgets

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
)

// TimestampWidget displays the current time.
type TimestampWidget struct{}

func (w *TimestampWidget) ID() string { return "timestamp" }

func (w *TimestampWidget) Render(ctx *Context) string {
	format := "15:04"

	if ctx.Config.Widgets.Timestamp.ShowSeconds {
		format = "15:04:05"
	}

	return ansi.Colorize(time.Now().Format(format), ctx.Theme.Muted)
}

var osNames = map[string]string{
	"darwin":  "macOS",
	"linux":   "Linux",
	"windows": "Windows",
	"freebsd": "FreeBSD",
}

// OSInfoWidget displays the OS and architecture.
type OSInfoWidget struct{}

func (w *OSInfoWidget) ID() string { return "os-info" }

func (w *OSInfoWidget) Render(ctx *Context) string {
	name := osNames[runtime.GOOS]

	if name == "" {
		name = runtime.GOOS
	}

	return ansi.Colorize(name+" "+runtime.GOARCH, ctx.Theme.Muted)
}

// HostnameWidget displays the machine hostname (useful over SSH).
type HostnameWidget struct{}

func (w *HostnameWidget) ID() string { return "hostname" }

func (w *HostnameWidget) Render(ctx *Context) string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return ""
	}

	if i := strings.IndexByte(host, '.'); i > 0 {
		host = host[:i]
	}

	return ansi.Colorize(host, ctx.Theme.Muted)
}

// SeparatorWidget displays a visual separator.
type SeparatorWidget struct{}

func (w *SeparatorWidget) ID() string { return "separator" }

func (w *SeparatorWidget) Render(ctx *Context) string {
	char := ctx.Config.Widgets.Separator.Char

	if char == "" {
		char = "│"
	}

	return ansi.Colorize(" "+char+" ", ctx.Theme.Separator)
}

// SpacerWidget adds a single space between widgets.
type SpacerWidget struct{}

func (w *SpacerWidget) ID() string { return "spacer" }

func (w *SpacerWidget) Render(_ *Context) string { return " " }

// VimModeWidget displays the current vim mode.
type VimModeWidget struct{}

func (w *VimModeWidget) ID() string { return "vim-mode" }

func (w *VimModeWidget) Render(ctx *Context) string {
	if ctx.Input.Vim == nil || ctx.Input.Vim.Mode == "" {
		return ""
	}

	mode := strings.ToUpper(ctx.Input.Vim.Mode)
	color := ctx.Theme.Info

	switch {
	case mode == "INSERT":
		color = ctx.Theme.Success
	case strings.HasPrefix(mode, "VISUAL"):
		color = ctx.Theme.Warning
	}

	return ansi.ColorBold(mode, color)
}

// LinesChangedWidget displays lines added/removed.
// Source "git" sums the working tree diff of the project and nested repos;
// source "session" uses the counters reported by Claude Code.
type LinesChangedWidget struct{}

func (w *LinesChangedWidget) ID() string { return "lines-changed" }

func (w *LinesChangedWidget) Render(ctx *Context) string {
	var added, removed int

	if ctx.Config.Widgets.LinesChanged.Source == "session" {
		added = ctx.Input.Cost.TotalLinesAdded
		removed = ctx.Input.Cost.TotalLinesRemoved
	} else {
		added, removed = workingTreeStats(projectDir(ctx))
	}

	if added == 0 && removed == 0 {
		return ansi.Colorize("+0 -0", ctx.Theme.Muted)
	}

	return ansi.Colorize(fmt.Sprintf("+%d", added), ctx.Theme.Success) +
		" " +
		ansi.Colorize(fmt.Sprintf("-%d", removed), ctx.Theme.Danger)
}

// workingTreeStats sums diff stats across the root repo and any nested repos.
func workingTreeStats(dir string) (int, int) {
	added, removed := gitDiffStats(dir)

	for _, nested := range findNestedRepos(dir) {
		a, r := gitDiffStats(nested)
		added += a
		removed += r
	}

	return added, removed
}

func gitDiffStats(dir string) (int, int) {
	out, err := gitCommand(dir, "diff", "--numstat")
	if err != nil || out == "" {
		return 0, 0
	}

	return parseNumstat(out)
}

// parseNumstat sums the added/removed columns of `git diff --numstat` output.
func parseNumstat(out string) (int, int) {
	added, removed := 0, 0

	for _, line := range strings.Split(out, "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		// Binary files show "-" instead of numbers.
		a, errA := strconv.Atoi(parts[0])
		r, errR := strconv.Atoi(parts[1])

		if errA != nil || errR != nil {
			continue
		}

		added += a
		removed += r
	}

	return added, removed
}
