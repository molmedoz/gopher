# Gopher Quick Reference

**One-page reference for Gopher commands and common workflows.**

---

## 📋 Essential Commands

| Command | Description | Example |
|---------|-------------|---------|
| `gopher list` | List installed versions (interactive) | `gopher list` |
| `gopher list-remote` | List available versions | `gopher list-remote` |
| `gopher install <version>` | Install a Go version | `gopher install 1.21.0` |
| `gopher use <version>` | Switch to a version | `gopher use 1.21.0` |
| `gopher use system` | Switch to system Go | `gopher use system` |
| `gopher current` | Show active version | `gopher current` |
| `gopher system` | Show system Go info | `gopher system` |
| `gopher uninstall <version>` | Remove a version | `gopher uninstall 1.20.0` |
| `gopher clean` | Remove download cache | `gopher clean` |
| `gopher purge` | Remove all Gopher data | `gopher purge` |
| `gopher version` | Show Gopher version | `gopher version` |
| `gopher help` | Show help | `gopher help` |

---

## 🚀 Common Workflows

### Install and Switch to New Version
```bash
gopher install 1.21.0
gopher use 1.21.0
go version  # Verify
```

### Return to System Go
```bash
gopher use system
go version  # Verify
```

### Check Current Setup
```bash
gopher list      # See installed versions
gopher current   # See active version
gopher system    # See system Go
```

### Clean Up Old Versions
```bash
gopher list                    # Find versions to remove
gopher uninstall 1.20.0        # Remove specific version
gopher clean                   # Clean download cache
gopher purge                   # Remove all Gopher data (requires confirmation)
```

### Find and Install Latest Version
```bash
gopher list-remote             # Browse available versions
gopher install 1.23.0          # Install latest
gopher use 1.23.0              # Switch to it
```

---

## 🎯 Command Flags

### Global Flags (Before Command)
| Flag | Description | Example |
|------|-------------|---------|
| `--json` | JSON output | `gopher --json list` |
| `--verbose`, `-v` | Verbose logging | `gopher -v install 1.21.0` |
| `--quiet`, `-q` | Minimal output | `gopher -q list` |
| `--no-interactive` | Disable pagination | `gopher --no-interactive list` |
| `--config <path>` | Custom config file | `gopher --config custom.json list` |
| `--page-size <n>` | Items per page | `gopher --page-size 5 list` |
| `--page <n>` | Go to specific page | `gopher --page 2 list` |

**Important:** Flags must come **before** the command name!

### Interactive Pagination Controls
When using `list` or `list-remote`:
- **Enter** or **n** - Next page
- **p** - Previous page
- **q** - Quit
- **<number>** - Jump to page
- **h** - Show help

---

## ⚙️ Environment Management

### Configuration Commands
```bash
gopher env list                          # List all settings
gopher env show go1.21.0                 # Show env for version
gopher env set gopath_mode=shared        # Set GOPATH mode
gopher env set custom_gopath=/path       # Set custom GOPATH
gopher env reset                         # Reset to defaults
```

### GOPATH Modes
| Mode | Command | Description |
|------|---------|-------------|
| Shared | `gopher env set gopath_mode=shared` | All versions share GOPATH (default) |
| Version-Specific | `gopher env set gopath_mode=version-specific` | Each version has own GOPATH |
| Custom | `gopher env set gopath_mode=custom` | Use custom GOPATH location |

---

## 📂 File Locations

### Configuration
| OS | Path |
|----|------|
| Linux/macOS | `~/.gopher/config.json` |
| Windows | `%USERPROFILE%\gopher\config.json` |

### Installation Directories
| OS | Path |
|----|------|
| Linux/macOS | `~/.gopher/versions/` |
| Windows | `%USERPROFILE%\gopher\versions\` |

### Download Cache
| OS | Path |
|----|------|
| Linux/macOS | `~/.gopher/downloads/` |
| Windows | `%USERPROFILE%\gopher\downloads\` |

---

## 🔧 JSON Output Examples

### List Versions (JSON)
```bash
gopher --json list
```
```json
[
  {
    "version": "go1.21.0",
    "os": "darwin",
    "arch": "arm64",
    "installed_at": "2025-08-27T08:49:40-07:00",
    "is_active": true,
    "is_system": false,
    "path": "/Users/user/.gopher/versions/go1.21.0/bin/go"
  }
]
```

### Current Version (JSON)
```bash
gopher --json current
```
```json
{
  "version": "go1.21.0",
  "path": "/Users/user/.gopher/versions/go1.21.0/bin/go",
  "is_system": false
}
```

---

## 🐛 Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| Command not found | Check PATH: `echo $PATH \| grep gopher` |
| Version not switching | Reload shell: `source ~/.bashrc` or `source ~/.zshrc` |
| Permission denied (Linux/macOS) | Fix permissions: `chmod -R u+w ~/.gopher` |
| Slow downloads | Try different mirror: `gopher env set mirror_url=<url>` |
| Can't uninstall active version | Switch first: `gopher use system` then uninstall |

### Enable Debug Logging
```bash
gopher --verbose install 1.21.0
gopher -v list
```

---

## 🌍 Platform-Specific Notes

### Windows
- **Developer Mode required** for symlinks (see [Windows Setup Guide](WINDOWS_SETUP_GUIDE.md))
- Use PowerShell or Command Prompt
- Paths use `%USERPROFILE%` instead of `~`

### macOS
- Apple Silicon (M1/M2/M3) fully supported
- Homebrew Go detected automatically
- Uses `~/.gopher/` for all data

### Linux
- All major distributions supported
- Multiple architectures (x86, x64, ARM, ARM64)
- Uses `~/.gopher/` for all data

---

## 📝 Configuration File Example

`~/.gopher/config.json`:
```json
{
  "install_dir": "~/.gopher/versions",
  "download_dir": "~/.gopher/downloads",
  "mirror_url": "https://go.dev/dl/",
  "auto_cleanup": true,
  "max_versions": 5,
  "gopath_mode": "shared",
  "custom_gopath": "",
  "goproxy": "https://proxy.golang.org,direct",
  "gosumdb": "sum.golang.org"
}
```

---

## 🔗 Quick Links

- **[Full User Guide](docs/USER_GUIDE.md)** - Complete documentation
- **[FAQ](docs/FAQ.md)** - Frequently asked questions
- **[Examples](docs/EXAMPLES.md)** - 50+ practical examples
- **[Windows Setup](docs/WINDOWS_SETUP_GUIDE.md)** - Windows-specific guide
- **[Documentation Index](DOCUMENTATION_INDEX.md)** - All documentation

---

## 💡 Pro Tips

1. **Use JSON for scripting**: `gopher --json list | jq '.[] | select(.is_active)'`
2. **Disable pagination for scripts**: `gopher --no-interactive list`
3. **Auto-cleanup saves space**: Enable in config with `"auto_cleanup": true`
4. **Check before switching**: Run `gopher list` to see what's installed
5. **System Go is always available**: `gopher use system` works even if you uninstall all Gopher versions

---

**Last Updated:** 2025-10-13  
**Version:** 1.0

**Need more details?** See the [Complete Documentation Index](DOCUMENTATION_INDEX.md)

