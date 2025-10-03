// Package errors provides standardized error handling for Gopher.
//
// This package implements a comprehensive error system with structured error
// types, error codes, and helper functions for consistent error handling
// throughout the application.
//
// Features:
//   - Structured error types with error codes
//   - Error wrapping and context preservation
//   - Validation helpers for common error types
//   - Recovery mechanisms for panic handling
//   - Logging integration for error reporting
//
// Error Categories:
//   - Validation errors: Invalid input, missing arguments, format errors
//   - Installation errors: Download, installation, and uninstallation failures
//   - System errors: System Go detection, symlink, and environment issues
//   - Configuration errors: Config loading, saving, and validation failures
//   - Alias errors: Alias management and validation failures
//   - Security errors: Path traversal and security validation failures
//
// Usage:
//
//	// Create a new error
//	err := errors.New(errors.ErrCodeInvalidVersion, "invalid version format")
//
//	// Wrap an existing error
//	err = errors.Wrapf(err, errors.ErrCodeInstallationFailed, "failed to install %s", version)
//
//	// Check error type
//	if errors.IsErrorCode(err, errors.ErrCodeInvalidVersion) {
//	    // Handle invalid version error
//	}
//
// For more details, see the individual function documentation.
package errors

import (
	"fmt"
	"runtime"
)

// ErrorCode represents a specific error type
type ErrorCode string

const (
	// Validation errors
	ErrCodeInvalidVersion   ErrorCode = "INVALID_VERSION"
	ErrCodeInvalidArgument  ErrorCode = "INVALID_ARGUMENT"
	ErrCodeInvalidFormat    ErrorCode = "INVALID_FORMAT"
	ErrCodeMissingArgument  ErrorCode = "MISSING_ARGUMENT"
	ErrCodeInvalidAliasName ErrorCode = "INVALID_ALIAS_NAME"
	ErrCodeReservedName     ErrorCode = "RESERVED_NAME"

	// Installation errors
	ErrCodeVersionNotInstalled     ErrorCode = "VERSION_NOT_INSTALLED"
	ErrCodeVersionAlreadyInstalled ErrorCode = "VERSION_ALREADY_INSTALLED"
	ErrCodeInstallationFailed      ErrorCode = "INSTALLATION_FAILED"
	ErrCodeUninstallationFailed    ErrorCode = "UNINSTALLATION_FAILED"
	ErrCodeDownloadFailed          ErrorCode = "DOWNLOAD_FAILED"
	ErrCodeExtractionFailed        ErrorCode = "EXTRACTION_FAILED"

	// System errors
	ErrCodeSystemGoNotAvailable   ErrorCode = "SYSTEM_GO_NOT_AVAILABLE"
	ErrCodeSymlinkFailed          ErrorCode = "SYMLINK_FAILED"
	ErrCodeEnvironmentSetupFailed ErrorCode = "ENVIRONMENT_SETUP_FAILED"
	ErrCodeShellDetectionFailed   ErrorCode = "SHELL_DETECTION_FAILED"

	// Configuration errors
	ErrCodeConfigLoadFailed    ErrorCode = "CONFIG_LOAD_FAILED"
	ErrCodeConfigSaveFailed    ErrorCode = "CONFIG_SAVE_FAILED"
	ErrCodeInvalidConfigValue  ErrorCode = "INVALID_CONFIG_VALUE"
	ErrCodeUnknownConfigOption ErrorCode = "UNKNOWN_CONFIG_OPTION"

	// Alias errors
	ErrCodeAliasAlreadyExists ErrorCode = "ALIAS_ALREADY_EXISTS"
	ErrCodeAliasNotFound      ErrorCode = "ALIAS_NOT_FOUND"
	ErrCodeAliasLoadFailed    ErrorCode = "ALIAS_LOAD_FAILED"
	ErrCodeAliasSaveFailed    ErrorCode = "ALIAS_SAVE_FAILED"
	ErrCodeAliasUpdateFailed  ErrorCode = "ALIAS_UPDATE_FAILED"
	ErrCodeAliasRemoveFailed  ErrorCode = "ALIAS_REMOVE_FAILED"

	// File system errors
	ErrCodeFileNotFound       ErrorCode = "FILE_NOT_FOUND"
	ErrCodeDirectoryNotFound  ErrorCode = "DIRECTORY_NOT_FOUND"
	ErrCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrCodeDiskSpaceExhausted ErrorCode = "DISK_SPACE_EXHAUSTED"

	// Network errors
	ErrCodeNetworkUnavailable ErrorCode = "NETWORK_UNAVAILABLE"
	ErrCodeTimeoutExceeded    ErrorCode = "TIMEOUT_EXCEEDED"
	ErrCodeServerUnavailable  ErrorCode = "SERVER_UNAVAILABLE"

	// Generic errors
	ErrCodeUnknown            ErrorCode = "UNKNOWN_ERROR"
	ErrCodeNotImplemented     ErrorCode = "NOT_IMPLEMENTED"
	ErrCodeOperationCancelled ErrorCode = "OPERATION_CANCELLED"
)

// GopherError represents a structured error with context
type GopherError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
	WrappedErr error                  `json:"-"`
	File       string                 `json:"file,omitempty"`
	Line       int                    `json:"line,omitempty"`
}

