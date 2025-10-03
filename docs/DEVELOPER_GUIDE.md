# Gopher Developer Guide

This guide is for developers who want to contribute to Gopher or understand its internal architecture.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Development Setup](#development-setup)
- [Building and Testing](#building-and-testing)
- [Code Style and Standards](#code-style-and-standards)
- [Adding New Features](#adding-new-features)
- [Debugging](#debugging)
- [Release Process](#release-process)

## Architecture Overview

Gopher follows a modular architecture with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Manager Layer  │    │  Storage Layer  │
│                 │    │                 │    │                 │
│  cmd/gopher/    │◄──►│ internal/version│◄──►│ internal/installer│
│  main.go        │    │  manager.go     │    │  installer.go   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐    ┌─────────────────┐
         │              │  System Layer   │    │ Download Layer  │
         └──────────────►│                 │    │                 │
                        │ internal/version│    │ internal/downloader│
                        │ system.go       │    │ downloader.go   │
                        └─────────────────┘    └─────────────────┘
                                 │
                        ┌─────────────────┐
                        │  Config Layer   │
                        │                 │
                        │ internal/config │
                        │ config.go       │
                        └─────────────────┘
```

### Key Components

1. **CLI Layer** (`cmd/gopher/`): Command-line interface and argument parsing
2. **Manager Layer** (`internal/version/`): Core version management logic
3. **System Layer** (`internal/version/system.go`): System Go detection and management
4. **Download Layer** (`internal/downloader/`): Go version downloading and verification
5. **Storage Layer** (`internal/installer/`): Go version installation and extraction
6. **Config Layer** (`internal/config/`): Configuration management

## Project Structure

```
gopher/
├── cmd/
│   └── gopher/
│       └── main.go              # CLI entry point
├── internal/
│   ├── version/
│   │   ├── types.go             # Core data structures
│   │   ├── manager.go           # Main manager implementation
│   │   ├── system.go            # System Go detection
│   │   ├── types_test.go        # Type tests
│   │   ├── manager_test.go      # Manager tests
│   │   └── system_test.go       # System tests
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── config_test.go       # Config tests
│   ├── downloader/
│   │   └── downloader.go        # Download and verification
│   └── installer/
│       └── installer.go         # Installation and extraction
├── docs/
│   ├── USER_GUIDE.md            # User documentation
│   ├── API_REFERENCE.md         # API documentation
│   └── DEVELOPER_GUIDE.md       # This file
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Build and development commands
└── README.md                    # Project overview
```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile commands)

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/molmedoz/gopher.git
cd gopher

# Install dependencies
go mod tidy

# Build the project
make build

# Run tests
make test
```

### Development Tools

Install recommended development tools:

```bash
# Install development tools
make install-tools

# This installs:
# - golangci-lint (linting)
# - goimports (import formatting)
# - air (hot reloading)
```

## Building and Testing

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build development version
make build-dev

# Clean build artifacts
make clean
```

### Test Commands

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/version -v
```

### Code Quality

```bash
# Format code
make format

# Run linter
make lint

# Run go vet
make vet

# Run all checks
make check
```

### Development Workflow

```bash
# Start development cycle
make dev

# This runs:
# 1. Format code
# 2. Run linter
# 3. Run tests
# 4. Build binary
```

## Code Style and Standards

### Go Code Style

Follow standard Go conventions:

```go
// Package comment
package version

import (
    "fmt"
    "time"
)

// Interface comment
type ManagerInterface interface {
    // Method comment
    ListInstalled() ([]Version, error)
}

// Struct comment
type Version struct {
    Version     string    `json:"version"`      // Field comment
    OS          string    `json:"os"`           // Field comment
    Arch        string    `json:"arch"`         // Field comment
    InstalledAt time.Time `json:"installed_at"`
    IsActive    bool      `json:"is_active"`
    IsSystem    bool      `json:"is_system"`
    Path        string    `json:"path"`
}

// Method comment
func (v Version) String() string {
    if v.IsSystem {
        return fmt.Sprintf("%s (%s/%s) [system]", v.Version, v.OS, v.Arch)
    }
    return fmt.Sprintf("%s (%s/%s)", v.Version, v.OS, v.Arch)
}
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `version`, `config`)
- **Interfaces**: descriptive name ending with `-er` (e.g., `ManagerInterface`)
- **Structs**: PascalCase (e.g., `Version`, `SystemDetector`)
- **Functions**: PascalCase for public, camelCase for private
- **Variables**: camelCase (e.g., `versionString`, `configPath`)

### Error Handling

Always handle errors explicitly:

```go
// Good
func (m *Manager) Install(version string) error {
    if err := ValidateVersion(version); err != nil {
        return fmt.Errorf("invalid version: %w", err)
    }
    
    // ... rest of implementation
}

// Bad
func (m *Manager) Install(version string) error {
    ValidateVersion(version) // Error ignored
    
    // ... rest of implementation
}
```

### Testing Standards

Write comprehensive tests:

```go
func TestValidateVersion(t *testing.T) {
    tests := []struct {
        name    string
        version string
        wantErr bool
    }{
        {"valid version", "1.21.0", false},
        {"invalid version", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateVersion(tt.version)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Adding New Features

### 1. Plan the Feature

- Define the user-facing interface
- Identify required internal changes
- Consider backward compatibility
- Plan test coverage

### 2. Update Data Structures

If new data is needed, update the relevant structs:

```go
// Add new fields to existing structs
type Version struct {
    // ... existing fields
    NewField string `json:"new_field"`
}

// Or create new structs
type NewFeature struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}
```

### 3. Update Interfaces

Add new methods to interfaces:

```go
type ManagerInterface interface {
    // ... existing methods
    NewMethod() (NewFeature, error)
}
```

### 4. Implement the Feature

Implement the feature in the appropriate layer:

```go
func (m *Manager) NewMethod() (NewFeature, error) {
    // Implementation
}
```

### 5. Update CLI

Add new commands or options:

```go
case "new-command":
    return newCommand(manager, args)
```

### 6. Add Tests

Write comprehensive tests:

```go
func TestNewMethod(t *testing.T) {
    // Test implementation
}
```

### 7. Update Documentation

- Update README.md
- Update USER_GUIDE.md
- Update API_REFERENCE.md
- Add examples

## Debugging

### Debug Mode

Enable debug output:

```bash
export GOPHER_DEBUG=1
gopher list
```

### Logging

Add debug logging:

```go
import "log"

func (m *Manager) Install(version string) error {
    if os.Getenv("GOPHER_DEBUG") == "1" {
        log.Printf("Installing version: %s", version)
    }
    
    // ... implementation
}
```

### Testing Specific Features

```bash
# Test specific functionality
go test -run TestSystemDetector ./internal/version -v

# Test with race detection
go test -race ./internal/version

# Test with coverage
go test -cover ./internal/version
```

### Common Debug Scenarios

#### System Go Not Detected

```go
// Add debug logging to system detection
func (sd *SystemDetector) DetectSystemGo() (*Version, error) {
    goPath, err := exec.LookPath("go")
    if err != nil {
        log.Printf("DEBUG: go not found in PATH: %v", err)
        return nil, err
    }
    log.Printf("DEBUG: found go at: %s", goPath)
    // ... rest of implementation
}
```

#### Download Issues

```go
// Add debug logging to downloader
func (d *Downloader) Download(version, downloadDir string) (string, error) {
    log.Printf("DEBUG: downloading version %s to %s", version, downloadDir)
    // ... implementation
}
```

## Release Process

### 1. Update Version

Update version in `cmd/gopher/main.go`:

```go
const versionString = "gopher v1.0.0"
```

### 2. Update Changelog

Create or update `CHANGELOG.md`:

```markdown
## [0.2.0] - 2024-01-01

### Added
- New feature X
- New command Y

### Changed
- Improved performance
- Updated dependencies

### Fixed
- Bug fix A
- Bug fix B
```

### 3. Run Full Test Suite

```bash
# Run all tests
make test

# Run with race detection
make race-test

# Run linting
make lint

# Run security checks
make security
```

### 4. Build Release

```bash
# Build for all platforms
make build-all

# Create distribution packages
make dist
```

### 5. Create Release

```bash
# Tag the release
git tag v1.0.0
git push origin v1.0.0

# Create GitHub release
gh release create v1.0.0 \
  --title "Gopher v1.0.0" \
  --notes "$(cat CHANGELOG.md)" \
  build/release/*
```

## Contributing Guidelines

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Run the full test suite
7. Submit a pull request

### Commit Messages

Use conventional commits:

```
feat: add system Go detection
fix: resolve download timeout issue
docs: update API reference
test: add tests for system detection
refactor: simplify version parsing
```

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests are comprehensive
- [ ] Documentation is updated
- [ ] No breaking changes (or documented)
- [ ] Error handling is proper
- [ ] Performance is acceptable
- [ ] Security considerations addressed

## Performance Considerations

### Memory Usage

- Use pointers for large structs
- Avoid unnecessary allocations
- Reuse buffers where possible

### Network Operations

- Use timeouts for downloads
- Implement retry logic
- Cache results when appropriate

### File Operations

- Use appropriate file permissions
- Clean up temporary files
- Handle concurrent access

## Security Considerations

### Input Validation

- Validate all user inputs
- Sanitize file paths
- Check file permissions

### Download Security

- Verify SHA256 checksums
- Use HTTPS for downloads
- Validate file sizes

### File System Security

- Use secure file permissions
- Avoid symlink attacks
- Validate file paths

## Troubleshooting

### Common Issues

#### Build Failures

```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version

# Update dependencies
go mod tidy
```

#### Test Failures

```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestSpecific ./internal/version -v

# Check for race conditions
go test -race ./...
```

#### Linting Issues

```bash
# Fix formatting
go fmt ./...

# Fix imports
goimports -w .

# Run linter
golangci-lint run
```

This developer guide provides comprehensive information for contributing to Gopher. For more specific questions, see the [API Reference](API_REFERENCE.md) or open an issue on GitHub.

## See Also

### For Developers

- **[API Reference](API_REFERENCE.md)** - Complete API documentation
- **[Test Strategy](TEST_STRATEGY.md)** - Testing approach and best practices
- **[Refactoring Summary](REFACTORING_SUMMARY.md)** - Recent architectural changes
- **[Contributing](../CONTRIBUTING.md)** - Contribution guidelines
- **[Roadmap](ROADMAP.md)** - Future features and enhancements

### For Users

- **[User Guide](USER_GUIDE.md)** - Complete user documentation
- **[Quick Reference](../QUICK_REFERENCE.md)** - One-page command reference
- **[FAQ](FAQ.md)** - Frequently asked questions
- **[Examples](EXAMPLES.md)** - Practical usage examples

### Testing & Quality

- **[Testing Guide](TESTING_GUIDE.md)** - Multi-platform testing strategy
- **[VM Setup Guide](VM_SETUP_GUIDE.md)** - VM setup for testing
- **[Docker Testing](../docker/README.md)** - Docker-based testing

### Additional Resources

- **[Documentation Index](../DOCUMENTATION_INDEX.md)** - Navigate all documentation
- **[README](../README.md)** - Project overview
- **[Changelog](../CHANGELOG.md)** - Version history
