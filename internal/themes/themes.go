package themes

import "github.com/EvanPluchart/claude-code-status-line/internal/ansi"

// Theme defines a color theme for the statusline.
type Theme struct {
	Name      string
	Primary   string
	Secondary string
	Muted     string
	Success   string
	Warning   string
	Danger    string
	Info      string
	Text      string
	Separator string
	BarFilled string
	BarEmpty  string
}

var registry = map[string]Theme{
	"default": {
		Name: "default", Primary: ansi.Cyan, Secondary: ansi.Blue,
		Muted: ansi.Gray, Success: ansi.Green, Warning: ansi.Yellow,
		Danger: ansi.Red, Info: ansi.Cyan, Text: ansi.White,
		Separator: ansi.Gray, BarFilled: ansi.Green, BarEmpty: ansi.Gray,
	},
	"minimal": {
		Name: "minimal", Primary: ansi.White, Secondary: ansi.Gray,
		Muted: ansi.Gray, Success: ansi.White, Warning: ansi.White,
		Danger: ansi.White, Info: ansi.Gray, Text: ansi.White,
		Separator: ansi.Gray, BarFilled: ansi.White, BarEmpty: ansi.Gray,
	},
	"neon": {
		Name: "neon", Primary: ansi.BrightMagenta, Secondary: ansi.BrightCyan,
		Muted: ansi.Gray, Success: ansi.BrightGreen, Warning: ansi.BrightYellow,
		Danger: ansi.BrightRed, Info: ansi.BrightCyan, Text: ansi.BrightWhite,
		Separator: ansi.BrightMagenta, BarFilled: ansi.BrightCyan, BarEmpty: ansi.Gray,
	},
	"dracula": {
		Name: "dracula", Primary: ansi.Magenta, Secondary: ansi.BrightMagenta,
		Muted: ansi.Gray, Success: ansi.Green, Warning: ansi.Yellow,
		Danger: ansi.Red, Info: ansi.Cyan, Text: ansi.White,
		Separator: ansi.Gray, BarFilled: ansi.Magenta, BarEmpty: ansi.Gray,
	},
	"catppuccin": {
		Name: "catppuccin", Primary: ansi.Magenta, Secondary: ansi.Yellow,
		Muted: ansi.Gray, Success: ansi.Green, Warning: ansi.Yellow,
		Danger: ansi.Red, Info: ansi.Blue, Text: ansi.White,
		Separator: ansi.Gray, BarFilled: ansi.Magenta, BarEmpty: ansi.Gray,
	},
	"nord": {
		Name: "nord", Primary: ansi.Cyan, Secondary: ansi.Blue,
		Muted: ansi.Gray, Success: ansi.Green, Warning: ansi.Yellow,
		Danger: ansi.Red, Info: ansi.Cyan, Text: ansi.White,
		Separator: ansi.Gray, BarFilled: ansi.Cyan, BarEmpty: ansi.Gray,
	},
}

// Get returns a theme by name. Falls back to default.
func Get(name string) Theme {
	if t, ok := registry[name]; ok {
		return t
	}

	return registry["default"]
}

// Names returns all available theme names.
func Names() []string {
	names := make([]string, 0, len(registry))

	for name := range registry {
		names = append(names, name)
	}

	return names
}
