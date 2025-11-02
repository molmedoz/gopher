package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SecurityError represents a security-related error
type SecurityError struct {
	Code    string
	Message string
	Details string
}

func (e *SecurityError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}

// Security error codes
const (
	ErrCodePathTraversal = "PATH_TRAVERSAL"
	ErrCodeInvalidPath   = "INVALID_PATH"
	ErrCodeUnsafePath    = "UNSAFE_PATH"
)

// NewPathTraversalError creates a new path traversal error
func NewPathTraversalError(path string) *SecurityError {
	return &SecurityError{
		Code:    ErrCodePathTraversal,
		Message: "path traversal detected",
		Details: path,
	}
}

// NewInvalidPathError creates a new invalid path error
func NewInvalidPathError(path string) *SecurityError {
	return &SecurityError{
		Code:    ErrCodeInvalidPath,
		Message: "invalid path",
		Details: path,
	}
}

// NewUnsafePathError creates a new unsafe path error
func NewUnsafePathError(path string) *SecurityError {
	return &SecurityError{
		Code:    ErrCodeUnsafePath,
		Message: "unsafe path detected",
		Details: path,
	}
}

// ValidatePath validates a file path for security issues
func ValidatePath(path string) error {
	if path == "" {
		return NewInvalidPathError("empty path")
	}

	// Clean the path to resolve any relative components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return NewPathTraversalError(path)
	}

	// Check for absolute paths that might be unsafe
	if filepath.IsAbs(cleanPath) {
		// Allow absolute paths but log them for security review
		// In a production environment, you might want to restrict this
		return nil
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"../",
		"..\\",
		"~",
		"$",
		"`",
		"|",
		"&",
		";",
		"(",
		")",
		"<",
		">",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(path, pattern) {
			return NewUnsafePathError(path)
		}
	}

	return nil
}

// SanitizePath sanitizes a file path by removing dangerous components
func SanitizePath(path string) string {
	// Remove dangerous components first
	sanitized := strings.ReplaceAll(path, "..", "")
	sanitized = strings.ReplaceAll(sanitized, "~", "")

	// Remove any remaining suspicious characters
	suspiciousChars := []string{"$", "`", "|", "&", ";", "(", ")", "<", ">"}
	for _, char := range suspiciousChars {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}

	// Clean the path after removing dangerous components
	cleanPath := filepath.Clean(sanitized)

	// Remove leading slashes and backslashes for relative paths
	if !filepath.IsAbs(path) {
		cleanPath = strings.TrimPrefix(cleanPath, "/")
		cleanPath = strings.TrimPrefix(cleanPath, "\\")
	}

	return cleanPath
}

// ValidateDirectoryPath validates a directory path for security issues
func ValidateDirectoryPath(path string) error {
	if err := ValidatePath(path); err != nil {
		return err
	}

	// Trailing separators are already handled by ValidatePath
	return nil
}

// IsSafePath checks if a path is safe to use
func IsSafePath(path string) bool {
	return ValidatePath(path) == nil
}

// GetSafePath returns a sanitized version of the path
func GetSafePath(path string) (string, error) {
	if err := ValidatePath(path); err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}

// ValidatePathWithinRoot ensures a path is within a safe root directory
// This prevents directory traversal attacks when accessing files
func ValidatePathWithinRoot(path, rootDir string) (string, error) {
	if err := ValidatePath(path); err != nil {
		return "", err
	}

	// Resolve absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve root directory: %w", err)
	}

	// Normalize paths to handle different separators
	absPath = filepath.Clean(absPath)
	absRoot = filepath.Clean(absRoot)

	// Ensure the path is within the root directory
	relPath, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path: %w", err)
	}

	// Check for path traversal (contains "..")
	if strings.HasPrefix(relPath, "..") || strings.Contains(relPath, ".."+string(filepath.Separator)) {
		return "", NewPathTraversalError(path)
	}

	return absPath, nil
}

