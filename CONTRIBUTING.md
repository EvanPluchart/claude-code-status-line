# Contributing to claude-code-status-line

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- Go 1.22+
- golangci-lint (optional, for linting)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/EvanPluchart/claude-code-status-line.git
cd claude-code-status-line

# Build
make build

# Run tests
make test

# Run linter
make lint

# Test locally
echo '{"model":{"id":"claude-opus-4-6","display_name":"Claude Opus 4.6"}}' | ./claude-code-status-line
```

### Project Structure

```
cmd/
└── claude-code-status-line/
    └── main.go           # CLI entry point
internal/
├── ansi/                 # ANSI color helpers
├── config/               # YAML config loading
├── engine/               # Main rendering engine
├── i18n/                 # Translations (en, fr)
├── parser/               # JSON input parser
├── themes/               # Color themes
└── widgets/              # Widget implementations
```

## How to Contribute

### Reporting Bugs

1. Check [existing issues](https://github.com/EvanPluchart/claude-code-status-line/issues)
2. Create a new issue using the **Bug Report** template
3. Include your OS, Go version, and Claude Code version

### Suggesting Features

1. Check [existing issues](https://github.com/EvanPluchart/claude-code-status-line/issues)
2. Create a new issue using the **Feature Request** template
3. Describe the widget/theme/feature and its use case

### Submitting Code

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Write your code following the [Code Conventions](docs/CODE_CONVENTIONS.md)
4. Add tests for new functionality
5. Run the full test suite: `make test`
6. Run the linter: `make lint`
7. Commit with a clear message: `Feature - Add my new feature`
8. Push and open a Pull Request

### Adding a Widget

1. Create a new file in `internal/widgets/` (e.g., `my-widget.go`)
2. Implement the `Widget` interface (`ID() string` and `Render(*Context) string`)
3. Register it in `internal/widgets/widgets.go` `init()` function
4. Add tests

### Adding a Theme

1. Add a new theme struct to `internal/themes/themes.go`
2. Register it in the `themes` map

### Adding a Translation

1. Add a new locale map in `internal/i18n/i18n.go`
2. Register it in the `locales` map

## Code Conventions

See [docs/CODE_CONVENTIONS.md](docs/CODE_CONVENTIONS.md) for the full coding standards.

### Quick Summary

- `gofmt` and `goimports` for formatting
- Always use braces for conditions
- Boolean naming: `is`/`has`/`should`/`can` prefix
- Imports grouped: stdlib > external > internal
- Commit format: `Type - Description`

## Pull Request Process

1. Fill in the PR template
2. Ensure CI passes (lint, tests, build)
3. Request a review
4. Address review feedback
5. Squash and merge

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
