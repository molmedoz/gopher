# Makefile Guide

**Version**: v1.0.0  
**Last Updated**: October 15, 2025

Complete guide to Gopher's Makefile commands for development and CI/CD.

---

## 🎯 **Quick Start**

```bash
# First time setup
make setup          # Install dependencies and tools

# Before committing
make ci             # Run all CI checks locally

# If all green, push!
git push
```

---

## 📋 **Command Reference**

### **Build Commands**

| Command | Description | Time |
|---------|-------------|------|
| `make build` | Build binary | ~3s |
| `make build-all` | Build for all platforms | ~15s |
| `make install` | Install to system (/usr/local/bin) | ~3s |
| `make clean` | Clean build artifacts | ~1s |

**Example**:
```bash
make build
./build/gopher version
```

---

### **Testing Commands**

| Command | Description | Race Detection | Time |
|---------|-------------|----------------|------|
| `make test` | Tests + coverage | ❌ No | ~9s |
| `make test-verbose` | Verbose output | ❌ No | ~9s |
| `make test-race` | With race detection | ✅ Yes | ~18s |

**Recommendation**: Use `make test` for fast feedback, `make test-race` when working on concurrent code.

**Example**:
```bash
# Fast test
make test

# Comprehensive (check for race conditions)
make test-race
```

---

### **Code Quality Commands**

#### **Formatting**:
| Command | Description | Fixes Issues |
|---------|-------------|--------------|
| `make fmt` | Format code | ✅ Yes |
| `make fmt-check` | Check formatting | ❌ No (CI mode) |

#### **Imports**:
| Command | Description | Fixes Issues |
|---------|-------------|--------------|
| `make imports` | Fix imports | ✅ Yes |
| `make imports-check` | Check imports | ❌ No (CI mode) |

#### **Linting**:
| Command | Description | Time |
|---------|-------------|------|
| `make staticcheck` | Run staticcheck | ~3s |
| `make lint` | Run golangci-lint | ~8s |
| `make lint-ci` | All lint checks | ~6s |
| `make vet` | Run go vet | ~2s |

**Example**:
```bash
# Fix formatting and imports
make fmt imports

# Check if ready for CI
make lint-ci
```

---

### **CI Commands** ⭐

| Command | Description | Includes | Time |
|---------|-------------|----------|------|
| `make ci` | **ALL CI checks** | lint + vet + test + build | ~20s |
| `make ci-race` | CI with race detection | lint + vet + test-race + build | ~40s |

**Recommendation**: Use `make ci` before every push.

**What it runs**:
1. ✅ Format check
2. ✅ Import check
3. ✅ Staticcheck
4. ✅ Go vet
5. ✅ Tests with coverage
6. ✅ Build

**Example**:
```bash
# Before committing
make ci

# If all green:
git add .
git commit -m "feat: add feature"
git push
```

---

### **Release Commands**

| Command | Description | Requires |
|---------|-------------|----------|
| `make release` | Build for all platforms | Tests pass |
| `make dist` | Create distribution packages | `make release` |
| `make goreleaser-check` | Validate GoReleaser config | GoReleaser installed |
| `make goreleaser-snapshot` | Test release build | GoReleaser installed |
| `make validate-release TAG=vX.Y.Z` | Validate release | Scripts in place |
| `make create-tag TAG=vX.Y.Z` | Create and push tag | Git clean |

**Example**:
```bash
# Prepare release
make release
make dist

# Create release tag
make create-tag TAG=v1.0.0
```

---

### **Utility Commands**

| Command | Description |
|---------|-------------|
| `make deps` | Download dependencies |
| `make tidy` | Tidy go.mod and go.sum |
| `make version` | Show version info |
| `make install-tools` | Install dev tools (goimports, staticcheck, etc.) |
| `make setup` | Complete development setup |

**Example**:
```bash
# New contributor setup
make setup

# Verify setup
make ci
```

---

### **Security Commands**

| Command | Description | Tool |
|---------|-------------|------|
| `make security-scan` | Security scan | staticcheck + vet |
| `make vuln-check` | Vulnerability check | govulncheck |
| `make security-all` | All security checks | All above |

**Example**:
```bash
# Check for vulnerabilities
make vuln-check

# Full security audit
make security-all
```

---

### **Docker Commands**

| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker images |
| `make docker-test` | Run tests in Docker |
| `make docker-clean` | Clean Docker images |

