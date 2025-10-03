<div align="center">
  <img src="gopher.jpeg" alt="Gopher Logo" width="200"/>
  
  # Gopher - Go Version Manager

  <!-- Uncomment when approved for GitHub Sponsors:
  [![GitHub Sponsors](https://img.shields.io/github/sponsors/molmedoz?style=social)](https://github.com/sponsors/molmedoz)
  -->
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
  [![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://go.dev)
</div>

A simple, fast, and lightweight Go version manager written in Go. Gopher enables you to install, manage, and switch between multiple Go versions, including your system-installed Go.

> 📚 **New to Gopher?** Start with the [**Quick Reference**](QUICK_REFERENCE.md) for a one-page overview, or explore the [**Complete Documentation Index**](DOCUMENTATION_INDEX.md) to find exactly what you need!

## Features

- 🚀 **Fast**: Minimal dependencies, built primarily with Go standard library
- 🔒 **Secure**: Cryptographic verification of downloaded Go binaries
- 🎯 **Simple**: Clean CLI interface with intuitive commands
- 🌍 **Cross-platform**: Works on Linux, macOS, and Windows
- 📦 **Lightweight**: Zero third-party dependencies (only official Go extended packages)
- 🔄 **Auto-cleanup**: Automatically manages old versions to save space
- 🏠 **System Integration**: Seamlessly manages both system and Gopher-installed Go versions
- 📊 **JSON Support**: Full JSON output for scripting and automation
- 🌐 **Environment Management**: Comprehensive GOPATH and GOROOT management
- ⚙️ **Flexible Configuration**: Customizable environment variables and workspace settings
- 🎨 **Visual Progress**: Animated spinners and progress bars with auto-sizing
- 🖥️ **Terminal Aware**: Smart terminal detection and cross-platform output handling
- 🔧 **Verbosity Control**: Command-line flags for detailed output control

## Compatibility

Gopher is designed to work across multiple platforms and Go versions:

| Gopher Version | Supported Go Versions | Platforms | Architecture |
|----------------|----------------------|-----------|--------------|
| v1.0.0+ | 1.20+ | Linux, macOS, Windows | x86, x64, ARM, ARM64 |

**Platform Support:**
- ✅ **Linux**: All major distributions (Ubuntu, Debian, Fedora, CentOS, Arch, etc.)
- ✅ **macOS**: Intel and Apple Silicon (M1/M2/M3/M4)
- ✅ **Windows**: Windows 10/11 with Developer Mode enabled

**Go Version Support:**
- Manages any Go version from 1.20 onwards
- Can detect and manage system-installed Go versions
- Supports pre-release and beta versions

## Dependencies

Gopher is designed to be lightweight with minimal dependencies:

**Runtime Dependencies:**
- ✅ **Zero third-party dependencies** - No external packages from unknown sources
- ✅ **Official Go packages only** - Uses golang.org/x/term (official Go extended library)

**Dependency Details:**
```
golang.org/x/term v0.36.0
  ├── Purpose: Terminal capabilities (TTY detection, width detection)
  ├── Maintained by: Go team (official)
  ├── Size: ~100KB
  └── Why: Better UX (auto-sizing progress bars, reliable cross-platform behavior)
```

**Why golang.org/x/term?**
- Industry standard (used by Docker, kubectl, cobra, and most Go CLI tools)
- Better Windows compatibility than stdlib alternatives
- Automatic terminal width detection (prevents line wrapping)
- Reliable TTY detection (progress bars only show when appropriate)

**Future:** We track "Pure Standard Library Implementation" as a very low priority item in our [roadmap](docs/ROADMAP.md#phase-6-pure-standard-library) for those who prefer zero external dependencies.

## Installation

### Package Managers (Recommended)

#### **macOS:**
```bash
# Using Homebrew (recommended)
brew install molmedoz/tap/gopher

# Or in two steps:
# brew tap molmedoz/tap
# brew install gopher
```

#### **Windows:**
```powershell
# Using Chocolatey
choco install gopher
```

#### **Linux:**
```bash
# Using Snap (all distributions)
sudo snap install gopher --classic

# Or download binary from releases
# See: https://github.com/molmedoz/gopher/releases
```

### Go Install

```bash
go install github.com/molmedoz/gopher/cmd/gopher@latest
```

### From Source

```bash
git clone https://github.com/molmedoz/gopher.git
cd gopher
go build -o gopher cmd/gopher/*.go
sudo mv gopher /usr/local/bin/
```

### Interactive Setup

Run the setup wizard to automatically configure gopher for your platform:

```bash
gopher init
```

This will:
- Detect your platform and current configuration
- Check for required dependencies (Developer Mode on Windows, Homebrew on macOS, etc.)
- Test symlink creation and PATH configuration
- Provide platform-specific next steps

### Platform-Specific Guides

- **Windows:** See [Windows Setup Guide](docs/WINDOWS_SETUP_GUIDE.md) for detailed instructions
- **All Platforms:** See [User Guide](docs/USER_GUIDE.md) for complete documentation

## Usage

### Basic Commands

```bash
# Interactive setup wizard
gopher init

# List installed Go versions (including system)
gopher list

# Install and switch to a Go version
gopher install 1.21.0
gopher use 1.21.0

# Switch to system Go
gopher use system

# Show current Go version
gopher current

# Uninstall a Go version
gopher uninstall 1.20.7

# Show detailed help
gopher help
```

### Environment Management

```bash
# Configure GOPATH mode
gopher env set gopath_mode=shared

# Show environment variables for a version
gopher env show go1.21.0

# Reset to defaults
gopher env reset
```

### Interactive Pagination

Both `list` and `list-remote` commands feature interactive pagination by default:

```bash
# Interactive mode (default)
gopher list
gopher list-remote

# Navigate with:
# - Enter or 'n': Next page
# - 'p': Previous page  
# - Number: Jump to specific page
# - 'q': Quit
# - 'h': Help

# Disable interactive mode
gopher --no-interactive list
gopher --no-interactive list-remote
```

### JSON Output

All commands support JSON output for scripting (flags must come before the command):

```bash
gopher --json list
gopher --json current
gopher --json system
```

### Visual Indicators

The `list` command shows clear visual indicators for the active version:

```bash
# Active version is highlighted with arrow and [active] flag
gopher list
# Output:
#   go1.24.7 (darwin/arm64) [system]
#   go1.22.0 (darwin/arm64)
#   go1.22.1 (darwin/arm64)
# → go1.23.0 (darwin/arm64) [active]
```

### Verbosity Control

Control output verbosity with command-line flags:

```bash
# Default: INFO level
gopher list

# Verbose: DEBUG level (shows detailed information)
gopher --verbose install 1.21.0
gopher -v install 1.21.0

# Quiet: ERROR level only (minimal output)
gopher --quiet list
gopher -q list
```

### GOPATH Modes

Gopher supports three GOPATH modes:

**Shared Mode (Default)**
```bash
gopher env set gopath_mode=shared
# All Go versions share the same GOPATH
```

**Version-Specific Mode**
```bash
gopher env set gopath_mode=version-specific
# Each Go version has its own isolated GOPATH
```

**Custom Mode**
```bash
gopher env set gopath_mode=custom
gopher env set custom_gopath=/path/to/workspace
# Use a custom GOPATH location
```

### Configuration

Gopher stores its configuration in:
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

## Examples

### Basic Workflow

```bash
# Run interactive setup (first time)
gopher init

# Install and switch to Go 1.21.0
gopher install 1.21.0
gopher use 1.21.0

# Verify installation
go version

# List all versions
gopher list

# Switch to system Go
gopher use system
```

### Scripting with JSON

```bash
# Get current version in JSON format
current_version=$(gopher current --json | jq -r '.version')
echo "Current Go version: $current_version"
```

## Configuration

Gopher stores its configuration in:
- **Linux/macOS**: `~/.gopher/config.json`
- **Windows**: `%USERPROFILE%\gopher\config.json`

### Environment Variables

- `GOPHER_CONFIG`: Path to custom configuration file
- `GOPHER_INSTALL_DIR`: Custom installation directory
- `GOPHER_DOWNLOAD_DIR`: Custom download directory

For detailed configuration options, see the [User Guide](docs/USER_GUIDE.md#configuration).

## Development

```bash
# Clone and build
git clone https://github.com/molmedoz/gopher.git
cd gopher
go build -o gopher cmd/gopher/main.go

# Run tests
go test ./...

# Run with custom config
gopher --config /path/to/config.json list
```

## Security

- All downloaded Go binaries are verified using SHA256 checksums
- Downloads are performed over HTTPS from official Go mirrors
- Minimal dependencies (only official Go packages) reduce the attack surface

## Documentation

📚 **[Complete Documentation Index](DOCUMENTATION_INDEX.md)** - Find all documentation quickly!

### Quick Start
- **[Quick Reference](QUICK_REFERENCE.md)** ⚡ - One-page command reference and workflows
- **[FAQ](docs/FAQ.md)** ❓ - Frequently asked questions
- **[User Guide](docs/USER_GUIDE.md)** 📖 - Comprehensive user documentation
- **[Examples](docs/EXAMPLES.md)** 💡 - 50+ practical usage examples

### Platform-Specific
- **[Windows Setup](docs/WINDOWS_SETUP_GUIDE.md)** 🪟 - Windows-specific setup guide
- **[Testing Guide](docs/TESTING_GUIDE.md)** 🧪 - Multi-platform testing strategy

### For Developers
- **[Developer Guide](docs/DEVELOPER_GUIDE.md)** 👨‍💻 - Contributing and development
- **[API Reference](docs/API_REFERENCE.md)** 📚 - Detailed API documentation

## 💖 Support This Project

If Gopher has been helpful to you, consider supporting its development!

<!-- 
When approved for GitHub Sponsors, uncomment this badge:
[![Sponsor Gopher](https://img.shields.io/badge/Sponsor-Gopher-pink?style=for-the-badge&logo=github)](https://github.com/sponsors/molmedoz)
-->

**GitHub Sponsors coming soon!** We're applying to the GitHub Sponsors program to make it easy for you to support Gopher's development.

In the meantime, you can support this project by:
- ⭐ **Starring the repository** on GitHub
- 🐛 **Reporting bugs** and issues
- 💡 **Suggesting features** and improvements
- 📝 **Contributing** code or documentation
- 📢 **Sharing** Gopher with others

Your support (code, feedback, or future sponsorship) helps:
- 🔧 Maintain and improve Gopher
- 🐛 Fix bugs and add new features
- 📚 Create better documentation
- ⏰ Dedicate more time to open source
- 🌟 Build more tools for the Go community

<!-- sponsors --><!-- sponsors -->

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [g](https://github.com/voidint/g) and [GVM](https://github.com/moovweb/gvm)
- Built with Go standard library + official Go extended packages (golang.org/x/term)
- Uses official Go download mirrors
- Thank you to all [contributors and sponsors](SPONSORS.md)!

---

<div align="center">
  <img src="molmedoz_hat.png" alt="molmedoz" width="80"/>
  
  **Created with ❤️ by [molmedoz](https://github.com/molmedoz)**
  
  *GitHub Sponsors coming soon! Star ⭐ the repo to show your support!*
</div>
