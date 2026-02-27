# claude-code-status-line

A fully customizable, multi-line statusline for [Claude Code](https://docs.anthropic.com/en/docs/claude-code).

```
 Opus 4.6 │ my-project │ feature/auth │ 12m34s │ $1.23
 ━━━━━━━━━━━━────────── 58% (116k/200k)
```

## Features

- **3 configurable lines** with any combination of widgets
- **18 built-in widgets**: model, directory, git branch, cost, tokens, progress bar, duration, vim mode, lines changed, and more
- **6 themes**: default, minimal, neon, dracula, catppuccin, nord
- **Multi-currency**: USD, EUR, GBP, JPY, CAD, and more
- **i18n**: English and French out of the box
- **Cross-platform**: macOS (Apple Silicon + Intel), Linux, Windows 10/11, WSL2
- **Blazing fast**: single Go binary, ~5ms execution (well within Claude Code's 300ms budget)
- **Zero-config start**: works with sensible defaults

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap EvanPluchart/tap
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
# Initialize config and register in Claude Code
claude-code-status-line init

# With options
claude-code-status-line init --locale fr --theme dracula --currency EUR
```

That's it! Restart Claude Code and your new statusline will appear.

## Available Widgets

| Widget | Description | Example |
|--------|-------------|---------|
| `model` | Current Claude model | `Opus 4.6` |
| `directory` | Working directory | `my-project` |
| `git-branch` | Current git branch | `feature/auth` |
| `git-status` | Dirty/clean indicator | `✓` / `✗` |
| `nested-repos` | Nested git repos count | `3 repos` |
| `cost` | Session cost (multi-currency) | `$1.23` / `1,23€` |
| `token-bar` | Context usage progress bar | `━━━━━━━━────` |
| `token-count` | Tokens used/max | `45.2k/200k` |
| `duration` | Session duration | `12m34s` |
| `context-percent` | Context usage % | `58%` |
| `timestamp` | Current time | `14:32` |
| `cache-ratio` | Cache hit ratio | `Cache: 87%` |
| `total-tokens` | Total I/O tokens | `↑120k ↓45k` |
| `os-info` | Platform info | `macOS arm64` |
| `vim-mode` | Vim mode indicator | `NORMAL` / `INSERT` |
| `lines-changed` | Lines added/removed | `+42 -7` |
| `separator` | Visual separator | `│` |
| `spacer` | Flexible space | |

## Configuration

Configuration file: `~/.claude-statusline/config.yml`

```yaml
# Language: en | fr
locale: "en"

# Theme: default | minimal | neon | dracula | catppuccin | nord
theme: "default"

# Lines configuration
lines:
  - widgets: [model, separator, directory, separator, git-branch, separator, duration, separator, cost]
  - widgets: [token-bar, context-percent, token-count]
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

# Color thresholds
thresholds:
  context:
    green: 0
    yellow: 50
    orange: 70
    red: 90
  cost:
    green: 0
    yellow: 0.25
    orange: 1.0
    red: 5.0
  duration:
    green: 0
    yellow: 60
    red: 1800
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

## Uninstall

```bash
claude-code-status-line uninstall
```

This restores your previous statusline configuration (if any) and removes the config directory.

## Requirements

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) CLI

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
