# Changelog

All notable changes to Gopher will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

*No unreleased changes yet*

## [v1.0.1] - 2025-11-01

### Added
- **Automatic PATH Validation**
  - Gopher now automatically checks if `GOPATH/bin` is in PATH after version switches
  - Provides helpful warnings with platform-specific fix instructions if missing
  - Supports Windows, macOS (zsh/bash), and Linux with OS-appropriate commands
  - Helps users discover when installed Go tools won't be accessible from command line
  - See [USER_GUIDE.md](docs/USER_GUIDE.md#gopathbin-not-in-path-warning) for details

### Security
- **Comprehensive Security Hardening**
  - Fixed all HIGH severity gosec findings (G115 integer overflow in archive extraction)
  - Fixed all MEDIUM severity findings:
    - G110: Added decompression bomb protection (1GB per-file limits)
    - G204: Secured subprocess execution with `runGoCommand` helper
    - G304: Added path validation for file operations
    - G301/G306: Reviewed and justified file/directory permissions
  - Fixed all LOW severity findings (G104: Handled all unhandled errors)
  - Added security validation helpers in `internal/security/`
  - All security fixes documented in [SECURITY.md](SECURITY.md)

### Fixed
- **Release Process Improvements**
  - Fixed archive format configuration (zip for Windows, tar.gz for Unix)
  - Fixed GoReleaser v2 `--skip` flags to use correct singular forms (homebrew, chocolatey, scoop)
  - Added tag rollback safety: tags are only pushed after successful GitHub release
  - Prevents orphaned tags if release process fails
  - Improved release workflow reliability and error handling

### Changed
- **Release Workflow**
  - Tags are created locally first, only pushed after successful GitHub release
  - Automatic cleanup of tags if GitHub release fails
  - Better error handling and rollback mechanisms
  - Enhanced release documentation and troubleshooting guides

## [v1.0.0] - 2025-10-15

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
  - Makefile-based local CI (`make ci`, `make test`, `make lint`)
  - Shell-based E2E tests (19 scenarios across platforms)
  - Comprehensive documentation suite
  - Detailed roadmap with prioritized features
  
- **Maintenance Commands**
  - `gopher clean` - Remove download cache to free disk space
  - `gopher purge` - Complete removal of all Gopher data with confirmation
  
- **Development Tools**
  - `make install-tools` - Automatically install goimports, golangci-lint, and gofumpt
  - `make check-format` - Check code format without modifying (CI mode)
  - `make check-imports` - Check imports without modifying (CI mode)
  - `make ci` - Run complete CI workflow locally (matches GitHub Actions)
  - Automatic tool path resolution (no PATH configuration needed)
  
### Distribution
- **Homebrew**: Available via `brew install molmedoz/tap/gopher` (macOS/Linux)
- **Linux Packages**: Debian (.deb), RedHat (.rpm), Alpine (.apk), Arch (.pkg.tar.zst)
- **GitHub Releases**: Pre-built binaries for all platforms (Linux, macOS, Windows, ARM)
- **Direct Download**: All platforms can download from GitHub Releases
- **Chocolatey**: Coming in v1.0.1+ (requires Windows build runner)
- **Snap**: Pending Snapcraft name approval (will be added in future release)

### Technical Improvements
- **Testing**: All tests include `-race` flag by default for race condition detection
- **Makefile**: Tools use full paths from GOPATH/bin (no PATH setup required)
- **CI/CD**: Version injection added to all GitHub Actions builds
- **CI/CD**: Race detection added to all test workflows
- **CI/Local Parity**: `make ci` matches GitHub Actions exactly
- **Dependencies**: Uses `golang.org/x/term` for enhanced terminal support
- **Testing**: 50+ unit tests with 35% coverage + 19 E2E test scenarios
- **CI/CD**: GitHub Actions workflows for lint, test, build, security, coverage
- **Architecture**: Modular design with reusable components
- **Thread Safety**: Mutex-protected concurrent operations in alias manager
- **Cross-platform**: Platform-specific optimizations for Windows, macOS, and Linux
- **Security**: Security scanning with gosec, vulnerability checks with govulncheck

### Commands Available
- `gopher list` - List installed Go versions (with aliases shown inline)
- `gopher list-remote` - List available Go versions with pagination
- `gopher install <version>` - Install a Go version
- `gopher uninstall <version>` - Uninstall a Go version
- `gopher use <version>` - Switch to a Go version
- `gopher use system` - Switch to system Go
- `gopher current` - Show current Go version
- `gopher system` - Show system Go information
- `gopher clean` - Remove download cache to free disk space
- `gopher purge` - Complete removal with confirmation
- `gopher alias` - Manage version aliases (create, list, remove, show, import, export, suggest, bulk)
- `gopher init` - Interactive setup wizard
- `gopher setup` - Shell integration setup
- `gopher status` - Show persistence and shell integration status
- `gopher env` - Environment configuration management (list, show, set, reset)
- `gopher debug` - Debug information for troubleshooting
- `gopher version` - Show gopher version

### User Experience Improvements
- Visual progress bars for downloads (auto-sized to terminal width)
- Animated spinners for long operations (extraction, uninstall)
- Aliases shown inline in `gopher list` output for quick reference
- Version normalization (omit 'go' prefix in all commands)
- Graceful handling of missing/corrupted alias files
- Clean, professional terminal output
- No interference between progress indicators and other output
- Better Windows terminal compatibility

### Fixed
- Windows progress bar issues (progress bars now update on same line)
- Windows spinner compatibility (animations work correctly)
- Terminal width detection prevents line wrapping
- TTY detection for better piped output handling
- Race conditions in concurrent alias operations (added mutex protection)
- Alias file format migration (old array to map format)
- Empty/missing alias file handling
- Flag parsing order in CLI
- JSON output for empty lists
- Test panic recovery output in CI

### Known Issues
None reported. Please report any issues on [GitHub](https://github.com/molmedoz/gopher/issues).
