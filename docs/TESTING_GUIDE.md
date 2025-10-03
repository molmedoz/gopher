# Gopher Testing Guide
**Multi-Platform Testing & Verification Strategy**

---

## 📊 Current Status

### ✅ Build Status
```bash
✓ Project builds successfully
✓ All tests pass (12 packages)
✓ No compilation errors
✓ No linter warnings
```

### 📁 Project Structure
```
gopher/
├── cmd/gopher/main.go          # CLI entry point (2232 lines)
├── internal/
│   ├── version/                # Version management (48 files)
│   ├── config/                 # Configuration
│   ├── downloader/             # Download logic
│   ├── installer/              # Installation logic
│   ├── security/               # Security checks
│   ├── errors/                 # Error handling
│   ├── logger/                 # Structured logging
│   └── ...
├── docker/                     # Docker test environments
└── test/                       # Integration tests
```

---

## 🎯 Testing Strategy Overview

### **Three-Tier Testing Approach**

1. **Automated Tests** (Already passing ✅)
   - Unit tests
   - Integration tests
   - Security tests

2. **Docker Container Tests** (Pre-configured ✅)
   - Linux (Ubuntu, Debian, Alpine)
   - macOS simulation
   - Windows simulation

3. **Native Platform Tests** (Manual/CI)
   - Real Linux machines
   - Real macOS machines
   - Real Windows machines (VM)

---

## 🐳 Docker Testing (Recommended First Step)

### **Available Docker Test Scenarios**

#### **Linux/Unix Tests**
```bash
# Test on Unix-like system WITHOUT pre-installed Go
make docker-test-unix-no-go

# Test on Unix-like system WITH pre-installed Go
make docker-test-unix-with-go
```

#### **macOS Tests** (Linux-based simulation)
```bash
# Test on macOS-like environment WITHOUT pre-installed Go
make docker-test-macos-no-go

# Test on macOS-like environment WITH pre-installed Go
make docker-test-macos-with-go
```

#### **Windows Tests** (Wine-based simulation)
```bash
# Test Windows build WITHOUT pre-installed Go
make docker-test-windows-no-go

# Test Windows build WITH pre-installed Go
make docker-test-windows-with-go

# Test Windows in simulated environment
make docker-test-windows-simulated
```

#### **Run All Docker Tests**
```bash
# Build all Docker images and run all tests
make docker-build
make docker-test
```

---

## 🔬 Comprehensive Testing Plan

### **Phase 1: Quick Verification (10 minutes)**

#### **1. Build Verification**
```bash
# Build the project
make build

# Verify binary exists
ls -lh build/gopher

# Test basic commands
./build/gopher version
./build/gopher help
```

#### **2. Unit Tests**
```bash
# Run all unit tests
make test

# Run with coverage
make test-coverage

# Run verbose tests
make test-verbose
```

#### **3. Security Checks**
```bash
# Run security scan
make security-scan

# Check for vulnerabilities
make vuln-check

# Run all security checks
make security-all
```

---

### **Phase 2: Docker Testing (30-60 minutes)**

#### **Step 1: Build Docker Images**
```bash
# Build all Docker test images
make docker-build

# This will create:
# - gopher-test-unix-no-go
# - gopher-test-unix-with-go
# - gopher-test-macos-no-go
# - gopher-test-macos-with-go
# - gopher-test-windows-no-go
# - gopher-test-windows-with-go
```

#### **Step 2: Run Automated Docker Tests**
```bash
# Test each scenario
make docker-test-unix-no-go
make docker-test-unix-with-go
make docker-test-macos-no-go
make docker-test-macos-with-go
make docker-test-windows-no-go
make docker-test-windows-with-go

# Or run all at once
make docker-test
```

#### **Step 3: Manual Docker Testing**

