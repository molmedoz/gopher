# Windows Testing Guide for Gopher

This guide explains how to test Gopher on Windows systems with different configurations.

## 🪟 **Windows Testing Options**

### **Option 1: Native Windows Testing (Recommended)**

This is the most accurate way to test Gopher on Windows.

#### Prerequisites
- Windows 10/11 or Windows Server 2019/2022
- Go installed on the system (optional, for testing system Go integration)
- PowerShell 5.1+ or PowerShell Core 7+

#### Steps

1. **Build Windows Binary**:
   ```bash
   # From your development machine (macOS/Linux)
   make build-all
   # This creates: build/gopher-windows-amd64.exe
   ```

2. **Transfer to Windows**:
   - Copy `build/gopher-windows-amd64.exe` to your Windows machine
   - Rename it to `gopher.exe`
   - Place it in a directory in your PATH (e.g., `C:\Windows\System32` or create a `C:\gopher\bin` directory)

3. **Test Basic Functionality**:
   ```cmd
   # Test gopher version
   gopher version
   
   # Test help
   gopher help
   
   # Test system Go detection (if Go is installed)
   gopher system
   
   # Test listing versions
   gopher list
   ```

4. **Test with System Go** (if Go is installed):
   ```cmd
   # Check current Go version
   go version
   
   # Test gopher system detection
   gopher system
   
   # Test switching to system Go
   gopher use system
   
   # Verify switch worked
   gopher current
   go version
   ```

5. **Test Installation and Switching**:
   ```cmd
   # Install a Go version
   gopher install 1.20.7
   
   # List installed versions
   gopher list
   
   # Switch to installed version
   gopher use 1.20.7
   
   # Verify switch
   gopher current
   go version
   
   # Switch back to system
   gopher use system
   ```

6. **Test Environment Setup**:
   ```cmd
   # Setup shell integration
   gopher setup
   
   # Check environment variables
   gopher env list
   
   # Show environment for specific version
   gopher env show go1.20.7
   ```

### **Option 2: Docker with Windows Containers**

This requires Docker Desktop with Windows container support enabled.

#### Prerequisites
- Docker Desktop with Windows containers enabled
- Windows 10 Pro/Enterprise or Windows Server 2016+

#### Steps

1. **Enable Windows Containers**:
   - Right-click Docker Desktop icon
   - Select "Switch to Windows containers"
   - Wait for Docker to restart

2. **Build Windows Binary**:
   ```bash
   make build-all
   ```

3. **Build Windows Container**:
   ```bash
   # Build the Windows native container
   docker build -f docker/Dockerfile.windows-native -t gopher-windows-native .
   ```

4. **Run Windows Tests**:
   ```bash
   # Run the Windows container
   docker run --rm gopher-windows-native
   
   # Or run interactively
   docker run -it gopher-windows-native cmd
   ```

### **Option 3: WSL2 with Windows Binary**

Test the Windows binary in WSL2 environment.

#### Prerequisites
- Windows 10/11 with WSL2 enabled
- Ubuntu or other Linux distribution in WSL2

#### Steps

1. **Build Windows Binary**:
   ```bash
   make build-all
   ```

2. **Copy to WSL2**:
   ```bash
   # From WSL2
   cp /mnt/c/path/to/gopher-windows-amd64.exe ./gopher.exe
   chmod +x gopher.exe
   ```

3. **Test in WSL2**:
   ```bash
   # Test basic functionality
   ./gopher.exe version
   ./gopher.exe help
   ```

## 🧪 **Test Scenarios**

### **Scenario 1: Windows with System Go**

Test Gopher when Go is already installed on Windows.

**Setup**:
- Install Go from https://golang.org/dl/
- Add Go to PATH
- Install Gopher

**Tests**:
```cmd
# 1. Verify system Go detection
gopher system

# 2. List versions (should show system Go)
gopher list

# 3. Test current version
gopher current

# 4. Install additional version
gopher install 1.20.7

# 5. Switch between versions
gopher use 1.20.7
gopher current
gopher use system
gopher current

# 6. Test environment setup
gopher setup
gopher env list
```

### **Scenario 2: Windows without System Go**

Test Gopher when no Go is installed on Windows.

**Setup**:
- Fresh Windows installation
- No Go installed
- Install Gopher

