package progress

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
)

func TestNewTerminalWriter(t *testing.T) {
	// Test with nil output (should use os.Stdout)
	tw := newTerminalWriter(nil)
	if tw.output == nil {
		t.Error("Expected output to be set when nil is provided")
	}
	if tw.lastLineLen != 0 {
		t.Errorf("Expected lastLineLen to be 0, got %d", tw.lastLineLen)
	}

	// Test with custom output
	var buf bytes.Buffer
	tw = newTerminalWriter(&buf)
	if tw.output != &buf {
		t.Error("Expected output to be the provided buffer")
	}
}

func TestTerminalWriter_ClearAndPrint(t *testing.T) {
	tests := []struct {
		name        string
		lines       []string
		wantContain string
	}{
		{
			name:        "single line",
			lines:       []string{"Hello World"},
			wantContain: "Hello World",
		},
		{
			name:        "multiple lines",
			lines:       []string{"Line 1", "Line 2", "Line 3"},
			wantContain: "Line 3",
		},
		{
			name:        "longer then shorter",
			lines:       []string{"This is a very long line", "Short"},
			wantContain: "Short",
		},
		{
			name:        "empty line",
			lines:       []string{"Content", ""},
			wantContain: "\r",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tw := newTerminalWriter(&buf)

			for _, line := range tt.lines {
				tw.clearAndPrint(line)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("Output should contain %q, got %q", tt.wantContain, output)
			}

			// Check for carriage return
			if !strings.Contains(output, "\r") {
				t.Error("Output should contain carriage return for line clearing")
			}
		})
	}
}

func TestTerminalWriter_PrintFinal(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{"simple message", "Done!"},
		{"with checkmark", "✓ Complete"},
		{"empty line", ""},
		{"long message", strings.Repeat("x", 150)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tw := newTerminalWriter(&buf)

			tw.printFinal(tt.line)

			output := buf.String()

			// Should contain the line content
			if !strings.Contains(output, tt.line) {
				t.Errorf("Output should contain %q, got %q", tt.line, output)
			}

			// Should end with newline
			if !strings.HasSuffix(output, "\n") {
				t.Error("Final output should end with newline")
			}
		})
	}
}

func TestTerminalWriter_Clear(t *testing.T) {
	var buf bytes.Buffer
	tw := newTerminalWriter(&buf)

	// Print something first
	tw.clearAndPrint("Some content")
	buf.Reset()

	// Clear the line
	tw.clear()

	output := buf.String()

	// Should contain carriage return
	if !strings.Contains(output, "\r") {
		t.Error("Clear should use carriage return")
	}

	// Should contain spaces for clearing
	if !strings.Contains(output, " ") {
		t.Error("Clear should use spaces to clear the line")
	}
}

func TestTerminalWriter_PlatformSpecific(t *testing.T) {
	var buf bytes.Buffer
	tw := newTerminalWriter(&buf)

	testLine := "Platform test line"
	tw.clearAndPrint(testLine)

	output := buf.String()

	// Should contain the test line
	if !strings.Contains(output, testLine) {
		t.Errorf("Output should contain %q", testLine)
	}

	// On Windows, should pad to fixed width
	if runtime.GOOS == "windows" {
		if len(output) < 120 {
			t.Log("Note: Windows output should be padded to 120 chars for clearing")
		}
	}

	// Should start with carriage return
	if !strings.HasPrefix(output, "\r") {
		t.Error("Output should start with carriage return")
	}
}

func TestTerminalWriter_UnixPadding(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	var buf bytes.Buffer
	tw := newTerminalWriter(&buf)

	// Print a long line
	longLine := strings.Repeat("x", 80)
	tw.clearAndPrint(longLine)
	buf.Reset()

	// Print a shorter line - should have padding
	shortLine := "short"
	tw.clearAndPrint(shortLine)

	output := buf.String()

	// Should contain more characters than just the short line
	// (accounting for padding)
	if len(output) <= len(shortLine)+2 { // +2 for \r and potential formatting
		t.Error("Unix systems should add padding when printing shorter lines")
	}
}

func TestTerminalWriter_Sequence(t *testing.T) {
	var buf bytes.Buffer
	tw := newTerminalWriter(&buf)

	// Simulate a typical progress sequence
	lines := []string{
		"Downloading... 0%",
		"Downloading... 25%",
		"Downloading... 50%",
		"Downloading... 75%",
		"Downloading... 100%",
	}

	for _, line := range lines {
		tw.clearAndPrint(line)
	}

	// Print final message
	buf.Reset()
	tw.printFinal("✓ Complete")

	output := buf.String()

	// Final output should contain completion message
	if !strings.Contains(output, "Complete") {
		t.Error("Final output should contain completion message")
	}

	// Should end with newline
	if !strings.HasSuffix(output, "\n") {
		t.Error("Final output should end with newline")
	}
}

// Benchmark tests
func BenchmarkTerminalWriter_ClearAndPrint(b *testing.B) {
	var buf bytes.Buffer
	tw := newTerminalWriter(&buf)
	line := "Benchmark test line with some content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tw.clearAndPrint(line)
	}
}

func BenchmarkTerminalWriter_PrintFinal(b *testing.B) {
	var buf bytes.Buffer
	tw := newTerminalWriter(&buf)
	line := "Final benchmark line"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		tw.printFinal(line)
	}
}
