package widgets

import (
	"fmt"
	"strings"

	"github.com/EvanPluchart/claude-code-status-line/internal/ansi"
	"github.com/EvanPluchart/claude-code-status-line/internal/exchange"
)

var currencySymbols = map[string]string{
	"USD": "$", "EUR": "\u20ac", "GBP": "\u00a3", "JPY": "\u00a5",
	"CAD": "CA$", "AUD": "A$", "CHF": "CHF", "CNY": "\u00a5",
	"KRW": "\u20a9", "INR": "\u20b9", "BRL": "R$", "MXN": "MX$",
	"SEK": "kr", "NOK": "kr", "DKK": "kr", "PLN": "z\u0142",
	"CZK": "K\u010d", "TRY": "\u20ba", "RUB": "\u20bd", "ZAR": "R",
}

// fallbackRates are used when the cached exchange rates are unavailable.
var fallbackRates = map[string]float64{
	"USD": 1, "EUR": 0.92, "GBP": 0.79, "JPY": 149.5,
	"CAD": 1.36, "AUD": 1.53, "CHF": 0.88, "CNY": 7.24,
	"KRW": 1320, "INR": 83.1, "BRL": 4.97, "MXN": 17.15,
	"SEK": 10.42, "NOK": 10.55, "DKK": 6.88, "PLN": 3.98,
	"CZK": 22.8, "TRY": 30.2, "RUB": 89.5, "ZAR": 18.6,
}

var europeanCurrencies = map[string]bool{
	"EUR": true, "CHF": true, "PLN": true, "CZK": true,
	"SEK": true, "NOK": true, "DKK": true,
}

// CostWidget displays the session cost.
type CostWidget struct{}

func (w *CostWidget) ID() string { return "cost" }

func (w *CostWidget) Render(ctx *Context) string {
	costUSD := ctx.Input.Cost.TotalCostUSD
	currency := ctx.Config.Widgets.Cost.Currency
	decimals := ctx.Config.Widgets.Cost.Decimals

	if currency == "" {
		currency = "USD"
	}

	if decimals == 0 {
		decimals = 2
	}

	symbol := currencySymbols[currency]
	if symbol == "" {
		symbol = currency
	}

	rate, ok := exchange.GetRate(currency)
	if !ok {
		rate = fallbackRates[currency]
	}

	if rate == 0 {
		rate = 1
	}

	converted := costUSD * rate
	format := fmt.Sprintf("%%.%df", decimals)
	value := fmt.Sprintf(format, converted)

	var formatted string

	if europeanCurrencies[currency] {
		formatted = strings.Replace(value, ".", ",", 1) + symbol
	} else {
		formatted = symbol + value
	}

	color := ctx.Theme.Success
	thresholds := ctx.Config.Thresholds.Cost

	if costUSD >= thresholds.Red {
		color = ctx.Theme.Danger + ansi.Bold
	} else if costUSD >= thresholds.Orange {
		color = ctx.Theme.Danger
	} else if costUSD >= thresholds.Yellow {
		color = ctx.Theme.Warning
	}

	return color + formatted + ansi.RST
}
