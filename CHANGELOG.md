# Changelog

All notable changes to Gopher will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

*No unreleased changes yet*

## [1.0.0] - 2025-10-15

### Added
- **Version Management**
  - Install, uninstall, and switch between Go versions
  - List installed and available remote versions
  - System Go detection and integration
  - Cross-platform support (Linux, macOS, Windows)
  
- **Alias System**
  - Create, list, and remove version aliases
  - Alias suggestions for common patterns
  - Import/export alias configurations
  - Show aliases by version
  
- **Environment Management**
  - Three GOPATH modes: shared (default), version-specific, and custom
  - Automatic environment variable configuration (GOROOT, GOPATH, GOPROXY, GOSUMDB)
  - Environment activation script generation
  - Shell integration and persistence
  - Comprehensive environment configuration commands (`gopher env`)
  
- **Progress System**
  - Visual progress bars for downloads with auto-sizing
  - Animated spinners for long operations (extraction, uninstall, metadata)
  - Cross-platform terminal handling (Windows, macOS, Linux)
  - Automatic terminal width detection (prevents line wrapping)
  - TTY detection (graceful degradation when piped)
  - Reusable formatters package (bytes, speed, percentage)
  
- **Developer Experience**
  - JSON output support for all commands
  - Verbose and quiet modes
  - Interactive setup wizard (`gopher init`)
  - Comprehensive documentation suite
  - Detailed roadmap with prioritized features
  
### Technical Improvements
- **Dependencies**: Uses `golang.org/x/term` for enhanced terminal support
- **Testing**: 11/11 test suites passing with 36% code coverage
- **Architecture**: Modular design with reusable components
- **Cross-platform**: Platform-specific optimizations for Windows, macOS, and Linux

### Commands Available
- `gopher list` - List installed Go versions (including system)
- `gopher list-remote` - List available Go versions with pagination
- `gopher install <version>` - Install a Go version
- `gopher uninstall <version>` - Uninstall a Go version
- `gopher use <version>` - Switch to a Go version
- `gopher use system` - Switch to system Go
- `gopher current` - Show current Go version
- `gopher system` - Show system Go information
- `gopher alias` - Manage version aliases (create, list, remove, show, import, export)
- `gopher init` - Interactive setup wizard
- `gopher setup` - Shell integration setup
- `gopher status` - Show persistence and shell integration status
- `gopher env` - Environment configuration management
- `gopher debug` - Debug information for troubleshooting
- `gopher version` - Show gopher version

### User Experience Improvements
- Visual progress bars for downloads (auto-sized to terminal width)
- Animated spinners for long operations
- Clean, professional terminal output
- No interference between progress indicators and other output
- Better Windows terminal compatibility

### Fixed
- Windows progress bar issues (progress bars now update on same line)
- Windows spinner compatibility (animations work correctly)
- Terminal width detection prevents line wrapping
- TTY detection for better piped output handling
- Flag parsing order in CLI
- JSON output for empty lists

## [0.1.0] - 2024-01-01

### Added
- Initial release of Gopher
- Basic Go version management functionality
- Install, uninstall, and switch Go versions
- List installed and available versions
- Configuration management
- Cross-platform support
- Zero external dependencies
- SHA256 verification of downloads
- Comprehensive test suite

### Features
- **Version Management**: Install, uninstall, and switch between Go versions
- **System Integration**: Detect and manage system-installed Go versions
- **Security**: Cryptographic verification of all downloads
- **Performance**: Fast, lightweight implementation
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Zero Dependencies**: Uses only Go standard library
- **JSON Support**: Full JSON output for scripting and automation
- **Auto-cleanup**: Automatically manages old versions to save space

### Commands
- `gopher list` - List installed Go versions (including system)
- `gopher list-remote` - List available Go versions
- `gopher install <version>` - Install a Go version
- `gopher uninstall <version>` - Uninstall a Go version
- `gopher use <version>` - Switch to a Go version
- `gopher use system` - Switch to system Go
- `gopher current` - Show current Go version
- `gopher system` - Show system Go information
- `gopher version` - Show gopher version

