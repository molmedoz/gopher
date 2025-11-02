package runtime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ============================================================================
// Alias Management - Advanced Features
// ============================================================================

// CreateAliasInteractive creates an alias with interactive conflict resolution
func (am *AliasManager) CreateAliasInteractive(name, version string, allowOverride, noOverride, force bool) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return fmt.Errorf("failed to load aliases: %w", err)
	}

	// Validate alias name
	if err := am.ValidateAliasName(name); err != nil {
		return fmt.Errorf("invalid alias name: %w", err)
	}

	// Check if version is installed
	if !am.isVersionInstalled(version) {
		return fmt.Errorf("version %s is not installed (use 'gopher install %s' first)", version, version)
	}

	// Check if alias already exists
	if existing, exists := am.aliases[name]; exists {
		// Handle conflict resolution
		if force {
			// Force mode - update without confirmation
			existing.Version = NormalizeVersion(version)
			existing.Updated = time.Now()
		} else if noOverride {
			// No override mode - return error
			return fmt.Errorf("alias '%s' already exists and points to %s (use 'gopher alias remove %s' first)", name, existing.Version, name)
		} else if allowOverride {
			// Allow override mode - update without confirmation
			existing.Version = NormalizeVersion(version)
			existing.Updated = time.Now()
		} else {
			// Interactive mode - ask for confirmation
			if err := am.handleAliasConflict(name, existing.Version, version); err != nil {
				return err
			}
			// If we get here, user confirmed the update
			existing.Version = NormalizeVersion(version)
			existing.Updated = time.Now()
		}
	} else {
		// Create new alias
		alias := &Alias{
			Name:    name,
			Version: NormalizeVersion(version),
			Created: time.Now(),
			Updated: time.Now(),
		}
		am.aliases[name] = alias
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		return fmt.Errorf("failed to save aliases: %w", err)
	}

	return nil
}

// UpdateAliasInteractive updates an alias with interactive conflict resolution
func (am *AliasManager) UpdateAliasInteractive(name, version string, allowOverride, noOverride, force bool) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return fmt.Errorf("failed to load aliases: %w", err)
	}

	// Check if alias exists
	existing, exists := am.aliases[name]
	if !exists {
		return fmt.Errorf("alias '%s' does not exist", name)
	}

	// Check if version is installed
	if !am.isVersionInstalled(version) {
		return fmt.Errorf("version %s is not installed (use 'gopher install %s' first)", version, version)
	}

	// Handle conflict resolution if version is different
	if existing.Version != NormalizeVersion(version) {
		if force {
			// Force mode - update without confirmation
			existing.Version = NormalizeVersion(version)
			existing.Updated = time.Now()
		} else if noOverride {
			// No override mode - return error
			return fmt.Errorf("alias '%s' already points to %s (use 'gopher alias remove %s' first)", name, existing.Version, name)
		} else if allowOverride {
			// Allow override mode - update without confirmation
			existing.Version = NormalizeVersion(version)
			existing.Updated = time.Now()
		} else {
			// Interactive mode - ask for confirmation
			if err := am.handleAliasConflict(name, existing.Version, version); err != nil {
				return err
			}
			// If we get here, user confirmed the update
			existing.Version = NormalizeVersion(version)
			existing.Updated = time.Now()
		}
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		return fmt.Errorf("failed to save aliases: %w", err)
	}

	return nil
}

