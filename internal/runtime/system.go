package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/errors"
)

// ============================================================================
// System Go Detection
// ============================================================================

// SystemDetectorImpl handles detection of system-installed Go versions
type SystemDetectorImpl struct{}

// NewSystemDetector creates a new system detector
func NewSystemDetector() *SystemDetectorImpl {
	return &SystemDetectorImpl{}
}

// DetectSystemGo detects the system-installed Go version
func (sd *SystemDetectorImpl) DetectSystemGo() (*Version, error) {
	// Try to find go binary in PATH
	goPath, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("go not found in PATH: %w", err)
	}

	// Get the version by running 'go version'
	cmd := exec.Command(goPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get go version: %w", err)
	}

	// Parse the version output
	versionStr := strings.TrimSpace(string(output))
	version, err := sd.parseGoVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go version: %w", err)
	}

	// Get file info for installation time
	fileInfo, err := os.Stat(goPath)
	var installedAt time.Time
	if err != nil {
		installedAt = time.Now()
	} else {
		installedAt = fileInfo.ModTime()
	}

	// Determine if this is a system installation
	isSystem := sd.isSystemInstallation(goPath)

	return &Version{
		Version:     version,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: installedAt,
		IsActive:    true,
		IsSystem:    isSystem,
		Path:        goPath,
	}, nil
}

// parseGoVersion parses the output of 'go version' command
func (sd *SystemDetectorImpl) parseGoVersion(output string) (string, error) {
	// Expected format: "go version go1.21.0 darwin/arm64"
	// Or for devel versions: "go version devel go1.22-abc123 darwin/arm64"
	parts := strings.Fields(output)
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected go version output format: %s", output)
	}

	// Handle devel versions
	if parts[2] == "devel" {
		if len(parts) < 4 {
			return "", fmt.Errorf("unexpected devel version format: %s", output)
		}
		return parts[2] + " " + parts[3], nil
	}

	// Extract version from "go1.21.0"
	versionPart := parts[2]
	if !strings.HasPrefix(versionPart, "go") {
		return "", fmt.Errorf("unexpected version format: %s", versionPart)
	}

	return versionPart, nil
}

// isSystemInstallation determines if the Go installation is a system installation
func (sd *SystemDetectorImpl) isSystemInstallation(goPath string) bool {
	// Common system installation paths
	systemPaths := []string{
		"/usr/bin/go",
		"/usr/local/bin/go",
		"/opt/go/bin/go",
		"/usr/local/go/bin/go",
		"C:\\Program Files\\Go\\bin\\go.exe",
		"C:\\Go\\bin\\go.exe",
	}

	for _, systemPath := range systemPaths {
		if filepath.Clean(goPath) == filepath.Clean(systemPath) {
			return true
		}
	}

	// Check if it's in a system directory
	dir := filepath.Dir(goPath)
	systemDirs := []string{
		"/usr/bin",
		"/usr/local/bin",
		"/opt/go/bin",
		"/usr/local/go/bin",
		"/opt/homebrew/opt/go/libexec/bin",
		"/usr/local/opt/go/libexec/bin",
		"C:\\Program Files\\Go\\bin",
		"C:\\Go\\bin",
	}

	for _, systemDir := range systemDirs {
		if filepath.Clean(dir) == filepath.Clean(systemDir) {
			return true
		}
	}

	// Check if it's a Homebrew installation
	if strings.Contains(goPath, "/opt/homebrew/") || strings.Contains(goPath, "/usr/local/opt/") {
		return true
	}

	return false
}

// GetSystemGoPath returns the path to the system Go binary
func (sd *SystemDetectorImpl) GetSystemGoPath() (string, error) {
	goPath, err := exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("go not found in PATH: %w", err)
	}
	return goPath, nil
}

// IsSystemGoAvailable checks if system Go is available
func (sd *SystemDetectorImpl) IsSystemGoAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

// GetSystemGoInfo returns detailed information about system Go
func (sd *SystemDetectorImpl) GetSystemGoInfo() (*SystemGoInfo, error) {
	goPath, err := sd.GetSystemGoPath()
	if err != nil {
		return nil, err
	}

	// Get version
	cmd := exec.Command(goPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get go version: %w", err)
	}

	// Get GOROOT
	cmd = exec.Command(goPath, "env", "GOROOT")
	gorootOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get GOROOT: %w", err)
	}

	// Get GOPATH
	cmd = exec.Command(goPath, "env", "GOPATH")
	gopathOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get GOPATH: %w", err)
	}

	// Get file info
	_, err = os.Stat(goPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &SystemGoInfo{
		Version:    strings.TrimSpace(string(output)),
		GOROOT:     strings.TrimSpace(string(gorootOutput)),
		GOPATH:     strings.TrimSpace(string(gopathOutput)),
		Executable: goPath,
		IsValid:    true,
	}, nil
}

// ============================================================================
// Version Utility Functions
// ============================================================================

// ValidateVersion validates a Go version string
func ValidateVersion(version string) error {
	return errors.ValidateVersion(version)
}

// NormalizeVersion normalizes a version string to include 'go' prefix
func NormalizeVersion(version string) string {
	if version == "" {
		return "go"
	}

	// If it already has 'go' prefix, return as is
	if strings.HasPrefix(version, "go") {
		return version
	}

	// Add 'go' prefix
	return "go" + version
}

// CompareVersions compares two version strings
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareVersions(v1, v2 string) int {
	// Normalize versions
	v1 = NormalizeVersion(v1)
	v2 = NormalizeVersion(v2)

	// Remove 'go' prefix for comparison
	v1 = strings.TrimPrefix(v1, "go")
	v2 = strings.TrimPrefix(v2, "go")

	// Split into parts
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare major version
	if len(parts1) > 0 && len(parts2) > 0 {
		if parts1[0] != parts2[0] {
			if parts1[0] < parts2[0] {
				return -1
			}
			return 1
		}
	}

	// Compare minor version
	if len(parts1) > 1 && len(parts2) > 1 {
		if parts1[1] != parts2[1] {
			if parts1[1] < parts2[1] {
				return -1
			}
			return 1
		}
	}

	// Compare patch version
	if len(parts1) > 2 && len(parts2) > 2 {
		if parts1[2] != parts2[2] {
			if parts1[2] < parts2[2] {
				return -1
			}
			return 1
		}
	}

	// If all parts are equal, versions are equal
	return 0
}
