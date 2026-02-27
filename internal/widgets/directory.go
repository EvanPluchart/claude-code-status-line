package widgets

import (
	"path/filepath"

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

	name := filepath.Base(dir)
	if name == "" || name == "." {
		name = "~"
	}

	return ansi.Colorize(name, ctx.Theme.Secondary)
}
