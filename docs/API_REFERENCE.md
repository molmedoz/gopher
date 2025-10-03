# Gopher API Reference

This document provides detailed information about Gopher's internal APIs and data structures.

## Table of Contents

- [Data Structures](#data-structures)
- [Interfaces](#interfaces)
- [Manager API](#manager-api)
- [System Detection API](#system-detection-api)
- [Configuration API](#configuration-api)
- [Error Handling](#error-handling)
- [JSON Schema](#json-schema)

## Data Structures

### Version

Represents a Go version with metadata.

```go
type Version struct {
    Version     string    `json:"version"`      // e.g., "1.21.0"
    OS          string    `json:"os"`           // e.g., "darwin", "linux", "windows"
    Arch        string    `json:"arch"`         // e.g., "amd64", "arm64"
    InstalledAt time.Time `json:"installed_at"`
    IsActive    bool      `json:"is_active"`
    IsSystem    bool      `json:"is_system"`    // true if system installation
    Path        string    `json:"path"`         // path to go binary
}
```

**Methods:**
- `String() string` - Basic string representation
- `FullString() string` - Detailed string representation
- `IsCompatible() bool` - Check if compatible with current system

### VersionInfo

Information about available Go versions.

```go
type VersionInfo struct {
    Version     string `json:"version"`
    Stable      bool   `json:"stable"`
    ReleaseDate string `json:"release_date"`
    Files       []File `json:"files"`
}
```

### File

Downloadable file information.

```go
type File struct {
    Filename string `json:"filename"`
    OS       string `json:"os"`
    Arch     string `json:"arch"`
    Size     int64  `json:"size"`
    SHA256   string `json:"sha256"`
}
```

### SystemGoInfo

Detailed system Go information.

```go
type SystemGoInfo struct {
    Path        string    `json:"path"`
    Version     string    `json:"version"`
    GOROOT      string    `json:"goroot"`
    GOPATH      string    `json:"gopath"`
    InstalledAt time.Time `json:"installed_at"`
    IsSystem    bool      `json:"is_system"`
}
```

### Config

Gopher configuration.

```go
type Config struct {
    InstallDir  string `json:"install_dir"`
    DownloadDir string `json:"download_dir"`
    MirrorURL   string `json:"mirror_url"`
    AutoCleanup bool   `json:"auto_cleanup"`
    MaxVersions int    `json:"max_versions"`
}
```

## Interfaces

### ManagerInterface

Main interface for version management operations.

```go
type ManagerInterface interface {
    // ListInstalled returns all installed Go versions
    ListInstalled() ([]Version, error)
    
    // ListAvailable returns all available Go versions
    ListAvailable() ([]VersionInfo, error)
    
    // Install downloads and installs a Go version
    Install(version string) error
    
    // Uninstall removes a Go version
    Uninstall(version string) error
    
    // Use switches to a Go version
    Use(version string) error
    
    // GetCurrent returns the currently active Go version
    GetCurrent() (*Version, error)
    
    // IsInstalled checks if a version is installed
    IsInstalled(version string) (bool, error)
}
```

## Manager API

### Manager

Main manager implementation.

```go
type Manager struct {
    config         *config.Config
    downloader     *downloader.Downloader
    installer      *installer.Installer
    systemDetector *SystemDetector
}
```

### Constructor

```go
func NewManager(cfg *config.Config) *Manager
```

Creates a new manager with the provided configuration.

### Core Methods

#### ListInstalled

```go
func (m *Manager) ListInstalled() ([]Version, error)
```

Returns all installed Go versions, including system versions.

**Returns:**
- `[]Version` - List of installed versions
- `error` - Any error that occurred

**Example:**
```go
versions, err := manager.ListInstalled()
if err != nil {
    log.Fatal(err)
}

for _, v := range versions {
    fmt.Printf("Version: %s, System: %v\n", v.Version, v.IsSystem)
}
```

#### Install

```go
func (m *Manager) Install(version string) error
```

Downloads and installs a specific Go version.

**Parameters:**
- `version` - Go version to install (e.g., "1.21.0", "go1.21.0")

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
err := manager.Install("1.21.0")
if err != nil {
    log.Fatal(err)
}
```

#### Use

```go
func (m *Manager) Use(version string) error
```

Switches to a specific Go version.

**Parameters:**
- `version` - Go version to use, or "system" for system Go

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
// Switch to specific version
err := manager.Use("1.21.0")
if err != nil {
    log.Fatal(err)
}

// Switch to system Go
err = manager.Use("system")
if err != nil {
    log.Fatal(err)
}
```

#### GetCurrent

```go
func (m *Manager) GetCurrent() (*Version, error)
```

Returns the currently active Go version.

**Returns:**
- `*Version` - Current version information
- `error` - Any error that occurred

**Example:**
```go
current, err := manager.GetCurrent()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Current version: %s\n", current.Version)
```

#### IsInstalled

```go
func (m *Manager) IsInstalled(version string) (bool, error)
```

Checks if a version is installed.

**Parameters:**
- `version` - Go version to check, or "system" for system Go

**Returns:**
- `bool` - True if installed
- `error` - Any error that occurred

**Example:**
```go
installed, err := manager.IsInstalled("1.21.0")
if err != nil {
    log.Fatal(err)
}
if installed {
    fmt.Println("Go 1.21.0 is installed")
}
```

### System Methods

#### GetSystemInfo

```go
func (m *Manager) GetSystemInfo() (*SystemGoInfo, error)
```

Returns detailed information about system Go.

**Returns:**
- `*SystemGoInfo` - System Go information
- `error` - Any error that occurred

**Example:**
```go
info, err := manager.GetSystemInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("System Go: %s\n", info.Version)
```

## System Detection API

### SystemDetector

Handles detection of system-installed Go versions.

```go
type SystemDetector struct{}
```

### Constructor

```go
func NewSystemDetector() *SystemDetector
```

Creates a new system detector.

### Core Methods

#### DetectSystemGo

```go
func (sd *SystemDetector) DetectSystemGo() (*Version, error)
```

Detects the system-installed Go version.

**Returns:**
- `*Version` - System Go version information
- `error` - Any error that occurred

**Example:**
```go
detector := NewSystemDetector()
version, err := detector.DetectSystemGo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("System Go: %s\n", version.Version)
```

#### IsSystemGoAvailable

```go
func (sd *SystemDetector) IsSystemGoAvailable() bool
```

Checks if system Go is available.

**Returns:**
- `bool` - True if system Go is available

**Example:**
```go
detector := NewSystemDetector()
if detector.IsSystemGoAvailable() {
    fmt.Println("System Go is available")
}
```

#### GetSystemGoPath

```go
func (sd *SystemDetector) GetSystemGoPath() (string, error)
```

Returns the path to the system Go binary.

**Returns:**
- `string` - Path to go binary
- `error` - Any error that occurred

**Example:**
```go
path, err := detector.GetSystemGoPath()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("System Go path: %s\n", path)
```

#### GetSystemGoInfo

```go
func (sd *SystemDetector) GetSystemGoInfo() (*SystemGoInfo, error)
```

Returns detailed system Go information.

**Returns:**
- `*SystemGoInfo` - Detailed system Go information
- `error` - Any error that occurred

**Example:**
```go
info, err := detector.GetSystemGoInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("GOROOT: %s\n", info.GOROOT)
```

## Configuration API

### Config

Configuration management.

```go
type Config struct {
    InstallDir  string `json:"install_dir"`
    DownloadDir string `json:"download_dir"`
    MirrorURL   string `json:"mirror_url"`
    AutoCleanup bool   `json:"auto_cleanup"`
    MaxVersions int    `json:"max_versions"`
}
```

### Constructor

```go
func DefaultConfig() *Config
```

Returns the default configuration.

### Methods

#### Load

```go
func Load(configPath string) (*Config, error)
```

Loads configuration from file.

**Parameters:**
- `configPath` - Path to configuration file

**Returns:**
- `*Config` - Loaded configuration
- `error` - Any error that occurred

**Example:**
```go
config, err := config.Load("~/.gopher/config.json")
if err != nil {
    log.Fatal(err)
}
```

#### Save

```go
func (c *Config) Save(configPath string) error
```

Saves configuration to file.

**Parameters:**
- `configPath` - Path to save configuration

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
err := config.Save("~/.gopher/config.json")
if err != nil {
    log.Fatal(err)
}
```

#### Validate

```go
func (c *Config) Validate() error
```

Validates the configuration.

**Returns:**
- `error` - Any validation errors

**Example:**
```go
err := config.Validate()
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

### Common Error Types

#### Validation Errors

```go
err := ValidateVersion("invalid")
// Returns: "invalid version format: invalid"
```

#### Installation Errors

```go
err := manager.Install("1.21.0")
// Possible errors:
// - "version 1.21.0 is already installed"
// - "failed to download version: ..."
// - "failed to install version: ..."
```

#### System Errors

```go
err := manager.Use("system")
// Possible errors:
// - "system Go is not available"
// - "failed to create symlink: permission denied"
```

### Error Handling Best Practices

```go
// Check for specific error types
if err != nil {
    if strings.Contains(err.Error(), "permission denied") {
        // Handle permission errors
    } else if strings.Contains(err.Error(), "not found") {
        // Handle not found errors
    } else {
        // Handle other errors
    }
}
```

## JSON Schema

### Version JSON Schema

```json
{
  "type": "object",
  "properties": {
    "version": {
      "type": "string",
      "description": "Go version (e.g., '1.21.0')"
    },
    "os": {
      "type": "string",
      "description": "Operating system (e.g., 'darwin', 'linux', 'windows')"
    },
    "arch": {
      "type": "string",
      "description": "Architecture (e.g., 'amd64', 'arm64')"
    },
    "installed_at": {
      "type": "string",
      "format": "date-time",
      "description": "Installation timestamp"
    },
    "is_active": {
      "type": "boolean",
      "description": "Whether this version is currently active"
    },
    "is_system": {
      "type": "boolean",
      "description": "Whether this is a system installation"
    },
    "path": {
      "type": "string",
      "description": "Path to the go binary"
    }
  },
  "required": ["version", "os", "arch", "installed_at", "is_active", "is_system", "path"]
}
```

### SystemGoInfo JSON Schema

```json
{
  "type": "object",
  "properties": {
    "path": {
      "type": "string",
      "description": "Path to the go binary"
    },
    "version": {
      "type": "string",
      "description": "Full go version output"
    },
    "goroot": {
      "type": "string",
      "description": "GOROOT environment variable value"
    },
    "gopath": {
      "type": "string",
      "description": "GOPATH environment variable value"
    },
    "installed_at": {
      "type": "string",
      "format": "date-time",
      "description": "Installation timestamp"
    },
    "is_system": {
      "type": "boolean",
      "description": "Whether this is a system installation"
    }
  },
  "required": ["path", "version", "goroot", "gopath", "installed_at", "is_system"]
}
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/molmedoz/gopher/internal/config"
    "github.com/molmedoz/gopher/internal/version"
)

func main() {
    // Load configuration
    cfg, err := config.Load(config.GetConfigPath())
    if err != nil {
        log.Fatal(err)
    }
    
    // Create manager
    manager := version.NewManager(cfg)
    
    // List installed versions
    versions, err := manager.ListInstalled()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, v := range versions {
        fmt.Printf("Version: %s, System: %v\n", v.Version, v.IsSystem)
    }
}
```

### System Go Detection

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/molmedoz/gopher/internal/version"
)

func main() {
    detector := version.NewSystemDetector()
    
    if detector.IsSystemGoAvailable() {
        info, err := detector.GetSystemGoInfo()
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("System Go: %s\n", info.Version)
        fmt.Printf("GOROOT: %s\n", info.GOROOT)
        fmt.Printf("GOPATH: %s\n", info.GOPATH)
    } else {
        fmt.Println("No system Go found")
    }
}
```

### Custom Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/molmedoz/gopher/internal/config"
    "github.com/molmedoz/gopher/internal/version"
)

func main() {
    // Create custom configuration
    cfg := &config.Config{
        InstallDir:  "/opt/go-versions",
        DownloadDir: "/tmp/gopher-downloads",
        MirrorURL:   "https://go.dev/dl/",
        AutoCleanup: true,
        MaxVersions: 10,
    }
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        log.Fatal(err)
    }
    
    // Create manager with custom config
    manager := version.NewManager(cfg)
    
    // Use manager...
}
```

This API reference provides comprehensive information about Gopher's internal APIs. For more examples and usage patterns, see the [User Guide](USER_GUIDE.md).