**For Linux/Unix:**
```bash
# Enter the container
docker run -it gopher-test-unix-no-go /bin/bash

# Inside container, test commands
/gopher/build/gopher version
/gopher/build/gopher list
/gopher/build/gopher list-remote --page-size 5
/gopher/build/gopher install 1.21.0  # Test installation
/gopher/build/gopher list            # Verify installed
/gopher/build/gopher use 1.21.0      # Test version switching
/gopher/build/gopher current         # Verify active version
/gopher/build/gopher system          # Check system Go
/gopher/build/gopher alias create stable 1.21.0  # Test aliases
/gopher/build/gopher alias list
/gopher/build/gopher uninstall 1.21.0  # Test uninstall
exit
```

**For macOS Simulation:**
```bash
# Enter the container
docker run -it gopher-test-macos-no-go /bin/bash

# Run the same test commands as above
```

**For Windows Simulation:**
```bash
# Enter the container (if interactive mode works)
docker run -it gopher-test-windows-simulated /bin/bash

# Test Windows-specific paths and behavior
```

---

### **Phase 3: Native Platform Testing**

#### **🐧 Linux Testing (Native or VM)**

**Recommended Distributions:**
- Ubuntu 22.04 LTS (most common)
- Debian 12
- Fedora 39
- Arch Linux (bleeding edge)

**Test Script for Linux:**
```bash
#!/bin/bash
# save as test-linux.sh

echo "=== Testing Gopher on Linux ==="

# 1. Build test
echo "Step 1: Building..."
make build
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"

# 2. Basic commands test
echo "Step 2: Testing basic commands..."
./build/gopher version
./build/gopher help | head -20
echo "✅ Basic commands work"

# 3. List remote versions
echo "Step 3: Testing list-remote..."
./build/gopher list-remote --no-interactive --page-size 5
echo "✅ List remote works"

# 4. Install a version
echo "Step 4: Testing install..."
./build/gopher install 1.21.0
if [ $? -ne 0 ]; then
    echo "❌ Install failed"
    exit 1
fi
echo "✅ Install successful"

# 5. List installed versions
echo "Step 5: Testing list..."
./build/gopher list --no-interactive
echo "✅ List installed works"

# 6. Switch version
echo "Step 6: Testing use..."
./build/gopher use 1.21.0
echo "✅ Version switch works"

# 7. Check current version
echo "Step 7: Testing current..."
./build/gopher current
echo "✅ Current version works"

# 8. Test aliases
echo "Step 8: Testing aliases..."
./build/gopher alias create stable 1.21.0
./build/gopher alias list
./build/gopher alias show stable
echo "✅ Aliases work"

# 9. System detection
echo "Step 9: Testing system detection..."
./build/gopher system
echo "✅ System detection works"

# 10. Cleanup
echo "Step 10: Testing uninstall..."
./build/gopher uninstall 1.21.0
./build/gopher alias remove stable
echo "✅ Cleanup successful"

echo ""
echo "=== ✅ All Linux tests passed ==="
```

**Run the test:**
```bash
chmod +x test-linux.sh
./test-linux.sh
```

---

#### **🍎 macOS Testing (Native)**

**Required:**
- macOS 11 (Big Sur) or later
- Xcode Command Line Tools installed

