// Package usage manages the per-model weekly rate limit cache used by the model-usage widget.
//
// Claude Code does not send per-model windows in the statusline payload (only five_hour and
// seven_day), so they are fetched from the same endpoint its own /usage panel calls. The render
// path never performs network calls nor reads credentials: it only reads the on-disk cache.
// When the cache is stale, a detached background process refreshes it.
package usage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/detach"
)

const (
	apiURL = "https://api.anthropic.com/api/oauth/usage?cedar_ember=1&skip_spend=1"
	// oauthBeta is the beta header Claude Code sends on this endpoint.
	oauthBeta = "oauth-2025-04-20"
	// maxAge is deliberately above the 180s the endpoint tolerates before rate limiting.
	maxAge       = 5 * time.Minute
	retryDelay   = 5 * time.Minute
	fetchTimeout = 5 * time.Second
	// weeklyScopedKind marks the per-model weekly windows among all returned limits.
	weeklyScopedKind = "weekly_scoped"
)

// ModelUsage is a weekly rate limit window scoped to a single model (e.g. Fable).
type ModelUsage struct {
	DisplayName string  `json:"display_name"`
	Percent     float64 `json:"percent"`
	ResetsAt    int64   `json:"resets_at"`
}

// CachedUsage is the on-disk cache structure.
type CachedUsage struct {
	UpdatedAt time.Time    `json:"updated_at"`
	Models    []ModelUsage `json:"models"`
}

var (
	once       sync.Once
	cachedData *CachedUsage
)

func cachePath() string {
	return filepath.Join(config.ConfigDir(), "usage.json")
}

// attemptPath marks the last background refresh attempt to avoid spawn storms.
func attemptPath() string {
	return filepath.Join(config.ConfigDir(), "usage.attempt")
}

// UserAgent returns the User-Agent the usage endpoint expects, built from the running
// Claude Code version. Unknown clients get rate limited aggressively.
func UserAgent(claudeVersion string) string {
	if claudeVersion == "" {
		return "claude-code"
	}

	return "claude-code/" + claudeVersion
}

// Models returns the per-model weekly limits from the local cache.
// When the cache is stale, a background refresh is scheduled for the given Claude Code version.
func Models(claudeVersion string) []ModelUsage {
	once.Do(func() {
		cachedData = loadCache()

		if cachedData == nil || time.Since(cachedData.UpdatedAt) > maxAge {
			scheduleRefresh(claudeVersion)
		}
	})

	if cachedData == nil {
		return nil
	}

	return cachedData.Models
}

func loadCache() *CachedUsage {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil
	}

	var cached CachedUsage

	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}

	return &cached
}

// scheduleRefresh spawns a detached "update-usage" process, at most once per retryDelay.
func scheduleRefresh(claudeVersion string) {
	if os.Getenv("CLAUDE_STATUSLINE_OFFLINE") != "" {
		return
	}

	if info, err := os.Stat(attemptPath()); err == nil && time.Since(info.ModTime()) < retryDelay {
		return
	}

	if err := os.MkdirAll(config.ConfigDir(), 0o755); err != nil {
		return
	}

	if err := os.WriteFile(attemptPath(), []byte(time.Now().Format(time.RFC3339)), 0o644); err != nil {
		return
	}

	_ = detach.Spawn("update-usage", "--quiet", UserAgent(claudeVersion))
}

// apiLimit is one window as returned by the usage endpoint.
type apiLimit struct {
	Kind     string  `json:"kind"`
	Percent  float64 `json:"percent"`
	ResetsAt string  `json:"resets_at"`
	Scope    struct {
		Model *struct {
			DisplayName string `json:"display_name"`
		} `json:"model"`
	} `json:"scope"`
}

// parseModels keeps the per-model weekly windows.
//
// The is_active flag is deliberately ignored: it marks the currently binding constraint,
// not whether a window should be shown (weekly_all reports false while still being enforced).
func parseModels(limits []apiLimit) []ModelUsage {
	models := make([]ModelUsage, 0, len(limits))

	for _, limit := range limits {
		if limit.Kind != weeklyScopedKind {
			continue
		}

		if limit.Scope.Model == nil || limit.Scope.Model.DisplayName == "" {
			continue
		}

		models = append(models, ModelUsage{
			DisplayName: limit.Scope.Model.DisplayName,
			Percent:     limit.Percent,
			ResetsAt:    parseResetsAt(limit.ResetsAt),
		})
	}

	return models
}

// parseResetsAt converts an ISO 8601 timestamp to Unix seconds, 0 when absent or unparsable.
func parseResetsAt(value string) int64 {
	if value == "" {
		return 0
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return 0
	}

	return parsed.Unix()
}

// fetchModels calls the usage endpoint with the local Claude Code credentials.
func fetchModels(userAgent string) ([]ModelUsage, error) {
	token, err := oauthToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build usage request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", oauthBeta)
	req.Header.Set("Content-Type", "application/json")
	// The endpoint rate limits unknown clients aggressively.
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: fetchTimeout}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch usage: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch usage: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Limits []apiLimit `json:"limits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode usage: %w", err)
	}

	return parseModels(result.Limits), nil
}

// Refresh forces a synchronous refresh of the per-model usage cache.
func Refresh(userAgent string) error {
	models, err := fetchModels(userAgent)
	if err != nil {
		return err
	}

	cached := &CachedUsage{
		UpdatedAt: time.Now(),
		Models:    models,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("encode usage: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath()), 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	if err := os.WriteFile(cachePath(), data, 0o600); err != nil {
		return fmt.Errorf("write cache: %w", err)
	}

	cachedData = cached

	return nil
}

// CacheAge returns how old the on-disk cache is, and false when there is no cache.
func CacheAge() (time.Duration, bool) {
	cached := loadCache()

	if cached == nil {
		return 0, false
	}

	return time.Since(cached.UpdatedAt), true
}
