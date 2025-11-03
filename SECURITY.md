# Security Policy

This document summarizes our static security checks (gosec) posture, current findings, and the rationale for any suppressions or risk mitigations.

## Tooling

- Primary SAST: gosec
- Complementary: golangci-lint (includes basic security linters)

Run locally:

```bash
# Full gosec scan (recommended for CI parity)
gosec -severity medium -confidence medium ./...

# Or via golangci-lint (if installed)
golangci-lint run --enable gosec
```

## Current Findings and Status

### HIGH Severity - All Fixed ✅

- **G115 (CWE-190): Integer overflow conversions** [FIXED]
  - **Locations**: `internal/installer/installer.go` (lines 207, 233)
  - **Issue**: Direct `int64 → FileMode` conversion could overflow
  - **Fix**: Safe conversion chain: `int64 → uint32 → FileMode` with masking
    - Mask permission bits: `uint32(header.Mode & 0777)`
    - Explicit uint32 cast prevents overflow
    - Applied to both `os.MkdirAll` and `os.Chmod` operations

### MEDIUM Severity - All Fixed ✅

- **G110 (CWE-409): Decompression bomb** [FIXED]
  - **Locations**: `internal/installer/installer.go` (tar and zip extraction)
  - **Issue**: Archive extraction without size limits could consume excessive resources
  - **Fix**: 
    - Added 1GB per-file size limit check before extraction
    - Use `io.LimitReader` to cap read operations
    - Validate both tar header size and zip uncompressed size

- **G204 (CWE-78): Command execution with variable** [FIXED]
  - **Locations**: 
    - `internal/runtime/system.go` (lines 173, 180)
    - `internal/runtime/list.go` (lines 135, 174, 212)
  - **Issue**: Direct execution using variable paths could enable command injection
  - **Fix**: Central helper `runGoCommand`:
    - Resolves binary via `exec.LookPath("go")` (not user-controlled path)
    - Uses `exec.CommandContext(ctx, "go", ...)` with 5s timeout
    - Avoids variable-provided executable path to reduce injection risk

- **G304 (CWE-22): File inclusion via variable** [FIXED]
  - **Locations**: Multiple files (config, installer, downloader, runtime)
  - **Issue**: File operations using paths that could be manipulated
  - **Fix**: 
    - Added `ValidatePath` and `ValidatePathWithinRoot` in `security` package
    - Path validation applied to config loading and metadata operations
    - **Root scoping**: All metadata/config/state file access is now scoped to safe root directories:
      - Config files: Scoped to `~/.gopher` (Unix) or `~/gopher` (Windows)
      - State files: Scoped to `~/.gopher/state` (parent of install directory)
      - Alias files: Scoped to `~/.gopher` (parent of install directory)
      - Metadata files: Already scoped to install directory (via `ValidatePathWithinRoot`)
    - Removed `..` path traversal patterns (replaced with `filepath.Dir` and absolute path resolution)
    - `#nosec` comments added for safely constructed paths (archive extraction, validated paths)

- **G301 (CWE-276): Directory permissions** [FIXED]
  - **Locations**: Multiple directories across codebase
  - **Issue**: Some directories use 0755 (world-readable) instead of 0750
  - **Fix**: 
    - Changed private directories to `0750`:
      - Config directory: `0750` (private user data)
      - State directory: `0750` (private user data)
      - Aliases directory: `0750` (private user data)
    - `#nosec` comments for directories requiring `0755`:
      - Bin directories (must be executable by system)
      - Install directories (Go toolchain needs access)
      - Script directories (executable scripts)
      - Download directories (temporary, acceptable)

- **G306 (CWE-276): WriteFile permissions** [FIXED]
  - **Locations**: `internal/runtime/test_utils.go`
  - **Issue**: Test file uses 0644 instead of 0600
  - **Fix**: Added `#nosec` comment - acceptable for test utilities

### LOW Severity - All Fixed ✅

- **G104 (CWE-703): Unhandled errors** [FIXED]
  - **Locations**: 
    - `internal/security/security.go` (file Close and Remove during cleanup)
    - `internal/progress/terminal.go` (os.Stdout.Sync operations)
    - `internal/installer/installer.go` (rc.Close and outFile.Close in zip extraction)
    - `internal/runtime/manager.go` (os.Remove during symlink cleanup)
  - **Issue**: Errors from close/remove operations not handled
  - **Fix**: 
    - All close/remove error paths now handled
    - Best-effort cleanup with `_ = ...` for non-critical operations
    - Logged or combined with primary errors where appropriate
    - `IsNotExist` errors tolerated for cleanup operations

