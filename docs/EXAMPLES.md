# Gopher Examples

This document provides practical examples of using Gopher in various scenarios.

## Table of Contents

- [Basic Usage Examples](#basic-usage-examples)
- [System Integration Examples](#system-integration-examples)
- [Scripting Examples](#scripting-examples)
- [CI/CD Integration](#cicd-integration)
- [Development Workflows](#development-workflows)
- [Visual Indicators Examples](#visual-indicators-examples)
- [Verbosity Control Examples](#verbosity-control-examples)
- [Troubleshooting Examples](#troubleshooting-examples)

## Basic Usage Examples

### Installing and Switching Go Versions

```bash
# Check current setup
gopher list
gopher current

# Install a new Go version
gopher install 1.21.0

# Switch to the new version
gopher use 1.21.0

# Verify the switch
go version

# Switch back to system Go
gopher use system
```

### Initial Setup and Configuration

```bash
# Run interactive setup wizard
gopher init

# Set up shell integration
gopher setup

# Check status and configuration
gopher status

# Get debug information
gopher debug
```

### Discovering Available Versions

```bash
# List available Go versions for installation (interactive by default)
gopher list-remote

# Disable interactive mode
gopher --no-interactive list-remote

# List with pagination (non-interactive)
gopher --page-size 10 --no-interactive list-remote

# Filter for specific versions
gopher --filter "1.21" list-remote

# Get JSON output for scripting
gopher --json list-remote
```

### Managing Multiple Versions

```bash
# Install multiple versions
gopher install 1.20.7
gopher install 1.21.0
gopher install 1.22.0

# List all installed versions (interactive by default)
gopher list

# Disable interactive mode
gopher --no-interactive list

# Switch between versions
gopher use 1.20.7
gopher use 1.21.0
gopher use 1.22.0

# Remove old versions
gopher uninstall 1.20.7

# Clean download cache to free disk space
gopher clean

# Remove all Gopher data (requires confirmation)
gopher purge
```

### Maintenance and Cleanup

```bash
# Clean download cache (safe, keeps installed versions)
gopher clean

# Example output:
# Cleaning download cache...
# ✓ Successfully cleaned download cache
#   Freed: 125.3 MB

# Purge all Gopher data (destructive, requires confirmation)
gopher purge

# Interactive confirmation:
# ⚠️  WARNING: This will permanently delete ALL Gopher data:
#   • All installed Go versions
#   • Download cache
#   • Configuration files
#   • State files and aliases
#   • Gopher-created symlinks
#
# This operation CANNOT be undone!
#
# Type 'yes' to confirm purge: yes
#
# Purging all Gopher data...
# ✓ Successfully purged all Gopher data

# After purge, you can reinitialize
gopher init
```

### Working with System Go

```bash
# Check system Go information
gopher system

# Switch to system Go
gopher use system

# Check if system Go is available
if gopher system > /dev/null 2>&1; then
    echo "System Go is available"
else
    echo "System Go not found"
fi
```

## System Integration Examples

### Shell Integration

#### Bash Profile Integration

Add to `~/.bashrc`:

```bash
# Gopher integration
export PATH="$HOME/.gopher/versions/current/bin:$PATH"

# Auto-switch Go versions based on directory
function cd() {
    builtin cd "$@"
    if [[ -f .gopher-version ]]; then
        local version=$(cat .gopher-version)
        if gopher list --json | jq -e ".[] | select(.version == \"go$version\")" > /dev/null; then
            gopher use "$version"
        else
            echo "Installing Go $version..."
            gopher install "$version"
            gopher use "$version"
        fi
    fi
}

# Go version in prompt
export PS1='$(gopher current --json | jq -r ".version") $ '
```

#### Zsh Integration

Add to `~/.zshrc`:

```zsh
# Gopher integration
export PATH="$HOME/.gopher/versions/current/bin:$PATH"

# Auto-switch Go versions
function chpwd() {
    if [[ -f .gopher-version ]]; then
        local version=$(cat .gopher-version)
        if gopher list --json | jq -e ".[] | select(.version == \"go$version\")" > /dev/null; then
            gopher use "$version"
        else
            echo "Installing Go $version..."
            gopher install "$version"
            gopher use "$version"
        fi
    fi
}

# Go version in prompt
autoload -U promptinit; promptinit
export PS1='$(gopher current --json | jq -r ".version") $ '
```

### IDE Integration

#### VS Code Settings

Create `.vscode/settings.json`:

```json
{
    "go.goroot": "~/.gopher/versions/go1.21.0",
    "go.toolsEnvVars": {
        "GOROOT": "~/.gopher/versions/go1.21.0"
    },
    "go.alternateTools": {
        "go": "~/.gopher/versions/go1.21.0/bin/go"
    }
}
```

#### GoLand Integration

1. Go to Settings → Go → GOROOT
2. Set GOROOT to `~/.gopher/versions/go1.21.0`
3. Or use the "Detect" button to auto-detect

### Docker Integration

#### Dockerfile Example

```dockerfile
FROM golang:1.21-alpine AS gopher

# Install gopher
RUN go install github.com/molmedoz/gopher/cmd/gopher@latest

# Install Go versions
RUN gopher install 1.20.7
RUN gopher install 1.21.0

FROM alpine:latest

# Copy gopher and Go versions
COPY --from=gopher /go/bin/gopher /usr/local/bin/gopher
COPY --from=gopher /root/.gopher /root/.gopher

# Set up environment
ENV PATH="/root/.gopher/versions/current/bin:$PATH"

# Use specific Go version
RUN gopher use 1.21.0
```

#### Docker Compose Example

```yaml
version: '3.8'
services:
  app:
    build: .
    environment:
      - GOPHER_VERSION=1.21.0
    volumes:
      - .:/app
    working_dir: /app
    command: |
      sh -c "
        gopher use $GOPHER_VERSION &&
        go mod tidy &&
        go run main.go
      "
```

## Scripting Examples

### Basic Scripting

#### Check Go Version

```bash
#!/bin/bash
# check-go-version.sh

current_version=$(gopher current --json | jq -r '.version')
echo "Current Go version: $current_version"

if [[ "$current_version" == "go1.21.0" ]]; then
    echo "✅ Correct Go version"
    exit 0
else
    echo "❌ Wrong Go version. Expected go1.21.0, got $current_version"
    exit 1
fi
```

#### Install Required Go Version

```bash
#!/bin/bash
# install-go-version.sh

VERSION=${1:-"1.21.0"}

echo "Installing Go $VERSION..."

# Check if already installed
if gopher list --json | jq -e ".[] | select(.version == \"go$VERSION\")" > /dev/null; then
    echo "Go $VERSION is already installed"
else
    gopher install "$VERSION"
fi

# Switch to the version
gopher use "$VERSION"

# Verify installation
go version
```

#### Project Setup Script

```bash
#!/bin/bash
# setup-project.sh

PROJECT_DIR=${1:-"."}
GO_VERSION=${2:-"1.21.0"}

cd "$PROJECT_DIR" || exit 1

echo "Setting up project in $PROJECT_DIR with Go $GO_VERSION"

# Create .gopher-version file
echo "$GO_VERSION" > .gopher-version

# Install and use Go version
gopher install "$GO_VERSION"
gopher use "$GO_VERSION"

# Initialize Go module if needed
if [[ ! -f go.mod ]]; then
    go mod init "$(basename "$PROJECT_DIR")"
fi

# Install dependencies
go mod tidy

echo "Project setup complete!"
echo "Go version: $(go version)"
echo "Project directory: $(pwd)"
```

### Advanced Scripting

#### Version Management Script

```bash
#!/bin/bash
# manage-go-versions.sh

case "$1" in
    "list")
        gopher list
        ;;
    "install")
        gopher install "$2"
        ;;
    "use")
        gopher use "$2"
        ;;
    "current")
        gopher current
        ;;
    "system")
        gopher system
        ;;
    "cleanup")
        # Remove old versions (keep last 3)
        versions=$(gopher list --json | jq -r '.[] | select(.is_system == false) | .version' | sort -V)
        count=$(echo "$versions" | wc -l)
        if [[ $count -gt 3 ]]; then
            to_remove=$(echo "$versions" | head -n $((count - 3)))
            echo "$to_remove" | while read -r version; do
                echo "Removing $version..."
                gopher uninstall "$version"
            done
        fi
        ;;
    *)
        echo "Usage: $0 {list|install|use|current|system|cleanup}"
        exit 1
        ;;
esac
```

#### Health Check Script

```bash
#!/bin/bash
# health-check.sh

echo "🔍 Gopher Health Check"
echo "====================="

# Check gopher installation
if ! command -v gopher >/dev/null 2>&1; then
    echo "❌ Gopher not found in PATH"
    exit 1
fi
echo "✅ Gopher installed: $(gopher version)"

# Check Go installation
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go not found in PATH"
    exit 1
fi
echo "✅ Go installed: $(go version)"

# Check gopher versions
echo ""
echo "📋 Installed Go versions:"
gopher list

# Check system Go
echo ""
echo "🏠 System Go information:"
if gopher system >/dev/null 2>&1; then
    gopher system
else
    echo "❌ No system Go found"
fi

# Check configuration
echo ""
echo "⚙️  Configuration:"
config_path=$(gopher --help 2>&1 | grep -o '--config [^ ]*' | head -1)
if [[ -n "$config_path" ]]; then
    echo "Config file: $config_path"
    if [[ -f "$config_path" ]]; then
        echo "✅ Config file exists"
    else
        echo "❌ Config file not found"
    fi
fi

echo ""
echo "✅ Health check complete!"
```

## CI/CD Integration

### GitHub Actions

#### Basic Go Setup

```yaml
name: Go CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go with Gopher
      run: |
        go install github.com/molmedoz/gopher/cmd/gopher@latest
        gopher install 1.21.0
        gopher use 1.21.0
        echo "$(gopher current --json | jq -r '.path')" >> $GITHUB_PATH
    
    - name: Verify Go version
      run: go version
    
    - name: Run tests
      run: go test ./...
```

#### Multi-version Testing

```yaml
name: Multi-version Testing

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.20.7, 1.21.0, 1.22.0]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Gopher
      run: go install github.com/molmedoz/gopher/cmd/gopher@latest
    
    - name: Install Go ${{ matrix.go-version }}
      run: |
        gopher install ${{ matrix.go-version }}
        gopher use ${{ matrix.go-version }}
        echo "$(gopher current --json | jq -r '.path')" >> $GITHUB_PATH
    
    - name: Test with Go ${{ matrix.go-version }}
      run: |
        go version
        go test ./...
```

#### Custom Go Version

```yaml
name: Custom Go Version

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Gopher
      run: go install github.com/molmedoz/gopher/cmd/gopher@latest
    
    - name: Install specific Go version
      run: |
        # Read Go version from .gopher-version file
        GO_VERSION=$(cat .gopher-version || echo "1.21.0")
        gopher install "$GO_VERSION"
        gopher use "$GO_VERSION"
        echo "$(gopher current --json | jq -r '.path')" >> $GITHUB_PATH
    
    - name: Test
      run: go test ./...
```

### GitLab CI

```yaml
stages:
  - test

variables:
  GO_VERSION: "1.21.0"

test:
  stage: test
  image: golang:1.21-alpine
  before_script:
    - go install github.com/molmedoz/gopher/cmd/gopher@latest
    - gopher install $GO_VERSION
    - gopher use $GO_VERSION
    - export PATH="$(gopher current --json | jq -r '.path' | xargs dirname):$PATH"
  script:
    - go version
    - go test ./...
```

### Jenkins Pipeline

```groovy
pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21.0'
    }
    
    stages {
        stage('Setup') {
            steps {
                sh 'go install github.com/molmedoz/gopher/cmd/gopher@latest'
                sh 'gopher install ${GO_VERSION}'
                sh 'gopher use ${GO_VERSION}'
                sh 'export PATH="$(gopher current --json | jq -r \'.path\' | xargs dirname):$PATH"'
            }
        }
        
        stage('Test') {
            steps {
                sh 'go version'
                sh 'go test ./...'
            }
        }
    }
}
```

## Development Workflows

### Project-specific Go Versions

#### Using .gopher-version Files

```bash
# In your project root
echo "1.21.0" > .gopher-version

# Auto-switch script
#!/bin/bash
if [[ -f .gopher-version ]]; then
    version=$(cat .gopher-version)
    if gopher list --json | jq -e ".[] | select(.version == \"go$version\")" > /dev/null; then
        gopher use "$version"
    else
        echo "Installing Go $version..."
        gopher install "$version"
        gopher use "$version"
    fi
fi
```

#### Makefile Integration

```makefile
# Makefile

GO_VERSION := $(shell cat .gopher-version 2>/dev/null || echo "1.21.0")

.PHONY: setup
setup:
	@echo "Setting up Go $(GO_VERSION)..."
	@gopher install $(GO_VERSION)
	@gopher use $(GO_VERSION)
	@go version

.PHONY: test
test: setup
	@go test ./...

.PHONY: build
build: setup
	@go build -o bin/app .

.PHONY: clean
clean:
	@rm -rf bin/
	@gopher uninstall $(GO_VERSION)
```

### Team Development

#### Shared Configuration

Create `.gopher-config.json` in your project:

```json
{
  "install_dir": "./.gopher/versions",
  "download_dir": "./.gopher/downloads",
  "mirror_url": "https://go.dev/dl/",
  "auto_cleanup": true,
  "max_versions": 3
}
```

Use with gopher:

```bash
gopher --config .gopher-config.json install 1.21.0
gopher --config .gopher-config.json use 1.21.0
```

#### Team Setup Script

```bash
#!/bin/bash
# team-setup.sh

echo "🚀 Setting up development environment..."

# Check if gopher is installed
if ! command -v gopher >/dev/null 2>&1; then
    echo "Installing gopher..."
    go install github.com/molmedoz/gopher/cmd/gopher@latest
fi

# Install required Go version
GO_VERSION=$(cat .gopher-version 2>/dev/null || echo "1.21.0")
echo "Installing Go $GO_VERSION..."

if gopher list --json | jq -e ".[] | select(.version == \"go$GO_VERSION\")" > /dev/null; then
    echo "Go $GO_VERSION already installed"
else
    gopher install "$GO_VERSION"
fi

# Switch to the version
gopher use "$GO_VERSION"

# Install project dependencies
go mod tidy

echo "✅ Development environment ready!"
echo "Go version: $(go version)"
```

## Visual Indicators Examples

### Understanding Version Display

Gopher provides clear visual indicators to help you identify the active version and understand the status of different Go versions.

#### Basic List Output

```bash
# Standard list command
gopher list
```

**Output with visual indicators:**
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

#### Switching Versions

```bash
# Switch to a different version
gopher use 1.22.0

# List to see the change
gopher list
```

**Output after switching:**
```
  go1.24.7 (darwin/arm64) [system]
→ go1.22.0 (darwin/arm64) [active]
  go1.22.1 (darwin/arm64)
  go1.23.0 (darwin/arm64)
```

#### System Go Detection

```bash
# Switch to system Go
gopher use system

# List to see system Go is now active
gopher list
```

**Output with system Go active:**
```
→ go1.24.7 (darwin/arm64) [system]
  go1.22.0 (darwin/arm64)
  go1.22.1 (darwin/arm64)
  go1.23.0 (darwin/arm64)
```

### Color Support Examples

#### Terminal with Color Support

```bash
# In a color-capable terminal
gopher list
# Shows colored output with green active version, cyan system version
```

#### Plain Text Output

```bash
# Redirect to file (no colors)
gopher list > versions.txt

# Use JSON output (no colors)
gopher --json list
```

### Scripting with Visual Indicators

#### Extract Active Version

```bash
# Get the active version using visual indicators
active_version=$(gopher list | grep "→" | sed 's/→ //' | awk '{print $1}')
echo "Active version: $active_version"
```

#### Check if System Go is Active

```bash
# Check if system Go is currently active
if gopher list | grep -q "→.*\[system\]"; then
    echo "System Go is active"
else
    echo "Gopher-managed Go is active"
fi
```

## Verbosity Control Examples

### Understanding Log Levels

Gopher provides three log levels to control the amount of output you see.

#### Quiet Mode (ERROR level only)

```bash
# Minimal output - only errors
gopher --quiet list
gopher -q install 1.21.0

# Useful for scripting
current_version=$(gopher -q current --json | jq -r '.version')
echo "Current version: $current_version"
```

#### Default Mode (INFO level)

```bash
# Standard output - general information
gopher list
gopher install 1.21.0
```

**Example output:**
```
📦 Installing Go 1.21.0...
✓ Successfully installed Go 1.21.0
```

#### Verbose Mode (DEBUG level)

```bash
# Detailed output with debugging information
gopher --verbose install 1.21.0
gopher -v list
```

**Example output:**
```
[DEBUG] Starting Go installation {version=go1.21.0 file=/tmp/go1.21.0.darwin-arm64.tar.gz}
[INFO] Downloading file {file=go1.21.0.darwin-arm64.tar.gz}
[DEBUG] Extracting files {count=1000}
[INFO] Successfully installed Go version {version=go1.21.0}
```

### Installation Examples

#### Quiet Installation

```bash
# Minimal output during installation
gopher -q install 1.21.0
# Only shows errors if they occur
```

#### Verbose Installation

```bash
# Detailed installation progress
gopher -v install 1.21.0
```

**Output:**
```
[DEBUG] Starting Go installation {version=go1.21.0 file=/tmp/go1.21.0.darwin-arm64.tar.gz}
[INFO] Downloading file {file=go1.21.0.darwin-arm64.tar.gz}
[DEBUG] Download response received {status_code=200 headers=map[Content-Type:[application/octet-stream]]}
[INFO] Download completed {file=go1.21.0.darwin-arm64.tar.gz}
[DEBUG] Extracting files {count=1000}
[INFO] Extraction completed {files_extracted=1000}
[INFO] Successfully installed Go version {version=go1.21.0}
```

#### Default Installation

```bash
# Standard installation output
gopher install 1.21.0
```

**Output:**
```
📦 Installing Go 1.21.0...
✓ Successfully installed Go 1.21.0
```

### Listing Examples

#### Quiet Listing

```bash
# Just the list, no extra information
gopher -q list
```

#### Verbose Listing

```bash
# Detailed listing with metadata
gopher -v list
```

**Output:**
```
[DEBUG] Scanning installation directory {path=/Users/user/.gopher/versions}
[DEBUG] Found installed version {version=go1.21.0 path=/Users/user/.gopher/versions/go1.21.0}
[DEBUG] Checking symlink {path=/Users/user/.local/bin/go}
[DEBUG] Active version detected {version=go1.21.0}
  go1.24.7 (darwin/arm64) [system]
→ go1.21.0 (darwin/arm64) [active]
```

#### Default Listing

```bash
# Standard listing format
gopher list
```

### System Detection Examples

#### Quiet System Info

```bash
# Minimal system information
gopher -q system
```

#### Verbose System Detection

```bash
# Detailed system detection process
gopher -v system
```

**Output:**
```
[DEBUG] Detecting system Go {path=/usr/local/go/bin/go}
[DEBUG] System Go found {version=go1.24.7 path=/usr/local/go/bin/go}
[DEBUG] Checking Homebrew Go {path=/opt/homebrew/bin/go}
[DEBUG] Homebrew Go not found
[DEBUG] System detection completed {version=go1.24.7 path=/usr/local/go/bin/go}
System Go: go1.24.7 (darwin/arm64)
Path: /usr/local/go/bin/go
```

#### Default System Info

```bash
# Standard system information
gopher system
```

**Output:**
```
System Go: go1.24.7 (darwin/arm64)
Path: /usr/local/go/bin/go
```

### Combining with Other Flags

#### Verbose JSON Output

```bash
# Detailed output in JSON format
gopher -v --json list
```

#### Quiet Installation with Custom Config

```bash
# Quiet installation using custom configuration
gopher -q --config /path/to/config.json install 1.21.0
```

#### Verbose Alias Operations

```bash
# Detailed alias creation
gopher -v alias create stable 1.21.0
```

**Output:**
```
[DEBUG] Creating alias {name=stable version=go1.21.0}
[DEBUG] Validating alias name {name=stable}
[DEBUG] Checking if version is installed {version=go1.21.0}
[DEBUG] Saving aliases {file=/Users/user/.gopher/aliases.json}
[INFO] Alias created {name=stable version=go1.21.0}
```

### Use Case Examples

#### Development and Debugging

```bash
# Debug installation issues
gopher -v install 1.21.0

# Debug system detection problems
gopher -v system

# Debug alias operations
gopher -v alias create stable 1.21.0
```

#### Scripting and Automation

```bash
#!/bin/bash
# Script that uses quiet mode for automation

# Get current version quietly
current_version=$(gopher -q current --json | jq -r '.version')
echo "Current version: $current_version"

# Install version quietly
if ! gopher -q install 1.21.0; then
    echo "Installation failed"
    exit 1
fi

# Switch to version quietly
gopher -q use 1.21.0
```

#### User-Friendly Output

```bash
# Standard output for interactive use
gopher list
gopher install 1.21.0
gopher use 1.21.0
```

## Troubleshooting Examples

### Common Issues and Solutions

#### Permission Denied

```bash
# Problem: Permission denied when switching versions
# Solution: Use sudo or change installation directory

# Option 1: Use sudo
sudo gopher use 1.21.0

# Option 2: Use local installation
export GOPHER_INSTALL_DIR=~/.local/go-versions
gopher use 1.21.0
```

#### Go Not Found After Switching

```bash
# Problem: go command not found
# Solution: Check PATH and symlinks

# Check if symlink exists
ls -la /usr/local/bin/go

# Recreate symlink
gopher use 1.21.0

# Check PATH
echo $PATH

# Add to PATH if needed
export PATH="$(gopher current --json | jq -r '.path' | xargs dirname):$PATH"
```

#### System Go Not Detected

```bash
# Problem: System Go not showing in gopher list
# Solution: Check Go installation and PATH

# Check if Go is in PATH
which go

# Check system detection
gopher system

# Verify Go installation
go version

# Check if it's a recognized system installation
gopher system --json | jq '.is_system'
```

#### Download Failures

```bash
# Problem: Failed to download Go version
# Solution: Check network and try different mirror

# Check internet connection
ping go.dev

# Try different mirror
gopher --config <(echo '{"mirror_url": "https://golang.org/dl/"}') install 1.21.0

# Check disk space
df -h

# Check download directory permissions
ls -la ~/.gopher/downloads
```

### Debug Scripts

#### Debug Installation

```bash
#!/bin/bash
# debug-install.sh

VERSION=${1:-"1.21.0"}

echo "🔍 Debugging Go $VERSION installation..."

# Check gopher
echo "Gopher version: $(gopher version)"

# Check system Go
echo "System Go:"
gopher system || echo "No system Go"

# Check installed versions
echo "Installed versions:"
gopher list

# Try to install
echo "Installing Go $VERSION..."
if gopher install "$VERSION"; then
    echo "✅ Installation successful"
    gopher use "$VERSION"
    go version
else
    echo "❌ Installation failed"
    echo "Check logs and try again"
fi
```

#### Debug System Detection

```bash
#!/bin/bash
# debug-system.sh

echo "🔍 Debugging system Go detection..."

# Check if go is in PATH
if which go >/dev/null 2>&1; then
    echo "✅ Go found in PATH: $(which go)"
    go version
else
    echo "❌ Go not found in PATH"
    exit 1
fi

# Check gopher system detection
echo "Gopher system detection:"
gopher system

# Check if detected as system
if gopher system --json | jq -e '.is_system' >/dev/null; then
    echo "✅ Detected as system installation"
else
    echo "❌ Not detected as system installation"
    echo "Path: $(gopher system --json | jq -r '.path')"
fi
```

These examples should help you get started with Gopher and integrate it into your development workflow. For more information, see the [User Guide](USER_GUIDE.md) or [API Reference](API_REFERENCE.md).
