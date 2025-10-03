package runtime

import (
	"encoding/json"
	"fmt"
	"github.com/molmedoz/gopher/internal/errors"
	"github.com/molmedoz/gopher/internal/security"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// Alias Management - Core Operations
// ============================================================================

// LoadAliases loads aliases from the aliases file
func (am *AliasManager) LoadAliases() error {
	// Check if aliases file exists
	if _, err := os.Stat(am.aliasesFile); os.IsNotExist(err) {
		// File doesn't exist, create empty aliases map
		am.aliases = make(map[string]*Alias)
		return nil
	}

	// Read aliases file
	data, err := os.ReadFile(am.aliasesFile)
	if err != nil {
		return fmt.Errorf("failed to read aliases file: %w", err)
	}

	// Parse JSON
	var aliases map[string]*Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		return fmt.Errorf("failed to parse aliases file: %w", err)
	}

	am.aliases = aliases
	return nil
}

// SaveAliases saves aliases to the aliases file
func (am *AliasManager) SaveAliases() error {
	// Ensure directory exists
	dir := filepath.Dir(am.aliasesFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create aliases directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(am.aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	// Write to file
	if err := os.WriteFile(am.aliasesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write aliases file: %w", err)
	}

	return nil
}

// ValidateAliasName validates an alias name
func (am *AliasManager) ValidateAliasName(name string) error {
	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	// Check for reserved names
	reservedNames := []string{"system", "sys", "current", "list", "install", "uninstall", "use", "version", "help", "init", "setup", "status", "debug", "env", "alias", "go", "git", "ls", "cd", "pwd"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) {
			return fmt.Errorf("'%s' is a reserved name and cannot be used as an alias", name)
		}
	}

	// Check for valid characters (alphanumeric, hyphens, underscores, dots)
	matched, err := regexp.MatchString("^[a-zA-Z0-9._-]+$", name)
	if err != nil {
		return fmt.Errorf("failed to validate alias name: %w", err)
	}
	if !matched {
		return fmt.Errorf("alias name contains invalid characters")
	}

	// Check length
	if len(name) > 50 {
		return fmt.Errorf("alias name is too long")
	}

	return nil
}

// isVersionInstalled checks if a version is installed
func (am *AliasManager) isVersionInstalled(version string) bool {
	if am.manager == nil {
		return true // Assume installed if manager reference is nil
	}

	installedVersions, err := am.manager.ListInstalled()
	if err != nil {
		// Fallback to directory check if ListInstalled fails
		versionDir := filepath.Join(am.manager.config.InstallDir, version)
		if _, err := os.Stat(versionDir); os.IsNotExist(err) {
			return false
		}
		goBinary := filepath.Join(versionDir, "bin", "go")
		if _, err := os.Stat(goBinary); os.IsNotExist(err) {
			return false
		}
		return true
	}

	for _, v := range installedVersions {
		if v.Version == version {
			return true
		}
	}

	return false
}

// CreateAlias creates a new alias
func (am *AliasManager) CreateAlias(name, version string) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeAliasLoadFailed, "failed to load aliases")
	}

	// Validate alias name
	if err := errors.ValidateAliasName(name); err != nil {
		return err
	}

	// Validate alias name for security (path traversal protection)
	if err := security.ValidatePath(name); err != nil {
		return errors.Newf(errors.ErrCodeInvalidAliasName, "invalid alias name: %v", err)
	}

	// Validate version format
	if err := errors.ValidateVersion(version); err != nil {
		return err
	}

	// Validate version for security (path traversal protection)
	if err := security.ValidatePath(version); err != nil {
		return errors.Newf(errors.ErrCodeInvalidVersion, "invalid version: %v", err)
	}

	// Check if alias already exists
	if _, exists := am.aliases[name]; exists {
		return errors.Newf(errors.ErrCodeAliasAlreadyExists, "alias '%s' already exists (use 'gopher alias remove %s' first)", name, name)
	}

	// Check if version is installed
	if !am.isVersionInstalled(version) {
		return errors.Newf(errors.ErrCodeVersionNotInstalled, "version %s is not installed (use 'gopher install %s' first)", version, version)
	}

	// Create new alias
	alias := &Alias{
		Name:    name,
		Version: NormalizeVersion(version),
		Created: time.Now(),
		Updated: time.Now(),
	}

	am.aliases[name] = alias

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeAliasSaveFailed, "failed to save aliases")
	}

	return nil
}

// GetAlias gets an alias by name
func (am *AliasManager) GetAlias(name string) (*Alias, bool) {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return nil, false
	}

	alias, exists := am.aliases[name]
	return alias, exists
}

// ListAliases returns all aliases
func (am *AliasManager) ListAliases() ([]*Alias, error) {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return nil, errors.Wrapf(err, errors.ErrCodeAliasLoadFailed, "failed to load aliases")
	}

	var result []*Alias
	for _, alias := range am.aliases {
		result = append(result, alias)
	}

	return result, nil
}

