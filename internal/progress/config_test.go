package progress

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()

	// Test default values
	if config.Width != 50 {
		t.Errorf("Expected default width 50, got %d", config.Width)
	}
	if config.UpdateThrottle != 100*time.Millisecond {
		t.Errorf("Expected default throttle 100ms, got %v", config.UpdateThrottle)
	}
	if !config.ShowSpeed {
		t.Error("Expected ShowSpeed to be true by default")
	}
	if !config.ShowBytes {
		t.Error("Expected ShowBytes to be true by default")
	}
	if config.FilledChar != "█" {
		t.Errorf("Expected filled char '█', got '%s'", config.FilledChar)
	}
	if config.EmptyChar != "░" {
		t.Errorf("Expected empty char '░', got '%s'", config.EmptyChar)
	}
	if config.Silent {
		t.Error("Expected Silent to be false by default")
	}
	if config.Minimal {
		t.Error("Expected Minimal to be false by default")
	}
}

func TestWithWidth(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected int
	}{
		{"valid width", 80, 80},
		{"small width", 10, 10},
		{"large width", 200, 200},
		{"zero width", 0, 50},      // Should keep default
		{"negative width", -1, 50}, // Should keep default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := applyOptions(WithWidth(tt.width))
			if config.Width != tt.expected {
				t.Errorf("Expected width %d, got %d", tt.expected, config.Width)
			}
		})
	}
}

func TestWithUpdateThrottle(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected time.Duration
	}{
		{"50ms", 50 * time.Millisecond, 50 * time.Millisecond},
		{"1 second", 1 * time.Second, 1 * time.Second},
		{"zero", 0, 100 * time.Millisecond},                         // Should keep default
		{"negative", -1 * time.Millisecond, 100 * time.Millisecond}, // Should keep default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := applyOptions(WithUpdateThrottle(tt.duration))
			if config.UpdateThrottle != tt.expected {
				t.Errorf("Expected throttle %v, got %v", tt.expected, config.UpdateThrottle)
			}
		})
	}
}

func TestWithSpeed(t *testing.T) {
	// Test enabling
	config := applyOptions(WithSpeed(true))
	if !config.ShowSpeed {
		t.Error("Expected ShowSpeed to be true")
	}

	// Test disabling
	config = applyOptions(WithSpeed(false))
	if config.ShowSpeed {
		t.Error("Expected ShowSpeed to be false")
	}
}

func TestWithBytes(t *testing.T) {
	// Test enabling
	config := applyOptions(WithBytes(true))
	if !config.ShowBytes {
		t.Error("Expected ShowBytes to be true")
	}

	// Test disabling
	config = applyOptions(WithBytes(false))
	if config.ShowBytes {
		t.Error("Expected ShowBytes to be false")
	}
}

func TestWithChars(t *testing.T) {
	tests := []struct {
		name          string
		filled        string
		empty         string
		expectedFill  string
		expectedEmpty string
	}{
		{"equals and dash", "=", "-", "=", "-"},
		{"hash and dot", "#", ".", "#", "."},
		{"emoji", "🟩", "⬜", "🟩", "⬜"},
		{"empty filled", "", "-", "█", "-"}, // Keep default filled
		{"empty empty", "=", "", "=", "░"},  // Keep default empty
		{"both empty", "", "", "█", "░"},    // Keep both defaults
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := applyOptions(WithChars(tt.filled, tt.empty))
			if config.FilledChar != tt.expectedFill {
				t.Errorf("Expected filled char '%s', got '%s'", tt.expectedFill, config.FilledChar)
			}
			if config.EmptyChar != tt.expectedEmpty {
				t.Errorf("Expected empty char '%s', got '%s'", tt.expectedEmpty, config.EmptyChar)
			}
		})
	}
}

func TestWithSilent(t *testing.T) {
	config := applyOptions(WithSilent())
	if !config.Silent {
		t.Error("Expected Silent to be true")
	}

	// Silent shouldn't affect other settings
	if config.Width != 50 {
		t.Error("Silent mode shouldn't change width")
	}
}

func TestWithMinimal(t *testing.T) {
	config := applyOptions(WithMinimal())

	if !config.Minimal {
		t.Error("Expected Minimal to be true")
	}
	if config.ShowSpeed {
		t.Error("Minimal mode should disable speed display")
	}
	if config.ShowBytes {
		t.Error("Minimal mode should disable bytes display")
	}
}

func TestWithCustom(t *testing.T) {
	config := applyOptions(WithCustom(func(c *ProgressConfig) {
		c.Width = 100
		c.ShowSpeed = false
		c.FilledChar = "X"
	}))

	if config.Width != 100 {
		t.Errorf("Expected width 100, got %d", config.Width)
	}
	if config.ShowSpeed {
		t.Error("Expected ShowSpeed to be false")
	}
	if config.FilledChar != "X" {
		t.Errorf("Expected filled char 'X', got '%s'", config.FilledChar)
	}
}

func TestMultipleOptions(t *testing.T) {
	config := applyOptions(
		WithWidth(80),
		WithSpeed(false),
		WithBytes(false),
		WithChars("=", "-"),
	)

	if config.Width != 80 {
		t.Errorf("Expected width 80, got %d", config.Width)
	}
	if config.ShowSpeed {
		t.Error("Expected ShowSpeed to be false")
	}
	if config.ShowBytes {
		t.Error("Expected ShowBytes to be false")
	}
	if config.FilledChar != "=" {
		t.Errorf("Expected filled char '=', got '%s'", config.FilledChar)
	}
	if config.EmptyChar != "-" {
		t.Errorf("Expected empty char '-', got '%s'", config.EmptyChar)
	}
}

func TestOptionOrder(t *testing.T) {
	// Last option should win
	config := applyOptions(
		WithWidth(80),
		WithWidth(100),
		WithWidth(60),
	)

	if config.Width != 60 {
		t.Errorf("Expected width 60 (last option), got %d", config.Width)
	}
}

func TestConflictingOptions(t *testing.T) {
	// Minimal should override previous settings
	config := applyOptions(
		WithSpeed(true),
		WithBytes(true),
		WithMinimal(), // Should disable both
	)

	if config.ShowSpeed {
		t.Error("WithMinimal should disable speed even if previously enabled")
	}
	if config.ShowBytes {
		t.Error("WithMinimal should disable bytes even if previously enabled")
	}

	// But can be overridden after
	config = applyOptions(
		WithMinimal(),
		WithSpeed(true),
	)

	if !config.ShowSpeed {
		t.Error("Options after WithMinimal should be able to re-enable features")
	}
}

func TestNoOptions(t *testing.T) {
	config := applyOptions()

	// Should be same as default
	defaultCfg := defaultConfig()
	if config.Width != defaultCfg.Width {
		t.Error("No options should produce default config")
	}
	if config.ShowSpeed != defaultCfg.ShowSpeed {
		t.Error("No options should produce default config")
	}
}

// Benchmark tests
func BenchmarkApplyOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		applyOptions(
			WithWidth(80),
			WithSpeed(false),
			WithChars("=", "-"),
		)
	}
}

func BenchmarkDefaultConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defaultConfig()
	}
}