### Documentation
- Comprehensive user guide
- API reference documentation
- Developer guide for contributors
- Practical usage examples
- Troubleshooting guide

### Installation
```bash
# Using Go install
go install github.com/molmedoz/gopher/cmd/gopher@latest

# From source
git clone https://github.com/molmedoz/gopher.git
cd gopher
go build -o gopher cmd/gopher/main.go
sudo mv gopher /usr/local/bin/
```

### Configuration
Default configuration location:
- **Linux/macOS**: `~/.gopher/config.json`
- **Windows**: `%USERPROFILE%\gopher\config.json`

Default configuration:
```json
{
  "install_dir": "~/.gopher/versions",
  "download_dir": "~/.gopher/downloads",
  "mirror_url": "https://go.dev/dl/",
  "auto_cleanup": true,
  "max_versions": 5
}
```

### System Go Support
Gopher automatically detects and manages system-installed Go versions:
- Homebrew installations (Intel and Apple Silicon)
- System package installations
- Manual installations in standard locations
- Windows installations

### JSON Output
All commands support JSON output for scripting:
```bash
gopher list --json
gopher current --json
gopher system --json
```

### Examples
```bash
# Basic usage
gopher list
gopher install 1.21.0
gopher use 1.21.0
gopher use system

# JSON output
gopher list --json | jq '.[] | select(.is_active) | .version'

# Project-specific versions
echo "1.21.0" > .gopher-version
gopher use $(cat .gopher-version)
```

### Testing
Comprehensive test suite with:
- Unit tests for all components
- Integration tests for system detection
- Error handling tests
- Cross-platform compatibility tests
- Performance tests

### Performance
- Fast installation and switching
- Minimal memory footprint
- Efficient disk usage with auto-cleanup
- Quick system Go detection

### Security
- SHA256 verification of all downloads
- HTTPS-only downloads from official mirrors
- No external dependencies reduce attack surface
- Secure file permissions and handling

---

## Release Notes

### v0.1.0 - Initial Release

This is the initial release of Gopher, a Go version manager that provides:

- **Simple and Fast**: Minimal dependencies, built with Go standard library only
- **System Integration**: Seamlessly manages both system and gopher-installed Go versions
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Secure**: Cryptographic verification of downloaded Go binaries
- **Scriptable**: Full JSON output for automation and scripting
- **Lightweight**: Zero external dependencies beyond Go standard library

Gopher is designed to be a simple, fast, and reliable alternative to existing Go version managers, with a focus on system integration and ease of use.

### Key Features

1. **Version Management**: Install, uninstall, and switch between Go versions
2. **System Go Support**: Automatically detect and manage system-installed Go
3. **JSON Output**: Full JSON support for scripting and automation
4. **Auto-cleanup**: Automatically manage old versions to save space
5. **Cross-platform**: Native support for Linux, macOS, and Windows
6. **Zero Dependencies**: Uses only Go standard library
7. **Security**: SHA256 verification of all downloads
8. **Performance**: Fast installation and switching

### Installation

```bash
go install github.com/molmedoz/gopher/cmd/gopher@latest
```

### Quick Start

```bash
# Check current setup
gopher list
gopher current

# Install a new Go version
gopher install 1.21.0

# Switch to the new version
gopher use 1.21.0

# Switch back to system Go
gopher use system
```

### Documentation

- [User Guide](docs/USER_GUIDE.md) - Comprehensive user documentation
- [API Reference](docs/API_REFERENCE.md) - Detailed API documentation
- [Developer Guide](docs/DEVELOPER_GUIDE.md) - Contributing and development
- [Examples](docs/EXAMPLES.md) - Practical usage examples

### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### License

MIT License - see [LICENSE](LICENSE) file for details.
