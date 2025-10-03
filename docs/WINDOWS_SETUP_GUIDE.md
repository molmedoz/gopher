# Windows Setup Guide

**Complete guide for setting up Gopher on Windows with automated setup and smart detection.**

**Version:** v1.0.0+  
**Last Updated:** October 2025

---

## 📋 Prerequisites

- **Windows 10 or later**
- **PowerShell 5.1+** (comes with Windows 10)
- **Administrator access** (for PATH configuration)

---

## ⚡ Quick Start (5 Minutes)

The **easiest** way to set up Gopher on Windows:

### Step 1: Download Gopher

```powershell
# Create bin directory
New-Item -ItemType Directory -Path "$env:USERPROFILE\bin" -Force

# Download latest release
$version = "1.0.0"
$url = "https://github.com/molmedoz/gopher/releases/download/v$version/gopher-windows-amd64.exe"
Invoke-WebRequest -Uri $url -OutFile "$env:USERPROFILE\bin\gopher.exe"

# Add gopher to PATH
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$env:USERPROFILE\bin*") {
    [Environment]::SetEnvironmentVariable("PATH", "$env:USERPROFILE\bin;$userPath", "User")
}
```

**Restart PowerShell**, then verify:
```powershell
gopher version
```

---

### Step 2: Run Interactive Setup

```powershell
gopher init
```

**This will:**
- ✅ Detect your Windows environment (PowerShell, paths, etc.)
- ✅ Create all required directories automatically
- ✅ Test symlink creation
- ✅ Show you **exact PowerShell commands** to complete setup
- ✅ Explain what each step does

**Example output:**
```powershell
🚀 Welcome to Gopher Setup Wizard!

📋 System Detection
===================
Platform: windows/amd64
Shell: PowerShell ✅
Symlink Directory: C:\Users\YourName\AppData\Local\bin

🎯 Setup Complete! Next Steps:
================================

STEP 1: Add Gopher's bin directory to PATH (CRITICAL!)
--------------------------------------------------------

Copy and run this PowerShell command as Administrator:

  $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
  $gopherBin = "C:\Users\YourName\AppData\Local\bin"
  if ($userPath -notlike "*$gopherBin*") {
    [Environment]::SetEnvironmentVariable("PATH", "$gopherBin;$userPath", "User")
  }

Then RESTART your PowerShell terminal.

STEP 2: Install a Go version
  gopher install 1.21.0

STEP 3: Switch to it
  gopher use 1.21.0

  ⚠️  Important: If you have system Go installed, Gopher will warn you
  about PATH order and provide exact commands to fix it.

STEP 4: Verify it works
  go version
  # Should show: go version go1.21.0 windows/amd64

📁 Directories created:
  Config:    C:\Users\YourName\gopher\config.json
  Versions:  C:\Users\YourName\gopher\versions
  Downloads: C:\Users\YourName\gopher\downloads
  State:     C:\Users\YourName\gopher\state
  Symlinks:  C:\Users\YourName\AppData\Local\bin
```

**Just copy and run the commands shown!** ✅

---

### Step 3: Follow the Instructions

Copy and run the PowerShell commands from `gopher init` output:

1. **Add Gopher's bin to PATH** (copy/paste from output)
2. **Restart PowerShell**
3. **Install a Go version**: `gopher install 1.21.0`
4. **Switch to it**: `gopher use 1.21.0`
5. **Verify**: `go version`

**That's it!** 🎉

---

## 📁 Directory Structure

Gopher uses a clean, visible directory structure on Windows (no hidden folders):

```
C:\Users\YourName\
├── bin\
│   └── gopher.exe                    ← Gopher executable
│
├── gopher\                            ← Data directory (visible, no dot)
│   ├── config.json                   ← Configuration
│   ├── versions\                     ← Installed Go versions
│   │   ├── go1.21.0\
│   │   │   ├── bin\
│   │   │   │   └── go.exe
│   │   │   ├── src\
│   │   │   └── pkg\
│   │   └── go1.22.0\
│   ├── downloads\                    ← Temporary downloads
│   │   └── go1.21.0.windows-amd64.zip
│   └── state\                        ← Active version tracking
│       └── active-version
│
└── AppData\Local\
    └── bin\
        └── go.exe                     ← Symlink to active version
```

**Why "gopher" (no dot)?**
- Windows users expect visible folders in their home directory
- Easy to find and manage
- Follows Windows conventions

