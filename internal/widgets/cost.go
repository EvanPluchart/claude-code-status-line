package widgets

import (
	"fmt"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/exchange"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
)

var currencySymbols = map[string]string{
	"USD": "$", "EUR": "€", "GBP": "£", "JPY": "¥",
	"CAD": "CA$", "AUD": "A$", "CHF": "CHF", "CNY": "¥",
	"KRW": "₩", "INR": "₹", "BRL": "R$", "MXN": "MX$",
	"SEK": "kr", "NOK": "kr", "DKK": "kr", "PLN": "zł",
	"CZK": "Kč", "TRY": "₺", "RUB": "₽", "ZAR": "R",
}

// fallbackRates are used when the cached exchange rates are unavailable (snapshot 2026-08).
var fallbackRates = map[string]float64{
	"USD": 1, "EUR": 0.86, "GBP": 0.74, "JPY": 159.2,
	"CAD": 1.39, "AUD": 1.39, "CHF": 0.81, "CNY": 6.74,
	"KRW": 1384.6, "INR": 95.5, "BRL": 5.15, "MXN": 16.95,
	"SEK": 9.52, "NOK": 9.34, "DKK": 6.41, "PLN": 3.69,
	"CZK": 20.67, "TRY": 48.15, "RUB": 84.3, "ZAR": 15.95,
}

// europeanCurrencies use a decimal comma and a trailing symbol.
var europeanCurrencies = map[string]bool{
	"EUR": true, "CHF": true, "PLN": true, "CZK": true,
	"SEK": true, "NOK": true, "DKK": true,
}

// formatCost converts a USD amount into the configured currency and formats it.
func formatCost(costUSD float64, ctx *Context) string {
	currency := strings.ToUpper(ctx.Config.Widgets.Cost.Currency)
	decimals := ctx.Config.Widgets.Cost.Decimals

	if currency == "" {
		currency = "USD"
	}

	if decimals <= 0 {
		decimals = 2
	}

	symbol := currencySymbols[currency]

	if symbol == "" {
		symbol = currency
	}

	rate := 1.0

	if currency != "USD" {
		cached, ok := exchange.GetRate(currency)

		if ok && cached > 0 {
			rate = cached
		} else if fallback := fallbackRates[currency]; fallback > 0 {
			rate = fallback
		}
	}

	value := fmt.Sprintf("%.*f", decimals, costUSD*rate)

	if europeanCurrencies[currency] {
		return strings.Replace(value, ".", ",", 1) + symbol
	}

	return symbol + value
}

// CostWidget displays the session cost.
type CostWidget struct{}

func (w *CostWidget) ID() string { return "cost" }

func (w *CostWidget) Render(ctx *Context) string {
	costUSD := ctx.Input.Cost.TotalCostUSD

	return ansi.Colorize(formatCost(costUSD, ctx), thresholdColor(costUSD, ctx.Config.Thresholds.Cost, ctx))
}

// BurnRateWidget displays the spending rate per hour of wall-clock time.
type BurnRateWidget struct{}

func (w *BurnRateWidget) ID() string { return "burn-rate" }

func (w *BurnRateWidget) Render(ctx *Context) string {
	durationMS := ctx.Input.Cost.TotalDurationMS

	// Too early in the session to give a meaningful figure.
	if durationMS < 60_000 {
		return ""
	}

	perHour := ctx.Input.Cost.TotalCostUSD / (float64(durationMS) / 3_600_000)
	t := i18n.Get(ctx.Config.Locale)

	return ansi.Colorize(formatCost(perHour, ctx)+t.PerHourSuffix, ctx.Theme.Muted)
}
