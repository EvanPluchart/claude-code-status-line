// Package exchange manages the USD exchange rates cache used by the cost widget.
//
// The render path never performs network calls: it only reads the on-disk cache.
// When the cache is stale, a detached background process is spawned to refresh it
// so that the statusline stays within its execution budget.
package exchange

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
	apiURL       = "https://open.er-api.com/v6/latest/USD"
	maxAge       = 24 * time.Hour
	retryDelay   = 30 * time.Minute
	fetchTimeout = 5 * time.Second
)

// CachedRates is the on-disk cache structure.
type CachedRates struct {
	UpdatedAt time.Time          `json:"updated_at"`
	Rates     map[string]float64 `json:"rates"`
}

var (
	once       sync.Once
	cachedData *CachedRates
)

func cachePath() string {
	return filepath.Join(config.ConfigDir(), "rates.json")
}

// attemptPath marks the last background refresh attempt to avoid spawn storms.
func attemptPath() string {
	return filepath.Join(config.ConfigDir(), "rates.attempt")
}

// GetRate returns the exchange rate for a currency from the local cache.
// When the cache is stale, a background refresh is scheduled.
func GetRate(currency string) (float64, bool) {
	once.Do(func() {
		cachedData = loadCache()

		if cachedData == nil || time.Since(cachedData.UpdatedAt) > maxAge {
			scheduleRefresh()
		}
	})

	if cachedData == nil {
		return 0, false
	}

	rate, ok := cachedData.Rates[currency]

	return rate, ok
}

func loadCache() *CachedRates {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil
	}

	var cached CachedRates

	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}

	return &cached
}

// scheduleRefresh spawns a detached "update-rates" process, at most once per retryDelay.
func scheduleRefresh() {
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

	_ = detach.Spawn("update-rates", "--quiet")
}

func fetchRates() (map[string]float64, error) {
	client := &http.Client{Timeout: fetchTimeout}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("fetch rates: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch rates: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode rates: %w", err)
	}

	if len(result.Rates) == 0 {
		return nil, fmt.Errorf("decode rates: empty response")
	}

	return result.Rates, nil
}

// Refresh forces a synchronous refresh of the exchange rates cache.
func Refresh() error {
	rates, err := fetchRates()
	if err != nil {
		return err
	}

	cached := &CachedRates{
		UpdatedAt: time.Now(),
		Rates:     rates,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("encode rates: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath()), 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	if err := os.WriteFile(cachePath(), data, 0o644); err != nil {
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
