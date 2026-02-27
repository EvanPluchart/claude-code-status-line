# Claude Code Status Line - Implementation Plan

> A fully customizable, multi-line statusline for Claude Code.

---

## 1. Project Overview

### 1.1 Vision

**claude-code-status-line** is an open-source, cross-platform, customizable statusline for [Claude Code](https://docs.anthropic.com/en/docs/claude-code). It replaces the default single-line status with a rich, configurable 3-line display featuring widgets that show real-time session information.

### 1.2 Key Features

- **3 configurable lines**: each line can contain any combination of widgets
- **Built-in widgets**: model, directory, git branch, git status, cost, tokens, progress bar, duration, vim mode, lines changed, and more
- **Themes**: predefined color themes (default, minimal, neon, dracula, catppuccin, nord)
- **i18n**: English and French out of the box, extensible for community translations
- **Currency support**: configurable currency for cost display (USD, EUR, GBP, JPY, etc.)
- **Cross-platform**: macOS (Apple Silicon + Intel), Linux, Windows 10/11, WSL2
- **Blazing fast**: single Go binary, ~5ms execution
- **Zero-config start**: works with sensible defaults
- **One-command uninstall**: cleanly removes config and restores previous statusline

### 1.3 Target Users

- Claude Code power users who want more visibility into their sessions
- Developers who want their terminal experience to feel polished and informative
- Teams that want consistent status information across members

---

## 2. Technical Architecture

### 2.1 How Claude Code Statusline Works

Claude Code supports custom statuslines via the `statusLine` setting in `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "claude-code-status-line"
  }
}
```

**Critical performance constraints:**
- Each statusline invocation **spawns a new process**
- Claude Code has a **300ms debounce** before running the command
- In-flight processes are **cancelled** if new input arrives
- The binary must complete well within the 300ms window

### 2.2 Input Data (from Claude Code)

The JSON input includes:

| Field | Type | Description |
|-------|------|-------------|
| `model.id` | string | Model identifier (e.g., `claude-opus-4-6`) |
| `model.display_name` | string | Human-readable model name |
| `workspace.project_dir` | string | Root project directory |
| `cwd` | string | Current working directory |
| `context_window.used_percentage` | number | Context usage percentage |
| `context_window.context_window_size` | number | Total context window size |
| `context_window.current_usage.*` | object | Token breakdown (input, output, cache) |
| `context_window.total_input_tokens` | number | Total input tokens for session |
| `context_window.total_output_tokens` | number | Total output tokens for session |
| `cost.total_cost_usd` | number | Total session cost in USD |
| `cost.total_duration_ms` | number | Total session duration in ms |
| `cost.total_lines_added` | number | Lines added in session |
| `cost.total_lines_removed` | number | Lines removed in session |
| `vim.mode` | string | Current vim mode (normal/insert) |
| `session_id` | string | Session identifier |
| `version` | string | Claude Code version |
| `exceeds_200k_tokens` | bool | Whether context exceeds 200k |

### 2.3 Tech Stack

| Technology | Purpose |
|-----------|---------|
| **Go 1.22+** | Fast, single-binary compilation |
| **gopkg.in/yaml.v3** | YAML config parsing (only dependency) |
| **goreleaser** | Cross-platform builds and releases |
| **golangci-lint** | Code quality |
| **GitHub Actions** | CI/CD |

---

## 3. Project Structure

```
claude-code-status-line/
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   ├── workflows/
│   │   ├── ci.yml              # Lint + test + build
│   │   └── release.yml         # goreleaser on tag push
│   └── PULL_REQUEST_TEMPLATE.md
├── cmd/
│   └── claude-code-status-line/
│       └── main.go             # CLI entry point
├── internal/
│   ├── ansi/
│   │   └── ansi.go             # ANSI color helpers
│   ├── config/
│   │   └── config.go           # YAML config loading
│   ├── engine/
│   │   └── engine.go           # Main rendering engine
│   ├── i18n/
│   │   └── i18n.go             # Translations (en, fr)
│   ├── parser/
│   │   └── parser.go           # JSON input parser
│   ├── themes/
│   │   └── themes.go           # Color themes
│   └── widgets/
│       ├── widgets.go          # Widget registry
│       ├── model.go            # Model name
│       ├── directory.go        # Working directory
│       ├── git.go              # Git branch, status, nested repos
│       ├── cost.go             # Session cost (multi-currency)
│       ├── tokens.go           # Token bar, count, %, cache, totals
│       ├── duration.go         # Session duration
│       └── misc.go             # Timestamp, OS, separator, spacer, vim, lines
├── docs/
│   ├── IMPLEMENTATION_PLAN.md  # This file
│   └── CODE_CONVENTIONS.md     # Coding standards
├── .golangci.yml
├── .goreleaser.yml
├── .gitignore
├── go.mod
├── go.sum
├── install.sh                  # Shell install script
├── install.ps1                 # PowerShell install script
├── Makefile
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── LICENSE                     # MIT
├── README.md
└── SECURITY.md
```

---

## 4. Widget System

### 4.1 Available Widgets

| Widget ID | Description | Example Output |
|-----------|-------------|----------------|
| `model` | Current Claude model | `Opus 4.6` |
| `directory` | Working directory name | `my-project` |
| `git-branch` | Current git branch | `feature/auth` |
| `git-status` | Dirty/clean indicator | `✓` or `✗` |
| `nested-repos` | Count of nested git repos | `3 repos` |
| `cost` | Session cost | `$1.23` / `1,23€` |
| `token-bar` | Progress bar for context | `━━━━━━━━────` |
| `token-count` | Tokens used/max | `45.2k/200k` |
| `duration` | Session duration | `12m34s` |
| `context-percent` | Context usage % | `23%` |
| `timestamp` | Current time | `14:32` |
| `cache-ratio` | Cache hit ratio | `Cache: 87%` |
| `total-tokens` | Total I/O tokens | `↑120k ↓45k` |
| `os-info` | Platform info | `macOS arm64` |
| `vim-mode` | Vim mode indicator | `NORMAL` / `INSERT` |
| `lines-changed` | Lines added/removed | `+42 -7` |
| `separator` | Visual separator | `│` |
| `spacer` | Flexible space | ` ` |

### 4.2 Widget Interface

```go
type Widget interface {
    ID() string
    Render(ctx *Context) string
}

type Context struct {
    Input  *parser.Input
    Config *config.Config
    Theme  themes.Theme
}
```

---

## 5. Distribution & Installation

### 5.1 Homebrew Tap (Primary for macOS/Linux)

```bash
brew tap EvanPluchart/tap
brew install claude-code-status-line
```

### 5.2 Go Install

```bash
go install github.com/EvanPluchart/claude-code-status-line/cmd/claude-code-status-line@latest
```

### 5.3 Shell Script

```bash
curl -sSL https://raw.githubusercontent.com/EvanPluchart/claude-code-status-line/main/install.sh | sh
```

### 5.4 GitHub Releases

Pre-built binaries for all platforms available on every release.

---

## 6. CI/CD & Release Process

### 6.1 GitHub Actions Workflows

#### ci.yml (on PR and push to main/develop)
1. Checkout code
2. Setup Go 1.22
3. Run golangci-lint
4. Run tests on ubuntu, macos, windows
5. Build binary

#### release.yml (on tag push v*)
1. Checkout code
2. Setup Go 1.22
3. Run goreleaser (builds all platforms, creates GitHub Release, updates Homebrew tap)

### 6.2 Versioning

- Follow [Semantic Versioning](https://semver.org/)
- Tags format: `v1.0.0`, `v1.1.0`, etc.
- Changelog auto-generated by goreleaser

---

## 7. Performance Targets

| Metric | Target |
|--------|--------|
| Binary execution | < 10ms |
| Binary size | < 5MB |
| Memory usage | < 5MB |

---

*Version: 1.0 - 2026-02-27*
