package errors

import (
	"testing"
)

func TestErrorLogger(t *testing.T) {
	logger := NewErrorLogger(LogLevelInfo)
	if logger == nil {
		t.Fatal("Expected error logger to be created")
	}

	// Test logging a regular error
	err := New(ErrCodeInvalidVersion, "invalid version")
	logger.LogError(err, map[string]interface{}{"test": "value"})

	// Test logging with format
	logger.LogErrorf(LogLevelInfo, "test error: %v", err)

	// Test logging GopherError
	gopherErr := New(ErrCodeInvalidVersion, "invalid version")
	logger.LogGopherError(gopherErr, map[string]interface{}{"test": "value"})
}

func TestLogLevel(t *testing.T) {
	logger := NewErrorLogger(LogLevelInfo)

	// Test setting and getting log level
	logger.SetLevel(LogLevelError)
	if logger.GetLevel() != LogLevelError {
		t.Errorf("Expected log level %v, got %v", LogLevelError, logger.GetLevel())
	}

	logger.SetLevel(LogLevelInfo)
	if logger.GetLevel() != LogLevelInfo {
		t.Errorf("Expected log level %v, got %v", LogLevelInfo, logger.GetLevel())
	}

	logger.SetLevel(LogLevelDebug)
	if logger.GetLevel() != LogLevelDebug {
		t.Errorf("Expected log level %v, got %v", LogLevelDebug, logger.GetLevel())
	}
}

func TestLogError(t *testing.T) {
	logger := NewErrorLogger(LogLevelInfo)

	// Test logging different types of errors
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "GopherError",
			err:  New(ErrCodeInvalidVersion, "invalid version"),
		},
		{
			name: "GopherError with context",
			err:  New(ErrCodeInvalidVersion, "invalid version").WithContext("version", "1.2.3"),
		},
		{
			name: "GopherError with details",
			err:  New(ErrCodeInvalidVersion, "invalid version").WithDetails("additional info"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These functions don't return errors, so we just test that they don't panic
			logger.LogError(tt.err, map[string]interface{}{"test": "value"})
			logger.LogErrorf(LogLevelInfo, "test error: %v", tt.err)
			if gopherErr, ok := tt.err.(*GopherError); ok {
				logger.LogGopherError(gopherErr, map[string]interface{}{"test": "value"})
			}
		})
	}
}

// TestErrorLogger_Comprehensive tests the error logger comprehensively
func TestErrorLogger_Comprehensive(t *testing.T) {
	// Test creating logger with different log levels
	logger1 := NewErrorLogger(LogLevelDebug)
	logger2 := NewErrorLogger(LogLevelInfo)
	logger3 := NewErrorLogger(LogLevelError)

	if logger1 == nil || logger2 == nil || logger3 == nil {
		t.Fatal("Expected error loggers to be created")
	}

	// Test logging different types of errors
	err := New(ErrCodeInvalidVersion, "invalid version")
	context := map[string]interface{}{"test": "value"}

	// Test that logging functions don't panic
	logger1.LogError(err, context)
	logger1.LogErrorf(LogLevelInfo, "test error: %v", err)
	logger1.LogGopherError(err, context)

	// Test setting and getting log levels
	logger1.SetLevel(LogLevelDebug)
	if logger1.GetLevel() != LogLevelDebug {
		t.Errorf("Expected log level %v, got %v", LogLevelDebug, logger1.GetLevel())
	}
}
