# Windows Usage Guide

Gopher supports version switching on Windows using ZIP downloads and symlinks.

## Quick Start

For detailed setup instructions, see [Windows Setup Guide](WINDOWS_SETUP_GUIDE.md).

### Usage

```bash
# Install a Go version
gopher install 1.21.0

# Switch to a version
gopher use 1.21.0

# List installed versions
gopher list

# Check current version
gopher current
```

### Symlink Locations

Gopher tries to create symlinks in these locations (in order):

1. `%USERPROFILE%\AppData\Local\bin\go.exe`
2. `%USERPROFILE%\bin\go.exe`

### If Symlink Creation Fails

If you get a "symlink creation not implemented for Windows" error, try these solutions:

#### Option 1: Enable Developer Mode (Recommended)
1. Open Windows Settings
2. Go to Update & Security > For developers
3. Turn on "Developer Mode"
4. Run `gopher use <version>` again

#### Option 2: Run as Administrator
1. Right-click Command Prompt or PowerShell
2. Select "Run as administrator"
3. Run `gopher use <version>`

#### Option 3: Manual PATH Setup
1. Add the symlink directory to your PATH environment variable
2. The directory will be shown in the error message
3. Restart your terminal
4. Run `gopher use <version>` again

### Troubleshooting

- **"Access is denied"**: Enable Developer Mode or run as Administrator
- **"Command not found"**: Add the symlink directory to your PATH
- **"Still using system version"**: Restart terminal after PATH changes
- **"No application is associated"**: Run `gopher debug` for diagnostics

### Example Session

```bash
C:\Users\YourName> gopher install 1.21.0
✓ Downloaded go1.21.0.windows-amd64.zip
✓ Installed go1.21.0

C:\Users\YourName> gopher use 1.21.0
✓ Created symlink in C:\Users\YourName\AppData\Local\bin\go.exe
✓ Added C:\Users\YourName\AppData\Local\bin to PATH for current session

C:\Users\YourName> go version
go version go1.21.0 windows/amd64
```

### Features

- ✅ Automatic PATH management
- ✅ No admin privileges required (with Developer Mode)
- ✅ Multiple Go versions can coexist
- ✅ Works with existing Go installations
