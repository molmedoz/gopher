# Gopher User Guide

This comprehensive guide covers everything you need to know about using Gopher, the Go version manager.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [System Go Management](#system-go-management)
- [Environment Management](#environment-management)
- [Configuration](#configuration)
- [Scripting and Automation](#scripting-and-automation)
- [Visual Indicators](#visual-indicators)
- [Verbosity Control](#verbosity-control)
- [Troubleshooting](#troubleshooting)
- [Advanced Usage](#advanced-usage)

## Installation

### Prerequisites

- **Windows:** Windows 10+ with PowerShell
- **macOS/Linux:** Standard shell (bash, zsh, fish)
- **All platforms:** Internet connection for downloading Go versions

### Installation Methods

#### Method 1: Package Managers (Recommended)

**macOS (Homebrew):**
```bash
brew install molmedoz/tap/gopher
```

**Windows (Chocolatey):**
```powershell
choco install gopher
```

**Linux (Snap):**
```bash
snap install gopher
```

#### Method 2: From GitHub Releases

**Windows:**
```powershell
# Download and install
$version = "1.0.0"
New-Item -ItemType Directory -Path "$env:USERPROFILE\bin" -Force
$url = "https://github.com/molmedoz/gopher/releases/download/v$version/gopher-windows-amd64.exe"
Invoke-WebRequest -Uri $url -OutFile "$env:USERPROFILE\bin\gopher.exe"

# Add to PATH
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$env:USERPROFILE\bin;$userPath", "User")

# Restart PowerShell
```

**macOS/Linux:**
```bash
# Download for your platform
curl -LO https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
chmod +x gopher-*
sudo mv gopher-* /usr/local/bin/gopher
```

#### Method 3: Using Go Install

```bash
go install github.com/molmedoz/gopher/cmd/gopher@latest
```

#### Method 4: From Source

```bash
git clone https://github.com/molmedoz/gopher.git
cd gopher
make build
sudo make install
```

### Verify Installation

```bash
gopher version
```

Expected output:
```
gopher v1.0.0
```

### First-Time Setup

**Windows:**
```powershell
gopher init
```
Then follow the on-screen instructions. See [Windows Setup Guide](WINDOWS_SETUP_GUIDE.md) for details.

**macOS/Linux:**
```bash
gopher init
```

**What `gopher init` does:**
- ✅ Detects your environment (shell, paths, etc.)
- ✅ Creates all required directories automatically
- ✅ Tests symlink creation
- ✅ Shows platform-specific setup instructions

**Directories created automatically:**
- **Windows:** `C:\Users\YourName\gopher\{versions,downloads,state}`
- **macOS/Linux:** `$HOME/.gopher/{versions,downloads,state}`

## Quick Start

### Workflow Overview

Here's a visual overview of common Gopher workflows:

```
┌─────────────────────────────────────────────────────────────────────┐
│                     GOPHER WORKFLOW DIAGRAM                         │
└─────────────────────────────────────────────────────────────────────┘

First Time Setup:
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Install  │ -> │  Check   │ -> │ Install  │ -> │   Use    │
│  Gopher  │    │  System  │    │ Go Ver.  │    │ Go Ver.  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
   (once)      gopher system   gopher install  gopher use

Regular Usage:
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│   List   │ -> │ Install  │ -> │  Switch  │ -> │  Verify  │
│ Versions │    │  Version │    │  Version │    │  Active  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
gopher list   gopher install  gopher use     go version

Cleanup:
┌──────────┐    ┌──────────┐    ┌──────────┐
│   List   │ -> │   Find   │ -> │ Uninstall│
│ Installed│    │  Unused  │    │  Version │
└──────────┘    └──────────┘    └──────────┘
gopher list                    gopher uninstall
```

### Version Switching Flow

```
Current State                Action                  New State
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   System    │         │   gopher    │         │   go1.21.0  │
│  go1.20.0   │  --->   │ use 1.21.0  │  --->   │   [active]  │
│  [active]   │         └─────────────┘         └─────────────┘
└─────────────┘                                        │
      ↑                                                │
      │         ┌─────────────┐                        │
      └─────────│   gopher    │  <─────────────────────┘
                │ use system  │
                └─────────────┘
```

### 1. Check Your Current Setup

```bash
# See what Go versions you have
gopher list

# Check your current Go version
gopher current

# View system Go information
gopher system
```

### 2. Install a New Go Version

```bash
# Install Go 1.21.0
gopher install 1.21.0

# Verify installation
gopher list
```

### 3. Switch Between Versions

```bash
# Switch to a specific version
gopher use 1.21.0

# Switch back to system Go
gopher use system

# Verify the switch
go version
```

### 4. Clean Up Old Versions

```bash
# Remove a version you no longer need
gopher uninstall 1.20.7
```

## Command Reference

### `gopher list`

Lists all installed Go versions, including system versions. **Interactive pagination is enabled by default.**

```bash
gopher list
```

**Interactive Mode (default):**
- Press **Enter** or **n** to go to the next page
- Press **p** to go to the previous page
- Press **q** to quit
- Type a page number to jump to that page
- Press **h** for help

**Output:**
```
Installed Go versions (page 1 of 1, showing 2 of 2 total):

  go1.25.1 (darwin/arm64) - active (installed: 2025-08-27 08:49:40) [system]
  go1.21.0 (darwin/arm64) - inactive (installed: 2025-08-27 09:15:22)

Page 1 of 1
Commands:
  n, next, →     - Next page
  p, prev, ←     - Previous page
  <number>       - Go to specific page
  q, quit, exit  - Quit
  h, help        - Show this help

Enter command:
```

**Options:**
- `--json`: Output in JSON format (disables interactive mode)
- `--no-interactive`: Disable interactive pagination
- `--page-size <number>`: Number of versions per page (default: 10)
- `--page <number>`: Page number to display (default: 1)

**Note:** Flags must be placed **before** the command name.

**Examples:**
```bash
# Interactive mode (default)
gopher list

# Disable interactive mode
gopher --no-interactive list

# JSON output
gopher --json list

# Change page size
gopher --page-size 5 list

# Go to specific page (non-interactive)
gopher --page 2 --no-interactive list
```

**JSON Output:**
```bash
gopher list --json
```

```json
[
  {
    "version": "go1.25.1",
    "os": "darwin",
    "arch": "arm64",
    "installed_at": "2025-08-27T08:49:40-07:00",
    "is_active": true,
    "is_system": true,
    "path": "/opt/homebrew/opt/go/libexec/bin/go"
  }
]
```

### `gopher list-remote`

Lists available Go versions for installation. **Interactive pagination is enabled by default.**

```bash
gopher list-remote
```

**Interactive Mode (default):**
- Press **Enter** or **n** to go to the next page
- Press **p** to go to the previous page
- Press **q** to quit
- Type a page number to jump to that page
- Press **h** for help

**Output:**
```
Available Go versions (page 1 of 33, showing 10 of 328 total):

  go1.25.1 (stable) - 
  go1.25.0 (stable) - 
  go1.24.7 (stable) - 
  ...

Page 1 of 33
Commands:
  n, next, →     - Next page
  p, prev, ←     - Previous page
  <number>       - Go to specific page
  q, quit, exit  - Quit
  h, help        - Show this help

Enter command:
```

**Options:**
- `--page-size <number>`: Number of versions per page (default: 10)
- `--page <number>`: Page number to display (default: 1)
- `--filter <text>`: Filter versions by text (e.g., '1.21', 'stable', 'rc')
- `--stable`: Show only stable versions
- `--no-interactive`: Disable interactive pagination
- `--json`: Output in JSON format (disables interactive mode)

**Note:** Flags must be placed **before** the command name.

**Examples:**
```bash
# Interactive listing (default)
gopher list-remote

# Disable interactive mode
gopher --no-interactive list-remote

# Pagination
gopher --page-size 5 list-remote
gopher --page 2 --page-size 10 --no-interactive list-remote

# Filtering
gopher --filter "1.21" list-remote
gopher --filter "rc" list-remote
gopher --stable list-remote

# JSON output
gopher --json list-remote
```

### `gopher install <version>`

Downloads and installs a specific Go version.

```bash
gopher install 1.21.0
```

**Examples:**
```bash
# Install latest stable
gopher install 1.21.0

# Install with go prefix
gopher install go1.21.0

# Install specific patch version
gopher install 1.21.1
```

**What happens during installation:**
1. Validates version format
2. Checks if already installed
3. Downloads from official Go mirrors
4. Verifies SHA256 checksums
5. Extracts to gopher directory
6. Creates version metadata
7. Cleans up downloaded files

### `gopher uninstall <version>`

Removes a Go version installed by gopher.

```bash
gopher uninstall 1.21.0
```

**Note:** Cannot uninstall system Go versions.

### `gopher use <version>`

Switches to a specific Go version.

```bash
gopher use 1.21.0
gopher use system
```

**Special versions:**
- `system` or `sys`: Switch to system Go

**What happens during switching:**
1. Validates version exists
2. Creates symlink to version's go binary
3. Updates PATH (requires appropriate permissions)

### `gopher current`

Shows the currently active Go version.

```bash
gopher current
```

**Output:**
```
Current Go version: go1.21.0 (darwin/arm64)
```

### `gopher system`

Shows detailed information about system Go.

```bash
gopher system
```

**Output:**
```
System Go Information:
  System Go: go version go1.25.1 darwin/arm64
  Path: /opt/homebrew/opt/go/libexec/bin/go
  GOROOT: /opt/homebrew/opt/go/libexec
  GOPATH: /Users/molmedoz/go
  Installed: 2025-08-27 08:49:40
  System: true
```

### `gopher version`

Shows gopher version information.

```bash
gopher version
```

### `gopher init`

Runs the interactive setup wizard for platform-specific configuration.

```bash
gopher init
```

**What it does:**
1. Detects your operating system and architecture
2. Prompts for installation preferences
3. Sets up shell integration
4. Configures environment variables
5. Creates initial configuration

**Example output:**
```
🚀 Welcome to Gopher Setup Wizard!

Detected system: darwin/arm64
Install directory: /Users/username/.gopher/versions
Download directory: /Users/username/.gopher/downloads

Setting up shell integration...
✓ Shell integration configured for zsh
✓ Configuration saved
```

### `gopher setup`

Sets up shell integration for persistent Go version switching.

```bash
gopher setup
```

**What it does:**
1. Detects your shell (bash, zsh, fish)
2. Adds gopher integration to shell profile
3. Enables automatic version switching
4. Sets up environment variables

**Supported shells:**
- Bash (`.bashrc`, `.bash_profile`)
- Zsh (`.zshrc`)
- Fish (`config.fish`)

### `gopher status`

Shows persistence status and shell integration information.

```bash
gopher status
```

**Output:**
```
Gopher Status:
  Shell Integration: ✓ Enabled (zsh)
  Profile File: /Users/username/.zshrc
  Current Version: go1.21.0
  Active Version: go1.21.0
  System Go: Available
```

### `gopher debug`

Shows debug information for troubleshooting.

```bash
gopher debug
```

**Output:**
```
Debug Information:
  Gopher Version: v1.0.0
  Go Version: go1.21.0
  OS: darwin
  Arch: arm64
  Install Dir: /Users/username/.gopher/versions
  Download Dir: /Users/username/.gopher/downloads
  Config File: /Users/username/.gopher/config.json
  Shell: zsh
  PATH: /Users/username/.gopher/versions/go1.21.0/bin:/usr/local/bin:...
```

**Useful for:**
- Troubleshooting installation issues
- Debugging shell integration problems
- Checking configuration settings
- Reporting bugs with system information

## System Go Management

Gopher automatically detects and manages your system-installed Go versions alongside Gopher-managed versions.

### Supported System Installations

Gopher recognizes these system Go installations:

- **macOS Homebrew**: `/opt/homebrew/opt/go/libexec/bin/go` (Apple Silicon)
- **macOS Homebrew**: `/usr/local/opt/go/libexec/bin/go` (Intel)
- **System packages**: `/usr/bin/go`, `/usr/local/bin/go`
- **Manual installations**: `/usr/local/go/bin/go`
- **Windows**: `C:\Program Files\Go\bin\go.exe`, `C:\Go\bin\go.exe`

### System Go Features

- **Automatic Detection**: Gopher finds your system Go automatically
- **Seamless Switching**: Switch between system and gopher versions easily
- **Version Information**: Get detailed info about your system Go installation
- **JSON Support**: System Go data available in JSON format

### Working with System Go

```bash
# Check if system Go is available
gopher system

# Switch to system Go
gopher use system

# List all versions (including system)
gopher list

# Get system Go info in JSON
gopher system --json
```

## Environment Management

Gopher provides comprehensive environment variable management to ensure proper Go development workflows. This section covers GOPATH and GOROOT management, environment configuration, and automatic script generation.

### GOPATH Management Modes

Gopher supports three different GOPATH management strategies to suit different development needs:

#### 1. Shared Mode (Default)

**Best for**: Most users, shared projects, simple setups

All Go versions share the same GOPATH directory, allowing packages to be shared across versions.

```bash
# Set shared mode (default)
gopher env set gopath_mode=shared

# Check current configuration
gopher env list

# Switch to a Go version
gopher use 1.21.0
# GOPATH=/home/user/go (shared across all versions)
```

**Advantages:**
- Simple and intuitive
- Packages shared across Go versions
- Smaller disk usage
- Matches most users' expectations

**Disadvantages:**
- Potential package compatibility issues
- Mixed dependencies from different Go versions

#### 2. Version-Specific Mode

**Best for**: Development work, testing, complete isolation

Each Go version has its own isolated GOPATH directory.

```bash
# Set version-specific mode
gopher env set gopath_mode=version-specific

# Switch to different versions
gopher use 1.21.0
# GOPATH=/home/user/.gopher/versions/go1.21.0/gopath

gopher use 1.22.0
# GOPATH=/home/user/.gopher/versions/go1.22.0/gopath
```

**Advantages:**
- Complete isolation between Go versions
- No package conflicts
- Safe to experiment with different Go versions
- Clean separation of dependencies

**Disadvantages:**
- Larger disk usage
- Need to reinstall packages for each version
- More complex setup

#### 3. Custom Mode

**Best for**: Specific project requirements, custom workspace setups

Use a custom GOPATH location specified by the user.

```bash
# Set custom mode
gopher env set gopath_mode=custom
gopher env set custom_gopath=/path/to/your/workspace

# Switch to a Go version
gopher use 1.21.0
# GOPATH=/path/to/your/workspace
```

**Advantages:**
- Full control over GOPATH location
- Can use existing project directories
- Flexible workspace management

**Disadvantages:**
- Manual configuration required
- Need to manage workspace location

### Environment Configuration

#### Viewing Configuration

```bash
# List all configuration options
gopher env list

# Show environment variables for a specific version
gopher env show go1.21.0

# Show environment variables in JSON format
gopher env show go1.21.0 --json
```

#### Setting Configuration

```bash
# Set GOPATH mode
gopher env set gopath_mode=version-specific

# Set custom GOPATH (when mode is custom)
gopher env set custom_gopath=/path/to/workspace

# Configure Go proxy
gopher env set goproxy=https://proxy.golang.org,direct

# Configure checksum database
gopher env set gosumdb=sum.golang.org

# Enable/disable environment variable setting
gopher env set set_environment=true
```

#### Resetting Configuration

```bash
# Reset all configuration to defaults
gopher env reset
```

### Automatic Environment Scripts

When switching Go versions, Gopher automatically generates environment activation scripts that set up all necessary environment variables.

#### Script Generation

```bash
# Switch to a Go version (generates script automatically)
gopher use 1.21.0

# Output:
# ✓ Environment variables configured for Go go1.21.0
#   To activate this environment, run:
#   source /home/user/.gopher/scripts/go-go1.21.0.env
#   Or add the following to your shell profile:
#   export GOROOT=/home/user/.gopher/versions/go1.21.0
#   export GOPATH=/home/user/go
#   export GOPROXY=https://proxy.golang.org,direct
#   export GOSUMDB=sum.golang.org
#   export PATH=/home/user/.gopher/versions/go1.21.0/bin:...
```

#### Script Location

Environment scripts are stored in `~/.gopher/scripts/`:

```
~/.gopher/scripts/
├── go-go1.21.0.env
├── go-go1.22.0.env
└── go-system.env
```

#### Activating Environment

**Method 1: Source the script**
```bash
source ~/.gopher/scripts/go-go1.21.0.env
```

**Method 2: Add to shell profile**
```bash
# Add to ~/.bashrc, ~/.zshrc, etc.
echo 'source ~/.gopher/scripts/go-go1.21.0.env' >> ~/.bashrc
```

**Method 3: Manual export**
```bash
export GOROOT=/home/user/.gopher/versions/go1.21.0
export GOPATH=/home/user/go
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
export PATH=/home/user/.gopher/versions/go1.21.0/bin:$PATH
```

### Environment Variables

Gopher manages the following environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `GOROOT` | Go installation directory | `/home/user/.gopher/versions/go1.21.0` |
| `GOPATH` | Go workspace directory | `/home/user/go` (shared) or `/home/user/.gopher/versions/go1.21.0/gopath` (version-specific) |
| `GOPROXY` | Go module proxy | `https://proxy.golang.org,direct` |
| `GOSUMDB` | Go checksum database | `sum.golang.org` |
| `PATH` | System PATH with Go binary | `/home/user/.gopher/versions/go1.21.0/bin:...` |

### Workflow Examples

#### Example 1: Shared Development Setup

```bash
# Use shared GOPATH for common development
gopher env set gopath_mode=shared
gopher env set set_environment=true

# Install and switch to Go 1.21.0
gopher install 1.21.0
gopher use 1.21.0

# Activate environment
source ~/.gopher/scripts/go-go1.21.0.env

# Verify setup
go version
echo $GOPATH  # /home/user/go
```

#### Example 2: Isolated Testing Setup

```bash
# Use version-specific GOPATH for testing
gopher env set gopath_mode=version-specific
gopher env set set_environment=true

# Install multiple Go versions
gopher install 1.21.0
gopher install 1.22.0

# Test with Go 1.21.0
gopher use 1.21.0
source ~/.gopher/scripts/go-go1.21.0.env
go mod init test-project
go get github.com/gin-gonic/gin

# Test with Go 1.22.0 (isolated environment)
gopher use 1.22.0
source ~/.gopher/scripts/go-go1.22.0.env
go mod init test-project  # Fresh start, no shared packages
```

#### Example 3: Project-Specific Setup

```bash
# Use custom GOPATH for specific project
gopher env set gopath_mode=custom
gopher env set custom_gopath=/path/to/my-project
gopher env set set_environment=true

# Switch to Go version
gopher use 1.21.0

# Activate environment
source ~/.gopher/scripts/go-go1.21.0.env

# Verify custom GOPATH
echo $GOPATH  # /path/to/my-project
```

### Troubleshooting Environment Issues

#### Environment Not Activated

**Problem**: Go commands not found after switching versions

**Solution**:
```bash
# Check if environment script exists
ls -la ~/.gopher/scripts/

# Activate environment manually
source ~/.gopher/scripts/go-go1.21.0.env

# Or recreate the environment
gopher use 1.21.0
```

#### Wrong GOPATH

**Problem**: GOPATH not set correctly

**Solution**:
```bash
# Check current configuration
gopher env list

# Check environment variables for version
gopher env show go1.21.0

# Reset configuration if needed
gopher env reset
```

#### Script Permission Issues

**Problem**: Cannot source environment script

**Solution**:
```bash
# Check script permissions
ls -la ~/.gopher/scripts/go-go1.21.0.env

# Fix permissions if needed
chmod +x ~/.gopher/scripts/go-go1.21.0.env

# Or use manual export
gopher env show go1.21.0
# Copy and paste the export commands
```

## Configuration

### Configuration File

Gopher stores its configuration in:

- **Linux/macOS**: `~/.gopher/config.json`
- **Windows**: `%USERPROFILE%\gopher\config.json`

### Default Configuration

```json
{
  "install_dir": "~/.gopher/versions",
  "download_dir": "~/.gopher/downloads",
  "mirror_url": "https://go.dev/dl/",
  "auto_cleanup": true,
  "max_versions": 5
}
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `install_dir` | Directory for Go versions | `~/.gopher/versions` |
| `download_dir` | Temporary download directory | `~/.gopher/downloads` |
| `mirror_url` | Go download mirror URL | `https://go.dev/dl/` |
| `auto_cleanup` | Auto-remove old versions | `true` |
| `max_versions` | Maximum versions to keep | `5` |

### Custom Configuration

```bash
# Use custom configuration file
gopher --config /path/to/config.json list

# Environment variables
export GOPHER_CONFIG=/path/to/config.json
export GOPHER_INSTALL_DIR=/opt/go-versions
export GOPHER_DOWNLOAD_DIR=/tmp/gopher-downloads
```

### Creating Custom Config

```bash
# Create config directory
mkdir -p ~/.gopher

# Create custom config
cat > ~/.gopher/config.json << EOF
{
  "install_dir": "/opt/go-versions",
  "download_dir": "/tmp/gopher-downloads",
  "mirror_url": "https://go.dev/dl/",
  "auto_cleanup": true,
  "max_versions": 10
}
EOF
```

## Scripting and Automation

### JSON Output

All commands support JSON output for scripting:

```bash
# Get current version in JSON
current_version=$(gopher current --json | jq -r '.version')
echo "Current Go version: $current_version"

# List versions as JSON
gopher list --json | jq '.[] | select(.is_active) | .version'

# Get system Go path
system_path=$(gopher system --json | jq -r '.path')
echo "System Go path: $system_path"
```

### Shell Integration

Add to your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
# Auto-switch Go versions based on directory
function cd() {
    builtin cd "$@"
    if [[ -f .gopher-version ]]; then
        gopher use $(cat .gopher-version)
    fi
}

# Go version prompt
export PS1='$(gopher current --json | jq -r ".version") $ '
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Setup Go with Gopher
  run: |
    go install github.com/molmedoz/gopher/cmd/gopher@latest
    gopher install 1.21.0
    gopher use 1.21.0
    echo "$(gopher current --json | jq -r '.path')" >> $GITHUB_PATH
```

### Project-Specific Versions

Create `.gopher-version` files in your projects:

```bash
# In your project directory
echo "1.21.0" > .gopher-version

# Auto-switch script
#!/bin/bash
if [[ -f .gopher-version ]]; then
    version=$(cat .gopher-version)
    if gopher list --json | jq -e ".[] | select(.version == \"go$version\")" > /dev/null; then
        gopher use $version
    else
        echo "Installing Go $version..."
        gopher install $version
        gopher use $version
    fi
fi
```

## Visual Indicators

Gopher provides clear visual indicators to help you identify the active Go version and distinguish between different types of versions.

### Active Version Indicators

The `list` command shows clear visual indicators for the currently active version:

```bash
gopher list
```

**Output with indicators:**
```
  go1.24.7 (darwin/arm64) [system]
  go1.22.0 (darwin/arm64)
  go1.22.1 (darwin/arm64)
→ go1.23.0 (darwin/arm64) [active]
```

**Indicator meanings:**
- `→` - Arrow pointing to the currently active version
- `[active]` - Text label indicating the active version
- `[system]` - Text label indicating system-installed Go
- Color coding (when supported):
  - **Green + Bold**: Active version
  - **Cyan**: System version
  - **Dim**: Inactive versions

### Color Support

Gopher automatically detects terminal color support and provides:
- **Automatic detection**: Works with modern terminals
- **Fallback support**: Plain text when colors aren't supported
- **Cross-platform**: Works on Linux, macOS, and Windows

### Disabling Visual Indicators

If you prefer plain text output:

```bash
# Use JSON output for scripting
gopher --json list

# Or redirect to a file
gopher list > versions.txt
```

## Verbosity Control

Gopher provides flexible verbosity control through command-line flags to help you get the right amount of information for your needs.

### Log Levels

Gopher supports multiple log levels:

- **ERROR**: Only error messages (--quiet flag)
- **INFO**: General information (default)
- **DEBUG**: Detailed debugging information (--verbose flag)

### Using Verbosity Flags

#### Quiet Mode (ERROR level only)

```bash
# Minimal output - only errors
gopher --quiet list
gopher -q install 1.21.0

# Useful for scripting and automation
current_version=$(gopher -q current --json | jq -r '.version')
```

#### Verbose Mode (DEBUG level)

```bash
# Detailed output with debugging information
gopher --verbose install 1.21.0
gopher -v list

# Shows detailed progress and internal operations
gopher -v install 1.21.0
# Output:
# [DEBUG] Starting Go installation {version=go1.21.0 file=/tmp/go1.21.0.darwin-arm64.tar.gz}
# [INFO] Downloading file {file=go1.21.0.darwin-arm64.tar.gz}
# [DEBUG] Extracting files {count=1000}
# [INFO] Successfully installed Go version {version=go1.21.0}
```

#### Default Mode (INFO level)

```bash
# Standard output - general information
gopher list
gopher install 1.21.0

# Shows important information without excessive detail
```

### Verbosity in Different Commands

#### Installation Commands

```bash
# Quiet installation (minimal output)
gopher -q install 1.21.0

# Verbose installation (detailed progress)
gopher -v install 1.21.0

# Default installation (standard progress)
gopher install 1.21.0
```

#### Listing Commands

```bash
# Quiet listing (just the list)
gopher -q list

# Verbose listing (with metadata)
gopher -v list

# Default listing (standard format)
gopher list
```

#### System Commands

```bash
# Quiet system info
gopher -q system

# Verbose system detection
gopher -v system

# Default system info
gopher system
```

### Combining with Other Flags

Verbosity flags work with other command-line options:

```bash
# Verbose JSON output
gopher -v --json list

# Quiet installation with custom config
gopher -q --config /path/to/config.json install 1.21.0

# Verbose alias creation
gopher -v alias create stable 1.21.0
```

### Use Cases

#### Development and Debugging

```bash
# Debug installation issues
gopher -v install 1.21.0

# Debug system detection
gopher -v system

# Debug alias operations
gopher -v alias create stable 1.21.0
```

#### Scripting and Automation

```bash
# Quiet mode for scripts
gopher -q current --json

# Error-only output for error handling
gopher -q install 1.21.0 || echo "Installation failed"
```

#### User-Friendly Output

```bash
# Standard output for interactive use
gopher list
gopher install 1.21.0
```

## Troubleshooting

### Troubleshooting Decision Tree

Use this decision tree to quickly identify and resolve common issues:

```
┌─────────────────────────────────────────────────────────────────┐
│                GOPHER TROUBLESHOOTING DECISION TREE             │
└─────────────────────────────────────────────────────────────────┘

START: What's the problem?
│
├─ [Gopher command not found]
│  │
│  ├─ Check if installed: which gopher
│  │  ├─ Not found → Reinstall (see Installation section)
│  │  └─ Found → Check PATH configuration
│  │             └─ Add to PATH: export PATH=$PATH:/path/to/gopher
│  │
│  └─ Still not working?
│     └─ Check shell config (~/.bashrc or ~/.zshrc)
│        └─ Reload: source ~/.bashrc
│
├─ [Version not switching / Wrong version active]
│  │
│  ├─ Verify switch: gopher current
│  │  └─ Check actual: go version
│  │
│  ├─ Not matching?
│  │  ├─ Check symlink: ls -la $(which go)
│  │  ├─ Reload shell: source ~/.bashrc
│  │  └─ Try switching again: gopher use <version>
│  │
│  └─ Still wrong?
│     └─ Check for multiple Go installations
│        └─ which -a go  # Shows all Go binaries in PATH
│
├─ [Installation failed / Download issues]
│  │
│  ├─ Check internet connection
│  │  └─ Test: curl -I https://go.dev/dl/
│  │
│  ├─ Try different mirror
│  │  └─ gopher env set mirror_url=https://golang.google.cn/dl/
│  │
│  ├─ Check disk space
│  │  └─ df -h ~/.gopher
│  │
│  └─ Enable verbose logging
│     └─ gopher --verbose install <version>
│
├─ [Permission denied errors]
│  │
│  ├─ Linux/macOS
│  │  ├─ Check permissions: ls -la ~/.gopher
│  │  ├─ Fix if needed: chmod -R u+w ~/.gopher
│  │  └─ May need sudo: sudo gopher use <version>
│  │
│  └─ Windows
│     ├─ Check Developer Mode: Settings → Privacy & Security
│     ├─ Enable Developer Mode (required for symlinks)
│     └─ Run as Administrator (if needed)
│
├─ [Go programs not finding packages]
│  │
│  ├─ Check GOPATH: echo $GOPATH
│  │  └─ Not set?
│  │     └─ gopher env show <version>
│  │        └─ Copy and run export commands
│  │
│  ├─ Check GOROOT: echo $GOROOT
│  │  └─ Should point to: ~/.gopher/versions/<version>
│  │
│  └─ Reset environment
│     └─ gopher env reset
│        └─ gopher use <version>
│
├─ [Slow downloads]
│  │
│  ├─ Try different mirror
│  │  └─ China: gopher env set mirror_url=https://golang.google.cn/dl/
│  │
│  ├─ Check proxy settings
│  │  └─ echo $GOPROXY
│  │
│  └─ Use direct connection
│     └─ gopher env set goproxy=direct
│
└─ [Other issues]
   │
   ├─ Enable debug logging
   │  └─ gopher --verbose <command>
   │
   ├─ Check system status
   │  └─ gopher system
   │
   ├─ Verify installation
   │  └─ gopher list
   │     └─ gopher current
   │
   └─ Still stuck?
      ├─ Check FAQ: docs/FAQ.md
      ├─ Search Issues: github.com/molmedoz/gopher/issues
      └─ Create New Issue with debug output
```

### Quick Diagnostic Commands

Run these commands to gather diagnostic information:

```bash
# Check Gopher installation
which gopher
gopher version

# Check current Go setup
which go
go version
gopher current
gopher system

# Check environment
echo $PATH
echo $GOROOT
echo $GOPATH
gopher env list

# List all installed versions
gopher list

# Check for multiple Go installations
which -a go   # Linux/macOS
where go      # Windows
```

### Common Issues

#### Permission Denied When Switching Versions

**Problem:**
```
Error: failed to create symlink: permission denied
```

**Solution:**
```bash
# On macOS/Linux, you may need sudo for system-wide symlinks
sudo gopher use 1.21.0

# Or use a local installation directory
export GOPHER_INSTALL_DIR=~/.local/go-versions
gopher use 1.21.0
```

#### Go Not Found After Switching

**Problem:**
```
go: command not found
```

**Solution:**
```bash
# Check if symlink was created
ls -la /usr/local/bin/go

# Recreate symlink
gopher use 1.21.0

# Check PATH
echo $PATH
```

#### System Go Not Detected

**Problem:**
System Go not showing in `gopher list`

**Solution:**
```bash
# Check if Go is in PATH
which go

# Check system detection
gopher system

# Verify Go installation
go version
```

#### Download Failures

**Problem:**
```
Error: failed to download version
```

**Solution:**
```bash
# Check internet connection
ping go.dev

# Try different mirror
gopher --config <(echo '{"mirror_url": "https://golang.org/dl/"}') install 1.21.0

# Check disk space
df -h
```

### Debug Mode

Enable verbose output for debugging:

```bash
# Set debug environment variable
export GOPHER_DEBUG=1

# Run commands with debug output
gopher list
gopher install 1.21.0
```

### Log Files

Gopher logs to:
- **Linux/macOS**: `~/.gopher/gopher.log`
- **Windows**: `%USERPROFILE%\gopher\gopher.log`

### Reset Gopher

If gopher gets into a bad state:

```bash
# Remove gopher directory
rm -rf ~/.gopher

# Reinstall
go install github.com/molmedoz/gopher/cmd/gopher@latest

# Reconfigure
gopher list
```

## Advanced Usage

### Custom Go Mirrors

```bash
# Use different mirror
gopher --config <(echo '{"mirror_url": "https://golang.org/dl/"}') install 1.21.0

# Corporate mirror
gopher --config <(echo '{"mirror_url": "https://internal-mirror.company.com/go/"}') install 1.21.0
```

### Multiple Gopher Instances

```bash
# Different configs for different projects
gopher --config ~/.gopher/project1.json use 1.21.0
gopher --config ~/.gopher/project2.json use 1.20.7
```

### Integration with Other Tools

#### VS Code

Add to VS Code settings:

```json
{
  "go.goroot": "~/.gopher/versions/go1.21.0",
  "go.toolsEnvVars": {
    "GOROOT": "~/.gopher/versions/go1.21.0"
  }
}
```

#### Docker

```dockerfile
# Use gopher in Docker
FROM golang:1.21-alpine AS gopher
RUN go install github.com/molmedoz/gopher/cmd/gopher@latest

FROM alpine:latest
COPY --from=gopher /go/bin/gopher /usr/local/bin/gopher
RUN gopher install 1.21.0
```

### Performance Tips

1. **Use SSD storage** for better performance
2. **Enable auto-cleanup** to save disk space
3. **Use local mirrors** for faster downloads
4. **Set appropriate max_versions** based on your needs

### Security Considerations

- Gopher verifies SHA256 checksums of all downloads
- Only downloads from official Go mirrors
- No external dependencies reduce attack surface
- All operations are local (no cloud sync)

## See Also

### Related Documentation

- **[Quick Reference](../QUICK_REFERENCE.md)** - One-page command reference and common workflows
- **[FAQ](FAQ.md)** - Frequently asked questions and quick answers
- **[Examples](EXAMPLES.md)** - 50+ practical usage examples
- **[API Reference](API_REFERENCE.md)** - Detailed API documentation
- **[Developer Guide](DEVELOPER_GUIDE.md)** - Contributing and development guide

### Platform-Specific Guides

- **[Windows Setup](WINDOWS_SETUP_GUIDE.md)** - Complete Windows setup with Developer Mode
- **[Windows Usage](WINDOWS_USAGE.md)** - Quick Windows reference
- **[Testing Guide](TESTING_GUIDE.md)** - Multi-platform testing strategy

### Additional Resources

- **[Documentation Index](../DOCUMENTATION_INDEX.md)** - Complete documentation navigation
- **[Roadmap](ROADMAP.md)** - Planned features and enhancements
- **[Changelog](../CHANGELOG.md)** - Version history and changes
- **[Contributing](../CONTRIBUTING.md)** - How to contribute to Gopher

## Getting Help

- **Documentation**: [README.md](../README.md)
- **Issues**: [GitHub Issues](https://github.com/molmedoz/gopher/issues)
- **Discussions**: [GitHub Discussions](https://github.com/molmedoz/gopher/discussions)

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for information on contributing to Gopher.