---

## 🔧 Advanced Setup

### Enable Developer Mode

**Why?** Allows symlink creation without admin privileges.

**How:**
1. Settings → Update & Security → For developers
2. Turn on "Developer Mode"
3. Restart terminal

**Alternative:** Run PowerShell as Administrator when using `gopher use`.

---

### Understanding PATH Order (CRITICAL!)

#### Automatic PATH Order Detection (v1.0.0+)

Gopher now automatically detects PATH order issues!

When you run `gopher use 1.21.0`, Gopher will:
- ✅ Create the symlink to the Go version
- ✅ Check if Gopher's bin directory is before system Go in PATH
- ⚠️ Warn you if system Go will take precedence
- 📋 Provide exact PowerShell commands to fix the PATH order

**Why this matters:** Windows searches PATH in order. If system Go comes first, it wins.

#### Example Warning:

```powershell
PS> gopher use 1.21.0
✓ Created symlink in C:\Users\YourName\AppData\Local\bin\go.exe

⚠️  WARNING: System Go will take precedence over Gopher-managed versions!

PATH Order Issue:
  Position 1: C:\Program Files\Go\bin (System Go) ← Found FIRST
  Position 2: C:\Users\YourName\AppData\Local\bin (Gopher) ← Found SECOND

This means 'go version' will still show system Go, not the version you just switched to.

TO FIX - Run this PowerShell command as Administrator:

  $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
  $gopherBin = "C:\Users\YourName\AppData\Local\bin"
  $pathArray = $userPath -split ';' | Where-Object { $_ -ne $gopherBin }
  $newUserPath = "$gopherBin;" + ($pathArray -join ';')
  [Environment]::SetEnvironmentVariable("PATH", $newUserPath, "User")
  
  # Restart your terminal for changes to take effect
```

**Just copy and run the command!** ✅

---

### Reload PATH in Current Session

After modifying PATH, reload it in your current PowerShell session:

```powershell
$env:PATH = [Environment]::GetEnvironmentVariable("PATH", "User") + ";" + [Environment]::GetEnvironmentVariable("PATH", "Machine")
```

Then verify:
```powershell
where.exe go
# Should show Gopher's symlink FIRST:
# C:\Users\YourName\AppData\Local\bin\go.exe
# C:\Program Files\Go\bin\go.exe
```

---

## 🧪 Testing Your Setup

### Complete Test

```powershell
# 1. Check gopher version
gopher version
# Should show: gopher v1.0.0

# 2. Check directories exist
Test-Path $env:USERPROFILE\gopher\versions
Test-Path $env:USERPROFILE\gopher\downloads
Test-Path $env:USERPROFILE\gopher\state
# All should return: True

# 3. Install a version
gopher install 1.21.0
# Should show single-line progress bar ✅

# 4. Switch to it
gopher use 1.21.0
# Should create symlink + warn about PATH if needed

# 5. Verify
go version
# Should show: go version go1.21.0 windows/amd64

# 6. Check symlink exists
Test-Path $env:LOCALAPPDATA\bin\go.exe
# Should return: True

# 7. Test switching
gopher use system
go version
# Should show system Go (if installed)

gopher use 1.21.0
go version
# Should show: go version go1.21.0
```

---

## 🔍 Troubleshooting

### Issue 1: "Access is denied" when creating symlinks

**Solution A: Enable Developer Mode** (see Step 1 above)

**Solution B: Run as Administrator**:
```powershell
# Right-click PowerShell and select "Run as administrator"
gopher use 1.21.0
```

---

### Issue 2: "Still using system version" after switching

**Cause:** System Go appears before Gopher in PATH.

**Solution:** Gopher v1.0.0+ automatically detects this and shows you the exact fix command!

Just run `gopher use <version>` and follow the instructions it provides.

---

### Issue 3: "gopher: command not found"

**Cause:** `$env:USERPROFILE\bin` not in PATH.

**Solution:**
```powershell
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$env:USERPROFILE\bin;$userPath", "User")

# Restart PowerShell
```

---

### Issue 4: "go: command not found" after switching

**Cause:** `%LOCALAPPDATA%\bin` not in PATH.

**Solution:** Run `gopher use <version>` again - it will show you the exact PATH command to run.

---

### Issue 5: Symlinks not working

