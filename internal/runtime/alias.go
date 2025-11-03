package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/errors"
	"github.com/molmedoz/gopher/internal/security"
)

// ============================================================================
// Alias Management - Core Operations
// ============================================================================

// loadAliasesOnce is the internal function that loads aliases exactly once
func (am *AliasManager) loadAliasesOnce() {
	// Validate aliases file path is within safe root
	aliasesFileAbs, err := filepath.Abs(am.aliasesFile)
	if err != nil {
		am.loadErr = fmt.Errorf("failed to resolve aliases file path: %w", err)
		return
	}

	// Get safe root directory (parent of InstallDir, e.g., ~/.gopher or ~/gopher)
	installDirAbs, err := filepath.Abs(am.config.InstallDir)
	if err != nil {
		am.loadErr = fmt.Errorf("failed to resolve install directory: %w", err)
		return
	}
	safeRoot := filepath.Dir(installDirAbs) // Parent of versions directory (e.g., ~/.gopher)

	// Validate aliases file is within safe root to prevent path traversal
	safeAliasesFile, err := security.ValidatePathWithinRoot(aliasesFileAbs, safeRoot)
	if err != nil {
		am.loadErr = fmt.Errorf("invalid aliases file path: %w", err)
		return
	}

	// Check if aliases file exists
	if _, err := os.Stat(safeAliasesFile); os.IsNotExist(err) {
		// File doesn't exist, create empty aliases map
		am.mu.Lock()
		am.aliases = make(map[string]*Alias)
		am.mu.Unlock()
		return
	}

	// Read aliases file
	// #nosec G304 -- path validated and scoped to safeRoot
	data, err := os.ReadFile(safeAliasesFile)
	if err != nil {
		am.loadErr = fmt.Errorf("failed to read aliases file: %w", err)
		return
	}

	// Parse JSON
	var aliases map[string]*Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		am.loadErr = fmt.Errorf("failed to parse aliases file: %w", err)
		return
	}

	am.mu.Lock()
	am.aliases = aliases
	am.mu.Unlock()
}

// LoadAliases loads aliases from the aliases file (uses sync.Once for efficiency)
func (am *AliasManager) LoadAliases() error {
	am.once.Do(am.loadAliasesOnce)
	return am.loadErr
}

// SaveAliases saves aliases to the aliases file
func (am *AliasManager) SaveAliases() error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Validate aliases file path is within safe root
	aliasesFileAbs, err := filepath.Abs(am.aliasesFile)
	if err != nil {
		return fmt.Errorf("failed to resolve aliases file path: %w", err)
	}

	// Get safe root directory (parent of InstallDir, e.g., ~/.gopher or ~/gopher)
	installDirAbs, err := filepath.Abs(am.config.InstallDir)
	if err != nil {
		return fmt.Errorf("failed to resolve install directory: %w", err)
	}
	safeRoot := filepath.Dir(installDirAbs) // Parent of versions directory (e.g., ~/.gopher)

	// Validate aliases file is within safe root to prevent path traversal
	safeAliasesFile, err := security.ValidatePathWithinRoot(aliasesFileAbs, safeRoot)
	if err != nil {
		return fmt.Errorf("invalid aliases file path: %w", err)
	}

	// Ensure directory exists
	// Use 0750 for aliases directory - private user data
	dir := filepath.Dir(safeAliasesFile)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create aliases directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(am.aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	// Write to file
	// #nosec G306 -- 0644 acceptable for aliases file (user-managed aliases)
	// #nosec G304 -- path validated and scoped to safeRoot
	if err := os.WriteFile(safeAliasesFile, data, 0644); err != nil {
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

	am.mu.Lock()
	// Check if alias already exists
	if _, exists := am.aliases[name]; exists {
		am.mu.Unlock()
		return errors.Newf(errors.ErrCodeAliasAlreadyExists, "alias '%s' already exists (use 'gopher alias remove %s' first)", name, name)
	}

	// Check if version is installed
	if !am.isVersionInstalled(version) {
		am.mu.Unlock()
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
	am.mu.Unlock()

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

	am.mu.RLock()
	defer am.mu.RUnlock()

	alias, exists := am.aliases[name]
	return alias, exists
}

// ListAliases returns all aliases
func (am *AliasManager) ListAliases() ([]*Alias, error) {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return nil, errors.Wrapf(err, errors.ErrCodeAliasLoadFailed, "failed to load aliases")
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

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

	am.mu.Lock()
	// Check if alias exists
	if _, exists := am.aliases[name]; !exists {
		am.mu.Unlock()
		return errors.Newf(errors.ErrCodeAliasNotFound, "alias '%s' does not exist", name)
	}

	// Remove alias
	delete(am.aliases, name)
	am.mu.Unlock()

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

	am.mu.Lock()
	// Check if alias exists
	alias, exists := am.aliases[name]
	if !exists {
		am.mu.Unlock()
		return errors.Newf(errors.ErrCodeAliasNotFound, "alias '%s' does not exist", name)
	}

	// Check if version is installed
	if !am.isVersionInstalled(version) {
		am.mu.Unlock()
		return errors.Newf(errors.ErrCodeVersionNotInstalled, "version %s is not installed (use 'gopher install %s' first)", version, version)
	}

	// Update alias
	alias.Version = NormalizeVersion(version)
	alias.Updated = time.Now()
	am.mu.Unlock()

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

	am.mu.RLock()
	defer am.mu.RUnlock()

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

	am.mu.RLock()
	var suggestions []string
	usedNames := make(map[string]bool)

	// Mark existing aliases as used
	for name := range am.aliases {
		usedNames[strings.ToLower(name)] = true
	}
	am.mu.RUnlock()

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