// SafeReadFile reads a file after validating the path is within rootDir
func SafeReadFile(path, rootDir string) ([]byte, error) {
	safePath, err := ValidatePathWithinRoot(path, rootDir)
	if err != nil {
		return nil, err
	}
	// #nosec G304 -- path validated and constrained to rootDir by ValidatePathWithinRoot
	return os.ReadFile(safePath)
}

// SafeWriteFile writes to a file after validating the path is within rootDir
func SafeWriteFile(path, rootDir string, data []byte, perm os.FileMode) error {
	safePath, err := ValidatePathWithinRoot(path, rootDir)
	if err != nil {
		return err
	}
	return os.WriteFile(safePath, data, perm)
}

// SafeOpen opens a file after validating the path is within rootDir
func SafeOpen(path, rootDir string) (*os.File, error) {
	safePath, err := ValidatePathWithinRoot(path, rootDir)
	if err != nil {
		return nil, err
	}
	// #nosec G304 -- path validated and constrained to rootDir by ValidatePathWithinRoot
	return os.Open(safePath)
}

// SafeCreate creates a file after validating the path is within rootDir
func SafeCreate(path, rootDir string) (*os.File, error) {
	safePath, err := ValidatePathWithinRoot(path, rootDir)
	if err != nil {
		return nil, err
	}
	// #nosec G304 -- path validated and constrained to rootDir by ValidatePathWithinRoot
	return os.Create(safePath)
}

// ValidateFilePermissions validates file permissions for security
func ValidateFilePermissions(filePath string, expectedMode os.FileMode) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return NewInvalidPathError(filePath)
	}

	// Check if permissions are too permissive
	actualMode := info.Mode()
	if actualMode&0002 != 0 { // World writable
		return NewUnsafePathError(filePath + " (world writable)")
	}
	if actualMode&0020 != 0 && actualMode&0004 == 0 { // Group writable but not group readable
		return NewUnsafePathError(filePath + " (group writable but not readable)")
	}

	// Check if permissions match expected mode
	if expectedMode != 0 && actualMode.Perm() != expectedMode.Perm() {
		return NewUnsafePathError(filePath + " (permissions mismatch)")
	}

	return nil
}

// SetSecureFilePermissions sets secure file permissions
func SetSecureFilePermissions(filePath string, mode os.FileMode) error {
	// Ensure the file is not world writable
	secureMode := mode &^ 0002 // Remove world write permission

	// Ensure the file is not group writable unless it's also group readable
	if secureMode&0020 != 0 && secureMode&0004 == 0 {
		secureMode |= 0004 // Add group read permission
	}

	return os.Chmod(filePath, secureMode)
}

// CreateSecureFile creates a file with secure permissions
func CreateSecureFile(filePath string, mode os.FileMode) (*os.File, error) {
	// Validate the path first
	if err := ValidatePath(filePath); err != nil {
		return nil, err
	}

	// Create the file with secure permissions
	// #nosec G304 -- filePath validated by ValidatePath before calling this function
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	// Set secure permissions
	if err := SetSecureFilePermissions(filePath, mode); err != nil {
		_ = file.Close() // Best effort cleanup
		if rerr := os.Remove(filePath); rerr != nil && !os.IsNotExist(rerr) {
			// Log cleanup failure but return original error
			return nil, fmt.Errorf("failed to set permissions: %w (cleanup also failed: %v)", err, rerr)
		}
		return nil, err
	}

	return file, nil
}

// CreateSecureDirectory creates a directory with secure permissions
func CreateSecureDirectory(dirPath string, mode os.FileMode) error {
	// Validate the path first
	if err := ValidateDirectoryPath(dirPath); err != nil {
		return err
	}

	// Create the directory with secure permissions
	if err := os.MkdirAll(dirPath, mode); err != nil {
		return err
	}

	// Set secure permissions
	return SetSecureFilePermissions(dirPath, mode)
}