**Test Script for macOS:**
```bash
#!/bin/bash
# save as test-macos.sh

echo "=== Testing Gopher on macOS ==="

# 1. Build test
echo "Step 1: Building..."
make build
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"

# 2. Check for Homebrew Go
echo "Step 2: Checking for Homebrew Go..."
which go
go version 2>/dev/null || echo "No system Go found"

# 3. Basic commands test
echo "Step 3: Testing basic commands..."
./build/gopher version
./build/gopher help | head -20
echo "✅ Basic commands work"

# 4. System detection (macOS specific)
echo "Step 4: Testing system detection..."
./build/gopher system --json
./build/gopher debug
echo "✅ System detection works"

# 5. Test Homebrew detection
echo "Step 5: Testing Homebrew Go detection..."
if command -v brew &> /dev/null; then
    echo "Homebrew detected"
    ./build/gopher system | grep -i homebrew || echo "No Homebrew Go detected"
fi

# 6. Install a version
echo "Step 6: Testing install..."
./build/gopher install 1.21.0
if [ $? -ne 0 ]; then
    echo "❌ Install failed"
    exit 1
fi
echo "✅ Install successful"

# 7. Test version switching
echo "Step 7: Testing version switching..."
./build/gopher use 1.21.0
./build/gopher current
echo "✅ Version switching works"

# 8. Test shell integration
echo "Step 8: Testing shell integration..."
./build/gopher init
./build/gopher status
echo "✅ Shell integration works"

# 9. Test persistence
echo "Step 9: Testing persistence..."
# This requires a new shell session
echo "Manual test required: Open new terminal and run 'go version'"

# 10. Cleanup
echo "Step 10: Cleanup..."
./build/gopher uninstall 1.21.0
echo "✅ Cleanup successful"

echo ""
echo "=== ✅ All macOS tests passed ==="
echo ""
echo "⚠️  MANUAL TEST REQUIRED:"
echo "1. Run: ./build/gopher install 1.21.0"
echo "2. Run: ./build/gopher use 1.21.0"
echo "3. Run: ./build/gopher setup"
echo "4. Open a NEW terminal window"
echo "5. Run: go version"
echo "6. Verify it shows go1.21.0"
```

**Run the test:**
```bash
chmod +x test-macos.sh
./test-macos.sh
```

---

#### **🪟 Windows Testing (Native or VM)**

**Required:**
- Windows 10/11
- Git Bash or WSL2 (for bash scripts)
- OR PowerShell (for PowerShell scripts)

**Test Script for Windows (PowerShell):**
```powershell
# save as test-windows.ps1

Write-Host "=== Testing Gopher on Windows ===" -ForegroundColor Cyan

# 1. Build test
Write-Host "Step 1: Building..." -ForegroundColor Yellow
go build -o build\gopher.exe cmd\gopher\main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Build successful" -ForegroundColor Green

# 2. Basic commands test
Write-Host "Step 2: Testing basic commands..." -ForegroundColor Yellow
.\build\gopher.exe version
.\build\gopher.exe help | Select-Object -First 20
Write-Host "✅ Basic commands work" -ForegroundColor Green

# 3. List remote versions
Write-Host "Step 3: Testing list-remote..." -ForegroundColor Yellow
.\build\gopher.exe list-remote --no-interactive --page-size 5
Write-Host "✅ List remote works" -ForegroundColor Green

# 4. Install a version
Write-Host "Step 4: Testing install..." -ForegroundColor Yellow
.\build\gopher.exe install 1.21.0
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Install failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Install successful" -ForegroundColor Green

# 5. List installed versions
Write-Host "Step 5: Testing list..." -ForegroundColor Yellow
.\build\gopher.exe list --no-interactive
Write-Host "✅ List installed works" -ForegroundColor Green

# 6. Switch version
Write-Host "Step 6: Testing use..." -ForegroundColor Yellow
.\build\gopher.exe use 1.21.0
Write-Host "✅ Version switch works" -ForegroundColor Green

# 7. Check current version
Write-Host "Step 7: Testing current..." -ForegroundColor Yellow
.\build\gopher.exe current
Write-Host "✅ Current version works" -ForegroundColor Green

# 8. Test aliases
Write-Host "Step 8: Testing aliases..." -ForegroundColor Yellow
.\build\gopher.exe alias create stable 1.21.0
.\build\gopher.exe alias list
.\build\gopher.exe alias show stable
Write-Host "✅ Aliases work" -ForegroundColor Green

# 9. System detection
Write-Host "Step 9: Testing system detection..." -ForegroundColor Yellow
.\build\gopher.exe system
Write-Host "✅ System detection works" -ForegroundColor Green

# 10. Test Windows paths
Write-Host "Step 10: Testing Windows-specific paths..." -ForegroundColor Yellow
.\build\gopher.exe debug
Write-Host "✅ Windows paths work" -ForegroundColor Green

# 11. Cleanup
Write-Host "Step 11: Cleanup..." -ForegroundColor Yellow
.\build\gopher.exe uninstall 1.21.0
.\build\gopher.exe alias remove stable
Write-Host "✅ Cleanup successful" -ForegroundColor Green

Write-Host ""
Write-Host "=== ✅ All Windows tests passed ===" -ForegroundColor Green
```

