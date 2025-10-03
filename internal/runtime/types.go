// Package version provides the core version management functionality for Gopher.
//
// This package implements the main version management system, including:
//   - Go version installation and uninstallation
//   - Version switching and symlink management
//   - System Go detection and integration
//   - Alias management for version shortcuts
//   - Environment variable management
//   - Download and installation of Go releases
//
// The package is designed with a modular architecture using interfaces
// for dependency injection, making it highly testable and maintainable.
//
// Core Components:
//   - Manager: Main version management orchestrator
//   - AliasManager: Handles version aliases and shortcuts
//   - Version: Represents a Go version with metadata
//   - Various interfaces for download, installation, and system detection
//
// Usage:
//
//	// Create a new manager instance
//	cfg := config.Load("/path/to/config.json")
//	manager := NewManager(cfg, envProvider)
//
//	// Install a Go version
//	err := manager.Install("1.21.0")
//
//	// Switch to a version
//	err := manager.Use("1.21.0")
//
//	// List installed versions
//	versions, err := manager.ListInstalled()
//
// For more detailed examples, see the individual function documentation.
package runtime

import (
	"fmt"
	"runtime"
	"time"

	"github.com/molmedoz/gopher/internal/color"
	"github.com/molmedoz/gopher/internal/config"
	"github.com/molmedoz/gopher/internal/downloader"
	"github.com/molmedoz/gopher/internal/env"
	"github.com/molmedoz/gopher/internal/installer"
)

// Manager is the main orchestrator for Go version management.
//
// It coordinates between different components to provide a unified interface
// for installing, managing, and switching between Go versions. The Manager
// uses dependency injection through interfaces to ensure testability and
// modularity.
//
// The Manager handles:
//   - Go version installation and uninstallation
//   - Version switching via symlinks
//   - System Go detection and integration
//   - Environment variable management
//   - Alias management through AliasManager
//
// Example:
//
//	cfg := config.Load("/path/to/config.json")
//	manager := NewManager(cfg, envProvider)
//	err := manager.Install("1.21.0")
type Manager struct {
	config       *config.Config
	downloader   *downloader.Downloader
	installer    *installer.Installer
	aliasManager *AliasManager
	envProvider  env.Provider
}

// Alias represents a version alias that provides a shortcut name for a Go version.
//
// Aliases allow users to create memorable names for Go versions, making it easier
// to switch between versions without remembering exact version numbers.
//
// Example:
//
//	alias := &Alias{
//	    Name:    "stable",
//	    Version: "1.21.0",
//	    Created: time.Now(),
//	    Updated: time.Now(),
//	    Tags:    []string{"production", "stable"},
//	    Group:   "main",
//	}
type Alias struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Tags    []string  `json:"tags,omitempty"`  // Tags for organization
	Group   string    `json:"group,omitempty"` // Group for organization
}

// AliasManager handles all alias-related operations including creation, deletion,
// listing, and management of version aliases.
//
// It provides methods to:
//   - Create, update, and delete aliases
//   - List all aliases with filtering and grouping
//   - Import and export aliases
//   - Validate alias names and versions
//   - Manage alias tags and groups
//
// The AliasManager persists aliases to a JSON file and provides in-memory
// caching for fast access.
type AliasManager struct {
	config      *config.Config
	aliases     map[string]*Alias
	aliasesFile string
	manager     *Manager // Reference to the main manager for version checking
}

// Version represents a Go version with its metadata and status information.
//
// It contains all the necessary information to identify and manage a Go version,
// including platform details, installation status, and whether it's currently active.
//
// The Version struct is used throughout the system to represent both
// Gopher-managed versions and system-installed Go versions.
//
// Example:
//
//	version := &Version{
//	    Version:     "go1.21.0",
//	    OS:          "darwin",
//	    Arch:        "arm64",
//	    InstalledAt: time.Now(),
//	    IsActive:    true,
//	    IsSystem:    false,
//	    Path:        "/Users/user/.gopher/versions/go1.21.0",
//	}
type Version struct {
	Version     string    `json:"version"`
	OS          string    `json:"os"`
	Arch        string    `json:"arch"`
	InstalledAt time.Time `json:"installed_at"`
	IsActive    bool      `json:"is_active"`
	IsSystem    bool      `json:"is_system"`
	Path        string    `json:"path,omitempty"`
}

// String returns the string representation of the version
func (v *Version) String() string {
	if v.IsSystem {
		return fmt.Sprintf("%s (%s/%s) [system]", v.Version, v.OS, v.Arch)
	}
	return fmt.Sprintf("%s (%s/%s)", v.Version, v.OS, v.Arch)
}

// FullString returns the full string representation including OS and arch
func (v *Version) FullString() string {
	if v.IsSystem {
		return fmt.Sprintf("%s (%s/%s) [system]", v.Version, v.OS, v.Arch)
	}
	return fmt.Sprintf("%s (%s/%s)", v.Version, v.OS, v.Arch)
}

// DisplayString returns the string representation with active indicator and colors
func (v *Version) DisplayString() string {
	base := v.FullString()

	// Add active indicator
	if v.IsActive {
		base = "→ " + base + " [active]"
	} else {
		base = "  " + base
	}

	return base
}

// ColoredDisplayString returns the string representation with active indicator and colors
func (v *Version) ColoredDisplayString() string {
	base := v.FullString()

	// Add active indicator with colors
	if v.IsActive {
		// Use green color for active version
		activeColor := color.ActiveVersion()
		base = activeColor("→ " + base + " [active]")
	} else if v.IsSystem {
		// Use cyan color for system version
		systemColor := color.SystemVersion()
		base = "  " + systemColor(base)
	} else {
		// Use dim color for inactive versions
		inactiveColor := color.InactiveVersion()
		base = "  " + inactiveColor(base)
	}

	return base
}

// IsCompatible checks if the version is compatible with the current platform
func (v *Version) IsCompatible() bool {
	// System versions are always compatible
	if v.IsSystem {
		return true
	}

	// Check if OS and architecture match
	return v.OS == runtime.GOOS && v.Arch == runtime.GOARCH
}

// VersionMetadata represents metadata for an installed version
type VersionMetadata struct {
	Version     string    `json:"version"`
	OS          string    `json:"os"`
	Arch        string    `json:"arch"`
	InstalledAt time.Time `json:"installed_at"`
	InstallDir  string    `json:"install_dir"`
}

// SystemGoInfo represents system Go installation information
type SystemGoInfo struct {
	Version    string `json:"version"`
	GOROOT     string `json:"goroot"`
	GOPATH     string `json:"gopath"`
	Executable string `json:"executable"`
	IsValid    bool   `json:"is_valid"`
}
