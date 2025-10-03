package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/downloader"
)

// ListInstalled returns all installed Go versions including system-installed Go.
//
// It scans the installation directory and system paths to find all available
// Go versions, including both Gopher-managed and system-installed versions.
// The returned list includes metadata about each version such as installation
// time, platform information, and whether it's currently active.
//
// The listing process:
//  1. Detects the currently active version via symlink
//  2. Scans for system-installed Go versions
//  3. Scans the Gopher installation directory
//  4. Deduplicates versions to avoid showing the same version multiple times
//  5. Marks the active version and system versions appropriately
//  6. Sorts versions with system Go first, then by version number
//
// Returns:
//   - []Version: Slice of Version structs with metadata
//   - error: Any error encountered during the listing process
//
// Example:
//
//	versions, err := manager.ListInstalled()
//	if err != nil {
//	    log.Fatal("Failed to list versions:", err)
//	}
//
//	for _, v := range versions {
//	    fmt.Printf("%s (%s/%s) %s\n", v.Version, v.OS, v.Arch, v.InstalledAt)
//	}
func (m *Manager) ListInstalled() ([]Version, error) {
	var result []Version
	seenVersions := make(map[string]bool)

	// Get currently active version by checking symlink
	activeVersion, err := m.getCurrentActiveVersion()
	if err != nil {
		// If we can't determine active version, continue without it
		activeVersion = ""
	}

	// Add system version if available (ALWAYS show it first, regardless of active detection)
	// Try to detect system Go using multiple methods
	systemVersion := m.detectSystemVersionRobust()
	if systemVersion != nil {
		// Check if system is the active version
		systemVersion.IsActive = (activeVersion == "system")
		result = append(result, *systemVersion)
		// Use a unique key for system version to avoid conflicts
		seenVersions["system-"+systemVersion.Version] = true
	}

	// Add gopher-managed versions
	versions, err := m.installer.ListInstalled()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed versions: %w", err)
	}

	for _, v := range versions {
		version, err := m.getVersionInfo(v)
		if err != nil {
			// Skip versions with invalid metadata
			continue
		}

		// Skip if we've already seen this version (avoid duplicates)
		// Check both the version string and the system- prefixed version
		if seenVersions[version.Version] || seenVersions["system-"+version.Version] {
			continue
		}

		// Check if this gopher-managed version is active
		// Only mark as active if it's a gopher-managed version (not system)
		version.IsActive = (version.Version == activeVersion && activeVersion != "system")

		// Add gopher-managed version
		result = append(result, *version)
		seenVersions[version.Version] = true
	}

	return result, nil
}

