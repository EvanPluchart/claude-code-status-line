package widgets

import "testing"

func TestModelName(t *testing.T) {
	cases := []struct {
		id    string
		name  string
		short string
	}{
		{"claude-opus-5", "Opus 5", "O5"},
		{"claude-fable-5", "Fable 5", "F5"},
		{"claude-sonnet-5", "Sonnet 5", "S5"},
		{"claude-opus-4-8", "Opus 4.8", "O4.8"},
		{"claude-opus-4-6", "Opus 4.6", "O4.6"},
		{"claude-sonnet-4-6", "Sonnet 4.6", "S4.6"},
		{"claude-sonnet-4-5-20250929", "Sonnet 4.5", "S4.5"},
		{"claude-haiku-4-5-20251001", "Haiku 4.5", "H4.5"},
		{"claude-opus-4-20250514", "Opus 4", "O4"},
		{"claude-3-5-sonnet-20241022", "Sonnet 3.5", "S3.5"},
		{"claude-3-haiku-20240307", "Haiku 3", "H3"},
		{"us.anthropic.claude-opus-5-v1:0", "Opus 5", "O5"},
		{"claude-opus-4-5@20251101", "Opus 4.5", "O4.5"},
		{"claude-sonnet-4-6[1m]", "Sonnet 4.6 1M", "S4.6"},
		{"Claude-Opus-5", "Opus 5", "O5"},
		{"gpt-4o", "", ""},
		{"", "", ""},
	}

	for _, c := range cases {
		if got := ModelName(c.id); got != c.name {
			t.Errorf("ModelName(%q) = %q, want %q", c.id, got, c.name)
		}

		if got := ModelShortName(c.id); got != c.short {
			t.Errorf("ModelShortName(%q) = %q, want %q", c.id, got, c.short)
		}
	}
}

func TestModelWidgetFallsBackToDisplayName(t *testing.T) {
	ctx := testContext()
	ctx.Input.Model.ID = "custom-model"
	ctx.Input.Model.DisplayName = "My Model"

	got := stripANSI((&ModelWidget{}).Render(ctx))

	if got != "My Model" {
		t.Errorf("got %q, want display name fallback", got)
	}
}

func TestModelWidgetShortName(t *testing.T) {
	ctx := testContext()
	ctx.Config.Widgets.Model.ShortName = true

	if got := stripANSI((&ModelWidget{}).Render(ctx)); got != "O5" {
		t.Errorf("got %q, want O5", got)
	}
}