**Tests**:
```cmd
# 1. Verify no system Go
gopher system

# 2. List versions (should be empty)
gopher list

# 3. Install Go version
gopher install 1.21.0

# 4. List versions (should show installed version)
gopher list

# 5. Switch to installed version
gopher use 1.21.0

# 6. Verify Go is working
go version
```

### **Scenario 3: Windows with Multiple Go Versions**

Test Gopher with multiple Go versions installed.

**Setup**:
- Install Go 1.20.7 and 1.21.0 via Gopher
- Have system Go installed

**Tests**:
```cmd
# 1. List all versions
gopher list

# 2. Switch between versions
gopher use 1.20.7
gopher current
go version

gopher use 1.21.0
gopher current
go version

gopher use system
gopher current
go version

# 3. Test environment variables
gopher env show go1.20.7
gopher env show go1.21.0
gopher env show system
```

## 🔧 **Troubleshooting Windows Issues**

### **Common Issues**

1. **"gopher is not recognized as an internal or external command"**
   - **Solution**: Add gopher.exe to your PATH or use full path
   - **Fix**: `set PATH=%PATH%;C:\path\to\gopher`

2. **"Access is denied" when creating symlinks**
   - **Solution**: Run Command Prompt as Administrator
   - **Alternative**: Use `gopher setup` to configure environment variables instead of symlinks

3. **"Go is not recognized" after switching versions**
   - **Solution**: Restart Command Prompt or run `gopher setup`
   - **Check**: Verify PATH includes Go binary location

4. **PowerShell execution policy error**
   - **Solution**: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`
   - **Alternative**: Use Command Prompt instead of PowerShell

### **Debug Commands**

```cmd
# Check gopher version
gopher version

# Check system Go
gopher system

# List all versions
gopher list

# Check current version
gopher current

# Check environment
gopher env list

# Check Go installation
go version
where go

# Check PATH
echo %PATH%
```

## 📊 **Expected Results**

### **Windows with System Go**

```
C:\> gopher list
Installed Go versions:
  go1.21.5    [system]  ✅ active
  go1.20.7    [gopher]  ⚪ inactive

C:\> gopher current
Current Go version: go1.21.5 (system)

C:\> go version
go version go1.21.5 windows/amd64
```

### **Windows without System Go**

```
C:\> gopher list
Installed Go versions:
  go1.21.0    [gopher]  ✅ active

C:\> gopher current
Current Go version: go1.21.0

C:\> go version
go version go1.21.0 windows/amd64
```

## 🚀 **Quick Test Script**

Create a `test-windows.bat` file:

```batch
@echo off
echo Testing Gopher on Windows...
echo.

echo 1. Gopher version:
gopher version
echo.

echo 2. System Go detection:
gopher system
echo.

echo 3. List versions:
gopher list
echo.

echo 4. Current version:
gopher current
echo.

echo 5. Go version:
go version
echo.

echo 6. Install Go 1.20.7:
gopher install 1.20.7
echo.

echo 7. List after installation:
gopher list
echo.

echo 8. Switch to 1.20.7:
gopher use 1.20.7
echo.

echo 9. Current after switch:
gopher current
echo.

echo 10. Go version after switch:
go version
echo.

echo 11. Switch back to system:
gopher use system
echo.

echo 12. Current after switch back:
gopher current
echo.

echo 13. Go version after switch back:
go version
echo.

echo Testing completed!
pause
```

## 📝 **Reporting Issues**

When reporting Windows-specific issues, include:

1. **Windows Version**: `winver` command output
2. **Go Version**: `go version` output
3. **Gopher Version**: `gopher version` output
4. **Error Messages**: Full error output
5. **Steps to Reproduce**: Exact commands run
6. **Environment**: PATH, GOROOT, GOPATH values

## 🎯 **Success Criteria**

Windows testing is successful when:

- ✅ Gopher installs and runs without errors
- ✅ System Go is detected correctly (if present)
- ✅ Go versions can be installed and switched
- ✅ Environment variables are set correctly
- ✅ Shell integration works (if configured)
- ✅ JSON output works for scripting
- ✅ All commands work in both CMD and PowerShell

---

**Last Updated**: January 2024  
**Tested On**: Windows 10/11, Windows Server 2019/2022