**Note**: For specific scenarios, use `./docker/build.sh <scenario>`

---

## 🔄 **Common Workflows**

### **1. New Contributor Setup**
```bash
git clone https://github.com/molmedoz/gopher
cd gopher
make setup          # Install dependencies and tools
make ci             # Verify everything works
```

### **2. Daily Development**
```bash
# Make changes
vim internal/runtime/manager.go

# Quick test
make test

# Format and check
make fmt imports
make ci

# If all green, commit
git add .
git commit -m "fix: improve manager"
git push
```

### **3. Before Creating PR**
```bash
# Run full CI locally
make ci

# Check for race conditions
make test-race

# If both pass, create PR
git push origin feature-branch
```

### **4. Debugging CI Failure**
```bash
# CI failed? Reproduce locally:
make ci

# Check specific step
make fmt-check      # Formatting issue?
make imports-check  # Import issue?
make lint-ci        # Lint issue?
make test           # Test failure?

# Fix and verify
make fmt imports
make ci
```

### **5. Preparing Release**
```bash
# Run all checks
make ci
make test-race

# Build release
make release

# Create tag
make create-tag TAG=v1.0.0
```

---

## 🎨 **Output Colors**

Makefile uses colored output:
- 🔵 **Blue**: Operation starting
- 🟢 **Green**: Success
- 🟡 **Yellow**: Warning/Info
- 🔴 **Red**: Error
- 🔷 **Cyan**: Additional info

---

## ⚡ **Performance Tips**

### **Fast Iteration**:
```bash
# Fastest check
make fmt test

# Skip slow linters during iteration
make fmt imports test
```

### **Comprehensive Before Push**:
```bash
# Full CI suite
make ci             # ~20s

# Or with race detection
make ci-race        # ~40s
```

### **Parallel Testing** (if needed):
```bash
# Run tests in parallel
go test -parallel 4 ./...
```

---

## 🐛 **Troubleshooting**

### **golangci-lint not installed**
```bash
# See: https://golangci-lint.run/usage/install/
# macOS:
brew install golangci-lint

# Linux:
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### **Tests failing on clean repo**
```bash
make clean
make deps
make test
```

### **Coverage report not generated**
```bash
make test
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

---

## 📊 **CI Pipeline Mapping**

| CI Step | Makefile Command | What It Does |
|---------|------------------|--------------|
| Lint | `make lint-ci` | fmt-check + imports-check + staticcheck |
| Test | `make test` | Tests with coverage (NO race) |
| Build | `make build` | Build binary |
| **Full CI** | `make ci` | All of the above |

**To reproduce CI locally**: `make ci`

---

## 🎯 **Recommended Workflow**

```bash
# 1. Setup (once)
make setup

# 2. Develop
vim cmd/gopher/main.go

# 3. Quick check
make test

# 4. Before commit
make fmt imports
make ci

# 5. Commit and push
git commit -am "feat: add feature"
git push

# 6. CI runs automatically with same commands!
```

---

## 📚 **Advanced Usage**

### **Custom Build**:
```bash
# With custom version
VERSION=1.0.0 make build

# With custom flags
LDFLAGS="-X main.version=custom" make build
```

### **Selective Testing**:
```bash
# Test specific package
go test ./internal/runtime -v

# Test specific function
go test ./internal/runtime -run TestAliasManager -v
```

### **Coverage Analysis**:
```bash
# Generate coverage
make test

# View in browser
open coverage.html

# View in terminal
go tool cover -func=coverage.out
```

---

## ✅ **CI vs Local Development**

| Use Case | Command | Race Detection | Time |
|----------|---------|----------------|------|
| Quick iteration | `make test` | ❌ No | ~9s |
| Before commit | `make ci` | ❌ No | ~20s |
| Working on concurrency | `make test-race` | ✅ Yes | ~18s |
| Comprehensive check | `make ci-race` | ✅ Yes | ~40s |
| **CI Pipeline** | `make ci` | ❌ No | ~20s |

**Philosophy**: Fast feedback by default, comprehensive when needed.

---

## 🔗 **Related Documentation**

- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Development Guidelines](docs/DEVELOPMENT_GUIDELINES.md)** - Best practices
- **[Testing Guide](docs/TESTING_GUIDE.md)** - How to write tests
- **[CI/CD Setup](docs/internal/CI_CD_SETUP.md)** - Pipeline documentation

---

**Version**: v1.0.0  
**Status**: Production ready ✅

