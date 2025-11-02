package progress

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"golang.org/x/term"
)

// terminalWriter handles cross-platform terminal output with proper line clearing
type terminalWriter struct {
	output      io.Writer // Output destination (default: os.Stdout)
	isStdout    bool      // True if output is os.Stdout
	lastLineLen int       // Track last line length for proper clearing
}

// newTerminalWriter creates a new terminal writer
func newTerminalWriter(output io.Writer) *terminalWriter {
	if output == nil {
		output = os.Stdout
	}

	// Check if output is stdout AND is a terminal (TTY)
	isStdout := false
	if output == os.Stdout {
		// Verify it's actually a TTY (not piped/redirected)
		if term.IsTerminal(int(os.Stdout.Fd())) {
			isStdout = true
		}
	} else if file, ok := output.(*os.File); ok {
		if file == os.Stdout && term.IsTerminal(int(file.Fd())) {
			isStdout = true
		}
	}

	return &terminalWriter{
		output:      output,
		isStdout:    isStdout,
		lastLineLen: 0,
	}
}

// clearAndPrint clears the current line and prints new content
// Uses platform-specific clearing strategies for Windows vs Unix-like systems
func (tw *terminalWriter) clearAndPrint(line string) {
	// Sanitize: Remove any newline or carriage return characters from the line
	line = strings.ReplaceAll(line, "\n", "")
	line = strings.ReplaceAll(line, "\r", "")

	if runtime.GOOS == "windows" {
		// Windows: Use fixed-width padding to avoid terminal artifacts
		// Generous width ensures complete clearing of previous content
		tw.printWindowsLine(line, 120)
	} else {
		// Unix/Linux/macOS: Use dynamic padding based on previous line length
		tw.printUnixLine(line)
	}

	// Don't flush for stdout - fmt.Printf auto-flushes
	// Only flush for other writers (like in tests)
	if !tw.isStdout {
		tw.flush()
	}
}

// printFinal prints the final line with a newline
func (tw *terminalWriter) printFinal(line string) {
	if runtime.GOOS == "windows" {
		tw.printWindowsLine(line, 120)
		fmt.Fprintln(tw.output)
	} else {
		tw.printUnixLine(line)
		fmt.Fprintln(tw.output)
	}

	// Only flush for non-stdout
	if !tw.isStdout {
		tw.flush()
	}
}

// printWindowsLine handles Windows-specific line printing
func (tw *terminalWriter) printWindowsLine(line string, maxWidth int) {
	// Pad line to fixed width to clear previous content
	if len(line) < maxWidth {
		line = line + strings.Repeat(" ", maxWidth-len(line))
	}

	// Write directly for stdout, use fmt for other writers
	if tw.isStdout {
		// Use fmt.Printf with explicit format - no automatic newline
		fmt.Printf("\r%s", line)
	} else {
		fmt.Fprintf(tw.output, "\r%s", line)
	}
}

// printUnixLine handles Unix-like system line printing
func (tw *terminalWriter) printUnixLine(line string) {
	// Calculate padding needed to clear previous line
	padding := ""
	if tw.lastLineLen > len(line) {
		padding = strings.Repeat(" ", tw.lastLineLen-len(line))
	}
	tw.lastLineLen = len(line)

	// Write directly for stdout, use fmt for other writers
	if tw.isStdout {
		// Use fmt.Printf with explicit format - no automatic newline
		// Explicitly flush stdout after write
		fmt.Printf("\r%s%s", line, padding)
		_ = os.Stdout.Sync() // Best effort - sync errors are rare and typically non-fatal
	} else {
		fmt.Fprintf(tw.output, "\r%s%s", line, padding)
	}
}

// clear clears the current line
func (tw *terminalWriter) clear() {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(tw.output, "\r%s\r", strings.Repeat(" ", 120))
	} else {
		if tw.lastLineLen > 0 {
			fmt.Fprintf(tw.output, "\r%s\r", strings.Repeat(" ", tw.lastLineLen))
			tw.lastLineLen = 0
		}
	}
	tw.flush()
}

// flush forces the output to be written immediately
func (tw *terminalWriter) flush() {
	// For stdout, sync immediately
	if tw.isStdout {
		_ = os.Stdout.Sync() // Best effort - sync errors are rare and typically non-fatal
		return
	}

	// Try to sync if the writer supports it
	if syncer, ok := tw.output.(interface{ Sync() error }); ok {
		_ = syncer.Sync()
	}
}
