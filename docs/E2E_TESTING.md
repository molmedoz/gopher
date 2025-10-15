# End-to-End Testing Guide

**Version**: v1.0.0  
**Last Updated**: October 15, 2025

Complete guide to Gopher's end-to-end testing suite.

---

## 🎯 **Overview**

E2E tests verify that the entire Gopher application works correctly by:
- Building the actual binary
- Running real commands
- Testing on multiple platforms
- Verifying output and behavior

---

## 🚀 **Quick Start**

### **Run E2E Tests Locally**:
```bash
# Build and run E2E tests
make test-e2e

# Run all tests (unit + e2e)
make test-all

# Run with verbose output
go test ./test -run TestE2E -v
```

### **Run on Specific Platform**:
```bash
# Tests run on your current OS/architecture automatically
# For cross-platform testing, use GitHub Actions test-matrix
```

---

## 📋 **Test Coverage**

### **TestE2E_FullWorkflow** - Complete User Journey

Tests all core commands:

| Command | What It Tests |
|---------|---------------|
| `gopher version` | Version display works |
| `gopher help` | Help text displays |
| `gopher list` | Lists installed versions |
| `gopher system` | Shows system Go info |
| `gopher current` | Shows active version |
| `gopher list-remote` | Remote version listing |
| `gopher list-remote --filter` | Version filtering |
| `gopher --json list` | JSON output format |
| `gopher env list` | Environment configuration |
| `gopher alias list` | Alias listing (empty state) |

**Tests**: 10 scenarios  
**Time**: ~2s  
**Platform**: All (Linux, macOS, Windows)

---

### **TestE2E_AliasWorkflow** - Alias System

Tests alias management:

| Command | What It Tests |
|---------|---------------|
| `gopher alias` | Alias help display |
| `gopher alias list` | List aliases (empty/populated) |
| `gopher alias suggest` | Alias name suggestions |

**Tests**: 3 scenarios  
**Time**: ~0.25s  
**Platform**: All

---

### **TestE2E_ErrorHandling** - Error Cases

Tests error handling:

| Scenario | What It Tests |
|----------|---------------|
| Invalid command | Unknown command handling |
| Invalid version | Version validation |
| Missing arguments | Argument validation |
| Reserved names | Name reservation |

**Tests**: 4 scenarios  
**Time**: ~0.16s  
**Platform**: All

---

## 🔧 **Test Architecture**

### **How It Works**:

```go
// 1. Build real gopher binary
buildGopher(t)

// 2. Get binary path (handles .exe on Windows)
gopherPath := getGopherPath(t)

// 3. Set up isolated environment
tmpDir := t.TempDir()
os.Setenv("GOPHER_CONFIG", filepath.Join(tmpDir, "config.json"))

// 4. Run actual commands
output := runGopher(t, gopherPath, "version")

// 5. Verify output
if !strings.Contains(output, "gopher") {
    t.Error("Unexpected output")
}
```

### **Isolation**:
- ✅ Each test uses temporary directory
- ✅ Environment variables isolate configuration
- ✅ No interference with user's Gopher installation
- ✅ Clean slate for every test run

---

## 📊 **Test Matrix (CI)**

### **Platforms Tested**:

| OS | Go Versions | Architecture | Total |
|----|-------------|--------------|-------|
| **Ubuntu** | 1.20, 1.21, 1.22, 1.23 | amd64 | 4 tests |
| **macOS** | 1.21, 1.22, 1.23 | amd64, arm64 | 3 tests |
| **Windows** | 1.21, 1.22, 1.23 | amd64 | 3 tests |
| **Total** | | | **10 combinations** |

**CI Runtime**: ~20 minutes (parallel)

---

## 🎯 **What E2E Tests Cover**

### **✅ Covered**:
- Command execution
- Help text display
- JSON output
- Version listing
- Alias management
- Environment configuration
- Error handling
- Cross-platform compatibility

### **❌ Not Covered** (Too Slow/Large for E2E):
- Actual Go version downloads
- Full version installation
- System integration (symlinks, PATH)
- Shell integration

**Note**: These are covered by integration tests and manual testing

---

## 🔍 **Running Tests**

### **All Tests**:
```bash
# Unit + E2E
make test-all

# Just unit tests
make test

# Just E2E tests
make test-e2e
```

### **Specific E2E Test**:
```bash
# Run specific test
go test ./test -run TestE2E_FullWorkflow -v

# Run specific subtest
go test ./test -run TestE2E_FullWorkflow/version -v
```

### **Skip E2E Tests**:
```bash
# Short mode skips E2E tests
go test -short ./test

# They're slow, so short mode is useful for quick iteration
```

