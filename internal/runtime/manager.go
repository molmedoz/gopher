// Package runtime provides Go runtime version management functionality for Gopher.
//
// This package implements a complete Go version management system with support for:
//   - Installing and uninstalling Go versions
//   - Switching between Go versions
//   - Managing system-installed Go
//   - Version aliases for easy switching
//   - Environment variable management
//   - Shell integration for persistence
//
// The package is organized into focused files:
//   - types.go: Type definitions (Manager, Version, Alias, AliasManager)
//   - constructor.go: Constructors and initialization
//   - manager.go: Core manager utilities and configuration (this file)
//   - install.go: Install, uninstall, and installation checks
//   - switch.go: Version switching (Use) and current version detection
//   - list.go: Listing installed and available versions
//   - environment.go: Environment setup, shell integration, and symlinks
//   - system.go: System Go detection and utilities
//   - alias_*.go: Alias management (5 focused files)
//
// Usage Example:
//
//	// Create a new manager
//	cfg := config.Load("/path/to/config.json")
//	envProvider := &env.DefaultProvider{}
//	manager := NewManager(cfg, envProvider)
//
//	// Install a Go version
//	err := manager.Install("1.21.0")
//
//	// Switch to the installed version
//	err = manager.Use("1.21.0")
//
//	// Create an alias
//	err = manager.AliasManager().CreateAlias("stable", "1.21.0")
//
// For more detailed examples, see individual function documentation.
package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/config"
)

// ============================================================================
// Core Manager Methods
// ============================================================================

// GetInstallDir returns the installation directory where Gopher manages Go versions.
func (m *Manager) GetInstallDir() string {
	return m.config.InstallDir
}

// GetDownloadDir returns the directory where Go archives are downloaded.
func (m *Manager) GetDownloadDir() string {
	return m.config.DownloadDir
}

// GetConfig returns the manager's configuration.
func (m *Manager) GetConfig() *config.Config {
	return m.config
}

// AliasManager returns the alias manager instance for managing version aliases.
func (m *Manager) AliasManager() *AliasManager {
	return m.aliasManager
}

// ============================================================================
// Utility Methods
// ============================================================================

// getCurrentActiveVersion determines which version is currently active by checking symlinks.
func (m *Manager) getCurrentActiveVersion() (string, error) {
	// Define potential symlink paths
	var symlinkPaths []string
	if runtime.GOOS == "windows" {
		symlinkPaths = []string{
			filepath.Join(m.envProvider.Getenv("LOCALAPPDATA"), "gopher", "bin", "go.exe"),
		}
	} else {
		symlinkPaths = []string{
			"/usr/local/bin/go",
			filepath.Join(m.envProvider.Getenv("HOME"), ".local", "bin", "go"),
			"/opt/gopher/bin/go",
			"/usr/bin/go",
			"/opt/gopher/go",
		}
	}

	// Check each symlink path
	for _, symlinkPath := range symlinkPaths {
		if target, err := os.Readlink(symlinkPath); err == nil {
			if version := m.extractVersionFromPath(target); version != "" {
				// Verify this is actually a gopher-managed version
				if m.installer.IsInstalled(version) {
					// Check if this symlink directory is in PATH
					if m.isDirectoryInPath(filepath.Dir(symlinkPath)) {
						// Check if this symlink is actually being used
						if m.isSymlinkActuallyUsed(symlinkPath) {
							return version, nil
						}
					}
				}
			}
		}
	}

	// Fallback: check if we have a saved active version
	if activeVersion, err := m.getActiveVersionFromState(); err == nil {
		if m.installer.IsInstalled(activeVersion) {
			if m.hasGopherSymlinkInPath() {
				return activeVersion, nil
			}
		}
	}

	// If no gopher-managed symlink found, check if system Go is active
	systemDetector := NewSystemDetector()
	if systemDetector.IsSystemGoAvailable() {
		if systemPath, err := systemDetector.GetSystemGoPath(); err == nil {
			if m.isDirectoryInPath(filepath.Dir(systemPath)) {
				return "system", nil
			}
		}
	}

	return "unknown", nil
}

// extractVersionFromPath extracts a Go version string from a file path.
//
// It looks for patterns like "go1.21.0" in the path components.
func (m *Manager) extractVersionFromPath(path string) string {
	// Look for version pattern in path (e.g., go1.21.0)
	parts := filepath.SplitList(path)
	for _, part := range parts {
		base := filepath.Base(part)
		// Check if this looks like a Go version
		if strings.HasPrefix(base, "go1.") || strings.HasPrefix(base, "go2.") {
			return base
		}
	}

	// Try splitting by directory separator
	pathParts := strings.Split(path, string(filepath.Separator))
	for _, part := range pathParts {
		if strings.HasPrefix(part, "go1.") || strings.HasPrefix(part, "go2.") {
			return part
		}
	}

	return ""
}

// getVersionInfo gets version information for a given installed version.
//
// It returns an error if:
//   - The version is not installed
//   - The version installation is corrupted (missing Go binary)
//
// For versions installed before the metadata feature was added, it creates
// basic metadata by inspecting the installation directory (backward compatibility).
func (m *Manager) getVersionInfo(version string) (*Version, error) {
	// First verify the version is actually installed
	if !m.installer.IsInstalled(version) {
		return nil, fmt.Errorf("version %s is not installed", version)
	}

	// Try to get metadata
	metadata, err := m.installer.GetVersionMetadata(version)
	if err != nil {
		// Metadata missing could mean:
		// 1. Version installed before metadata feature
		// 2. Metadata file corrupted/deleted
		// 3. Version directory exists but incomplete
		//
		// For backward compatibility with old installations,
		// create basic metadata by inspecting the installation directory
		versionPath := filepath.Join(m.config.InstallDir, version)

		// Verify the Go binary actually exists
		goBinaryPath := filepath.Join(versionPath, "bin", "go")
		if runtime.GOOS == "windows" {
			goBinaryPath += ".exe"
		}

		if _, err := os.Stat(goBinaryPath); err != nil {
			return nil, fmt.Errorf("version %s appears corrupted: go binary not found", version)
		}

		// Get installation time from directory mod time
		installedAt := time.Now()
		if dirInfo, err := os.Stat(versionPath); err == nil {
			installedAt = dirInfo.ModTime()
		}

		// Return basic version info (backward compatibility)
		return &Version{
			Version:     version,
			OS:          runtime.GOOS,
			Arch:        runtime.GOARCH,
			InstalledAt: installedAt,
			IsActive:    false,
			IsSystem:    false,
			Path:        versionPath,
		}, nil
	}

	// Parse metadata fields
	installedAt := time.Now()
	if timeStr, ok := metadata["installed_at"]; ok {
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			installedAt = t
		}
	}

	// Use metadata from file
	return &Version{
		Version:     version,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: installedAt,
		IsActive:    false,
		IsSystem:    false,
		Path:        filepath.Join(m.config.InstallDir, version),
	}, nil
}

// autoCleanup removes old versions if the configured limit is exceeded.
//
// It keeps only the most recent versions up to the MaxVersions limit.
func (m *Manager) autoCleanup() error {
	versions, err := m.installer.ListInstalled()
	if err != nil {
		return fmt.Errorf("failed to list installed versions: %w", err)
	}

	if len(versions) <= m.config.MaxVersions {
		return nil
	}

	// Keep only the most recent versions
	toRemove := len(versions) - m.config.MaxVersions
	for i := 0; i < toRemove; i++ {
		if err := m.Uninstall(versions[i]); err != nil {
			return fmt.Errorf("failed to cleanup version %s: %w", versions[i], err)
		}
	}

	return nil
}
