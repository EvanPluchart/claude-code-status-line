package widgets

import (
	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/themes"
)

// Context holds everything a widget needs to render.
type Context struct {
	Input  *parser.Input
	Config *config.Config
	Theme  themes.Theme
}

// Widget is the interface all widgets must implement.
type Widget interface {
	ID() string
	Render(ctx *Context) string
}

var registry = map[string]Widget{}

func register(w Widget) {
	registry[w.ID()] = w
}

func init() {
	register(&ModelWidget{})
	register(&DirectoryWidget{})
	register(&GitBranchWidget{})
	register(&GitStatusWidget{})
	register(&NestedReposWidget{})
	register(&CostWidget{})
	register(&TokenBarWidget{})
	register(&TokenCountWidget{})
	register(&DurationWidget{})
	register(&ContextPercentWidget{})
	register(&TimestampWidget{})
	register(&CacheRatioWidget{})
	register(&TotalTokensWidget{})
	register(&OSInfoWidget{})
	register(&SeparatorWidget{})
	register(&SpacerWidget{})
	register(&VimModeWidget{})
	register(&LinesChangedWidget{})
	register(&SessionUsageWidget{})
	register(&WeeklyUsageWidget{})
}

// Get returns a widget by ID, or nil if not found.
func Get(id string) Widget {
	return registry[id]
}

// IDs returns all registered widget IDs.
func IDs() []string {
	ids := make([]string, 0, len(registry))

	for id := range registry {
		ids = append(ids, id)
	}

	return ids
}
