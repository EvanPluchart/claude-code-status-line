package widgets

import (
	"fmt"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

// AgentWidget displays the active agent name (when running with --agent).
type AgentWidget struct{}

func (w *AgentWidget) ID() string { return "agent" }

func (w *AgentWidget) Render(ctx *Context) string {
	if ctx.Input.Agent == nil || ctx.Input.Agent.Name == "" {
		return ""
	}

	return ansi.Colorize("@"+ctx.Input.Agent.Name, ctx.Theme.Secondary)
}

// OutputStyleWidget displays the active output style (hidden for "default").
type OutputStyleWidget struct{}

func (w *OutputStyleWidget) ID() string { return "output-style" }

func (w *OutputStyleWidget) Render(ctx *Context) string {
	if ctx.Input.OutputStyle == nil {
		return ""
	}

	name := ctx.Input.OutputStyle.Name

	if name == "" || strings.EqualFold(name, "default") {
		return ""
	}

	return ansi.Colorize(name, ctx.Theme.Muted)
}

// effortColors ranks effort levels from cheapest to most expensive.
var effortRank = map[string]int{
	"low": 0, "medium": 1, "high": 2, "xhigh": 3, "max": 4,
}

// EffortWidget displays the reasoning effort level.
type EffortWidget struct{}

func (w *EffortWidget) ID() string { return "effort" }

func (w *EffortWidget) Render(ctx *Context) string {
	if ctx.Input.Effort == nil || ctx.Input.Effort.Level == "" {
		return ""
	}

	level := strings.ToLower(ctx.Input.Effort.Level)
	color := ctx.Theme.Muted

	switch effortRank[level] {
	case 0, 1:
		color = ctx.Theme.Success
	case 2:
		color = ctx.Theme.Info
	case 3:
		color = ctx.Theme.Warning
	case 4:
		color = ctx.Theme.Danger
	}

	return ansi.Colorize("⚙ "+level, color)
}

// ThinkingWidget shows whether extended thinking is enabled.
type ThinkingWidget struct{}

func (w *ThinkingWidget) ID() string { return "thinking" }

func (w *ThinkingWidget) Render(ctx *Context) string {
	if ctx.Input.Thinking == nil || !ctx.Input.Thinking.Enabled {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)

	return ansi.Colorize("✦ "+t.ThinkingLabel, ctx.Theme.Info)
}

// FastModeWidget shows a marker when fast mode is active.
type FastModeWidget struct{}

func (w *FastModeWidget) ID() string { return "fast-mode" }

func (w *FastModeWidget) Render(ctx *Context) string {
	if !ctx.Input.FastMode {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)

	return ansi.ColorBold("⚡ "+t.FastLabel, ctx.Theme.Warning)
}

// PRWidget displays the open pull/merge request for the current branch.
type PRWidget struct{}

func (w *PRWidget) ID() string { return "pr" }

func (w *PRWidget) Render(ctx *Context) string {
	pr := ctx.Input.PR

	if pr == nil || pr.Number <= 0 {
		return ""
	}

	t := i18n.Get(ctx.Config.Locale)
	prefix := "PR"

	if pr.Kind == "mr" {
		prefix = "MR"
	}

	text := fmt.Sprintf("%s #%d", prefix, pr.Number)
	color := ctx.Theme.Info

	switch strings.ToLower(pr.ReviewState) {
	case "approved":
		color = ctx.Theme.Success
		text += " ✓"
	case "changes_requested":
		color = ctx.Theme.Danger
		text += " " + t.PRChangesRequested
	case "draft":
		color = ctx.Theme.Muted
		text += " " + t.PRDraft
	case "pending":
		color = ctx.Theme.Warning
	}

	return ansi.Colorize(text, color)
}

// WorktreeWidget displays the active worktree name.
type WorktreeWidget struct{}

func (w *WorktreeWidget) ID() string { return "worktree" }

func (w *WorktreeWidget) Render(ctx *Context) string {
	name := ""

	if ctx.Input.Worktree != nil {
		name = ctx.Input.Worktree.Name
	}

	if name == "" {
		name = ctx.Input.Workspace.GitWorktree
	}

	if name == "" {
		return ""
	}

	return ansi.Colorize("⌥ "+name, ctx.Theme.Secondary)
}

// SessionNameWidget displays the session name or title.
type SessionNameWidget struct{}

func (w *SessionNameWidget) ID() string { return "session-name" }

func (w *SessionNameWidget) Render(ctx *Context) string {
	name := strings.TrimSpace(ctx.Input.SessionName)

	if name == "" {
		return ""
	}

	const maxLen = 40

	if len([]rune(name)) > maxLen {
		name = string([]rune(name)[:maxLen-1]) + "…"
	}

	return ansi.Colorize(name, ctx.Theme.Text)
}

// RepoWidget displays the remote repository as owner/name.
type RepoWidget struct{}

func (w *RepoWidget) ID() string { return "repo" }

func (w *RepoWidget) Render(ctx *Context) string {
	repo := ctx.Input.Workspace.Repo

	if repo == nil || repo.Name == "" {
		return ""
	}

	text := repo.Name

	if repo.Owner != "" {
		text = repo.Owner + "/" + repo.Name
	}

	return ansi.Colorize(text, ctx.Theme.Secondary)
}

// ClaudeVersionWidget displays the Claude Code version.
type ClaudeVersionWidget struct{}

func (w *ClaudeVersionWidget) ID() string { return "claude-version" }

func (w *ClaudeVersionWidget) Render(ctx *Context) string {
	if ctx.Input.Version == "" {
		return ""
	}

	return ansi.Colorize("v"+strings.TrimPrefix(ctx.Input.Version, "v"), ctx.Theme.Muted)
}