// CreateAliasesBulk creates multiple aliases with conflict resolution
func (am *AliasManager) CreateAliasesBulk(aliases map[string]string, allowOverride, noOverride, force bool) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return fmt.Errorf("failed to load aliases: %w", err)
	}

	// Validate all aliases first
	for name, version := range aliases {
		if err := am.ValidateAliasName(name); err != nil {
			return fmt.Errorf("invalid alias name '%s': %w", name, err)
		}
		if !am.isVersionInstalled(version) {
			return fmt.Errorf("version %s is not installed for alias \"%s\" (use 'gopher install %s' first)", version, name, version)
		}
	}

	// Process each alias
	for name, version := range aliases {
		normalizedVersion := NormalizeVersion(version)

		if existing, exists := am.aliases[name]; exists {
			// Handle conflict resolution
			if force {
				// Force mode - update without confirmation
				existing.Version = normalizedVersion
				existing.Updated = time.Now()
			} else if noOverride {
				// No override mode - return error
				return fmt.Errorf("alias '%s' already exists and points to %s (use 'gopher alias remove %s' first)", name, existing.Version, name)
			} else if allowOverride {
				// Allow override mode - update without confirmation
				existing.Version = normalizedVersion
				existing.Updated = time.Now()
			} else {
				// Interactive mode - ask for confirmation
				if err := am.handleAliasConflict(name, existing.Version, version); err != nil {
					return err
				}
				// If we get here, user confirmed the update
				existing.Version = normalizedVersion
				existing.Updated = time.Now()
			}
		} else {
			// Create new alias
			alias := &Alias{
				Name:    name,
				Version: normalizedVersion,
				Created: time.Now(),
				Updated: time.Now(),
			}
			am.aliases[name] = alias
		}
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		return fmt.Errorf("failed to save aliases: %w", err)
	}

	return nil
}

// handleAliasConflict handles interactive conflict resolution
func (am *AliasManager) handleAliasConflict(name, currentVersion, newVersion string) error {
	fmt.Printf("\n⚠️  Alias Conflict Detected\n")
	fmt.Printf("   Alias: %s\n", name)
	fmt.Printf("   Current: %s\n", currentVersion)
	fmt.Printf("   New:     %s\n", newVersion)
	fmt.Print("\nUpdate alias? (y/yes/n/no): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	switch response {
	case "y", "yes":
		return nil // User confirmed, continue with update
	case "n", "no":
		return fmt.Errorf("operation cancelled by user")
	default:
		return fmt.Errorf("invalid response: %s (expected y/yes/n/no)", response)
	}
}

// ExportAliases exports aliases to a file
func (am *AliasManager) ExportAliases(filename string, format string, tags []string) error {
	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		return fmt.Errorf("failed to load aliases: %w", err)
	}

	// Filter aliases by tags if specified
	var aliasesToExport map[string]*Alias
	if len(tags) > 0 {
		aliasesToExport = make(map[string]*Alias)
		for name, alias := range am.aliases {
			// Check if alias has any of the specified tags
			for _, tag := range tags {
				for _, aliasTag := range alias.Tags {
					if strings.EqualFold(aliasTag, tag) {
						aliasesToExport[name] = alias
						break
					}
				}
			}
		}
	} else {
		aliasesToExport = am.aliases
	}

	// Determine file format
	if format == "" {
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".json":
			format = "json"
		default:
			format = "json" // Default to JSON
		}
	}

	// Export based on format
	switch strings.ToLower(format) {
	case "json":
		return am.exportToJSON(filename, aliasesToExport)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// ImportAliases imports aliases from a file
func (am *AliasManager) ImportAliases(filename string, allowOverride, noOverride, force bool) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Determine file format
	ext := strings.ToLower(filepath.Ext(filename))
	var format string
	switch ext {
	case ".json":
		format = "json"
	default:
		format = "json" // Default to JSON
	}

	// Import based on format
	var aliases map[string]string
	var err error
	switch strings.ToLower(format) {
	case "json":
		aliases, err = am.importFromJSON(filename)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to import aliases: %w", err)
	}

	// Create aliases using bulk creation
	return am.CreateAliasesBulk(aliases, allowOverride, noOverride, force)
}

// exportToJSON exports aliases to JSON file
func (am *AliasManager) exportToJSON(filename string, aliases map[string]*Alias) error {
	// Marshal to JSON
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	// Write to file
	// #nosec G306 -- 0644 acceptable for aliases file (user-managed aliases)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// importFromJSON imports aliases from JSON file
func (am *AliasManager) importFromJSON(filename string) (map[string]string, error) {
	// Read file
	// #nosec G304 -- filename is user-provided alias file path (validated by caller)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var aliases map[string]*Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to map[string]string for bulk creation
	result := make(map[string]string)
	for name, alias := range aliases {
		result[name] = alias.Version
	}

	return result, nil
}
