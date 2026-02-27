# Code Conventions - Claude Code Status Line

> Coding standards for the claude-code-status-line project.

---

## General Rules

### 1. Always Use Braces for Conditions

Never use inline returns or single-line conditions without braces.

```go
// Good
if user == nil {
    return nil
}

// Bad (not applicable in Go, but keep the mindset)
```

### 2. Whitespace - Aerated Code

Add blank lines to improve readability:

- **Before and after** each condition (`if`, `switch`)
- **Before and after** each loop (`for`, `range`)
- **Before** each `return`, `break`, `continue`

```go
// Good
func processInput(data *parser.Input) string {
    if data == nil {
        return ""
    }

    model := data.Model.DisplayName

    if model == "" {
        return "Unknown"
    }

    for _, widget := range widgets {
        widget.Render(ctx)
    }

    return result
}
```

### 3. Early Return Pattern

Use early returns to avoid deep nesting.

```go
// Good
func validateConfig(cfg *Config) error {
    if len(cfg.Lines) == 0 {
        return fmt.Errorf("no lines configured")
    }

    if len(cfg.Lines) > 3 {
        return fmt.Errorf("maximum 3 lines")
    }

    return nil
}
```

### 4. Comments

- **DO** comment complex logic, business rules, or non-obvious behavior
- **DON'T** comment self-explanatory code
- Code should be self-documenting through good naming

```go
// Good - explains WHY
// Use --no-optional-locks to avoid blocking other git operations
cmd := exec.Command("git", "--no-optional-locks", "status", "--porcelain")

// Bad - states the obvious
// Get the model name
modelName := getModelName(input)
```

---

## Go Rules

### 1. Formatting

Use `gofmt` and `goimports`. No exceptions. The CI enforces this.

### 2. Import Organization

Separate imports into groups with blank lines between:

```go
import (
    // stdlib
    "fmt"
    "os"
    "path/filepath"

    // external
    "gopkg.in/yaml.v3"

    // internal
    "github.com/EvanPluchart/claude-code-status-line/internal/config"
    "github.com/EvanPluchart/claude-code-status-line/internal/parser"
)
```

### 3. Boolean Naming

Always use prefixes: `is`, `has`, `should`, `can`, `will`

```go
// Good
isLoading := true
hasError := false
shouldRefresh := true

// Bad
loading := true
error := false
```

### 4. Error Handling

- Check errors immediately after the call
- Use early returns for error paths
- Wrap errors with context when propagating

```go
// Good
data, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}
```

### 5. Constants

```go
// Good
const (
    MaxLines       = 3
    DefaultBarWidth = 16
)

// Bad
const maxLines = 3
```

### 6. Struct Tags

Use `yaml` and `json` tags consistently with snake_case:

```go
type Config struct {
    Locale string `yaml:"locale" json:"locale"`
    Theme  string `yaml:"theme"  json:"theme"`
}
```

---

## Naming Conventions

### Case Styles

| Style | Usage | Example |
|-------|-------|---------|
| `PascalCase` | Exported types, functions | `Widget`, `RenderContext` |
| `camelCase` | Unexported variables, functions | `renderWidget`, `isActive` |
| `SCREAMING_SNAKE_CASE` | Not used in Go | - |
| `kebab-case` | Widget IDs, file names | `git-branch`, `token-bar` |

### File Naming

| Type | Pattern | Example |
|------|---------|---------|
| Widget file | `purpose.go` | `git.go`, `tokens.go` |
| Test file | `*_test.go` | `model_test.go` |
| Package | lowercase, no underscores | `widgets`, `config` |

---

## Git Conventions

### Branches

```
type/description
```

| Type | Usage | Example |
|------|-------|---------|
| `feature/` | New functionality | `feature/add-git-branch-widget` |
| `bugfix/` | Bug fix | `bugfix/fix-token-bar-overflow` |
| `hotfix/` | Urgent fix | `hotfix/fix-crash-on-empty-input` |
| `release/` | Release preparation | `release/v1.0.0` |
| `chore/` | Maintenance | `chore/update-dependencies` |

### Commits

```
Type - Description
```

| Type | Usage | Example |
|------|-------|---------|
| `Feature` | New functionality | `Feature - Add git branch widget` |
| `Bugfix` | Bug fix | `Bugfix - Fix token bar overflow on small terminals` |
| `Hotfix` | Urgent fix | `Hotfix - Fix crash when stdin is empty` |
| `Refactor` | Refactoring | `Refactor - Extract ANSI helpers to package` |
| `Test` | Tests | `Test - Add unit tests for cost widget` |
| `Docs` | Documentation | `Docs - Add widget reference guide` |
| `Chore` | Maintenance | `Chore - Update Go to 1.23` |
| `Style` | Formatting | `Style - Fix import ordering` |

### Rules

1. **One action per commit** - don't mix multiple changes
2. **Clear description** - explain "what", not "how"
3. **Imperative mood** - "Add" not "Added" or "Adds"
4. **Max 72 characters** for the main line

### Protected Branches

| Branch | Protection |
|--------|-----------|
| `main` | Merge via PR only, CI must pass |
| `develop` | Merge via PR only |

---

## Code Quality

### Pre-commit Checklist

- [ ] `gofmt` passes
- [ ] `goimports` passes
- [ ] `golangci-lint` passes
- [ ] All conditions use braces `{}`
- [ ] Blank lines before/after conditions, loops
- [ ] Blank lines before `return`, `break`
- [ ] Boolean variables use `is`/`has`/`should` prefix
- [ ] Imports are grouped and organized
- [ ] No obvious/unnecessary comments
- [ ] Complex logic is commented
- [ ] Errors are handled
- [ ] Tests pass

---

*Version: 1.0 - 2026-02-27*
