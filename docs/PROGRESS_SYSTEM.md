# Progress System

The progress system provides user feedback for long-running operations through progress bars and spinners.

## Overview

The `internal/progress` package provides:

- **ProgressBar** - Visual progress indicator for operations with known total size
- **Spinner** - Animated indicator for indeterminate operations
- **Formatters** - Reusable utilities for formatting bytes, speed, and percentages
- **Terminal Writer** - Cross-platform terminal output handling

## Architecture

### Modular Design

The progress system uses a modular architecture with clear separation:

```
ProgressBar / Spinner
  ├── Configuration (ProgressConfig with functional options)
  ├── Calculation (percentage, speed, etc.)
  ├── Rendering (buildLine - format to string)
  └── Output (terminalWriter - cross-platform I/O)
```

**Benefits:**
- ✅ Clear separation of concerns
- ✅ Reusable components (formatters, terminal writer, config)
- ✅ Easy to test each layer
- ✅ Extensible for future needs

## Components

### 1. ProgressBar

Displays a visual progress bar with percentage, bytes, and speed information.

**Basic Usage:**
```go
// Create a progress bar
pb := progress.NewProgressBar(totalBytes, "Downloading file")

// Update progress
pb.Update(currentBytes)

// Finish (shows 100% and adds newline)
pb.Finish()
```

**With Options:**
```go
// Minimal mode (percentage only)
pb := progress.NewProgressBar(total, "Downloading", progress.WithMinimal())

// Custom width
pb := progress.NewProgressBar(total, "Downloading", progress.WithWidth(80))

// Silent mode (no output)
pb := progress.NewProgressBar(total, "Downloading", progress.WithSilent())

// Multiple options
pb := progress.NewProgressBar(total, "Downloading",
    progress.WithWidth(60),
    progress.WithChars("=", "-"),
    progress.WithSpeed(false),
)
```

**Features:**
- ✅ Auto-detects terminal width to prevent line wrapping
- ✅ Updates on same line using carriage return (`\r`)
- ✅ Throttled updates (100ms default) to prevent flickering
- ✅ Automatic speed calculation
- ✅ Human-readable byte formatting
- ✅ Cross-platform (Windows, macOS, Linux)

### 2. Spinner

Animated spinner for operations without known progress.

**Usage:**
```go
// Create and start spinner
spinner := progress.NewSpinner("Loading data")
spinner.Start()

// Do work...
time.Sleep(5 * time.Second)

// Stop (shows checkmark)
spinner.Stop()
```

**With defer pattern:**
```go
spinner := progress.NewSpinner("Processing")
spinner.Start()
defer spinner.Stop()

// Do work... spinner automatically stops at end
```

**Features:**
- ✅ Animated rotating characters (⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏)
- ✅ Runs in background goroutine
- ✅ Shows checkmark (✓) on completion
- ✅ Configurable update frequency
- ✅ Silent mode support

### 3. ProgressWriter

Wraps an io.Writer to automatically update progress during write operations.

**Usage:**
```go
file, _ := os.Create("output.dat")
defer file.Close()

pb := progress.NewProgressBar(totalBytes, "Writing")
pw := progress.NewProgressWriter(file, pb)

// Writes automatically update progress
io.Copy(pw, source)
pb.Finish()
```

**Perfect for:**
- File downloads
- File copying
- Data streaming

### 4. Formatters Package

Reusable formatting utilities in `internal/formatters`:

```go
// Format bytes
formatters.FormatBytes(1048576)  // "1.0 MB"

// Format speed
formatters.FormatSpeed(10485760)  // "10.0 MB/s"

// Format percentage
formatters.FormatPercentage(0.753)  // "75.3%"
```

## Configuration Options

All options use the functional options pattern for flexibility:

