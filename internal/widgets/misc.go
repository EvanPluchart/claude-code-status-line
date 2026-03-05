package widgets

import (
	"fmt"
	"os/exec"
	"path/filepath"
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

// LinesChangedWidget displays lines added/removed from git diff.
// It sums changes across the root repo and any nested git repos.
type LinesChangedWidget struct{}

func (w *LinesChangedWidget) ID() string { return "lines-changed" }

func (w *LinesChangedWidget) Render(ctx *Context) string {
	dir := projectDir(ctx)
	added, removed := gitDiffStats(dir)

	// Find nested repos and sum their stats
	nestedDirs := findNestedRepos(dir)

	for _, nested := range nestedDirs {
		a, r := gitDiffStats(nested)
		added += a
		removed += r
	}

	if added == 0 && removed == 0 {
		return ansi.Colorize("+0", ctx.Theme.Muted) + " " + ansi.Colorize("-0", ctx.Theme.Muted)
	}

	return ansi.Colorize(fmt.Sprintf("+%d", added), ctx.Theme.Success) +
		" " +
		ansi.Colorize(fmt.Sprintf("-%d", removed), ctx.Theme.Danger)
}

func gitDiffStats(dir string) (int, int) {
	out, err := gitCommand(dir, "--no-optional-locks", "diff", "--numstat")
	if err != nil || out == "" {
		return 0, 0
	}

	added := 0
	removed := 0

	for _, line := range strings.Split(out, "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		// Binary files show "-" instead of numbers
		if parts[0] == "-" || parts[1] == "-" {
			continue
		}

		a := atoi(parts[0])
		r := atoi(parts[1])
		added += a
		removed += r
	}

	return added, removed
}

func findNestedRepos(rootDir string) []string {
	cmd := exec.Command("find", ".", "-maxdepth", "3", "-name", ".git", "-type", "d")
	cmd.Dir = rootDir

	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var dirs []string

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" || line == "./.git" {
			continue
		}

		// Strip "/.git" suffix to get the repo dir
		repoDir := strings.TrimSuffix(line, "/.git")

		if repoDir == "." {
			continue
		}

		dirs = append(dirs, filepath.Join(rootDir, repoDir))
	}

	return dirs
}

func atoi(s string) int {
	n := 0

	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}

		n = n*10 + int(c-'0')
	}

	return n
}