// Error implements the error interface
func (e *GopherError) Error() string {
	if e.WrappedErr != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.WrappedErr)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error for error chain inspection
func (e *GopherError) Unwrap() error {
	return e.WrappedErr
}

// WithContext adds context information to the error
func (e *GopherError) WithContext(key string, value interface{}) *GopherError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails adds detailed information to the error
func (e *GopherError) WithDetails(details string) *GopherError {
	e.Details = details
	return e
}

// New creates a new GopherError
func New(code ErrorCode, message string) *GopherError {
	_, file, line, _ := runtime.Caller(1)
	return &GopherError{
		Code:    code,
		Message: message,
		File:    file,
		Line:    line,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *GopherError {
	_, file, line, _ := runtime.Caller(1)
	return &GopherError{
		Code:       code,
		Message:    message,
		WrappedErr: err,
		File:       file,
		Line:       line,
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *GopherError {
	_, file, line, _ := runtime.Caller(1)
	return &GopherError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		WrappedErr: err,
		File:       file,
		Line:       line,
	}
}

// Newf creates a new GopherError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *GopherError {
	_, file, line, _ := runtime.Caller(1)
	return &GopherError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		File:    file,
		Line:    line,
	}
}

// IsGopherError checks if an error is a GopherError
func IsGopherError(err error) bool {
	_, ok := err.(*GopherError)
	return ok
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if gopherErr, ok := err.(*GopherError); ok {
		return gopherErr.Code
	}
	return ErrCodeUnknown
}

// IsErrorCode checks if an error has a specific code
func IsErrorCode(err error, code ErrorCode) bool {
	return GetErrorCode(err) == code
}

// Common error constructors for frequently used errors

// Validation errors
func NewInvalidVersion(version string) *GopherError {
	return Newf(ErrCodeInvalidVersion, "invalid version format: %s", version)
}

func NewMissingArgument(command string) *GopherError {
	return Newf(ErrCodeMissingArgument, "%s command requires additional arguments", command)
}

func NewInvalidFormat(expected string) *GopherError {
	return Newf(ErrCodeInvalidFormat, "invalid format: expected %s", expected)
}

func NewInvalidAliasName(name string) *GopherError {
	return Newf(ErrCodeInvalidAliasName, "invalid alias name: %s", name)
}

func NewReservedName(name string) *GopherError {
	return Newf(ErrCodeReservedName, "name '%s' is reserved", name)
}

// Installation errors
func NewVersionNotInstalled(version string) *GopherError {
	return Newf(ErrCodeVersionNotInstalled, "version %s is not installed", version)
}

func NewVersionAlreadyInstalled(version string) *GopherError {
	return Newf(ErrCodeVersionAlreadyInstalled, "version %s is already installed", version)
}

func NewInstallationFailed(version string, err error) *GopherError {
	return Wrapf(err, ErrCodeInstallationFailed, "failed to install version %s", version)
}

func NewDownloadFailed(version string, err error) *GopherError {
	return Wrapf(err, ErrCodeDownloadFailed, "failed to download version %s", version)
}

// System errors
func NewSystemGoNotAvailable() *GopherError {
	return New(ErrCodeSystemGoNotAvailable, "system Go is not available")
}

func NewSymlinkFailed(target, link string, err error) *GopherError {
	return Wrapf(err, ErrCodeSymlinkFailed, "failed to create symlink from %s to %s", target, link)
}

// Configuration errors
func NewConfigLoadFailed(path string, err error) *GopherError {
	return Wrapf(err, ErrCodeConfigLoadFailed, "failed to load configuration from %s", path)
}

func NewConfigSaveFailed(path string, err error) *GopherError {
	return Wrapf(err, ErrCodeConfigSaveFailed, "failed to save configuration to %s", path)
}

func NewUnknownConfigOption(option string) *GopherError {
	return Newf(ErrCodeUnknownConfigOption, "unknown configuration option: %s", option)
}

// File system errors
func NewFileNotFound(path string) *GopherError {
	return Newf(ErrCodeFileNotFound, "file not found: %s", path)
}

func NewDirectoryNotFound(path string) *GopherError {
	return Newf(ErrCodeDirectoryNotFound, "directory not found: %s", path)
}

func NewAliasNotFound(name string) *GopherError {
	return Newf(ErrCodeAliasNotFound, "alias not found: %s", name)
}

func NewPermissionDenied(path string) *GopherError {
	return Newf(ErrCodePermissionDenied, "permission denied: %s", path)
}

// Network errors
func NewNetworkUnavailable(err error) *GopherError {
	return Wrap(err, ErrCodeNetworkUnavailable, "network unavailable")
}

func NewTimeoutExceeded(operation string) *GopherError {
	return Newf(ErrCodeTimeoutExceeded, "timeout exceeded for operation: %s", operation)
}

// Generic errors
func NewNotImplemented(feature string) *GopherError {
	return Newf(ErrCodeNotImplemented, "feature not implemented: %s", feature)
}

func NewOperationCancelled(operation string) *GopherError {
	return Newf(ErrCodeOperationCancelled, "operation cancelled: %s", operation)
}
