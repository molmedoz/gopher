package progress

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/formatters"
)

// ProgressBar represents a progress bar
type ProgressBar struct {
	total      int64
	current    int64
	startTime  time.Time
	lastUpdate time.Time
	label      string
	config     *ProgressConfig
	terminal   *terminalWriter
}

// NewProgressBar creates a new progress bar with optional configuration
//
// Examples:
//
//	// Basic usage (backward compatible)
//	pb := NewProgressBar(1024*1024, "Downloading")
//
//	// With options
//	pb := NewProgressBar(1024*1024, "Downloading", WithWidth(80), WithMinimal())
func NewProgressBar(total int64, label string, options ...Option) *ProgressBar {
	config := applyOptions(options...)

	return &ProgressBar{
		total:      total,
		current:    0,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
		label:      label,
		config:     config,
		terminal:   newTerminalWriter(nil),
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(current int64) {
	pb.current = current
	now := time.Now()

	// Only update display at configured throttle rate to avoid flickering
	if now.Sub(pb.lastUpdate) < pb.config.UpdateThrottle && pb.current < pb.total {
		return
	}
	pb.lastUpdate = now

	pb.display()
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.current = pb.total
	pb.display()
	// Print final line with newline
	if !pb.config.Silent {
		fmt.Println()
	}
}

// display renders the progress bar
func (pb *ProgressBar) display() {
	// Skip display if silent mode
	if pb.config.Silent {
		return
	}

	if pb.total == 0 {
		return
	}

	// Calculate percentage
	percentage := float64(pb.current) / float64(pb.total)

	// Build the progress line
	line := pb.buildLine(percentage)

	// Output via terminal writer
	pb.terminal.clearAndPrint(line)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// buildLine builds the progress line based on configuration
func (pb *ProgressBar) buildLine(percentage float64) string {
	// Calculate filled width
	filled := int(percentage * float64(pb.config.Width))
	if filled < 0 {
		filled = 0
	}
	if filled > pb.config.Width {
		filled = pb.config.Width
	}

	// Create progress bar string
	bar := strings.Repeat(pb.config.FilledChar, filled) +
		strings.Repeat(pb.config.EmptyChar, pb.config.Width-filled)

	// Calculate speed
	elapsed := time.Since(pb.startTime)
	var speed float64
	if elapsed.Seconds() > 0 {
		speed = float64(pb.current) / elapsed.Seconds()
	}

	// Build line parts
	parts := []string{
		pb.label,
		fmt.Sprintf("[%s]", bar),
		formatters.FormatPercentage(percentage),
	}

	// Add bytes if configured and not minimal
	if pb.config.ShowBytes && !pb.config.Minimal {
		parts = append(parts, fmt.Sprintf("%s/%s",
			formatters.FormatBytes(pb.current),
			formatters.FormatBytes(pb.total)))
	}

	// Add speed if configured and not minimal
	if pb.config.ShowSpeed && !pb.config.Minimal {
		parts = append(parts, formatters.FormatSpeed(speed))
	}

	return strings.Join(parts, " ")
}

// ProgressWriter wraps an io.Writer to track progress
type ProgressWriter struct {
	writer io.Writer
	bar    *ProgressBar
}

// NewProgressWriter creates a new progress writer
func NewProgressWriter(writer io.Writer, bar *ProgressBar) *ProgressWriter {
	return &ProgressWriter{
		writer: writer,
		bar:    bar,
	}
}

// Write implements io.Writer
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.writer.Write(p)
	if pw.bar != nil {
		pw.bar.Update(pw.bar.current + int64(n))
	}
	return n, err
}

// Spinner represents a simple spinner for indeterminate progress
type Spinner struct {
	frames   []string
	index    int
	label    string
	ticker   *time.Ticker
	done     chan bool
	config   *ProgressConfig
	terminal *terminalWriter
}

// NewSpinner creates a new spinner with optional configuration
func NewSpinner(label string, options ...Option) *Spinner {
	config := applyOptions(options...)
	return &Spinner{
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:    0,
		label:    label,
		ticker:   time.NewTicker(config.UpdateThrottle),
		done:     make(chan bool),
		config:   config,
		terminal: newTerminalWriter(nil),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				if !s.config.Silent {
					// Build the spinner line
					line := fmt.Sprintf("%s %s", s.frames[s.index], s.label)
					s.terminal.clearAndPrint(line)
				}
				s.index = (s.index + 1) % len(s.frames)
			case <-s.done:
				return
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.ticker.Stop()
	s.done <- true

	if !s.config.Silent {
		// Build final line with checkmark
		line := fmt.Sprintf("%s %s", "✓", s.label)
		s.terminal.printFinal(line)
	}
}

// SimpleProgress shows simple progress messages
func SimpleProgress(message string) {
	fmt.Printf("⏳ %s...\n", message)
}

// CompleteProgress shows completion message
func CompleteProgress(message string) {
	fmt.Printf("✓ %s\n", message)
}
