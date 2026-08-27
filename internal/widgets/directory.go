package widgets

import (
	"path/filepath"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
)

// DirectoryWidget displays the working directory name.
type DirectoryWidget struct{}

func (w *DirectoryWidget) ID() string { return "directory" }

func (w *DirectoryWidget) Render(ctx *Context) string {
	dir := ctx.Input.Workspace.ProjectDir

	if dir == "" {
		dir = ctx.Input.CWD
	}

	return ansi.Colorize(shortenPath(dir, ctx.Config.Widgets.Directory.Depth), ctx.Theme.Secondary)
}

// shortenPath keeps the last `depth` segments of a path (basename when depth <= 1).
func shortenPath(dir string, depth int) string {
	dir = filepath.Clean(dir)
	segments := strings.Split(filepath.ToSlash(dir), "/")

	var kept []string

	for _, s := range segments {
		if s != "" && s != "." {
			kept = append(kept, s)
		}
	}

	if len(kept) == 0 {
		return "~"
	}

	if depth <= 1 {
		return kept[len(kept)-1]
	}

	if depth > len(kept) {
		depth = len(kept)
	}

	return strings.Join(kept[len(kept)-depth:], "/")
}