// RemoveAlias removes an alias
func (am *AliasManager) RemoveAlias(name string) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeAliasLoadFailed, "failed to load aliases")
	}

	// Validate alias name for security (path traversal protection)
	if err := security.ValidatePath(name); err != nil {
		return errors.Newf(errors.ErrCodeInvalidAliasName, "invalid alias name: %v", err)
	}

	// Check if alias exists
	if _, exists := am.aliases[name]; !exists {
		return errors.Newf(errors.ErrCodeAliasNotFound, "alias '%s' does not exist", name)
	}

	// Remove alias
	delete(am.aliases, name)

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeAliasSaveFailed, "failed to save aliases")
	}

	return nil
}

// UpdateAlias updates an existing alias
func (am *AliasManager) UpdateAlias(name, version string) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeAliasLoadFailed, "failed to load aliases")
	}

	// Validate alias name for security (path traversal protection)
	if err := security.ValidatePath(name); err != nil {
		return errors.Newf(errors.ErrCodeInvalidAliasName, "invalid alias name: %v", err)
	}

	// Validate version for security (path traversal protection)
	if err := security.ValidatePath(version); err != nil {
		return errors.Newf(errors.ErrCodeInvalidVersion, "invalid version: %v", err)
	}

	// Check if alias exists
	alias, exists := am.aliases[name]
	if !exists {
		return errors.Newf(errors.ErrCodeAliasNotFound, "alias '%s' does not exist", name)
	}

	// Check if version is installed
	if !am.isVersionInstalled(version) {
		return errors.Newf(errors.ErrCodeVersionNotInstalled, "version %s is not installed (use 'gopher install %s' first)", version, version)
	}

	// Update alias
	alias.Version = NormalizeVersion(version)
	alias.Updated = time.Now()

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeAliasSaveFailed, "failed to save aliases")
	}

	return nil
}

// GetAliasesByVersion returns all aliases pointing to a specific version
func (am *AliasManager) GetAliasesByVersion(version string) ([]*Alias, error) {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return nil, errors.Wrapf(err, errors.ErrCodeAliasLoadFailed, "failed to load aliases")
	}

	normalizedVersion := NormalizeVersion(version)
	var result []*Alias
	for _, alias := range am.aliases {
		if alias.Version == normalizedVersion {
			result = append(result, alias)
		}
	}

	return result, nil
}

// Standard alias patterns for suggestions
var standardAliasPatterns = []string{
	"stable", "latest", "dev", "development", "prod", "production",
	"test", "testing", "staging", "preview", "beta", "alpha", "rc",
	"lts", "current", "main", "master", "release", "nightly",
}

// Standard alias groups for suggestions
var standardAliasGroups = []string{
	"production", "development", "testing", "staging", "personal",
}

// SuggestAliases suggests common alias names for a given version
func (am *AliasManager) SuggestAliases(version string) []string {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return []string{}
	}

	var suggestions []string
	usedNames := make(map[string]bool)

	// Mark existing aliases as used
	for name := range am.aliases {
		usedNames[strings.ToLower(name)] = true
	}

	// Add standard patterns that are not in use
	for _, pattern := range standardAliasPatterns {
		if !usedNames[strings.ToLower(pattern)] {
			suggestions = append(suggestions, pattern)
		}
	}

	// Add version-specific suggestions
	versionSuggestions := am.getVersionSpecificSuggestions(version)
	for _, suggestion := range versionSuggestions {
		if !usedNames[strings.ToLower(suggestion)] {
			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

// getVersionSpecificSuggestions generates version-specific alias suggestions
func (am *AliasManager) getVersionSpecificSuggestions(version string) []string {
	var suggestions []string

	// Extract version parts
	version = strings.TrimPrefix(version, "go")
	parts := strings.Split(version, ".")

	if len(parts) >= 2 {
		major := parts[0]
		minor := parts[1]

		// Add major.minor suggestions
		suggestions = append(suggestions, major+"."+minor)
		suggestions = append(suggestions, "go"+major+"."+minor)

		// Add major version suggestions
		suggestions = append(suggestions, major)
		suggestions = append(suggestions, "go"+major)

		// Add specific version suggestions
		suggestions = append(suggestions, version)
		suggestions = append(suggestions, "go"+version)
	}

	// Add common patterns with version
	suggestions = append(suggestions, "v"+version)
	suggestions = append(suggestions, "go-v"+version)

	return suggestions
}

// GetStandardAliasPatterns returns the standard alias patterns
func (am *AliasManager) GetStandardAliasPatterns() []string {
	return standardAliasPatterns
}

// GetStandardAliasGroups returns the standard alias groups
func (am *AliasManager) GetStandardAliasGroups() []string {
	return standardAliasGroups
}
