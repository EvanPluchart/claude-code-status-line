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

		// Collect widget outputs, keeping track of which are separators.
		var entries []entry

		for _, widgetID := range line.Widgets {
			w := widgets.Get(widgetID)
			if w == nil {
				continue
			}

			rendered := w.Render(ctx)
			if rendered == "" {
				continue
			}

			entries = append(entries, entry{
				output:      rendered,
				isSeparator: widgetID == "separator",
			})
		}

		// Strip leading, trailing, and consecutive separators.
		parts := cleanSeparators(entries)

		if len(parts) > 0 {
			lines = append(lines, strings.Join(parts, ""))
		}
	}

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n") + "\n"
}

// cleanSeparators removes leading, trailing, and consecutive separator entries.
func cleanSeparators(entries []entry) []string {
	var out []string
	lastWasSep := true // treat start as separator to strip leading ones

	for _, e := range entries {
		if e.isSeparator {
			if lastWasSep {
				continue // skip consecutive or leading separator
			}
			lastWasSep = true
		} else {
			lastWasSep = false
		}

		out = append(out, e.output)
	}

	// Strip trailing separator
	if len(out) > 0 && lastWasSep {
		out = out[:len(out)-1]
	}

	return out
}

type entry struct {
	output      string
	isSeparator bool
}
