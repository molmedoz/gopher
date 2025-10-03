package progress

import "time"

// ProgressConfig holds configuration for progress display
type ProgressConfig struct {
	// Width is the width of the progress bar in characters
	Width int

	// UpdateThrottle is the minimum time between display updates
	UpdateThrottle time.Duration

	// ShowSpeed determines whether to show download speed
	ShowSpeed bool

	// ShowBytes determines whether to show byte counts
	ShowBytes bool

	// FilledChar is the character used for filled portion of bar
	FilledChar string

	// EmptyChar is the character used for empty portion of bar
	EmptyChar string

	// Silent disables all output (useful for scripting)
	Silent bool

	// Minimal shows minimal output (just percentage, no speed/bytes)
	Minimal bool
}

// Option is a functional option for configuring progress
type Option func(*ProgressConfig)

// defaultConfig returns the default configuration
func defaultConfig() *ProgressConfig {
	// Try to detect terminal width and adjust bar width accordingly
	barWidth := 50 // default

	// Get terminal width if available
	width, _, err := getTerminalSize()
	if err == nil && width > 0 {
		// Reserve space for label, brackets, percentage, bytes, speed
		// Roughly: "Label [bar] 100.0% 999.9 MB/999.9 MB 999.9 MB/s"
		reserved := 80 // approximate space for text
		barWidth = width - reserved

		// Constrain bar width to reasonable limits
		if barWidth < 20 {
			barWidth = 20
		} else if barWidth > 100 {
			barWidth = 100
		}
	}

	return &ProgressConfig{
		Width:          barWidth,
		UpdateThrottle: 100 * time.Millisecond,
		ShowSpeed:      true,
		ShowBytes:      true,
		FilledChar:     "█",
		EmptyChar:      "░",
		Silent:         false,
		Minimal:        false,
	}
}

// WithWidth sets the progress bar width
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithWidth(80))
func WithWidth(width int) Option {
	return func(c *ProgressConfig) {
		if width > 0 {
			c.Width = width
		}
	}
}

// WithUpdateThrottle sets the minimum time between updates
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithUpdateThrottle(50*time.Millisecond))
func WithUpdateThrottle(duration time.Duration) Option {
	return func(c *ProgressConfig) {
		if duration > 0 {
			c.UpdateThrottle = duration
		}
	}
}

// WithSpeed enables or disables speed display
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithSpeed(false))
func WithSpeed(show bool) Option {
	return func(c *ProgressConfig) {
		c.ShowSpeed = show
	}
}

// WithBytes enables or disables byte count display
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithBytes(false))
func WithBytes(show bool) Option {
	return func(c *ProgressConfig) {
		c.ShowBytes = show
	}
}

// WithChars sets custom characters for filled and empty portions
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithChars("=", "-"))
func WithChars(filled, empty string) Option {
	return func(c *ProgressConfig) {
		if filled != "" {
			c.FilledChar = filled
		}
		if empty != "" {
			c.EmptyChar = empty
		}
	}
}

// WithSilent disables all output
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithSilent())
func WithSilent() Option {
	return func(c *ProgressConfig) {
		c.Silent = true
	}
}

// WithMinimal enables minimal output (percentage only, no speed/bytes)
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithMinimal())
func WithMinimal() Option {
	return func(c *ProgressConfig) {
		c.Minimal = true
		c.ShowSpeed = false
		c.ShowBytes = false
	}
}

// WithCustom allows setting multiple config values at once
//
// Example:
//
//	pb := NewProgressBar(100, "Download", WithCustom(func(c *ProgressConfig) {
//	    c.Width = 80
//	    c.ShowSpeed = false
//	}))
func WithCustom(fn func(*ProgressConfig)) Option {
	return fn
}

// applyOptions applies functional options to a config
func applyOptions(options ...Option) *ProgressConfig {
	config := defaultConfig()
	for _, opt := range options {
		opt(config)
	}
	return config
}