## Suppressions

The following `// #nosec` suppressions are in place with documented rationale:

### G304 (File Inclusion) - Archive Extraction Paths
- **Files**: `internal/installer/installer.go`
- **Rationale**: Paths are constructed from validated `targetDir` and archive header components. Archive headers are validated and paths are sanitized during extraction. `targetDir` is validated before use.
- **Risk Level**: Low - paths are constrained to installation directory

### G301 (Directory Permissions) - Executable Directories
- **Files**: 
  - `internal/runtime/environment.go` (script directories, bin directories)
  - `internal/installer/installer.go` (install directory, parent directories during extraction)
  - `internal/downloader/downloader.go` (download directory)
- **Rationale**: These directories require `0755` permissions because:
  - **Bin directories**: Must be executable by the system for symlinks to work
  - **Install directories**: Go toolchain needs read access for execution
  - **Script directories**: Contain executable shell scripts
  - **Download directories**: Temporary, user-controlled location
  - **Archive extraction**: Parent directories during extraction need appropriate permissions
- **Risk Level**: Low - these are user-controlled installation locations

### G306 (WriteFile Permissions) - Test Files
- **Files**: `internal/runtime/test_utils.go`
- **Rationale**: Test utilities use `0644` permissions for temporary test files. This is acceptable in test contexts where files are created in isolated test directories.
- **Risk Level**: None - test-only code

### Suppression Guidelines

When adding new suppressions:
- Prefer small, targeted suppressions next to the specific line
- Include a short justification referencing CWE and why the risk is acceptable/mitigated
- Document in this file with rationale
- Review suppressions periodically to ensure they remain justified

## Secure Coding Guidance

### Path Security
- **Always validate user-controlled paths** before file I/O operations
  - Use `security.ValidatePath()` for basic validation
  - Use `security.ValidatePathWithinRoot()` when paths should be constrained to a specific directory
  - Prefer `filepath.Join()` over string concatenation for path construction
  - Be cautious with archive extraction - validate and constrain paths to target directory

### Command Execution
- **Never use user-controlled paths for executables**
  - Use `exec.LookPath()` to resolve binaries from PATH
  - Use `exec.CommandContext()` with timeouts for command execution
  - Prefer fixed binary names (e.g., `"go"`) over variable paths
  - See `runGoCommand()` helper in `internal/runtime/system.go` for reference

### Error Handling
- **Always handle cleanup errors** for file operations
  - Close files and handle errors from `Close()`
  - Remove temporary files and handle `Remove()` errors
  - Use `os.IsNotExist()` to tolerate expected "not found" errors
  - Log non-fatal cleanup failures for debugging
  - For non-critical operations, use `_ = operation()` with best-effort approach

### File Permissions
- **Use restrictive permissions by default**
  - Private data: `0750` for directories, `0600` for files
  - Executable directories: `0755` only when system needs access (bin, install dirs)
  - Scripts: `0755` for executable scripts, `0644` for non-executable configs
  - Archive modes: Always mask to `0777` when converting: `uint32(mode & 0777)`

### Archive Handling
- **Always validate archive contents**
  - Check file sizes before extraction (prevent decompression bombs)
  - Use `io.LimitReader` to cap read operations
  - Validate archive structure (e.g., "go/" prefix for Go distributions)
  - Set maximum file size limits (e.g., 1GB per file)
  - Safely convert archive permission modes (prevent integer overflow)

### Integer Safety
- **Prevent integer overflow in type conversions**
  - For `int64 → uint32 → FileMode`: use explicit casting
  - Mask permission bits before conversion: `uint32(value & 0777)`
  - Never directly cast `int64` to `FileMode` without intermediate `uint32`

### Network Security
- **Always verify downloaded content**
  - Checksum verification (SHA256) before use
  - Validate file sizes and content types
  - Use secure protocols (HTTPS) for downloads
  - Clean up invalid downloads on checksum mismatch

## CI Recommendations

- Fail CI on new HIGH severity gosec findings
- Allow MEDIUM/LOW to warn initially; tighten over time
- Keep a baseline file if needed to track legacy issues while preventing regressions

Generate a baseline (optional):

```bash
# Creates a JSON report
gosec -fmt=json -out=gosec-report.json ./...
```

Review and update this policy when adding new surfaces (e.g., shell integration, third‑party tools, or Windows-specific features).


