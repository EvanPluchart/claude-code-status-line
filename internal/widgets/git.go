package widgets

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

// nestedRepoMaxDepth bounds the directory walk used to find nested repositories.
const nestedRepoMaxDepth = 3

// skippedDirs are never descended into when looking for nested repositories.
var skippedDirs = map[string]bool{
	"node_modules": true, "vendor": true, ".cache": true, "dist": true,
	"build": true, "target": true, ".venv": true, "venv": true,
}

func gitCommand(cwd string, args ...string) (string, error) {
	fullArgs := append([]string{"--no-optional-locks"}, args...)
	cmd := exec.Command("git", fullArgs...)
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

// findNestedRepos returns nested git repository directories (excluding the root itself).
// It is implemented with a bounded filesystem walk so that it works on every platform.
func findNestedRepos(rootDir string) []string {
	var dirs []string

	_ = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}

		rel, relErr := filepath.Rel(rootDir, path)
		if relErr != nil || rel == "." {
			return nil
		}

		name := d.Name()

		if name == ".git" {
			if repo := filepath.Dir(path); repo != rootDir {
				dirs = append(dirs, repo)
			}

			return filepath.SkipDir
		}

		if skippedDirs[name] || (strings.HasPrefix(name, ".") && name != ".claude") {
			return filepath.SkipDir
		}

		if strings.Count(rel, string(filepath.Separator))+1 >= nestedRepoMaxDepth {
			return filepath.SkipDir
		}

		return nil
	})

	return dirs
}

// truncateBranch shortens long branch names, keeping the beginning and the end.
func truncateBranch(branch string, maxLen int) string {
	if maxLen <= 0 || len(branch) <= maxLen {
		return branch
	}

	if maxLen <= 3 {
		return branch[:maxLen]
	}

	keep := maxLen - 1
	head := keep / 2
	tail := keep - head

	return branch[:head] + "…" + branch[len(branch)-tail:]
}

// GitBranchWidget displays the current git branch.
type GitBranchWidget struct{}

func (w *GitBranchWidget) ID() string { return "git-branch" }

func (w *GitBranchWidget) Render(ctx *Context) string {
	branch, err := gitCommand(projectDir(ctx), "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil || branch == "" {
		t := i18n.Get(ctx.Config.Locale)

		return ansi.Colorize(t.NoGitRepo, ctx.Theme.Muted)
	}

	if branch == "HEAD" {
		sha, shaErr := gitCommand(projectDir(ctx), "rev-parse", "--short", "HEAD")

		if shaErr == nil && sha != "" {
			branch = "@" + sha
		}
	}

	return ansi.Colorize(truncateBranch(branch, ctx.Config.Widgets.Git.MaxBranchLength), ctx.Theme.Info)
}

// GitStatusWidget shows a dirty/clean indicator.
type GitStatusWidget struct{}

func (w *GitStatusWidget) ID() string { return "git-status" }

func (w *GitStatusWidget) Render(ctx *Context) string {
	status, err := gitCommand(projectDir(ctx), "status", "--porcelain")
	if err != nil {
		return ansi.Colorize("-", ctx.Theme.Muted)
	}

	if status == "" {
		return ansi.Colorize("✓", ctx.Theme.Success)
	}

	return ansi.Colorize("✗", ctx.Theme.Warning)
}

// GitChangesWidget shows staged / modified / untracked file counts (●2 ✚1 …3).
type GitChangesWidget struct{}

func (w *GitChangesWidget) ID() string { return "git-changes" }

func (w *GitChangesWidget) Render(ctx *Context) string {
	status, err := gitCommand(projectDir(ctx), "status", "--porcelain")
	if err != nil {
		return ""
	}

	staged, modified, untracked := countPorcelain(status)

	if staged == 0 && modified == 0 && untracked == 0 {
		return ansi.Colorize("✓", ctx.Theme.Success)
	}

	var parts []string

	if staged > 0 {
		parts = append(parts, ansi.Colorize(fmt.Sprintf("●%d", staged), ctx.Theme.Success))
	}

	if modified > 0 {
		parts = append(parts, ansi.Colorize(fmt.Sprintf("✚%d", modified), ctx.Theme.Warning))
	}

	if untracked > 0 {
		parts = append(parts, ansi.Colorize(fmt.Sprintf("…%d", untracked), ctx.Theme.Muted))
	}

	return strings.Join(parts, " ")
}

// countPorcelain parses `git status --porcelain` output into staged, modified and untracked counts.
func countPorcelain(status string) (int, int, int) {
	staged, modified, untracked := 0, 0, 0

	for _, line := range strings.Split(status, "\n") {
		if len(line) < 2 {
			continue
		}

		index, worktree := line[0], line[1]

		if index == '?' && worktree == '?' {
			untracked++

			continue
		}

		if index != ' ' && index != '?' {
			staged++
		}

		if worktree != ' ' && worktree != '?' {
			modified++
		}
	}

	return staged, modified, untracked
}

// GitAheadBehindWidget shows commits ahead/behind the upstream branch (↑2 ↓1).
type GitAheadBehindWidget struct{}

func (w *GitAheadBehindWidget) ID() string { return "git-ahead-behind" }

func (w *GitAheadBehindWidget) Render(ctx *Context) string {
	out, err := gitCommand(projectDir(ctx), "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if err != nil {
		return ""
	}

	fields := strings.Fields(out)
	if len(fields) != 2 {
		return ""
	}

	ahead, _ := strconv.Atoi(fields[0])
	behind, _ := strconv.Atoi(fields[1])

	if ahead == 0 && behind == 0 {
		return ansi.Colorize("≡", ctx.Theme.Muted)
	}

	var parts []string

	if ahead > 0 {
		parts = append(parts, ansi.Colorize(fmt.Sprintf("↑%d", ahead), ctx.Theme.Success))
	}

	if behind > 0 {
		parts = append(parts, ansi.Colorize(fmt.Sprintf("↓%d", behind), ctx.Theme.Warning))
	}

	return strings.Join(parts, " ")
}

// NestedReposWidget counts nested git repositories.
type NestedReposWidget struct{}

func (w *NestedReposWidget) ID() string { return "nested-repos" }

func (w *NestedReposWidget) Render(ctx *Context) string {
	t := i18n.Get(ctx.Config.Locale)
	count := len(findNestedRepos(projectDir(ctx)))

	if count <= 0 {
		return ansi.Colorize(t.NoNestedRepos, ctx.Theme.Muted)
	}

	label := t.RepoPlural

	if count == 1 {
		label = t.RepoSingular
	}

	return ansi.Colorize(fmt.Sprintf("%d %s", count, label), ctx.Theme.Muted)
}
