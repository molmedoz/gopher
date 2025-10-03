package errors

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator provides common validation functions
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateVersion validates a Go version string
func (v *Validator) ValidateVersion(version string) error {
	if version == "" {
		return New(ErrCodeInvalidVersion, "version cannot be empty")
	}

	// Remove 'go' prefix if present
	version = strings.TrimPrefix(version, "go")

	// Basic version format validation (semantic versioning)
	versionRegex := regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?(?:-([a-zA-Z0-9\-]+))?(?:\+([a-zA-Z0-9\-]+))?$`)
	if !versionRegex.MatchString(version) {
		return NewInvalidVersion(version)
	}

	// Additional validation for Go versions
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return NewInvalidVersion(version)
	}

	// Check for valid major.minor format
	if parts[0] == "0" && parts[1] == "0" {
		return NewInvalidVersion(version)
	}

	return nil
}

// ValidateAliasName validates an alias name
func (v *Validator) ValidateAliasName(name string) error {
	if name == "" {
		return New(ErrCodeInvalidAliasName, "alias name cannot be empty")
	}

	if len(name) > 50 {
		return New(ErrCodeInvalidAliasName, "alias name cannot exceed 50 characters")
	}

	// Check for reserved names
	reservedNames := map[string]bool{
		"system":      true,
		"sys":         true,
		"install":     true,
		"uninstall":   true,
		"use":         true,
		"list":        true,
		"list-remote": true,
		"alias":       true,
		"init":        true,
		"setup":       true,
		"status":      true,
		"debug":       true,
		"help":        true,
		"version":     true,
		"config":      true,
		"env":         true,
		"current":     true,
		"switch":      true,
		"remove":      true,
		"delete":      true,
		"add":         true,
		"create":      true,
		"update":      true,
		"export":      true,
		"import":      true,
		"bulk":        true,
	}

	if reservedNames[strings.ToLower(name)] {
		return NewReservedName(name)
	}

	// Check for valid characters (letters, numbers, hyphens, underscores, dots)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	if !validNameRegex.MatchString(name) {
		return New(ErrCodeInvalidAliasName, "alias name contains invalid characters. Only letters, numbers, hyphens, underscores, and dots are allowed")
	}

	// Check that name doesn't start or end with special characters
	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
		return New(ErrCodeInvalidAliasName, "alias name cannot start with special characters")
	}

	if strings.HasSuffix(name, "-") || strings.HasSuffix(name, "_") || strings.HasSuffix(name, ".") {
		return New(ErrCodeInvalidAliasName, "alias name cannot end with special characters")
	}

	return nil
}

// ValidateConfigValue validates a configuration value
func (v *Validator) ValidateConfigValue(key, value string) error {
	switch key {
	case "gopath_mode":
		validModes := []string{"shared", "version-specific", "custom"}
		for _, mode := range validModes {
			if value == mode {
				return nil
			}
		}
		return New(ErrCodeInvalidConfigValue, fmt.Sprintf("gopath_mode must be one of: %s", strings.Join(validModes, ", ")))

	case "set_environment":
		if value != "true" && value != "false" {
			return New(ErrCodeInvalidConfigValue, "set_environment must be 'true' or 'false'")
		}
		return nil

	case "auto_cleanup":
		if value != "true" && value != "false" {
			return New(ErrCodeInvalidConfigValue, "auto_cleanup must be 'true' or 'false'")
		}
		return nil

	case "max_versions":
		// This would need to be parsed as an integer, but we'll do basic validation here
		if value == "" {
			return New(ErrCodeInvalidConfigValue, "max_versions cannot be empty")
		}
		return nil

	case "mirror_url":
		if value == "" {
			return New(ErrCodeInvalidConfigValue, "mirror_url cannot be empty")
		}
		// Basic URL validation
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return New(ErrCodeInvalidConfigValue, "mirror_url must be a valid HTTP/HTTPS URL")
		}
		return nil

	case "custom_gopath":
		if value == "" {
			return New(ErrCodeInvalidConfigValue, "custom_gopath cannot be empty when gopath_mode is 'custom'")
		}
		return nil

	default:
		return NewUnknownConfigOption(key)
	}
}

// ValidatePath validates a file or directory path
func (v *Validator) ValidatePath(path string) error {
	if path == "" {
		return New(ErrCodeInvalidArgument, "path cannot be empty")
	}

	// Check for invalid characters in path (but allow : for Windows drive letters)
	invalidChars := []string{"<", ">", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return New(ErrCodeInvalidArgument, fmt.Sprintf("path contains invalid character: %s", char))
		}
	}

	// Check for : but allow it only for Windows drive letters (C:, D:, etc.)
	if strings.Contains(path, ":") {
		// Allow : only if it's followed by \ or / (Windows drive letter)
		colonIndex := strings.Index(path, ":")
		if colonIndex == -1 || colonIndex == len(path)-1 ||
			(path[colonIndex+1] != '\\' && path[colonIndex+1] != '/') {
			return New(ErrCodeInvalidArgument, "path contains invalid character: :")
		}
	}

	return nil
}

// ValidateShell validates a shell name
func (v *Validator) ValidateShell(shell string) error {
	validShells := []string{"bash", "zsh", "fish", "powershell", "cmd"}
	for _, validShell := range validShells {
		if shell == validShell {
			return nil
		}
	}
	return New(ErrCodeInvalidArgument, fmt.Sprintf("unsupported shell: %s. Supported shells: %s", shell, strings.Join(validShells, ", ")))
}

// ValidateKeyValuePair validates a key=value pair
func (v *Validator) ValidateKeyValuePair(keyValue string) error {
	if keyValue == "" {
		return New(ErrCodeInvalidFormat, "key=value pair cannot be empty")
	}

	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return NewInvalidFormat("key=value")
	}

	key, value := parts[0], parts[1]
	if key == "" {
		return New(ErrCodeInvalidFormat, "key cannot be empty")
	}

	if value == "" {
		return New(ErrCodeInvalidFormat, "value cannot be empty")
	}

	return nil
}

// ValidateCommand validates a command and its arguments
func (v *Validator) ValidateCommand(command string, args []string) error {
	if command == "" {
		return New(ErrCodeInvalidArgument, "command cannot be empty")
	}

	// Validate commands that require arguments
	commandsRequiringArgs := map[string]string{
		"install":   "version",
		"uninstall": "version",
		"use":       "version or alias",
		"show":      "version",
		"set":       "key=value",
		"alias":     "subcommand",
		"env":       "subcommand",
	}

	if requiredArg, exists := commandsRequiringArgs[command]; exists {
		if len(args) == 0 {
			return NewMissingArgument(command + " (requires " + requiredArg + ")")
		}
	}

	return nil
}

// Global validator instance
var (
	DefaultValidator = NewValidator()
)

// Convenience functions using the default validator

// ValidateVersion is a convenience function
func ValidateVersion(version string) error {
	return DefaultValidator.ValidateVersion(version)
}

// ValidateAliasName is a convenience function
func ValidateAliasName(name string) error {
	return DefaultValidator.ValidateAliasName(name)
}

// ValidateConfigValue is a convenience function
func ValidateConfigValue(key, value string) error {
	return DefaultValidator.ValidateConfigValue(key, value)
}

// ValidatePath is a convenience function
func ValidatePath(path string) error {
	return DefaultValidator.ValidatePath(path)
}

// ValidateShell is a convenience function
func ValidateShell(shell string) error {
	return DefaultValidator.ValidateShell(shell)
}

// ValidateKeyValuePair is a convenience function
func ValidateKeyValuePair(keyValue string) error {
	return DefaultValidator.ValidateKeyValuePair(keyValue)
}

// ValidateCommand is a convenience function
func ValidateCommand(command string, args []string) error {
	return DefaultValidator.ValidateCommand(command, args)
}
