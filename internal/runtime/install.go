package runtime

import (
	"fmt"

	"github.com/molmedoz/gopher/internal/errors"
	"github.com/molmedoz/gopher/internal/security"
)

// ============================================================================
// Installation & Uninstallation Operations
// ============================================================================

// Install downloads and installs a specific Go version.
//
// It performs a complete installation workflow including:
//   - Version validation and normalization
//   - Security validation (path traversal protection)
//   - Duplicate installation check
//   - Directory preparation
//   - Download from official Go mirrors
//   - Installation and extraction
//   - Auto-cleanup of old versions (if enabled)
//
// Parameters:
//   - version: The Go version to install (e.g., "1.21.0", "go1.21.0")
//
// Returns an error if the installation fails at any step.
//
// Example:
//
//	err := manager.Install("1.21.0")
//	if err != nil {
//	    log.Fatal("Installation failed:", err)
//	}
func (m *Manager) Install(version string) error {
	// Validate version format
	if err := ValidateVersion(version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Validate version for security (path traversal protection)
	if err := security.ValidatePath(version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Normalize version
	version = NormalizeVersion(version)

	// Check if already installed
	installed, err := m.IsInstalled(version)
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to check if version is installed")
	}
	if installed {
		return errors.NewVersionAlreadyInstalled(version)
	}

	// Ensure directories exist
	if err := m.config.EnsureDirectories(); err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to ensure directories")
	}

	// Download the version
	filePath, err := m.downloader.Download(version, m.config.DownloadDir)
	if err != nil {
		return errors.NewDownloadFailed(version, err)
	}

	// Install the version
	if err := m.installer.Install(version, filePath); err != nil {
		// Clean up downloaded file on failure (ignore errors on cleanup)
		_ = m.downloader.Cleanup(filePath)
		return errors.NewInstallationFailed(version, err)
	}

	// Clean up downloaded file
	if err := m.downloader.Cleanup(filePath); err != nil {
		// Log warning but don't fail the installation
		fmt.Printf("Warning: failed to clean up downloaded file: %v\n", err)
	}

	// Auto-cleanup if enabled
	if m.config.AutoCleanup {
		if err := m.autoCleanup(); err != nil {
			fmt.Printf("Warning: failed to auto-cleanup: %v\n", err)
		}
	}

	return nil
}

// Uninstall removes a specific Go version.
//
// It validates the version, ensures it's installed, and then removes it
// from the installation directory.
//
// Parameters:
//   - version: The Go version to uninstall (e.g., "1.21.0", "go1.21.0")
//
// Returns an error if the uninstallation fails.
//
// Example:
//
//	err := manager.Uninstall("1.21.0")
func (m *Manager) Uninstall(version string) error {
	// Validate version format
	if err := ValidateVersion(version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Validate version for security (path traversal protection)
	if err := security.ValidatePath(version); err != nil {
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

	// Uninstall the version
	if err := m.installer.Uninstall(version); err != nil {
		return errors.Wrapf(err, errors.ErrCodeUninstallationFailed, "failed to uninstall version %s", version)
	}

	return nil
}

// IsInstalled checks if a Go version is currently installed.
//
// It normalizes the version string and checks the installation directory.
//
// Parameters:
//   - version: The Go version to check (e.g., "1.21.0", "go1.21.0")
//
// Returns:
//   - bool: true if installed, false otherwise
//   - error: Any error encountered during the check
//
// Example:
//
//	installed, err := manager.IsInstalled("1.21.0")
//	if err != nil {
//	    log.Fatal("Check failed:", err)
//	}
//	if installed {
//	    fmt.Println("Version is installed!")
//	}
func (m *Manager) IsInstalled(version string) (bool, error) {
	// Normalize version
	version = NormalizeVersion(version)

	// Check if it's installed via installer
	return m.installer.IsInstalled(version), nil
}
