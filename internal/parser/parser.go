package parser

import (
	"encoding/json"
	"strings"
)

// DefaultContextWindowSize is used when Claude Code does not provide a size.
const DefaultContextWindowSize = 200000

// CurrentUsage represents the current context window usage.
type CurrentUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// ContextWindow represents context window metrics.
type ContextWindow struct {
	TotalInputTokens  int           `json:"total_input_tokens"`
	TotalOutputTokens int           `json:"total_output_tokens"`
	ContextWindowSize int           `json:"context_window_size"`
	UsedPercentage    float64       `json:"used_percentage"`
	RemainingPct      float64       `json:"remaining_percentage"`
	CurrentUsage      *CurrentUsage `json:"current_usage"`
}

// UsedTokens returns the number of tokens currently in the context window.
// It prefers the exact usage breakdown and falls back to the percentage.
func (cw ContextWindow) UsedTokens() int {
	if cw.CurrentUsage != nil {
		used := cw.CurrentUsage.InputTokens + cw.CurrentUsage.CacheCreationInputTokens + cw.CurrentUsage.CacheReadInputTokens

		if used > 0 {
			return used
		}
	}

	return int(float64(cw.ContextWindowSize) * cw.UsedPercentage / 100)
}

// Cost represents session cost information.
type Cost struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMS    int64   `json:"total_duration_ms"`
	TotalAPIDurationMS int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}

// Model represents the Claude model information.
type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// Repo represents the remote repository of the workspace.
type Repo struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// Workspace represents workspace paths.
type Workspace struct {
	CurrentDir  string   `json:"current_dir"`
	ProjectDir  string   `json:"project_dir"`
	AddedDirs   []string `json:"added_dirs"`
	GitWorktree string   `json:"git_worktree"`
	Repo        *Repo    `json:"repo,omitempty"`
}

// OutputStyle represents the active output style.
type OutputStyle struct {
	Name string `json:"name"`
}

// Effort represents the reasoning effort level.
type Effort struct {
	Level string `json:"level"`
}

// Thinking represents the extended thinking state.
type Thinking struct {
	Enabled bool `json:"enabled"`
}

// Vim represents vim mode state.
type Vim struct {
	Mode string `json:"mode"`
}

// Agent represents the active agent.
type Agent struct {
	Name string `json:"name"`
}

// PullRequest represents the open PR/MR for the current branch.
type PullRequest struct {
	Number      int    `json:"number"`
	URL         string `json:"url"`
	ReviewState string `json:"review_state"`
	Kind        string `json:"kind"`
}

// Worktree represents an active worktree session.
type Worktree struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	Branch         string `json:"branch"`
	OriginalCWD    string `json:"original_cwd"`
	OriginalBranch string `json:"original_branch"`
}

// RateLimit represents a single rate limit window.
type RateLimit struct {
	UsedPercentage float64 `json:"used_percentage"`
	ResetsAt       int64   `json:"resets_at"`
}

// RateLimits holds the session and weekly rate limits.
type RateLimits struct {
	FiveHour *RateLimit `json:"five_hour,omitempty"`
	SevenDay *RateLimit `json:"seven_day,omitempty"`
}

// Input is the full JSON payload from Claude Code.
type Input struct {
	CWD            string        `json:"cwd"`
	SessionID      string        `json:"session_id"`
	SessionName    string        `json:"session_name"`
	PromptID       string        `json:"prompt_id"`
	TranscriptPath string        `json:"transcript_path"`
	Version        string        `json:"version"`
	Model          Model         `json:"model"`
	Workspace      Workspace     `json:"workspace"`
	OutputStyle    *OutputStyle  `json:"output_style,omitempty"`
	Cost           Cost          `json:"cost"`
	ContextWindow  ContextWindow `json:"context_window"`
	Exceeds200K    bool          `json:"exceeds_200k_tokens"`
	FastMode       bool          `json:"fast_mode"`
	Effort         *Effort       `json:"effort,omitempty"`
	Thinking       *Thinking     `json:"thinking,omitempty"`
	RateLimits     *RateLimits   `json:"rate_limits,omitempty"`
	Vim            *Vim          `json:"vim,omitempty"`
	Agent          *Agent        `json:"agent,omitempty"`
	PR             *PullRequest  `json:"pr,omitempty"`
	Worktree       *Worktree     `json:"worktree,omitempty"`
}

// Parse reads raw JSON bytes and returns a parsed Input with defaults applied.
func Parse(data []byte) (*Input, error) {
	var input Input

	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	applyDefaults(&input)

	return &input, nil
}

func applyDefaults(input *Input) {
	if input.ContextWindow.ContextWindowSize <= 0 {
		input.ContextWindow.ContextWindowSize = DefaultContextWindowSize
	}

	if input.Model.DisplayName == "" {
		input.Model.DisplayName = "Claude"
	}

	if input.Workspace.ProjectDir == "" {
		input.Workspace.ProjectDir = input.CWD
	}

	if input.CWD == "" {
		input.CWD = input.Workspace.CurrentDir
	}

	if input.ContextWindow.CurrentUsage == nil {
		input.ContextWindow.CurrentUsage = &CurrentUsage{}
	}

	// used_percentage may be null early in the session: derive it from the usage breakdown.
	if input.ContextWindow.UsedPercentage == 0 && input.ContextWindow.ContextWindowSize > 0 {
		used := input.ContextWindow.UsedTokens()

		if used > 0 {
			input.ContextWindow.UsedPercentage = float64(used) / float64(input.ContextWindow.ContextWindowSize) * 100
		}
	}

	if input.ContextWindow.RemainingPct == 0 && input.ContextWindow.UsedPercentage > 0 {
		input.ContextWindow.RemainingPct = 100 - input.ContextWindow.UsedPercentage
	}

	if input.Vim != nil {
		input.Vim.Mode = strings.ToUpper(strings.TrimSpace(input.Vim.Mode))
	}
}
