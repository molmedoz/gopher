package errors

import (
	"errors"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	handler := NewErrorHandler(false)
	if handler == nil {
		t.Fatal("Expected error handler to be created")
	}

	// Test handling a regular error
	err := errors.New("test error")
	handled := handler.HandleError(err)
	if handled == "" {
		t.Error("Expected handled error to be returned")
	}

	// Test handling a GopherError
	gopherErr := New(ErrCodeInvalidVersion, "invalid version")
	handled = handler.HandleError(gopherErr)
	if handled == "" {
		t.Error("Expected handled GopherError to be returned")
	}
}

// TestErrorHandler_Comprehensive tests the error handler comprehensively
func TestErrorHandler_Comprehensive(t *testing.T) {
	// Test creating error handler with different options
	handler1 := NewErrorHandler(false)
	handler2 := NewErrorHandler(true)

	if handler1 == nil || handler2 == nil {
		t.Fatal("Expected error handlers to be created")
	}

	// Test handling different types of errors
	tests := []struct {
		name string
		err  error
	}{
		{"regular error", errors.New("regular error")},
		{"GopherError", New(ErrCodeInvalidVersion, "invalid version")},
		{"GopherError with context", New(ErrCodeInvalidVersion, "invalid version").WithContext("version", "1.2.3")},
		{"GopherError with details", New(ErrCodeInvalidVersion, "invalid version").WithDetails("additional info")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handled := handler1.HandleError(tt.err)
			if handled == "" {
				t.Error("Expected handled error to be returned")
			}
		})
	}
}
