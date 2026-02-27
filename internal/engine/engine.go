package engine

import (
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/themes"
	"github.com/EvanPluchart/claude-code-status-line/internal/widgets"
)

// Render produces the full statusline output from input and config.
func Render(input *parser.Input, cfg *config.Config) string {
	theme := themes.Get(cfg.Theme)
	ctx := &widgets.Context{
		Input:  input,
		Config: cfg,
		Theme:  theme,
	}

	var lines []string

	for _, line := range cfg.Lines {
		if len(line.Widgets) == 0 {
			continue
		}

		var parts []string

		for _, widgetID := range line.Widgets {
			w := widgets.Get(widgetID)
			if w == nil {
				continue
			}

			rendered := w.Render(ctx)
			if rendered != "" {
				parts = append(parts, rendered)
			}
		}

		if len(parts) > 0 {
			lines = append(lines, strings.Join(parts, ""))
		}
	}

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n") + "\n"
}
