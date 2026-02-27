package parser

import "encoding/json"

// CurrentUsage represents the current context window usage.
type CurrentUsage struct {
	InputTokens                int `json:"input_tokens"`
	OutputTokens               int `json:"output_tokens"`
	CacheCreationInputTokens   int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens       int `json:"cache_read_input_tokens"`
}

// ContextWindow represents context window metrics.
type ContextWindow struct {
	TotalInputTokens   int           `json:"total_input_tokens"`
	TotalOutputTokens  int           `json:"total_output_tokens"`
	ContextWindowSize  int           `json:"context_window_size"`
	UsedPercentage     float64       `json:"used_percentage"`
	RemainingPct       float64       `json:"remaining_percentage"`
	CurrentUsage       *CurrentUsage `json:"current_usage"`
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

// Workspace represents workspace paths.
type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

// Vim represents vim mode state.
type Vim struct {
	Mode string `json:"mode"`
}

// Agent represents the active agent.
type Agent struct {
	Name string `json:"name"`
}

// Input is the full JSON payload from Claude Code.
type Input struct {
	CWD              string        `json:"cwd"`
	SessionID        string        `json:"session_id"`
	TranscriptPath   string        `json:"transcript_path"`
	Version          string        `json:"version"`
	Model            Model         `json:"model"`
	Workspace        Workspace     `json:"workspace"`
	Cost             Cost          `json:"cost"`
	ContextWindow    ContextWindow `json:"context_window"`
	Exceeds200K      bool          `json:"exceeds_200k_tokens"`
	Vim              *Vim          `json:"vim,omitempty"`
	Agent            *Agent        `json:"agent,omitempty"`
}

// Parse reads raw JSON bytes and returns a parsed Input.
func Parse(data []byte) (*Input, error) {
	var input Input

	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	// Defaults
	if input.ContextWindow.ContextWindowSize == 0 {
		input.ContextWindow.ContextWindowSize = 200000
	}

	if input.Model.DisplayName == "" {
		input.Model.DisplayName = "Claude"
	}

	if input.Workspace.ProjectDir == "" {
		input.Workspace.ProjectDir = input.CWD
	}

	if input.ContextWindow.CurrentUsage == nil {
		input.ContextWindow.CurrentUsage = &CurrentUsage{}
	}

	return &input, nil
}
