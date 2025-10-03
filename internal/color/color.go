// Package color provides ANSI color support for terminal output in Gopher.
//
// This package implements a simple color utility system that automatically
// detects terminal capabilities and provides fallback support for environments
// that don't support colors.
//
// Features:
//   - ANSI color code support
//   - Automatic terminal detection
//   - Fallback to plain text when colors aren't supported
//   - Predefined color functions for common use cases
//   - Thread-safe color functions
//
// Usage:
//
//	// Check if colors are supported
//	if color.IsColorEnabled() {
//	    // Use colors
//	}
//
//	// Create color functions
//	red := color.RedColor()
//	green := color.GreenColor()
//	bold := color.BoldColor()
//
//	// Apply colors to text
//	coloredText := red("Error message")
//	successText := green("Success message")
//	boldText := bold("Important text")
//
//	// Use predefined functions for common cases
//	activeColor := color.ActiveVersion()
//	systemColor := color.SystemVersion()
//	inactiveColor := color.InactiveVersion()
//
// The package automatically detects if the output is going to a terminal
// and disables colors when appropriate (e.g., when redirecting to a file).
package color

import (
	"os"
	"runtime"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
)

// ColorFunc represents a function that applies color to text
type ColorFunc func(string) string

// Disabled returns a function that doesn't apply any color
func Disabled() ColorFunc {
	return func(text string) string {
		return text
	}
}

// RedColor returns a function that applies red color
func RedColor() ColorFunc {
	return func(text string) string {
		return Red + text + Reset
	}
}

// GreenColor returns a function that applies green color
func GreenColor() ColorFunc {
	return func(text string) string {
		return Green + text + Reset
	}
}

// YellowColor returns a function that applies yellow color
func YellowColor() ColorFunc {
	return func(text string) string {
		return Yellow + text + Reset
	}
}

// BlueColor returns a function that applies blue color
func BlueColor() ColorFunc {
	return func(text string) string {
		return Blue + text + Reset
	}
}

// CyanColor returns a function that applies cyan color
func CyanColor() ColorFunc {
	return func(text string) string {
		return Cyan + text + Reset
	}
}

// BoldColor returns a function that applies bold formatting
func BoldColor() ColorFunc {
	return func(text string) string {
		return Bold + text + Reset
	}
}

// DimColor returns a function that applies dim formatting
func DimColor() ColorFunc {
	return func(text string) string {
		return Dim + text + Reset
	}
}

// IsColorEnabled checks if color output is enabled
func IsColorEnabled() bool {
	// Check if we're in a terminal
	if runtime.GOOS == "windows" {
		// On Windows, check for ANSI support
		return os.Getenv("TERM") != "" || os.Getenv("ANSICON") != ""
	}

	// On Unix-like systems, check if stdout is a terminal
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Check if it's a character device (terminal)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// NewColorFunc creates a color function that respects color settings
func NewColorFunc(colorCode string) ColorFunc {
	if !IsColorEnabled() {
		return Disabled()
	}

	return func(text string) string {
		return colorCode + text + Reset
	}
}

// ActiveVersion creates a color function for active version highlighting
func ActiveVersion() ColorFunc {
	if !IsColorEnabled() {
		return Disabled()
	}

	return func(text string) string {
		return Bold + Green + text + Reset
	}
}

// SystemVersion creates a color function for system version highlighting
func SystemVersion() ColorFunc {
	if !IsColorEnabled() {
		return Disabled()
	}

	return func(text string) string {
		return Cyan + text + Reset
	}
}

// InactiveVersion creates a color function for inactive version highlighting
func InactiveVersion() ColorFunc {
	if !IsColorEnabled() {
		return Disabled()
	}

	return func(text string) string {
		return Dim + text + Reset
	}
}
