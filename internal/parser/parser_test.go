package parser

import "testing"

func TestParseFullPayload(t *testing.T) {
	data := []byte(`{
		"cwd": "/tmp/project",
		"session_id": "abc",
		"session_name": "My session",
		"version": "2.1.90",
		"model": {"id": "claude-opus-5", "display_name": "Opus"},
		"workspace": {"current_dir": "/tmp/project", "project_dir": "/tmp/project", "repo": {"host": "github.com", "owner": "acme", "name": "app"}},
		"output_style": {"name": "Explanatory"},
		"cost": {"total_cost_usd": 1.5, "total_duration_ms": 60000, "total_api_duration_ms": 20000, "total_lines_added": 10, "total_lines_removed": 2},
		"context_window": {"context_window_size": 1000000, "used_percentage": 12.5, "remaining_percentage": 87.5,
			"current_usage": {"input_tokens": 100000, "output_tokens": 500, "cache_creation_input_tokens": 5000, "cache_read_input_tokens": 20000}},
		"exceeds_200k_tokens": false,
		"fast_mode": true,
		"effort": {"level": "xhigh"},
		"thinking": {"enabled": true},
		"rate_limits": {"five_hour": {"used_percentage": 23.5, "resets_at": 1738425600}},
		"vim": {"mode": "insert"},
		"agent": {"name": "reviewer"},
		"pr": {"number": 42, "url": "https://github.com/acme/app/pull/42", "review_state": "approved"},
		"worktree": {"name": "feat", "branch": "worktree-feat"}
	}`)

	input, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if input.Model.ID != "claude-opus-5" {
		t.Errorf("model id = %q", input.Model.ID)
	}

	if input.Workspace.Repo == nil || input.Workspace.Repo.Owner != "acme" {
		t.Errorf("repo not parsed: %+v", input.Workspace.Repo)
	}

	if input.Effort == nil || input.Effort.Level != "xhigh" {
		t.Errorf("effort not parsed: %+v", input.Effort)
	}

	if input.Thinking == nil || !input.Thinking.Enabled {
		t.Errorf("thinking not parsed: %+v", input.Thinking)
	}

	if !input.FastMode {
		t.Error("fast_mode should be true")
	}

	if input.PR == nil || input.PR.Number != 42 || input.PR.ReviewState != "approved" {
		t.Errorf("pr not parsed: %+v", input.PR)
	}

	if input.Worktree == nil || input.Worktree.Name != "feat" {
		t.Errorf("worktree not parsed: %+v", input.Worktree)
	}

	if input.Vim.Mode != "INSERT" {
		t.Errorf("vim mode should be normalized to upper case, got %q", input.Vim.Mode)
	}

	if input.RateLimits == nil || input.RateLimits.FiveHour == nil || input.RateLimits.SevenDay != nil {
		t.Errorf("rate limits not parsed as expected: %+v", input.RateLimits)
	}

	if got := input.ContextWindow.UsedTokens(); got != 125000 {
		t.Errorf("UsedTokens = %d, want 125000", got)
	}
}

func TestParseDefaults(t *testing.T) {
	input, err := Parse([]byte(`{"cwd": "/tmp/x"}`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if input.ContextWindow.ContextWindowSize != DefaultContextWindowSize {
		t.Errorf("default context size = %d", input.ContextWindow.ContextWindowSize)
	}

	if input.Model.DisplayName != "Claude" {
		t.Errorf("default display name = %q", input.Model.DisplayName)
	}

	if input.Workspace.ProjectDir != "/tmp/x" {
		t.Errorf("project dir should fall back to cwd, got %q", input.Workspace.ProjectDir)
	}

	if input.ContextWindow.CurrentUsage == nil {
		t.Error("current usage should never be nil after Parse")
	}

	if input.Effort != nil || input.PR != nil || input.Vim != nil {
		t.Error("optional sections should stay nil when absent")
	}
}

func TestParseNullFields(t *testing.T) {
	data := []byte(`{
		"cwd": "/tmp/x",
		"context_window": {"context_window_size": 200000, "used_percentage": null, "remaining_percentage": null,
			"current_usage": {"input_tokens": 50000, "cache_read_input_tokens": 50000}}
	}`)

	input, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if input.ContextWindow.UsedPercentage != 50 {
		t.Errorf("used percentage should be derived from usage, got %v", input.ContextWindow.UsedPercentage)
	}

	if input.ContextWindow.RemainingPct != 50 {
		t.Errorf("remaining percentage should be derived, got %v", input.ContextWindow.RemainingPct)
	}
}

func TestParseNullCurrentUsage(t *testing.T) {
	input, err := Parse([]byte(`{"context_window": {"context_window_size": 200000, "used_percentage": 30, "current_usage": null}}`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if input.ContextWindow.CurrentUsage == nil {
		t.Fatal("current usage should be defaulted")
	}

	if got := input.ContextWindow.UsedTokens(); got != 60000 {
		t.Errorf("UsedTokens fallback from percentage = %d, want 60000", got)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	if _, err := Parse([]byte(`{not json`)); err == nil {
		t.Error("expected an error for invalid JSON")
	}
}