// detectSystemVersionRobust tries multiple methods to detect system Go version
func (m *Manager) detectSystemVersionRobust() *Version {
	// Method 1: Check common system Go locations directly (bypass PATH entirely)
	// Platform-specific paths
	var systemGoPaths []string
	switch runtime.GOOS {
	case "linux":
		systemGoPaths = []string{
			"/usr/local/go/bin/go",
			"/usr/lib/go/bin/go",
			"/opt/go/bin/go",
			"/usr/bin/go",
			"/usr/local/bin/go", // Common on some Linux distros
		}
	case "darwin":
		systemGoPaths = []string{
			"/usr/local/go/bin/go",
			"/opt/homebrew/bin/go", // Apple Silicon Homebrew
			"/usr/local/bin/go",    // Intel Homebrew
			"/usr/bin/go",
		}
	case "windows":
		systemGoPaths = []string{
			"C:\\Program Files\\Go\\bin\\go.exe",
			"C:\\Go\\bin\\go.exe",
		}
	default:
		systemGoPaths = []string{
			"/usr/local/go/bin/go",
			"/usr/lib/go/bin/go",
			"/opt/go/bin/go",
			"/usr/bin/go",
		}
	}

	for _, goPath := range systemGoPaths {
		if _, err := os.Stat(goPath); err == nil {
			// Found a system Go installation, try to get version
			cmd := exec.Command(goPath, "version")
			if output, err := cmd.Output(); err == nil {
				versionStr := strings.TrimSpace(string(output))
				// Parse version from output like "go version go1.21.0 linux/arm64"
				parts := strings.Fields(versionStr)
				if len(parts) >= 3 && parts[0] == "go" && parts[1] == "version" {
					version := parts[2]
					// Get file info for installation time
					fileInfo, err := os.Stat(goPath)
					var installedAt time.Time
					if err == nil {
						installedAt = fileInfo.ModTime()
					} else {
						installedAt = time.Now()
					}

					return &Version{
						Version:     version,
						OS:          runtime.GOOS,
						Arch:        runtime.GOARCH,
						InstalledAt: installedAt,
						IsActive:    false, // Will be set by caller
						IsSystem:    true,
						Path:        goPath,
					}
				}
			}
		}
	}

	// Method 2: Try to find Go in PATH but avoid gopher-managed versions
	if systemPath, err := exec.LookPath("go"); err == nil {
		// Check if this is not a gopher-managed version
		if !strings.Contains(systemPath, ".gopher") && !strings.Contains(systemPath, "gopher") {
			// Check if it's a symlink (gopher creates symlinks)
			if _, err := os.Readlink(systemPath); err == nil {
				// If it's a symlink, it's likely created by gopher, skip it
			} else {
				// Try to get version
				cmd := exec.Command(systemPath, "version")
				if output, err := cmd.Output(); err == nil {
					versionStr := strings.TrimSpace(string(output))
					parts := strings.Fields(versionStr)
					if len(parts) >= 3 && parts[0] == "go" && parts[1] == "version" {
						version := parts[2]
						// Get file info for installation time
						fileInfo, err := os.Stat(systemPath)
						var installedAt time.Time
						if err == nil {
							installedAt = fileInfo.ModTime()
						} else {
							installedAt = time.Now()
						}

						return &Version{
							Version:     version,
							OS:          runtime.GOOS,
							Arch:        runtime.GOARCH,
							InstalledAt: installedAt,
							IsActive:    false, // Will be set by caller
							IsSystem:    true,
							Path:        systemPath,
						}
					}
				}
			}
		}
	}

	// Method 3: Try to find the original system Go by checking GOROOT
	if goroot := m.envProvider.Getenv("GOROOT"); goroot != "" {
		goPath := filepath.Join(goroot, "bin", "go")
		if runtime.GOOS == "windows" {
			goPath = filepath.Join(goroot, "bin", "go.exe")
		}

		if _, err := os.Stat(goPath); err == nil {
			cmd := exec.Command(goPath, "version")
			if output, err := cmd.Output(); err == nil {
				versionStr := strings.TrimSpace(string(output))
				parts := strings.Fields(versionStr)
				if len(parts) >= 3 && parts[0] == "go" && parts[1] == "version" {
					version := parts[2]
					// Get file info for installation time
					fileInfo, err := os.Stat(goPath)
					var installedAt time.Time
					if err == nil {
						installedAt = fileInfo.ModTime()
					} else {
						installedAt = time.Now()
					}

					return &Version{
						Version:     version,
						OS:          runtime.GOOS,
						Arch:        runtime.GOARCH,
						InstalledAt: installedAt,
						IsActive:    false, // Will be set by caller
						IsSystem:    true,
						Path:        goPath,
					}
				}
			}
		}
	}

	// Method 4: Create a fallback system version entry
	// This ensures we always show something for system Go
	return &Version{
		Version:     "system",
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: time.Now(),
		IsActive:    false, // Will be set by caller
		IsSystem:    true,
		Path:        "system",
	}
}

// ListAvailable returns all available Go versions from official releases
func (m *Manager) ListAvailable() ([]downloader.VersionInfo, error) {
	// Fetch from the Go releases API
	return m.downloader.ListAvailableVersions()
}
