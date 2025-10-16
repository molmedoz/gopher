package errors

import (
	"errors"
	"testing"
)

func TestGopherError(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		message  string
		expected string
	}{
		{
			name:     "basic error",
			code:     ErrCodeInvalidVersion,
			message:  "invalid version",
			expected: "INVALID_VERSION: invalid version",
		},
		{
			name:     "error with wrapped error",
			code:     ErrCodeInstallationFailed,
			message:  "installation failed",
			expected: "INSTALLATION_FAILED: installation failed: wrapped error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message)
			if tt.name == "error with wrapped error" {
				wrappedErr := errors.New("wrapped error")
				err = Wrap(wrappedErr, tt.code, tt.message)
			}

			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}

			if err.Code != tt.code {
				t.Errorf("Code = %v, want %v", err.Code, tt.code)
			}

			if err.Message != tt.message {
				t.Errorf("Message = %v, want %v", err.Message, tt.message)
			}
		})
	}
}

func TestErrorHelpers(t *testing.T) {
	t.Run("IsGopherError", func(t *testing.T) {
		gopherErr := New(ErrCodeInvalidVersion, "test")
		regularErr := errors.New("regular error")

		if !IsGopherError(gopherErr) {
			t.Error("Expected gopherErr to be identified as GopherError")
		}

		if IsGopherError(regularErr) {
			t.Error("Expected regularErr to not be identified as GopherError")
		}
	})

	t.Run("GetErrorCode", func(t *testing.T) {
		gopherErr := New(ErrCodeInvalidVersion, "test")
		regularErr := errors.New("regular error")

		if GetErrorCode(gopherErr) != ErrCodeInvalidVersion {
			t.Error("Expected correct error code for GopherError")
		}

		if GetErrorCode(regularErr) != ErrCodeUnknown {
			t.Error("Expected ErrCodeUnknown for regular error")
		}
	})

	t.Run("IsErrorCode", func(t *testing.T) {
		gopherErr := New(ErrCodeInvalidVersion, "test")
		regularErr := errors.New("regular error")

		if !IsErrorCode(gopherErr, ErrCodeInvalidVersion) {
			t.Error("Expected gopherErr to match ErrCodeInvalidVersion")
		}

		if IsErrorCode(regularErr, ErrCodeInvalidVersion) {
			t.Error("Expected regularErr to not match ErrCodeInvalidVersion")
		}
	})
}

func TestErrorConstructors(t *testing.T) {
	t.Run("NewInvalidVersion", func(t *testing.T) {
		err := NewInvalidVersion("1.2.3")
		expected := "INVALID_VERSION: invalid version format: 1.2.3"
		if err.Error() != expected {
			t.Errorf("NewInvalidVersion() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("NewMissingArgument", func(t *testing.T) {
		err := NewMissingArgument("install")
		expected := "MISSING_ARGUMENT: install command requires additional arguments"
		if err.Error() != expected {
			t.Errorf("NewMissingArgument() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("NewVersionNotInstalled", func(t *testing.T) {
		err := NewVersionNotInstalled("1.21.0")
		expected := "VERSION_NOT_INSTALLED: version 1.21.0 is not installed"
		if err.Error() != expected {
			t.Errorf("NewVersionNotInstalled() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("NewSystemGoNotAvailable", func(t *testing.T) {
		err := NewSystemGoNotAvailable()
		expected := "SYSTEM_GO_NOT_AVAILABLE: system Go is not available"
		if err.Error() != expected {
			t.Errorf("NewSystemGoNotAvailable() = %v, want %v", err.Error(), expected)
		}
	})
}

func TestErrorWithContext(t *testing.T) {
	err := New(ErrCodeInvalidVersion, "test error")
	err = err.WithContext("version", "1.2.3")
	err = err.WithContext("command", "install")
	err = err.WithDetails("Additional details here")

	if err.Context["version"] != "1.2.3" {
		t.Error("Expected version context to be set")
	}

	if err.Context["command"] != "install" {
		t.Error("Expected command context to be set")
	}

	if err.Details != "Additional details here" {
		t.Error("Expected details to be set")
	}
}

func TestErrorUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, ErrCodeInstallationFailed, "wrapped message")

	if wrappedErr.Unwrap() != originalErr {
		t.Error("Expected Unwrap to return the original error")
	}
}
