package usage

import (
	"encoding/json"
	"testing"
	"time"
)

// TestParseModelsKeepsScopedWindows covers the filtering the widget depends on: the endpoint
// returns every window, only the per-model weekly ones must survive. The payload below is the
// real shape returned by the endpoint, including is_active:false on windows that are shown.
func TestParseModelsKeepsScopedWindows(t *testing.T) {
	raw := []byte(`{"limits": [
		{"kind": "session", "percent": 16, "is_active": true,
		 "resets_at": "2026-09-16T13:10:00.744050+00:00", "scope": null},
		{"kind": "weekly_all", "percent": 14, "is_active": false,
		 "resets_at": "2026-09-21T22:00:00.744068+00:00", "scope": null},
		{"kind": "weekly_scoped", "percent": 6, "is_active": false,
		 "resets_at": "2026-09-22T00:00:00Z", "scope": {"model": {"id": null, "display_name": "Fable"}}},
		{"kind": "weekly_scoped", "percent": 50, "scope": {"model": null}}
	]}`)

	var payload struct {
		Limits []apiLimit `json:"limits"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	models := parseModels(payload.Limits)

	if len(models) != 1 {
		t.Fatalf("models = %+v, want only the Fable window", models)
	}

	if models[0].DisplayName != "Fable" || models[0].Percent != 6 {
		t.Errorf("model = %+v, want Fable at 6%%", models[0])
	}

	if want := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC).Unix(); models[0].ResetsAt != want {
		t.Errorf("resets_at = %d, want %d", models[0].ResetsAt, want)
	}
}

func TestParseResetsAtIsDefensive(t *testing.T) {
	for _, value := range []string{"", "not-a-date", "1789564200"} {
		if got := parseResetsAt(value); got != 0 {
			t.Errorf("parseResetsAt(%q) = %d, want 0", value, got)
		}
	}
}

func TestUserAgent(t *testing.T) {
	if got := UserAgent("2.1.273"); got != "claude-code/2.1.273" {
		t.Errorf("UserAgent = %q", got)
	}

	if got := UserAgent(""); got != "claude-code" {
		t.Errorf("UserAgent fallback = %q", got)
	}
}
