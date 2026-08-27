# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- 17 new widgets: `effort`, `thinking`, `fast-mode`, `agent`, `output-style`, `session-name`,
  `claude-version`, `repo`, `worktree`, `git-changes`, `git-ahead-behind`, `pr`, `burn-rate`,
  `api-time`, `context-remaining`, `exceeds-200k`, `hostname`
- Parsing of the new Claude Code payload fields (`effort`, `thinking`, `fast_mode`, `pr`,
  `worktree`, `workspace.repo`, `output_style`, `session_name`, `prompt_id`)
- `preview` command (sample data or `--stdin` JSON payload) and `widgets` command
- `widgets.lines_changed.source` (`git` | `session`), `widgets.directory.depth`,
  `widgets.git.max_branch_length` and `thresholds.rate_limit` options
- `CLAUDE_STATUSLINE_DIR` / `CLAUDE_CONFIG_DIR` environment overrides
- Unit tests for parser, config, engine, widgets, i18n and themes

### Changed

- Model names are derived from the model ID (Claude 5 family, Bedrock/Vertex IDs, `[1m]` suffix)
  instead of a static table
- Adjacent widgets are now separated by a space; `spacer` behaves like a separator
- `thresholds.context`, `thresholds.cost` and `thresholds.duration` are now honored by every widget
- Configuration merge now keeps every default for keys absent from the file (thresholds included)
- `init` registers `padding: 0` and `hideVimModeIndicator` when the `vim-mode` widget is used
- Default context window fallback, fallback exchange rates and French translations (accents) refreshed
- golangci-lint configuration migrated to v2, CI pinned to the go.mod Go version

### Fixed

- Exchange rates were never refreshed in the background (goroutine killed at exit); a detached
  `update-rates --quiet` process now handles it, at most once every 30 minutes
- `nested-repos` and `lines-changed` relied on `find` and did not work on Windows
- `used_percentage: null` early in a session no longer renders 0%: it is derived from the usage
- Detached HEAD is displayed as `@<sha>` instead of `HEAD`