| Option | Description | Example |
|--------|-------------|---------|
| `WithWidth(int)` | Set bar width in characters | `WithWidth(80)` |
| `WithUpdateThrottle(duration)` | Minimum time between updates | `WithUpdateThrottle(50*time.Millisecond)` |
| `WithSpeed(bool)` | Show/hide speed display | `WithSpeed(false)` |
| `WithBytes(bool)` | Show/hide byte counts | `WithBytes(false)` |
| `WithChars(filled, empty)` | Custom bar characters | `WithChars("=", "-")` |
| `WithSilent()` | Disable all output | `WithSilent()` |
| `WithMinimal()` | Show only percentage | `WithMinimal()` |
| `WithCustom(func)` | Custom configuration | `WithCustom(func(c *ProgressConfig) {...})` |

## Cross-Platform Behavior

### Terminal Width Detection
- Automatically detects terminal width using `golang.org/x/term`
- Adjusts bar width to prevent line wrapping
- Falls back to 50 character width if detection fails

### Platform-Specific Handling

**Windows:**
- Uses fixed-width padding (120 chars) for reliable clearing
- Prevents terminal artifacts and flickering
- Progress bars and spinners work in CMD and PowerShell

**Unix/Linux/macOS:**
- Uses dynamic padding based on previous line length
- More efficient terminal output
- Works in all common terminal emulators

### TTY Detection
- Only uses carriage return (`\r`) when stdout is a TTY
- Falls back to simple output when piped or redirected
- Prevents broken output in scripts and logs

## Implementation Details

### Terminal Writer
The `terminalWriter` component handles all cross-platform terminal output:

```go
type terminalWriter struct {
    output      io.Writer
    isStdout    bool
    lastLineLen int
}
```

**Key features:**
- Sanitizes lines (removes accidental `\n` or `\r`)
- Platform-specific clearing strategies
- Proper flush handling
- Testable with custom io.Writer

### Throttling
Progress updates are throttled to prevent performance issues:

- Default: 100ms between updates
- Configurable via `WithUpdateThrottle()`
- Always shows final update (100%)

## Usage in Gopher

**Current usage:**

1. **Download progress** (`internal/downloader`)
   ```go
   progressBar := progress.NewProgressBar(fileSize, "Downloading "+filename)
   progressWriter := progress.NewProgressWriter(file, progressBar)
   io.Copy(progressWriter, resp.Body)
   progressBar.Finish()
   ```

2. **Uninstall spinner** (`cmd/gopher`)
   ```go
   spinner := progress.NewSpinner("Uninstalling Go "+version)
   spinner.Start()
   manager.Uninstall(version)
   spinner.Stop()
   ```

3. **Extraction spinner** (`internal/installer`)
   ```go
   spinner := progress.NewSpinner("Extracting archive")
   spinner.Start()
   defer spinner.Stop()
   // extraction logic...
   ```

4. **Metadata spinner** (`internal/installer`)
   ```go
   spinner := progress.NewSpinner("Creating version metadata")
   spinner.Start()
   createMetadata(...)
   spinner.Stop()
   ```

## Testing

### Unit Tests
```bash
# Test all progress components
go test ./internal/progress/...

# Test formatters
go test ./internal/formatters/...
```

### Manual Testing
```bash
# Test progress bar
./build/gopher install 1.23.0

# Test spinners
./build/gopher uninstall 1.23.0
```

## Future Enhancements

See `docs/ROADMAP.md` for planned features:

- **Multi-bar support** - Multiple concurrent progress bars
- **Progress manager** - Coordinate multiple operations
- **More progress types** - Pie chart, dots, percentage-only
- **Structured logging integration** - Log progress to files without interfering with terminal

## Troubleshooting

### Progress Bar Not Updating
- Ensure stdout is a TTY (not piped)
- Check terminal supports ANSI/VT100
- Verify TERM environment variable is set

### Line Wrapping Issues
- Progress bar auto-detects terminal width
- If wrapping occurs, try `WithWidth(30)` for narrower bar
- Or use `WithMinimal()` for percentage-only

### Windows-Specific Issues
- Ensure terminal supports Unicode (for spinner characters)
- CMD: Use `chcp 65001` for UTF-8 support
- PowerShell: Usually works out of the box

---

**Last Updated:** October 15, 2025  
**Version:** 1.0.0  
**Status:** Production Ready

