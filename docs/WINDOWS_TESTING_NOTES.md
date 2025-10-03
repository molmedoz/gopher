# Windows Testing Guide

## Overview

This guide helps verify that Gopher's progress indicators (progress bars and spinners) work correctly on Windows.

## What's Included in v1.0.0

### Progress Features
- **Progress bars** - Visual download progress with auto-sizing
- **Spinners** - Animated feedback for long operations
- **Cross-platform terminal handling** - Works on Windows CMD, PowerShell, and Windows Terminal

### Technical Implementation
- Uses `golang.org/x/term` for terminal width detection and TTY detection
- Platform-specific output handling (Windows uses fixed-width padding)
- Line sanitization to prevent newline issues
- Smart terminal detection (disables animations when piped)

## 🧪 Windows Testing Checklist

### Prerequisites
```powershell
# Ensure UTF-8 support for Unicode spinner characters
chcp 65001
```

### Test 1: Progress Bar (Download)
```powershell
# Test download with progress bar
.\gopher.exe uninstall 1.23.0
.\gopher.exe install 1.23.0
```

**Expected:**
- ✅ Progress bar updates on **same line** (not new lines)
- ✅ Shows: `[████████░░░░] 60.5% 85.3 MB/141.0 MB 10.2 MB/s`
- ✅ Bar fills up smoothly from left to right
- ✅ No flickering or artifacts

**Potential Issues:**
- ❌ New lines for each update → Width detection may have failed
- ❌ No progress bar → Check if stdout is redirected
- ❌ Garbled output → Try `chcp 65001` for UTF-8

### Test 2: Spinners (Uninstall, Extract, Metadata)
```powershell
.\gopher.exe uninstall 1.23.0
.\gopher.exe install 1.23.0
```

**Expected:**
```
⠹ Uninstalling Go 1.23.0
✓ Uninstalling Go 1.23.0

Downloading... [progress bar]

Installing Go 1.23.0
⠹ Extracting archive
✓ Extracting archive
⠹ Creating version metadata
✓ Creating version metadata
✓ Successfully installed Go 1.23.0
```

**Spinner should:**
- ✅ Animate (rotate through: ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏)
- ✅ Update on same line
- ✅ Show checkmark (✓) on completion
- ✅ Each spinner on its own line

**Potential Issues:**
- ❌ Static spinner (doesn't rotate) → Terminal might not support Unicode
- ❌ Square boxes instead of spinner → Use `chcp 65001`
- ❌ Overlapping spinners → Report as bug

### Test 3: Terminal Width Detection
```powershell
# Test in different window sizes
# Resize your terminal and run:
.\gopher.exe install 1.23.0

# Should adjust bar width automatically
```

### Test 4: Piped Output (Silent Mode)
```powershell
# Progress bar should gracefully degrade when piped
.\gopher.exe install 1.23.0 > output.txt
type output.txt

# Should see simple text, not ANSI codes
```

## 🐛 Known Windows-Specific Behaviors

### CMD vs PowerShell
- **PowerShell**: ✅ Better Unicode support, recommended
- **CMD**: ⚠️ Requires `chcp 65001` for spinner characters
- **Windows Terminal**: ✅ Best experience, full Unicode support

### Terminal Emulators
- **Windows Terminal**: ✅ Fully supported
- **ConEmu**: ✅ Should work
- **Git Bash**: ✅ Should work (Unix-like)
- **Old CMD**: ⚠️ May not support Unicode spinners

## 📊 Performance

### Expected Behavior
- **Progress updates**: Every 100ms (throttled)
- **Spinner updates**: Every 100ms (configurable)
- **CPU usage**: Negligible (<1%)
- **Memory**: ~100KB for progress structures

### If Performance Issues
- Use `WithUpdateThrottle(500*time.Millisecond)` to update less frequently
- Use `WithSilent()` to disable all output
- Use `WithMinimal()` for lighter output

## 🔍 Debugging

### Enable Verbose Output
```powershell
# See what's happening
.\gopher.exe -v install 1.23.0
```

### Check Terminal Capabilities
```powershell
# Check terminal size
mode con

# Check code page (should be 65001 for UTF-8)
chcp
```

### If Progress Bar Creates New Lines
1. Check terminal width: `mode con` - if very narrow, bar may wrap
2. Try minimal mode: modify code to use `WithMinimal()`
3. Check if output is redirected
4. Verify `golang.org/x/term` is installed: `go list -m golang.org/x/term`

## 📝 Reporting Issues

If you encounter issues on Windows, please report with:

1. **Windows Version**: Windows 10/11, build number
2. **Terminal**: CMD, PowerShell, Windows Terminal, etc.
3. **Code Page**: Output of `chcp`
4. **Terminal Size**: Output of `mode con`
5. **Issue Type**: New lines, static spinner, garbled output, etc.
6. **Screenshot or recording** if possible

## ✨ What Should Work

Based on the macOS testing:

✅ **Progress bar**:
- Updates on same line
- Shows progress, speed, bytes
- Auto-adjusts to terminal width
- No line wrapping

✅ **Spinners**:
- Animate smoothly
- Show completion checkmark
- Each on separate line
- No overlap

✅ **Overall UX**:
- Clean, professional output
- No debug messages
- Proper spacing
- Clear feedback

---

## 🎯 Windows-Specific Code

The progress system includes Windows-specific handling:

**Fixed-width padding** (`terminal.go`):
```go
if runtime.GOOS == "windows" {
    // On Windows, use fixed 120-char width for reliable clearing
    tw.printWindowsLine(line, 120)
}
```

**Benefits**:
- Prevents terminal artifacts
- More reliable line clearing
- Works with CMD's quirks

---

**Happy Testing on Windows!** 🎉

If the progress bar and spinners work on Windows as they do on macOS, the refactoring is a complete success!

