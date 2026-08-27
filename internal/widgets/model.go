package widgets

import (
	"regexp"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
)

// modelIDPattern matches Claude model IDs such as:
//
//	claude-opus-5, claude-sonnet-4-6, claude-haiku-4-5-20251001,
//	claude-3-5-sonnet-20241022, us.anthropic.claude-opus-5-v1:0, claude-opus-4-5@20251101
var modelIDPattern = regexp.MustCompile(`claude-(?:(\d+)-(?:(\d+)-)?)?(opus|sonnet|haiku|fable|mythos)(?:-(\d+)(?:-(\d+))?)?`)

var familyShortPrefix = map[string]string{
	"opus":   "O",
	"sonnet": "S",
	"haiku":  "H",
	"fable":  "F",
	"mythos": "M",
}

// ModelName derives a human-readable name ("Opus 4.6") from a model ID.
// Returns an empty string when the ID is not recognized.
func ModelName(id string) string {
	family, version := parseModelID(id)

	if family == "" {
		return ""
	}

	name := strings.ToUpper(family[:1]) + family[1:]

	if version != "" {
		name += " " + version
	}

	if strings.Contains(id, "[1m]") {
		name += " 1M"
	}

	return name
}

// ModelShortName derives a compact name ("O4.6") from a model ID.
// Returns an empty string when the ID is not recognized.
func ModelShortName(id string) string {
	family, version := parseModelID(id)

	if family == "" {
		return ""
	}

	return familyShortPrefix[family] + version
}

// parseModelID extracts the family and version ("4.6") from a model ID.
func parseModelID(id string) (string, string) {
	m := modelIDPattern.FindStringSubmatch(strings.ToLower(id))

	if m == nil {
		return "", ""
	}

	family := m[3]
	major := m[1]
	minor := m[2]

	// Modern IDs carry the version after the family name (claude-opus-4-6).
	if major == "" {
		major = m[4]
		minor = m[5]
	}

	// A trailing 8-digit group is a date snapshot, not a minor version.
	if len(minor) == 8 {
		minor = ""
	}

	if len(major) == 8 {
		major = ""
	}

	version := major

	if minor != "" {
		version += "." + minor
	}

	return family, version
}

// ModelWidget displays the current Claude model name.
type ModelWidget struct{}

func (w *ModelWidget) ID() string { return "model" }

func (w *ModelWidget) Render(ctx *Context) string {
	var name string

	if ctx.Config.Widgets.Model.ShortName {
		name = ModelShortName(ctx.Input.Model.ID)
	} else {
		name = ModelName(ctx.Input.Model.ID)
	}

	if name == "" {
		name = ctx.Input.Model.DisplayName
	}

	if name == "" {
		name = "Claude"
	}

	return ansi.ColorBold(name, contextColor(ctx.Input.ContextWindow.UsedPercentage, ctx))
}
