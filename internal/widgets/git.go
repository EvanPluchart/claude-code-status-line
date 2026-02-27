package widgets

import (
	"os/exec"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

func gitCommand(cwd string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func projectDir(ctx *Context) string {
	if ctx.Input.Workspace.ProjectDir != "" {
		return ctx.Input.Workspace.ProjectDir
	}

	return ctx.Input.CWD
}

// GitBranchWidget displays the current git branch.
type GitBranchWidget struct{}

func (w *GitBranchWidget) ID() string { return "git-branch" }

func (w *GitBranchWidget) Render(ctx *Context) string {
	branch, err := gitCommand(projectDir(ctx), "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil || branch == "" {
		return ""
	}

	return ansi.Colorize(branch, ctx.Theme.Info)
}

// GitStatusWidget shows a dirty/clean indicator.
type GitStatusWidget struct{}

func (w *GitStatusWidget) ID() string { return "git-status" }

func (w *GitStatusWidget) Render(ctx *Context) string {
	status, err := gitCommand(projectDir(ctx), "status", "--porcelain", "--no-optional-locks")
	if err != nil {
		return ""
	}

	if status == "" {
		return ansi.Colorize("\u2713", ctx.Theme.Success)
	}

	return ansi.Colorize("\u2717", ctx.Theme.Warning)
}

// NestedReposWidget counts nested git repositories.
type NestedReposWidget struct{}

func (w *NestedReposWidget) ID() string { return "nested-repos" }

func (w *NestedReposWidget) Render(ctx *Context) string {
	out, err := gitCommand(projectDir(ctx), "rev-parse", "--show-toplevel")
	if err != nil || out == "" {
		return ""
	}

	cmd := exec.Command("find", ".", "-maxdepth", "3", "-name", ".git", "-type", "d")
	cmd.Dir = projectDir(ctx)

	findOut, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(findOut)), "\n")
	count := len(lines) - 1 // subtract root repo

	if count <= 0 {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)
	label := t.RepoPlural

	if count == 1 {
		label = t.RepoSingular
	}

	return ansi.Colorize(strings.Join([]string{itoa(count), label}, " "), ctx.Theme.Muted)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}

	return s
}