**Check if Developer Mode is enabled:**
```powershell
# If this fails, Developer Mode is OFF:
gopher use 1.21.0

# Enable it in Settings or run as Administrator
```

---

### Issue 6: "Directories not created"

**Solution:** Gopher v1.0.0+ automatically creates directories on first run!

Just run any gopher command and directories will be created:
```powershell
gopher init
# or
gopher install 1.21.0
```

---

## 📚 Common Commands

```powershell
# List installed versions
gopher list

# List available versions
gopher list-remote

# Install a version
gopher install 1.21.0

# Switch to a version
gopher use 1.21.0

# Switch to system Go
gopher use system

# Show current version
gopher current

# Show system Go info
gopher system

# Create an alias
gopher alias create stable 1.21.0

# Use an alias
gopher use stable

# Check setup status
gopher status

# Debug information
gopher debug
```

---

## 🎯 Best Practices

### 1. Use Developer Mode
- Enables symlinks without admin privileges
- Makes version switching seamless
- Required for best experience

### 2. Fix PATH Order
- If you have system Go, put Gopher's bin directory **first** in PATH
- Follow the warning message from `gopher use`
- Verify with `where.exe go`

### 3. Use Aliases
- Create aliases for frequently used versions
- Example: `gopher alias create stable 1.21.0`
- Switch with: `gopher use stable`

### 4. Keep Clean
- Gopher auto-cleans old downloads
- Uninstall unused versions: `gopher uninstall 1.20.0`

---

## 🆘 Getting Help

### Built-in Help

```powershell
# General help
gopher help

# Command-specific help
gopher help install
gopher help use

# Check system status
gopher status

# Debug information
gopher debug
```

### Check PATH Configuration

```powershell
# Show current PATH
$env:PATH -split ';'

# Check where go is found
where.exe go

# Check Gopher's symlink
Test-Path $env:LOCALAPPDATA\bin\go.exe
```

### Reset and Start Fresh

See: `docs/internal/WINDOWS_RESET_GUIDE.md` for complete reset instructions.

**Quick reset:**
```powershell
Remove-Item "$env:USERPROFILE\bin\gopher.exe" -Force
Remove-Item "$env:USERPROFILE\gopher" -Recurse -Force
Remove-Item "$env:LOCALAPPDATA\bin\go.exe" -Force

# Clean PATH
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$cleanedPath = ($userPath -split ';' | Where-Object { 
    $_ -ne "$env:USERPROFILE\bin" -and 
    $_ -ne "$env:LOCALAPPDATA\bin"
}) -join ';'
[Environment]::SetEnvironmentVariable("PATH", $cleanedPath, "User")

# Restart PowerShell and start fresh
```

---

## 🔒 Security Notes

### Symlink Security

- Symlinks are created in `%LOCALAPPDATA%\bin` (user-specific)
- No system-wide changes
- Safe and isolated

### Download Verification

- All downloads are from official Go sources
- SHA256 checksums verified
- Secure HTTPS connections

---

## 📖 Additional Resources

- **User Guide:** `docs/USER_GUIDE.md`
- **FAQ:** `docs/FAQ.md`
- **Windows Reset:** `docs/internal/WINDOWS_RESET_GUIDE.md`
- **Windows PATH Reload:** `docs/internal/WINDOWS_PATH_RELOAD.md`
- **GitHub:** https://github.com/molmedoz/gopher

---

## 💡 Pro Tips

### Tip 1: Use Windows Terminal

Windows Terminal provides better Unicode support and performance:
- Download from Microsoft Store
- Better progress bar display
- Nicer colors and formatting

### Tip 2: Reload PATH Without Restarting

```powershell
$env:PATH = [Environment]::GetEnvironmentVariable("PATH", "User") + ";" + [Environment]::GetEnvironmentVariable("PATH", "Machine")
```

### Tip 3: Check Your Setup Anytime

```powershell
gopher status
gopher debug
```

### Tip 4: Verify PATH Order

```powershell
where.exe go
# Gopher's symlink should be FIRST
```

---

## 🎉 Success!

You're all set! Gopher is now installed and configured on Windows.

**Next steps:**
1. Install Go versions: `gopher install 1.21.0`
2. Switch between them: `gopher use 1.21.0`
3. Enjoy seamless Go version management! ✅

---

**For more help:** Run `gopher help` or visit https://github.com/molmedoz/gopher

**Last Updated:** October 2025  
**Version:** v1.0.0+
