package errors

import (
	"errors"
	"testing"
)

func TestRecoverer(t *testing.T) {
	logger := NewErrorLogger(LogLevelInfo)
	recoverer := NewRecoverer(logger)
	if recoverer == nil {
		t.Fatal("Expected recoverer to be created")
	}

	// Test recovering from panic
	func() {
		defer recoverer.Recover()
		panic("test panic")
	}()

	// Test that recoverer was created successfully
	t.Log("Recoverer created and basic panic recovery tested")
}

// TestRecoverer_Comprehensive tests the recoverer comprehensively
func TestRecoverer_Comprehensive(t *testing.T) {
	logger := NewErrorLogger(LogLevelInfo)
	recoverer := NewRecoverer(logger)

	if recoverer == nil {
		t.Fatal("Expected recoverer to be created")
	}

	// Test recovering from different types of panics
	tests := []struct {
		name  string
		panic interface{}
	}{
		{"string panic", "test panic"},
		{"error panic", errors.New("test error panic")},
		{"nil panic", nil},
		{"number panic", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer recoverer.Recover()
				panic(tt.panic)
			}()
		})
	}

	// Test that recoverer was created successfully
	t.Log("Recoverer comprehensive test completed")
}
