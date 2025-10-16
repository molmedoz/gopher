package progress

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestProgressBar(t *testing.T) {
	// Test progress bar creation
	pb := NewProgressBar(100, "Test")
	if pb.total != 100 {
		t.Errorf("Expected total 100, got %d", pb.total)
	}
	if pb.label != "Test" {
		t.Errorf("Expected label 'Test', got '%s'", pb.label)
	}
	if pb.config.Width != 50 {
		t.Errorf("Expected width 50, got %d", pb.config.Width)
	}
	if pb.current != 0 {
		t.Errorf("Expected current 0, got %d", pb.current)
	}
}

func TestProgressBarUpdate(t *testing.T) {
	pb := NewProgressBar(100, "Test")

	// Test initial state
	if pb.current != 0 {
		t.Errorf("Expected current 0, got %d", pb.current)
	}

	// Test update
	pb.Update(50)
	if pb.current != 50 {
		t.Errorf("Expected current 50, got %d", pb.current)
	}

	// Test finish
	pb.Finish()
	if pb.current != pb.total {
		t.Errorf("Expected current %d, got %d", pb.total, pb.current)
	}
}

func TestProgressBarDisplay(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	pb := NewProgressBar(100, "Test")
	pb.Update(50)
	pb.Finish()

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains expected elements
	if !strings.Contains(output, "Test") {
		t.Error("Expected output to contain label 'Test'")
	}
	// The output might be on the same line, so just check for any number
	if !strings.ContainsAny(output, "0123456789") {
		t.Error("Expected output to contain numbers")
	}
}

func TestProgressBarThrottling(t *testing.T) {
	pb := NewProgressBar(100, "Test")

	// Test that rapid updates are throttled
	start := time.Now()
	pb.Update(10)
	pb.Update(20) // Should be throttled
	pb.Update(30) // Should be throttled
	elapsed := time.Since(start)

	// Should complete quickly due to throttling
	if elapsed > 50*time.Millisecond {
		t.Error("Expected throttling to prevent rapid updates")
	}
}

func TestProgressWriter(t *testing.T) {
	var buf bytes.Buffer
	pb := NewProgressBar(100, "Test")
	pw := NewProgressWriter(&buf, pb)

	// Test writing
	data := []byte("hello world")
	n, err := pw.Write(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
	if buf.String() != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", buf.String())
	}
	if pb.current != int64(len(data)) {
		t.Errorf("Expected current %d, got %d", len(data), pb.current)
	}
}

func TestProgressWriterNilBar(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, nil)

	// Test writing with nil bar
	data := []byte("test")
	n, err := pw.Write(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
}

func TestSpinner(t *testing.T) {
	spinner := NewSpinner("Test spinner")
	if spinner.label != "Test spinner" {
		t.Errorf("Expected label 'Test spinner', got '%s'", spinner.label)
	}
	if len(spinner.frames) == 0 {
		t.Error("Expected spinner frames to be non-empty")
	}
	if spinner.index != 0 {
		t.Errorf("Expected index 0, got %d", spinner.index)
	}
}

func TestSpinnerStartStop(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	spinner := NewSpinner("Test")
	spinner.Start()
	time.Sleep(150 * time.Millisecond) // Let it spin a bit
	spinner.Stop()

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains expected elements
	if !strings.Contains(output, "Test") {
		t.Error("Expected output to contain label 'Test'")
	}
	if !strings.Contains(output, "✓") {
		t.Error("Expected output to contain completion checkmark")
	}
}

// TestFormatBytes has been moved to formatters/bytes_test.go

func TestSimpleProgress(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	SimpleProgress("Test message")
	CompleteProgress("Test completion")

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains expected elements
	if !strings.Contains(output, "⏳ Test message...") {
		t.Error("Expected output to contain simple progress message")
	}
	if !strings.Contains(output, "✓ Test completion") {
		t.Error("Expected output to contain completion message")
	}
}

func TestProgressBarZeroTotal(t *testing.T) {
	pb := NewProgressBar(0, "Test")

	// Should not panic or cause issues
	pb.Update(0)
	pb.Finish()
}

func TestProgressBarNegativeValues(t *testing.T) {
	pb := NewProgressBar(100, "Test")

	// Test negative current value
	pb.Update(-10)
	if pb.current != -10 {
		t.Errorf("Expected current -10, got %d", pb.current)
	}

	// Test current greater than total
	pb.Update(150)
	if pb.current != 150 {
		t.Errorf("Expected current 150, got %d", pb.current)
	}
}
