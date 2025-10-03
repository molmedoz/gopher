package runtime

import (
	"fmt"
	"runtime"
	"time"

	"github.com/molmedoz/gopher/internal/errors"
)

// ============================================================================
// Version Switching Operations
// ============================================================================

// Use switches to a specific Go version by creating a symlink.
//
// It handles switching between different Go versions by creating a symlink
// that points to the selected version's Go binary. This allows seamless
// switching between versions without modifying PATH or environment variables.
//
// Special version values:
//   - "system" or "sys": Switch to system-installed Go
//   - Alias names: Resolve to the actual version
//   - Version strings: Direct version switching
//
// Parameters:
//   - version: The Go version to switch to (e.g., "1.21.0", "system", "stable")
//
// Returns an error if the switching fails at any step.
//
// Example:
//
//	// Switch to a specific version
//	err := manager.Use("1.21.0")
//
//	// Switch to system Go
//	err := manager.Use("system")
//
//	// Switch using an alias
//	err := manager.Use("stable")
func (m *Manager) Use(version string) error {
	// Handle special case for system version
	if version == "system" || version == "sys" {
		return m.useSystemVersion()
	}

	// Check if version is an alias
	if alias, exists := m.aliasManager.GetAlias(version); exists {
		fmt.Printf("Using alias '%s' -> %s\n", version, alias.Version)
		version = alias.Version
	}

	// Validate version format
	if err := ValidateVersion(version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Normalize version
	version = NormalizeVersion(version)

	// Check if installed
	installed, err := m.IsInstalled(version)
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to check if version is installed")
	}
	if !installed {
		return errors.NewVersionNotInstalled(version)
	}

	// Get the go binary path
	binaryPath, err := m.installer.GetGoBinaryPath(version)
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to get go binary path")
	}

	// Create symlink or update PATH
	if err := m.createSymlink(binaryPath); err != nil {
		return errors.NewSymlinkFailed(binaryPath, "", err)
	}

	// Try to add symlink directory to PATH for current session
	if err := m.addSymlinkToPath(binaryPath); err != nil {
		fmt.Printf("Warning: failed to add symlink to PATH: %v\n", err)
		fmt.Printf("  You may need to manually add the symlink directory to your PATH\n")
	}

	// Set up environment variables
	if err := m.setupEnvironment(version); err != nil {
		return errors.Wrapf(err, errors.ErrCodeEnvironmentSetupFailed, "failed to setup environment")
	}

	// Save the active version for persistence
	if err := m.saveActiveVersion(version); err != nil {
		fmt.Printf("Warning: failed to save active version: %v\n", err)
	}

	// Set up shell integration for persistence
	if err := m.setupShellIntegration(); err != nil {
		fmt.Printf("Warning: failed to setup shell integration: %v\n", err)
	}

	// On Windows, check PATH order and warn if system Go will take precedence
	if runtime.GOOS == "windows" {
		if err := m.checkWindowsPathOrder(); err != nil {
			// Non-fatal: show warning but don't fail
			fmt.Printf("\n%v\n", err)
		}
	}

	return nil
}

// GetCurrent returns the currently active Go version.
//
// It checks multiple sources to determine which Go version is active:
//   - State file (saved active version)
//   - Symlinks in common locations
//   - System Go as fallback
//
// Returns the active Version with metadata, or an "unknown" version if
// no active version can be determined.
//
// Example:
//
//	current, err := manager.GetCurrent()
//	if err != nil {
//	    log.Fatal("Failed to get current version:", err)
//	}
//	fmt.Printf("Active version: %s\n", current.Version)
func (m *Manager) GetCurrent() (*Version, error) {
	systemDetector := NewSystemDetector()

	// First try to get the active version from state file
	if activeVersion, err := m.getActiveVersionFromState(); err == nil {
		// Check if it's a system version
		if activeVersion == "system" {
			if systemDetector.IsSystemGoAvailable() {
				return systemDetector.DetectSystemGo()
			}
		} else {
			// It's a gopher-managed version, get its info
			if version, err := m.getVersionInfo(activeVersion); err == nil {
				return version, nil
			}
		}
	}

	// Try to detect from symlinks
	if activeVersion, err := m.getCurrentActiveVersion(); err == nil {
		// Check if it's a system version
		if systemDetector.IsSystemGoAvailable() {
			if systemVersion, err := systemDetector.DetectSystemGo(); err == nil {
				if systemVersion.Version == activeVersion {
					return systemVersion, nil
				}
			}
		}

		// It's a gopher-managed version
		if version, err := m.getVersionInfo(activeVersion); err == nil {
			return version, nil
		}
	}

	// If no active version found, try system Go as fallback
	if systemDetector.IsSystemGoAvailable() {
		return systemDetector.DetectSystemGo()
	}

	// If no system Go is available, return unknown
	return &Version{
		Version:     "unknown",
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		IsActive:    true,
		InstalledAt: time.Now(),
	}, nil
}

// useSystemVersion switches to the system Go version.
//
// This is called internally when Use("system") is invoked.
// It handles platform-specific switching logic.
func (m *Manager) useSystemVersion() error {
	systemDetector := NewSystemDetector()
	if !systemDetector.IsSystemGoAvailable() {
		return fmt.Errorf("system Go not available")
	}

	// Get system Go path
	systemPath, err := systemDetector.GetSystemGoPath()
	if err != nil {
		return fmt.Errorf("failed to get system Go path: %w", err)
	}

	// On Windows, remove gopher symlinks to let system Go be found naturally
	if runtime.GOOS == "windows" {
		if err := m.removeGopherSymlinks(); err != nil {
			fmt.Printf("Warning: failed to remove gopher symlinks: %v\n", err)
		}
		fmt.Printf("✓ Switched to system Go version\n")
		fmt.Printf("  System Go path: %s\n", systemPath)
	} else {
		// On Unix systems, create symlink to system Go
		if err := m.createSymlink(systemPath); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
	}

	// Set up environment for system Go
	if err := m.setupSystemEnvironment(); err != nil {
		return fmt.Errorf("failed to setup system environment: %w", err)
	}

	// Save the system version as active
	if err := m.saveActiveVersion("system"); err != nil {
		fmt.Printf("Warning: failed to save active version: %v\n", err)
	}

	// Set up shell integration for persistence
	if err := m.setupShellIntegration(); err != nil {
		fmt.Printf("Warning: failed to setup shell integration: %v\n", err)
	}

	return nil
}

// GetSystemInfo returns detailed information about system Go.
//
// Returns an error if system Go is not available.
//
// Example:
//
//	info, err := manager.GetSystemInfo()
//	if err != nil {
//	    fmt.Println("System Go not available")
//	} else {
//	    fmt.Printf("System Go: %s at %s\n", info.Version, info.GOROOT)
//	}
func (m *Manager) GetSystemInfo() (*SystemGoInfo, error) {
	systemDetector := NewSystemDetector()
	if !systemDetector.IsSystemGoAvailable() {
		return nil, fmt.Errorf("system Go not available")
	}
	return systemDetector.GetSystemGoInfo()
}
