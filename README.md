# claude-code-status-line

A fully customizable, multi-line statusline for [Claude Code](https://code.claude.com/docs/en/statusline).

```
 Opus 5 │ my-project │ feature/auth ✓ │ 12m34s │ $1.23
 ━━━━━━━━━━━━──── 58% (116k/200k) │ ⚙ high ✦ think
```

## Features

- **3 configurable lines** with any combination of widgets
- **37 built-in widgets**: model, effort, thinking, fast mode, git (branch, changes, ahead/behind, PR), cost, burn rate, tokens, rate limits, worktree, agent, and more
- **6 themes**: default, minimal, neon, dracula, catppuccin, nord
- **Multi-currency**: USD, EUR, GBP, JPY, CAD, and more (rates refreshed daily in the background)
- **i18n**: English and French out of the box
- **Cross-platform**: macOS (Apple Silicon + Intel), Linux, Windows 10/11, WSL2
- **Blazing fast**: single Go binary, a few ms execution (well within Claude Code's 300ms budget)
- **Zero-config start**: works with sensible defaults
- **Live preview**: `claude-code-status-line preview` renders your config without launching Claude Code

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap EvanPluchart/tap
brew trust EvanPluchart/tap   # recent Homebrew versions require approving third-party taps
brew install claude-code-status-line
```

### Go Install

```bash
go install github.com/EvanPluchart/claude-code-status-line/cmd/claude-code-status-line@latest
```

### Shell Script (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/EvanPluchart/claude-code-status-line/main/install.sh | sh
```

### PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/EvanPluchart/claude-code-status-line/main/install.ps1 | iex
```

### GitHub Releases

Download pre-built binaries from [GitHub Releases](https://github.com/EvanPluchart/claude-code-status-line/releases).

## Quick Start

```bash
# Interactive wizard: config + registration in Claude Code
claude-code-status-line init

# Non-interactive
claude-code-status-line init --locale fr --theme dracula --currency EUR

# Preview the result right away
claude-code-status-line preview
```

That's it! Restart Claude Code and your new statusline will appear.

## Available Widgets

Run `claude-code-status-line widgets` to list them from the CLI.

### Session

| Widget | Description | Example |
|--------|-------------|---------|
| `model` | Current Claude model (derived from the model ID) | `Opus 5` / `O5` |
| `effort` | Reasoning effort level | `⚙ high` |
| `thinking` | Extended thinking enabled | `✦ think` |
| `fast-mode` | Fast mode active | `⚡ fast` |
| `agent` | Active agent (`--agent`) | `@reviewer` |
| `output-style` | Output style (hidden when `default`) | `Explanatory` |
| `session-name` | Session name or title | `Add OAuth login` |
| `claude-version` | Claude Code version | `v2.1.90` |
| `vim-mode` | Vim mode indicator | `NORMAL` / `INSERT` |

### Workspace & git

| Widget | Description | Example |
|--------|-------------|---------|
| `directory` | Project directory (`widgets.directory.depth` segments) | `my-project` |
| `repo` | Remote repository | `acme/my-project` |
| `worktree` | Active worktree | `⌥ feature-x` |
| `git-branch` | Current git branch (truncated past `max_branch_length`) | `feature/auth` |
| `git-status` | Dirty/clean indicator | `✓` / `✗` |
| `git-changes` | Staged / modified / untracked counts | `●2 ✚1 …3` |
| `git-ahead-behind` | Commits ahead/behind upstream | `↑2 ↓1` / `≡` |
| `pr` | Open PR/MR and review state | `PR #42 ✓` |
| `nested-repos` | Nested git repos count | `3 repos` |
| `lines-changed` | Lines added/removed (`source: git` or `session`) | `+42 -7` |

### Cost & time

| Widget | Description | Example |
|--------|-------------|---------|
| `cost` | Session cost (multi-currency) | `$1.23` / `1,23€` |
| `burn-rate` | Spending rate per hour | `$2.40/h` |
| `duration` | Session duration | `12m34s` |
| `api-time` | Time spent waiting for the API | `API 4m10s (33%)` |
| `timestamp` | Current time | `14:32` |

### Context & tokens

| Widget | Description | Example |
|--------|-------------|---------|
| `token-bar` | Context usage progress bar | `━━━━━━━━────` |
| `context-percent` | Context usage % | `58%` |
| `token-count` | Tokens used/max | `(116k/200k)` |
| `context-remaining` | Tokens left in the context | `84k left` |
| `total-tokens` | Total I/O tokens | `↑120k ↓45k` |
| `cache-ratio` | Cache hit ratio | `Cache: 87%` |
| `exceeds-200k` | Marker when the conversation exceeds 200k tokens | `>200k` |
| `session-usage` | 5-hour rate limit (Pro/Max) | `5h ━━━━──── 23% (reset 3h42m)` |
| `weekly-usage` | 7-day rate limit (Pro/Max) | `7d ━━━━━━── 41% (reset 4d08h)` |
| `model-usage` | Per-model weekly limits (see note below) | `Fable ━─────── 6% (reset 5d12h)` |

> **`model-usage` requires an extra opt-in.** Claude Code only sends the 5-hour and 7-day windows
> to statusline commands — the per-model windows (the separate *Fable* bar on
> [claude.ai usage](https://claude.ai/settings/usage)) are not part of the payload. When you add this
> widget to your config, the statusline reads your local Claude Code OAuth token (macOS Keychain, or
> `~/.claude/.credentials.json`) and queries `api.anthropic.com/api/oauth/usage` — the same endpoint
> Claude Code's own `/usage` panel calls — from a detached background process, at most once every
> 5 minutes. The token never leaves your machine except to Anthropic, and is never written to the
> cache (`~/.claude-statusline/usage.json`, mode `0600`) nor logged. No other widget reads
> credentials, and the render path itself never makes network calls. Set `CLAUDE_STATUSLINE_OFFLINE=1`
> to disable all background refreshes.

### Layout & system

| Widget | Description | Example |
|--------|-------------|---------|
| `separator` | Visual separator | `│` |
| `spacer` | Single space | |
| `os-info` | Platform info | `macOS arm64` |
| `hostname` | Machine hostname | `macbook` |

Widgets that have nothing to show (no PR, no agent, thinking off…) render nothing, and the surrounding separators collapse automatically.

## Configuration

Configuration file: `~/.claude-statusline/config.yml` (override the directory with `CLAUDE_STATUSLINE_DIR`).

```yaml
# Language: en | fr
locale: "en"

# Theme: default | minimal | neon | dracula | catppuccin | nord
theme: "default"

# Lines configuration (3 lines max)
lines:
  - widgets: [model, separator, directory, separator, git-branch, git-status, separator, duration, separator, cost]
  - widgets: [token-bar, context-percent, token-count, separator, effort, thinking, fast-mode]
  - widgets: []

# Widget options
widgets:
  cost:
    currency: "USD"
    decimals: 2
  token_bar:
    width: 16
    filled_char: "━"
    empty_char: "─"
  separator:
    char: "│"
  model:
    short_name: false
  timestamp:
    show_seconds: false
  directory:
    depth: 1
  git:
    max_branch_length: 40
  lines_changed:
    source: "git"   # git = working tree diff, session = lines changed by Claude Code

# Color thresholds (green below yellow, then yellow / orange / red)
thresholds:
  context:
    yellow: 60
    orange: 80
    red: 90
  cost:
    yellow: 0.25
    orange: 1.0
    red: 5.0
  duration:
    yellow: 60
    orange: 600
    red: 1800
  rate_limit:
    yellow: 60
    orange: 80
    red: 90
```

Only the keys you set are overridden; everything else keeps its default.

### CLI

```bash
claude-code-status-line config              # interactive wizard with live preview
claude-code-status-line config edit         # open config in $EDITOR
claude-code-status-line config set theme nord
claude-code-status-line config get currency
claude-code-status-line config reset
claude-code-status-line preview             # render with sample data
claude-code-status-line preview --stdin < payload.json
claude-code-status-line widgets             # list widget IDs
claude-code-status-line update              # self-update (non-Homebrew installs)
claude-code-status-line update-rates        # force an exchange rates refresh
```

## Themes

| Theme | Style |
|-------|-------|
| `default` | Adaptive colors based on context usage |
| `minimal` | Monochrome, subtle |
| `neon` | Bright, vibrant |
| `dracula` | Dracula color palette |
| `catppuccin` | Catppuccin Mocha |
| `nord` | Nord palette |

## How it works

Claude Code runs the binary with a JSON payload on stdin (model, cost, context window, rate limits, PR, worktree…). The binary reads `~/.claude-statusline/config.yml`, renders the configured widgets and prints up to 3 lines. The `init` command registers it in `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "claude-code-status-line",
    "padding": 0
  }
}
```

When the `vim-mode` widget is enabled, `hideVimModeIndicator` is also set so the mode is not displayed twice.

Rendering never performs network calls, so the statusline stays within its execution budget. Two
caches are refreshed by detached background processes instead: exchange rates (`update-rates`, daily)
and, only when the `model-usage` widget is enabled, per-model usage limits (`update-usage`, every
5 minutes). Both can be run by hand, and both are disabled by `CLAUDE_STATUSLINE_OFFLINE=1`.

## Uninstall

```bash
claude-code-status-line uninstall
```

This restores your previous statusline configuration (if any) and removes the config directory.

## Requirements

- [Claude Code](https://code.claude.com/docs/en/overview) CLI
- `git` on the PATH for the git widgets

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
