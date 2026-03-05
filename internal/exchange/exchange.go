package exchange

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	apiURL       = "https://open.er-api.com/v6/latest/USD"
	maxAge       = 24 * time.Hour
	fetchTimeout = 3 * time.Second
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
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".claude-statusline", "rates.json")
}

// GetRate returns the exchange rate for a currency.
// It loads from cache, refreshes in the background if stale, and returns 0 if unknown.
func GetRate(currency string) (float64, bool) {
	once.Do(func() {
		cachedData = loadCache()

		if cachedData != nil && time.Since(cachedData.UpdatedAt) > maxAge {
			go refreshCache()
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

func refreshCache() {
	rates, err := fetchRates()
	if err != nil {
		return
	}

	cached := &CachedRates{
		UpdatedAt: time.Now(),
		Rates:     rates,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return
	}

	dir := filepath.Dir(cachePath())
	_ = os.MkdirAll(dir, 0o755)
	_ = os.WriteFile(cachePath(), data, 0o644)

	cachedData = cached
}

func fetchRates() (map[string]float64, error) {
	client := &http.Client{Timeout: fetchTimeout}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Rates, nil
}

// Refresh forces a synchronous refresh of the exchange rates cache.
// Returns an error if the fetch fails.
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
		return err
	}

	dir := filepath.Dir(cachePath())
	_ = os.MkdirAll(dir, 0o755)
	_ = os.WriteFile(cachePath(), data, 0o644)

	cachedData = cached

	return nil
}
