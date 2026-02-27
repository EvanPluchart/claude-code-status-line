package ansi

import "os"

const (
	RST  = "\033[0m"
	Bold = "\033[1m"
	Dim  = "\033[2m"

	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Magenta      = "\033[35m"
	Cyan         = "\033[36m"
	White        = "\033[37m"
	Gray         = "\033[90m"
	BrightRed    = "\033[91m"
	BrightGreen  = "\033[92m"
	BrightYellow = "\033[93m"
	BrightBlue   = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"
)

// Colorize wraps text with an ANSI color and reset.
func Colorize(text, color string) string {
	if text == "" {
		return ""
	}

	return color + text + RST
}

// ColorBold wraps text with color + bold and reset.
func ColorBold(text, color string) string {
	if text == "" {
		return ""
	}

	return color + Bold + text + RST
}

// SupportsColor checks if the terminal supports ANSI colors.
func SupportsColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}

	// Check if stdout is a terminal
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return fi.Mode()&os.ModeCharDevice != 0
}
