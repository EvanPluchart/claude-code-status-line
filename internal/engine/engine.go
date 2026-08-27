package engine

import (
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/themes"
	"github.com/EvanPluchart/claude-code-status-line/internal/widgets"
)

// Placeholder returns a substitute rendering for a widget ID (used by previews).
type Placeholder func(theme themes.Theme) string

// Render produces the full statusline output from input and config.
func Render(input *parser.Input, cfg *config.Config) string {
	return RenderWith(input, cfg, nil)
}

// RenderWith renders the statusline, substituting widgets that have a placeholder.
func RenderWith(input *parser.Input, cfg *config.Config, placeholders map[string]Placeholder) string {
	theme := themes.Get(cfg.Theme)
	ctx := &widgets.Context{
		Input:  input,
		Config: cfg,
		Theme:  theme,
	}

	var lines []string

	for i, line := range cfg.Lines {
		if i >= config.MaxLines {
			break
		}

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

			rendered := ""

			if ph, ok := placeholders[widgetID]; ok {
				rendered = ph(theme)
			} else {
				rendered = w.Render(ctx)
			}

			if rendered == "" {
				continue
			}

			entries = append(entries, entry{
				output:      rendered,
				isSeparator: widgetID == "separator" || widgetID == "spacer",
			})
		}

		// Strip leading, trailing, and consecutive separators.
		parts := cleanSeparators(entries)

		if len(parts) > 0 {
			lines = append(lines, joinParts(parts))
		}
	}

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n") + "\n"
}

// cleanSeparators removes leading, trailing, and consecutive separator entries.
func cleanSeparators(entries []entry) []entry {
	var out []entry
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

		out = append(out, e)
	}

	// Strip trailing separator
	if len(out) > 0 && lastWasSep {
		out = out[:len(out)-1]
	}

	return out
}

// joinParts concatenates entries, inserting a single space between two
// adjacent widgets that are not separated by a separator/spacer.
func joinParts(entries []entry) string {
	var b strings.Builder

	for i, e := range entries {
		if i > 0 && !e.isSeparator && !entries[i-1].isSeparator {
			b.WriteByte(' ')
		}

		b.WriteString(e.output)
	}

	return b.String()
}

type entry struct {
	output      string
	isSeparator bool
}