**Run the test:**
```powershell
# In PowerShell
.\test-windows.ps1
```

---

## 🧪 Test Scenarios by Platform

### **Core Functionality Tests (All Platforms)**

| Test | Command | Expected Result |
|------|---------|-----------------|
| Version | `gopher version` | Shows version number |
| Help | `gopher help` | Shows help text |
| List Remote | `gopher list-remote --page-size 5` | Shows 5 versions |
| Install | `gopher install 1.21.0` | Downloads and installs |
| List Installed | `gopher list` | Shows installed versions |
| Use Version | `gopher use 1.21.0` | Switches to 1.21.0 |
| Current | `gopher current` | Shows active version |
| System | `gopher system` | Shows system Go info |
| Alias Create | `gopher alias create test 1.21.0` | Creates alias |
| Alias List | `gopher alias list` | Lists aliases |
| Alias Show | `gopher alias show test` | Shows alias details |
| Alias Remove | `gopher alias remove test` | Removes alias |
| Uninstall | `gopher uninstall 1.21.0` | Removes version |

---

### **Platform-Specific Tests**

#### **Linux-Specific**
- [ ] Test with system package manager Go (apt, yum, dnf)
- [ ] Test with manually installed Go
- [ ] Test bash/zsh shell integration
- [ ] Test with different distributions
- [ ] Test symlink creation
- [ ] Test PATH management

#### **macOS-Specific**
- [ ] Test with Homebrew Go detection
- [ ] Test with system Go (Xcode)
- [ ] Test zsh shell integration (default in modern macOS)
- [ ] Test bash shell integration
- [ ] Test fish shell integration
- [ ] Test PATH management in different shells
- [ ] Test with Apple Silicon (ARM64)
- [ ] Test with Intel (AMD64)

#### **Windows-Specific**
- [ ] Test with official Go installer
- [ ] Test with Chocolatey Go
- [ ] Test with Scoop Go
- [ ] Test PowerShell integration
- [ ] Test CMD integration
- [ ] Test Git Bash integration
- [ ] Test WSL2 integration
- [ ] Test Windows PATH management
- [ ] Test Windows Registry entries (if any)
- [ ] Test UNC paths
- [ ] Test spaces in paths

---

## 🔍 What to Verify on Each Platform

### **Installation Verification**
```bash
# After installing a version
1. Check directory exists: ~/.gopher/versions/go1.21.0/
2. Check binary exists: ~/.gopher/versions/go1.21.0/bin/go
3. Check binary works: ~/.gopher/versions/go1.21.0/bin/go version
4. Check symlink created: ~/.gopher/current -> versions/go1.21.0
```

### **Version Switching Verification**
```bash
# After switching version
1. Check symlink updated: ~/.gopher/current points to correct version
2. Check go command works: go version shows correct version
3. Check GOROOT set correctly: echo $GOROOT
4. Check GOPATH set correctly: echo $GOPATH
```

### **Persistence Verification**
```bash
# Test persistence across shell sessions
1. Install and activate a version
2. Close terminal
3. Open NEW terminal
4. Run: go version
5. Verify it shows the correct version
```

### **System Go Detection Verification**
```bash
# Test system Go detection
1. Run: gopher system
2. Verify it detects your system Go installation
3. Verify it shows correct path
4. Verify it shows correct version
5. Test: gopher use system
6. Run: go version
7. Verify it uses system Go
```

---

## 🐛 Common Issues to Test