---

## 📝 **Adding New E2E Tests**

### **1. Create Test Function**:
```go
func TestE2E_MyNewFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    buildGopher(t)
    gopherPath := getGopherPath(t)

    // Set up isolated environment
    tmpDir := t.TempDir()
    os.Setenv("GOPHER_CONFIG", filepath.Join(tmpDir, "config.json"))

    t.Run("my_scenario", func(t *testing.T) {
        output := runGopher(t, gopherPath, "my-command", "--my-flag")
        
        if !strings.Contains(output, "expected text") {
            t.Errorf("Unexpected output: %s", output)
        }
    })
}
```

### **2. Run Test**:
```bash
go test ./test -run TestE2E_MyNewFeature -v
```

### **3. Add to Documentation**:
- Update this file with new test coverage
- Add examples if complex

---

## 🐛 **Debugging E2E Tests**

### **Test Fails Locally**:
```bash
# Run with verbose output
go test ./test -run TestE2E_FullWorkflow -v

# Check specific subtest
go test ./test -run TestE2E_FullWorkflow/version -v

# Run gopher command manually
./build/gopher version
```

### **Test Fails in CI**:
```bash
# Reproduce CI environment
make ci

# Check test-matrix logs on GitHub
# Look for the specific OS/Go version combination that failed
```

### **Binary Not Found**:
```bash
# Ensure build directory exists
make build

# Check binary was created
ls -la build/gopher*
```

---

## 📊 **Performance**

| Test Suite | Tests | Time | Can Skip? |
|------------|-------|------|-----------|
| Unit tests | ~50 | ~9s | ❌ No |
| E2E tests | ~17 | ~3s | ✅ Yes (with -short) |
| **Total** | **~67** | **~12s** | |

**CI Impact**: +3s per platform (acceptable)

---

## ✅ **CI Integration**

### **test-matrix.yml** runs E2E tests on:
- ✅ Ubuntu (Linux) - 4 Go versions
- ✅ macOS (Intel + Apple Silicon) - 3 Go versions
- ✅ Windows - 3 Go versions

**Total**: 10 platform/version combinations

### **Workflow**:
```yaml
- name: Run unit tests
  run: go test -v ./...

- name: Build
  run: go build -v -o build/gopher ./cmd/gopher

- name: Run E2E tests
  run: go test ./test -run TestE2E -v
```

---

## 🎯 **Best Practices**

### **DO**:
- ✅ Use isolated environments (tmpDir)
- ✅ Check for expected output patterns
- ✅ Test error cases
- ✅ Use subtests for organization
- ✅ Skip in short mode (`testing.Short()`)
- ✅ Clean up after tests (defer, tmpDir)

### **DON'T**:
- ❌ Download real Go versions (too slow)
- ❌ Modify user's actual Gopher installation
- ❌ Hard-code absolute paths
- ❌ Depend on internet connectivity for basic tests
- ❌ Leave test artifacts behind

---

## 📚 **Test Scenarios**

### **Happy Path**:
```go
// All commands work as expected
gopher version      → Shows version
gopher help         → Shows help
gopher list         → Lists versions
gopher alias list   → Lists aliases
```

### **Edge Cases**:
```go
// Empty states
gopher alias list   → "No aliases found"

// Invalid input
gopher install bad  → Error
gopher alias create system 1.22 → Error (reserved name)
```

### **Cross-Platform**:
```go
// Binary name varies by OS
Windows: gopher.exe
Unix:    gopher

// Paths vary by OS
Windows: %USERPROFILE%\gopher
Unix:    ~/.gopher
```

---

## 🚀 **Future E2E Tests**

Consider adding:
- Full installation workflow (download + install)
- Version switching workflow
- Alias creation → use workflow
- Environment variable configuration
- Shell integration verification

**Tradeoff**: More coverage vs slower tests

---

## ✅ **Summary**

**E2E Tests Provide**:
- ✅ Real-world command execution
- ✅ Cross-platform verification
- ✅ Integration with CI
- ✅ Confidence in release quality

**Current Coverage**:
- ✅ All major commands tested
- ✅ Error handling verified
- ✅ JSON output validated
- ✅ Alias system checked
- ✅ Environment management tested

**Status**: Production ready for v1.0.0 ✅

---

## 🔗 **Related Documentation**

- [Testing Guide](TESTING_GUIDE.md) - Overall testing strategy
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute
- [Makefile Guide](MAKEFILE_GUIDE.md) - All make commands
- [Test Strategy](TEST_STRATEGY.md) - Testing philosophy

---

**Version**: v1.0.0  
**Status**: Production ready ✅

