package errors

import (
	"fmt"
	"os"
	"strings"
)

// ErrorHandler provides common error handling utilities
type ErrorHandler struct {
	// Add any configuration for error handling
	Verbose bool
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(verbose bool) *ErrorHandler {
	return &ErrorHandler{
		Verbose: verbose,
	}
}

// HandleError processes an error and returns a user-friendly message
func (h *ErrorHandler) HandleError(err error) string {
	if err == nil {
		return ""
	}

	// Check if it's already a GopherError
	if gopherErr, ok := err.(*GopherError); ok {
		return h.formatGopherError(gopherErr)
	}

	// Handle common system errors
	if os.IsNotExist(err) {
		return "File or directory not found"
	}
	if os.IsPermission(err) {
		return "Permission denied - you may need to run with elevated privileges"
	}

	// Return the original error message for unknown errors
	return err.Error()
}

// formatGopherError formats a GopherError for display
func (h *ErrorHandler) formatGopherError(err *GopherError) string {
	var parts []string

	// Add the main message
	parts = append(parts, err.Message)

	// Add details if available
	if err.Details != "" {
		parts = append(parts, fmt.Sprintf("Details: %s", err.Details))
	}

	// Add context if verbose mode is enabled
	if h.Verbose && len(err.Context) > 0 {
		contextParts := make([]string, 0, len(err.Context))
		for key, value := range err.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%v", key, value))
		}
		parts = append(parts, fmt.Sprintf("Context: %s", strings.Join(contextParts, ", ")))
	}

	// Add file and line if verbose mode is enabled
	if h.Verbose && err.File != "" {
		parts = append(parts, fmt.Sprintf("Location: %s:%d", err.File, err.Line))
	}

	return strings.Join(parts, "\n")
}

// ShouldRetry determines if an operation should be retried based on the error
func (h *ErrorHandler) ShouldRetry(err error) bool {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeNetworkUnavailable, ErrCodeTimeoutExceeded, ErrCodeServerUnavailable:
			return true
		default:
			return false
		}
	}
	return false
}

// GetRetryDelay returns the suggested delay before retrying an operation
func (h *ErrorHandler) GetRetryDelay(err error) int {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeNetworkUnavailable:
			return 5 // 5 seconds
		case ErrCodeTimeoutExceeded:
			return 10 // 10 seconds
		case ErrCodeServerUnavailable:
			return 30 // 30 seconds
		default:
			return 0
		}
	}
	return 0
}

// IsUserError determines if an error is caused by user input
func (h *ErrorHandler) IsUserError(err error) bool {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeInvalidVersion, ErrCodeInvalidArgument, ErrCodeInvalidFormat,
			ErrCodeMissingArgument, ErrCodeInvalidAliasName, ErrCodeReservedName,
			ErrCodeUnknownConfigOption, ErrCodeInvalidConfigValue:
			return true
		default:
			return false
		}
	}
	return false
}

// IsSystemError determines if an error is caused by system issues
func (h *ErrorHandler) IsSystemError(err error) bool {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeSystemGoNotAvailable, ErrCodeSymlinkFailed, ErrCodeEnvironmentSetupFailed,
			ErrCodeShellDetectionFailed, ErrCodePermissionDenied, ErrCodeDiskSpaceExhausted:
			return true
		default:
			return false
		}
	}
	return false
}

// IsNetworkError determines if an error is caused by network issues
func (h *ErrorHandler) IsNetworkError(err error) bool {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeNetworkUnavailable, ErrCodeTimeoutExceeded, ErrCodeServerUnavailable,
			ErrCodeDownloadFailed:
			return true
		default:
			return false
		}
	}
	return false
}

// GetErrorCategory returns a human-readable category for the error
func (h *ErrorHandler) GetErrorCategory(err error) string {
	if h.IsUserError(err) {
		return "User Input Error"
	}
	if h.IsSystemError(err) {
		return "System Error"
	}
	if h.IsNetworkError(err) {
		return "Network Error"
	}
	return "Unknown Error"
}

// SuggestSolution provides a suggested solution for common errors
func (h *ErrorHandler) SuggestSolution(err error) string {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeInvalidVersion:
			return "Please use a valid Go version format (e.g., '1.21.0' or 'go1.21.0')"
		case ErrCodeMissingArgument:
			return "Please provide the required arguments. Use 'gopher help' for usage information"
		case ErrCodeVersionNotInstalled:
			return "Use 'gopher list' to see installed versions, or 'gopher install <version>' to install a version"
		case ErrCodeVersionAlreadyInstalled:
			return "The version is already installed. Use 'gopher list' to see installed versions"
		case ErrCodeSystemGoNotAvailable:
			return "No system Go installation found. Install Go from https://golang.org/dl/ or use 'gopher install <version>'"
		case ErrCodePermissionDenied:
			return "Try running with elevated privileges (sudo on Unix, Run as Administrator on Windows)"
		case ErrCodeNetworkUnavailable:
			return "Check your internet connection and try again"
		case ErrCodeTimeoutExceeded:
			return "The operation timed out. Try again with a better internet connection"
		case ErrCodeSymlinkFailed:
			return "Symlink creation failed. You may need to enable Developer Mode on Windows or run with elevated privileges"
		case ErrCodeInvalidAliasName:
			return "Use only letters, numbers, hyphens, underscores, and dots. Avoid reserved names"
		case ErrCodeReservedName:
			return "Choose a different name that is not reserved by gopher"
		case ErrCodeUnknownConfigOption:
			return "Use 'gopher config list' to see available configuration options"
		default:
			return "Please check the error details and try again"
		}
	}
	return "Please check the error details and try again"
}