### **Path Issues**
- [ ] Spaces in installation path
- [ ] Unicode characters in path
- [ ] Very long paths (>260 chars on Windows)
- [ ] Symlink loops
- [ ] Permission issues

### **Network Issues**
- [ ] Slow connection
- [ ] Download interruption
- [ ] Proxy settings
- [ ] Firewall blocking

### **Version Conflicts**
- [ ] Multiple Go installations
- [ ] System Go vs Gopher Go
- [ ] Homebrew Go vs Gopher Go (macOS)
- [ ] WSL Go vs Windows Go

### **Shell Integration**
- [ ] Shell profile not found
- [ ] Multiple shell profiles
- [ ] Existing Go PATH entries
- [ ] Permission issues writing to profile

---

## 📝 Test Report Template

After testing, document results using this template:

```markdown
# Gopher Test Report

## Test Information
- **Date:** YYYY-MM-DD
- **Tester:** Your Name
- **Gopher Version:** v1.0.0
- **Platform:** [Linux/macOS/Windows]
- **OS Version:** [e.g., Ubuntu 22.04, macOS 14, Windows 11]

## Test Results Summary
- **Total Tests:** XX
- **Passed:** XX
- **Failed:** XX
- **Skipped:** XX

## Detailed Results

### Core Functionality
- [ ] ✅ Build successful
- [ ] ✅ Version command works
- [ ] ✅ Help command works
- [ ] ✅ List remote works
- [ ] ✅ Install works
- [ ] ✅ List installed works
- [ ] ✅ Use version works
- [ ] ✅ Current works
- [ ] ✅ System detection works
- [ ] ✅ Alias create works
- [ ] ✅ Alias list works
- [ ] ✅ Alias show works
- [ ] ✅ Alias remove works
- [ ] ✅ Uninstall works

### Platform-Specific
- [ ] ✅ Shell integration works
- [ ] ✅ PATH management works
- [ ] ✅ Persistence works
- [ ] ✅ System Go detection works

## Issues Found
1. [Issue description]
   - **Severity:** [Critical/High/Medium/Low]
   - **Steps to reproduce:**
   - **Expected behavior:**
   - **Actual behavior:**

## Notes
[Any additional observations or comments]
```

---

## 🚀 Quick Start Testing Checklist

### **Minimum Viable Test (5 minutes)**
```bash
# 1. Build
make build

# 2. Run basic commands
./build/gopher version
./build/gopher help

# 3. Test one install
./build/gopher install 1.21.0

# 4. Test version switch
./build/gopher use 1.21.0

# 5. Verify
./build/gopher current

# 6. Cleanup
./build/gopher uninstall 1.21.0
```

### **Standard Test (15 minutes)**
Run the platform-specific test script for your OS (see above).

### **Comprehensive Test (30-60 minutes)**
1. Run all Docker tests: `make docker-test`
2. Run platform-specific test script
3. Test manual scenarios (shell integration, persistence)
4. Document results

---

## 📚 Additional Resources

### **Makefile Commands**
```bash
make help              # Show all available commands
make test              # Run unit tests
make test-coverage     # Run tests with coverage
make security-all      # Run all security checks
make docker-test       # Run all Docker tests
make build-all         # Build for all platforms
```

### **Environment Variables**
```bash
GOPHER_INSTALL_DIR     # Override install directory
GOPHER_DOWNLOAD_DIR    # Override download directory
GOPHER_LOG_LEVEL       # Set log level (DEBUG, INFO, WARN, ERROR)
```

---

## ✅ Success Criteria

Your Gopher installation is working correctly if:

1. ✅ All unit tests pass
2. ✅ Binary builds without errors
3. ✅ Can install a Go version
4. ✅ Can switch between versions
5. ✅ Can detect system Go
6. ✅ Aliases work correctly
7. ✅ Shell integration persists across sessions
8. ✅ Uninstall removes versions cleanly

---

**Last Updated:** 2025-10-11
**Version:** 1.0


