package widgets

import (
	"sort"

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
	// Session
	register(&ModelWidget{})
	register(&EffortWidget{})
	register(&ThinkingWidget{})
	register(&FastModeWidget{})
	register(&AgentWidget{})
	register(&OutputStyleWidget{})
	register(&SessionNameWidget{})
	register(&ClaudeVersionWidget{})
	register(&VimModeWidget{})

	// Workspace & git
	register(&DirectoryWidget{})
	register(&RepoWidget{})
	register(&WorktreeWidget{})
	register(&GitBranchWidget{})
	register(&GitStatusWidget{})
	register(&GitChangesWidget{})
	register(&GitAheadBehindWidget{})
	register(&PRWidget{})
	register(&NestedReposWidget{})
	register(&LinesChangedWidget{})

	// Cost & time
	register(&CostWidget{})
	register(&BurnRateWidget{})
	register(&DurationWidget{})
	register(&APITimeWidget{})
	register(&TimestampWidget{})

	// Context & tokens
	register(&TokenBarWidget{})
	register(&ContextPercentWidget{})
	register(&TokenCountWidget{})
	register(&ContextRemainingWidget{})
	register(&TotalTokensWidget{})
	register(&CacheRatioWidget{})
	register(&Exceeds200KWidget{})

	// Rate limits
	register(&SessionUsageWidget{})
	register(&WeeklyUsageWidget{})
	register(&ModelUsageWidget{})

	// Layout & system
	register(&OSInfoWidget{})
	register(&HostnameWidget{})
	register(&SeparatorWidget{})
	register(&SpacerWidget{})
}

// Get returns a widget by ID, or nil if not found.
func Get(id string) Widget {
	return registry[id]
}

// IDs returns all registered widget IDs, sorted.
func IDs() []string {
	ids := make([]string, 0, len(registry))

	for id := range registry {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	return ids
}
