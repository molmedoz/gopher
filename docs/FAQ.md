# Gopher - Frequently Asked Questions (FAQ)

Quick answers to common questions about Gopher, the Go version manager.

---

## 📚 Table of Contents

- [Installation & Setup](#installation--setup)
- [Usage & Commands](#usage--commands)
- [System Integration](#system-integration)
- [Environment Management](#environment-management)
- [Troubleshooting](#troubleshooting)
- [Platform-Specific](#platform-specific)
- [Advanced Topics](#advanced-topics)

---

## Installation & Setup

### Q: Do I need to uninstall my system Go before using Gopher?

**A:** No! Gopher works **alongside** your system Go installation. You can easily switch between Gopher-managed versions and your system Go using `gopher use system`.

### Q: Where does Gopher install Go versions?

**A:** Gopher installs Go versions in:
- **Linux/macOS**: `~/.gopher/versions/`
- **Windows**: `%USERPROFILE%\gopher\versions\`

All versions are isolated from your system Go installation.

### Q: How much disk space does Gopher need?

**A:** Each Go version requires approximately 150-200 MB. Gopher supports auto-cleanup to manage disk space automatically. You can configure the maximum number of versions to keep in your config file.

### Q: Can I customize Gopher's installation directory?

**A:** Yes! Set the `GOPHER_INSTALL_DIR` environment variable or modify the `install_dir` setting in your configuration file (`~/.gopher/config.json`).

### Q: How do I completely uninstall Gopher?

**A:**
```bash
# 1. Switch back to system Go
gopher use system

# 2. Remove Gopher binary
sudo rm /usr/local/bin/gopher  # or wherever it's installed

# 3. Remove Gopher data (optional)
rm -rf ~/.gopher
```

---

## Usage & Commands

### Q: How do I see what Go version is currently active?

**A:** Use any of these commands:
```bash
gopher current      # Shows active version via Gopher
go version          # Shows actual Go version being used
which go            # Shows path to active Go binary
```

### Q: How do I list all installed Go versions?

**A:** Use:
```bash
gopher list                    # Interactive pagination (default)
gopher --no-interactive list   # Non-interactive list
gopher --json list             # JSON output for scripting
```

### Q: How do I find available Go versions to install?

**A:**
```bash
gopher list-remote             # Interactive pagination (default)
gopher --no-interactive list-remote  # Non-interactive list
```

### Q: Can I install multiple Go versions at once?

**A:** Currently, you need to install them one at a time:
```bash
gopher install 1.21.0
gopher install 1.22.0
gopher install 1.23.0
```

### Q: How do I uninstall a Go version I no longer need?

**A:**
```bash
gopher uninstall 1.20.0
```
Note: You cannot uninstall the currently active version. Switch to another version first.

---

## System Integration

### Q: Can I use Gopher alongside my system Go?

**A:** Absolutely! Gopher detects your system Go and allows you to switch between Gopher-managed versions and your system installation:
```bash
gopher install 1.21.0
gopher use 1.21.0     # Use Gopher version
gopher use system     # Switch back to system Go
```

### Q: How does Gopher detect my system Go?

**A:** Gopher searches for Go in:
1. Homebrew installation (`/opt/homebrew/opt/go` on macOS)
2. Standard installation paths (`/usr/local/go`, `/usr/bin/go`)
3. PATH environment variable

### Q: Will Gopher modify my system Go installation?

**A:** No! Gopher never modifies your system Go. It uses symlinks to manage version switching without touching your system installation.

### Q: Can I use different Go versions for different projects?

**A:** Yes! Simply switch versions when working on different projects:
```bash
cd ~/project1
gopher use 1.21.0

cd ~/project2
gopher use 1.22.0
```

Future versions will support automatic project-based version switching.

---

## Environment Management

### Q: What is GOPATH mode and which should I use?

**A:** Gopher supports three GOPATH modes:

1. **Shared (Default)**: All Go versions share the same GOPATH
   ```bash
   gopher env set gopath_mode=shared
   ```

2. **Version-Specific**: Each Go version has its own isolated GOPATH
   ```bash
   gopher env set gopath_mode=version-specific
   ```

3. **Custom**: Use a custom GOPATH location
   ```bash
   gopher env set gopath_mode=custom
   gopher env set custom_gopath=/path/to/workspace
   ```

For most users, **shared mode** is recommended for easier package management.

### Q: How do I see my current environment configuration?

**A:**
```bash
gopher env list              # List all configuration options
gopher env show go1.21.0     # Show environment for specific version
```

### Q: How do I reset my environment to defaults?

**A:**
```bash
gopher env reset
```

### Q: Does Gopher set GOROOT automatically?

**A:** Yes! Gopher automatically sets `GOROOT` to point to the active Go version's installation directory.

---

## Troubleshooting

### Q: Gopher command not found - what should I do?

**A:**
1. Check if Gopher is installed: `which gopher`
2. If not found, reinstall following the [installation guide](USER_GUIDE.md#installation)
3. If installed, ensure it's in your PATH:
   ```bash
   echo $PATH | grep gopher
   ```

### Q: Version switching doesn't work - I still see the old version

**A:**
1. Check current version: `gopher current`
2. Verify the switch: `go version`
3. If mismatch, check shell configuration:
   ```bash
   # Reload shell configuration
   source ~/.bashrc    # or ~/.zshrc
   ```
4. Check symlink: `ls -la $(which go)`

### Q: Downloads are slow or failing

**A:**
1. Check internet connection
2. Try using a different mirror:
   ```bash
   gopher env set mirror_url=https://golang.google.cn/dl/
   ```
3. Check if you can access the Go downloads page in your browser

### Q: I get "permission denied" errors

**A:**
On Linux/macOS, ensure you have write permissions:
```bash
# Check permissions
ls -la ~/.gopher

# Fix if needed
chmod -R u+w ~/.gopher
```

On Windows, ensure Developer Mode is enabled (see [Windows Setup Guide](WINDOWS_SETUP_GUIDE.md)).

### Q: How do I enable debug logging?

**A:**
```bash
gopher --verbose install 1.21.0   # Verbose output
gopher -v list                     # Short form
```

For detailed debugging, check Gopher's log files in `~/.gopher/logs/`.

---

## Platform-Specific

### Q: Do I need Developer Mode enabled on Windows?

**A:** Yes, for the best experience. Developer Mode allows Gopher to create symlinks without administrator privileges. See the [Windows Setup Guide](WINDOWS_SETUP_GUIDE.md) for instructions.

### Q: Can I run Gopher on macOS with Apple Silicon (M1/M2/M3)?

**A:** Yes! Gopher fully supports Apple Silicon. It will automatically download ARM64 versions of Go when running on Apple Silicon Macs.

### Q: Does Gopher work on Linux ARM?

**A:** Yes! Gopher supports multiple architectures including ARM, ARM64, x86, and x86_64.

### Q: Can I use Gopher in WSL (Windows Subsystem for Linux)?

**A:** Yes! Gopher works great in WSL. Follow the Linux installation instructions.

### Q: Does Gopher work in Docker containers?

**A:** Yes! See [EXAMPLES.md](EXAMPLES.md#docker-integration) for Docker integration examples.

---

## Advanced Topics

### Q: Can I use Gopher in CI/CD pipelines?

**A:** Absolutely! Gopher provides JSON output for easy scripting:
```bash
gopher --json list
gopher --json current
```

See [EXAMPLES.md](EXAMPLES.md#cicd-integration) for detailed CI/CD examples.

### Q: How does Gopher verify downloaded Go binaries?

**A:** Gopher uses SHA256 checksums to verify all downloads. The checksums are fetched from the official Go downloads page and verified before installation.

### Q: Can I configure Gopher to auto-cleanup old versions?

**A:** Yes! Set in your config file:
```json
{
  "auto_cleanup": true,
  "max_versions": 5
}
```

Gopher will automatically keep only the 5 most recent versions.

### Q: Can I use Gopher with go.mod version directives?

**A:** Yes! While Gopher doesn't automatically detect `go.mod` versions yet (planned for future release), you can manually switch to the required version:
```bash
# Check go.mod for required version
cat go.mod | grep "^go "
# Output: go 1.21

# Install and switch to required version
gopher install 1.21.0
gopher use 1.21.0
```

### Q: Go is automatically downloading toolchains and overriding my PATH. How do I prevent this?

**A:** Go 1.21+ has automatic toolchain management that can download newer Go versions when your `go.mod` requires them. To prevent this behavior:

**Option 1: Set GOTOOLCHAIN=local (Recommended)**
```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export GOTOOLCHAIN=local

# Apply immediately
source ~/.bashrc  # or ~/.zshrc
```

This tells Go to only use your locally installed version and never auto-download.

**Option 2: Align go.mod with installed Go version**
```bash
# Check your installed Go version
go version  # e.g., go version go1.24.9

# Update go.mod to match
# Change: go 1.25.1
# To: go 1.24.9
```

**Why this happens:**
- Go automatically downloads toolchains when `go.mod` requires a newer version
- These downloads go to `$GOPATH/pkg/mod/golang.org/toolchain@.../`
- The downloaded toolchain is prepended to PATH for that process
- IDEs like Cursor may inherit this behavior

**More info:** See [Go toolchain documentation](https://go.dev/doc/toolchain)

### Q: How do I script with Gopher?

**A:** Use JSON output for easy parsing:
```bash
#!/bin/bash
current=$(gopher --json current | jq -r '.version')
echo "Current Go version: $current"
```

See [EXAMPLES.md](EXAMPLES.md#scripting-examples) for more scripting examples.

### Q: Can I create aliases for Go versions?

**A:** Yes! Gopher v1.0+ supports version aliases:
```bash
gopher alias create stable 1.21.0
gopher use stable
```

See the [Roadmap](ROADMAP.md) for alias feature details.

---

## Still Have Questions?

### Can't find your answer?

1. Check the [User Guide](USER_GUIDE.md) for comprehensive documentation
2. Browse [Examples](EXAMPLES.md) for practical use cases
3. Review [Troubleshooting](USER_GUIDE.md#troubleshooting) section
4. Search [GitHub Issues](https://github.com/molmedoz/gopher/issues)
5. Create a new issue if you've found a bug or have a feature request

### Documentation Links

- [User Guide](USER_GUIDE.md) - Complete feature documentation
- [Examples](EXAMPLES.md) - 50+ practical examples
- [Developer Guide](DEVELOPER_GUIDE.md) - For contributors
- [API Reference](API_REFERENCE.md) - API documentation
- [Windows Setup](WINDOWS_SETUP_GUIDE.md) - Windows-specific setup
- [Documentation Index](../DOCUMENTATION_INDEX.md) - All documentation

---

**Last Updated:** 2025-10-13  
**Version:** 1.0  
**Maintainer:** Gopher Development Team

