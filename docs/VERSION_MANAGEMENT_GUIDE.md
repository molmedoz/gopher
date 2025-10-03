# Version Management Guide

**How Gopher automatically manages version information**

---

## ✅ **Current Status: IMPLEMENTED**

Gopher v1.0.0 uses **automated version management** with build-time ldflags injection. Version information is automatically extracted from git tags during build - **no manual updates required!**

```go
// cmd/gopher/main.go - Version variables injected at build time
var (
    appVersion = "dev"     // Set via: -X main.appVersion=v1.0.0
    appCommit  = "none"    // Set via: -X main.appCommit=abc123
    appDate    = "unknown" // Set via: -X main.appDate=2025-10-13
    appBuiltBy = "source"  // Set via: -X main.appBuiltBy=goreleaser
)
```

**Benefits of this approach:**
- ✅ Zero manual updates - version from git tags automatically
- ✅ Zero human errors - impossible to mismatch version and tag
- ✅ Detailed build info - commit hash, build date, builder
- ✅ Dev-friendly - shows "dev" when built locally
- ✅ Production-ready - integrated with GoReleaser and Makefile

---

## ✅ **Solutions Comparison**

| Solution | Pros | Cons | Recommended |
|----------|------|------|-------------|
| **1. Build-time ldflags** | ✅ No code changes needed<br>✅ Works with GoReleaser<br>✅ Zero runtime overhead | ⚠️ Requires build flags | ✅ **BEST** |
| **2. VERSION file + embed** | ✅ Human readable<br>✅ Easy to update | ❌ Extra file to maintain<br>❌ Can still forget to update | ⚠️ OK |
| **3. Git tags at runtime** | ✅ Always accurate | ❌ Requires .git directory<br>❌ Slower<br>❌ Doesn't work in dist | ❌ Not recommended |
| **4. Auto-generate VERSION** | ✅ Fully automated | ❌ Complex setup | ⚠️ Overkill |

---

## 🚀 **How It Works: Build-time ldflags**

### **Implementation Overview:**

Gopher uses **ldflags** to inject version information at build time:

```bash
go build -ldflags "-X main.appVersion=v1.0.0" ./cmd/gopher
```

The binary contains the version - **no code changes needed for releases!**

### **Current Implementation:**

#### **main.go Variables:**

```go
// cmd/gopher/main.go - Version information injected at build time
// Example: go build -ldflags "-X main.appVersion=v1.0.0 -X main.appCommit=abc123"
var (
    appVersion = "dev"     // Version (from git tag)
    appCommit  = "none"    // Git commit hash
    appDate    = "unknown" // Build date
    appBuiltBy = "source"  // Built by (goreleaser, manual, etc.)
)

// getVersionString returns the formatted version string
func getVersionString() string {
    if appVersion == "dev" {
        return "gopher dev (built from source)"
    }
    return fmt.Sprintf("gopher %s", appVersion)
}

// getFullVersionInfo returns detailed version information
func getFullVersionInfo() string {
    return fmt.Sprintf(`gopher %s
  commit: %s
  built: %s
  by: %s
  go: %s
  platform: %s/%s`, 
        appVersion, appCommit, appDate, appBuiltBy, 
        inruntime.Version(), inruntime.GOOS, inruntime.GOARCH)
}
```

#### **Makefile Integration:**

```makefile
# Extract version from git tags
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILD_BY := manual

# Build flags
LDFLAGS := -X main.appVersion=$(VERSION) \
           -X main.appCommit=$(COMMIT) \
           -X main.appDate=$(BUILD_DATE) \
           -X main.appBuiltBy=$(BUILD_BY)

# Build command
build:
	go build -ldflags "$(LDFLAGS)" -o build/gopher ./cmd/gopher
```

### **GoReleaser Integration:**

`.goreleaser.yml` automatically sets these variables:

```yaml
builds:
  - main: ./cmd/gopher
    binary: gopher
    ldflags:
      - -s -w
      - -X main.appVersion={{.Version}}
      - -X main.appCommit={{.FullCommit}}
      - -X main.appDate={{.Date}}
      - -X main.appBuiltBy=goreleaser
```

### **Benefits:**

✅ **No manual updates** - Version from git tag automatically  
✅ **Zero errors** - Impossible to mismatch version and tag  
✅ **Extra info** - Commit hash, build date, builder  
✅ **Dev builds** - Shows "dev" when built locally  
✅ **Production ready** - GoReleaser + Makefile integrated

---

## ✅ **Testing Different Build Methods**

### **Local Development Build:**
```bash
go build ./cmd/gopher
./gopher version
# Output: gopher dev (built from source)
```

### **Manual Build with Version:**
```bash
go build -ldflags "-X main.appVersion=v1.0.0 -X main.appCommit=abc123" ./cmd/gopher
./gopher version
# Output: gopher v1.0.0
```

### **GoReleaser Build:**
```bash
goreleaser release --snapshot --clean
./dist/gopher_linux_amd64_v1/gopher version
# Output: gopher v1.0.1-next (snapshot version)
```

### **Production Release:**
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# GoReleaser sets: version=v1.0.0, commit=<hash>, date=<timestamp>
```

---

## 🔒 **Validation & Safety**

### **Option 1: Pre-release Validation Script**

Create `scripts/validate-release.sh`:

```bash
#!/bin/bash
# Validates that a release is ready

set -e

# Get the latest tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "none")

echo "🔍 Validating release readiness..."
echo "Latest tag: $LATEST_TAG"

# Check if working directory is clean
if [[ -n $(git status -s) ]]; then
    echo "❌ Error: Working directory not clean. Commit all changes first."
    exit 1
fi

# Check if CHANGELOG.md mentions the version
if ! grep -q "$LATEST_TAG" CHANGELOG.md; then
    echo "⚠️  Warning: CHANGELOG.md doesn't mention $LATEST_TAG"
    echo "   Update CHANGELOG.md before releasing!"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if tests pass
echo "🧪 Running tests..."
if ! go test ./... ; then
    echo "❌ Error: Tests failed. Fix tests before releasing."
    exit 1
fi

echo "✅ Release validation passed!"
echo ""
echo "To create release, run:"
echo "  git tag -a v1.0.0 -m 'Release v1.0.0'"
echo "  git push origin v1.0.0"
```

### **Option 2: GitHub Actions Validation**

Add to `.github/workflows/validate-tag.yml`:

```yaml
name: Validate Tag

on:
  push:
    tags:
      - 'v*'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Validate tag format
        run: |
          TAG=${GITHUB_REF#refs/tags/}
          if [[ ! $TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
            echo "❌ Invalid tag format: $TAG"
            echo "Expected: v1.0.0 or v1.0.0-beta.1"
            exit 1
          fi
          echo "✅ Tag format valid: $TAG"
      
      - name: Check CHANGELOG
        run: |
          TAG=${GITHUB_REF#refs/tags/}
          if ! grep -q "$TAG" CHANGELOG.md; then
            echo "⚠️ Warning: $TAG not found in CHANGELOG.md"
          else
            echo "✅ CHANGELOG.md updated"
          fi
      
      - name: Run tests
        run: go test ./...
```

### **Option 3: Git Hook (Pre-push)**

Create `.git/hooks/pre-push`:

```bash
#!/bin/bash
# Validates before pushing tags

while read local_ref local_sha remote_ref remote_sha; do
    if [[ $local_ref =~ refs/tags/v ]]; then
        TAG=${local_ref#refs/tags/}
        
        echo "🔍 Validating tag: $TAG"
        
        # Check format
        if [[ ! $TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
            echo "❌ Invalid tag format. Use: v1.0.0 or v1.0.0-beta.1"
            exit 1
        fi
        
        # Check CHANGELOG
        if ! grep -q "$TAG" CHANGELOG.md; then
            echo "⚠️  Warning: $TAG not in CHANGELOG.md"
            read -p "Continue? (y/n) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
        
        echo "✅ Tag validation passed"
    fi
done

exit 0
```

---

## 📋 **All Solutions Detailed**

### **Solution 1: Build-time ldflags ⭐ RECOMMENDED**

**How it works:**
```
Developer           Build Process              Binary
─────────────────   ───────────────────────   ─────────────
                    
git tag v1.0.0  →   go build -ldflags      →  gopher version
                    "-X main.version=v1.0.0"   Output: v1.0.0
                    
                    Version injected at 
                    compile time
```

**Pros:**
- ✅ Zero runtime overhead
- ✅ No extra files needed
- ✅ Works perfectly with GoReleaser
- ✅ Fully automated
- ✅ Never gets out of sync

**Cons:**
- None! This is the standard Go approach

**Implementation:**
```go
var (
    appVersion = "dev"  // Injected at build time via ldflags
)
```

---

### **Solution 2: VERSION File + embed**

**How it works:**
```
VERSION file    →    Embedded in binary    →    Read at runtime
────────────         ────────────────────        ──────────────
v1.0.0               //go:embed VERSION           gopher version
                     var version string           Output: v1.0.0
```

**Pros:**
- ✅ Human-readable VERSION file
- ✅ Easy to see current version
- ✅ Can be auto-generated

**Cons:**
- ❌ Extra file to maintain
- ❌ Can forget to update
- ❌ Requires go:embed

**Implementation:**
```go
package main

import _ "embed"

//go:embed VERSION
var embeddedVersion string

func init() {
    if appVersion == "dev" {
        appVersion = strings.TrimSpace(embeddedVersion)
    }
}
```

---

### **Solution 3: Git Tags at Runtime**

**How it works:**
```
Binary runs    →    Executes git command    →    Returns version
────────────        ────────────────────         ──────────────
gopher version      git describe --tags          Output: v1.0.0
```

**Pros:**
- ✅ Always accurate if .git exists
- ✅ No build-time injection needed

**Cons:**
- ❌ Requires .git directory (doesn't work for distributed binaries!)
- ❌ Slower (spawns process)
- ❌ Won't work for installed binaries
- ❌ Requires git to be installed

**Not suitable for distributed binaries** ❌

---

### **Solution 4: Automated VERSION File Updates**

**How it works:**
```
git tag v1.0.0  →  GitHub Action  →  Update VERSION file  →  Commit
```

**Pros:**
- ✅ Fully automated
- ✅ Visible in repository

**Cons:**
- ❌ Creates extra commits
- ❌ Complex to set up
- ❌ Can cause sync issues

---

## ⭐ **Why We Chose: Build-time ldflags**

This is the **standard Go approach** used by major projects (kubectl, docker, hugo, etc.).

### **What You Get:**

1. ✅ **Automatic version from git tags** - No manual updates
2. ✅ **Detailed version info** - Version, commit, date, builder
3. ✅ **Dev vs production** - Shows "dev" when built locally
4. ✅ **Zero maintenance** - Works forever with no changes
5. ✅ **Validation** - Optional pre-push hooks prevent mistakes

### **How Users See It:**

**Production build:**
```bash
$ gopher version
gopher v1.0.0

$ gopher --verbose version
gopher v1.0.0
  commit: abc123def456
  built: 2025-10-13T10:30:00Z
  by: goreleaser
  go: go1.21.0
  platform: darwin/arm64
```

**Local dev build:**
```bash
$ go build ./cmd/gopher
$ ./gopher version
gopher dev (built from source)
```

---

## 📝 **Summary**

Gopher v1.0.0 implements fully automated version management using:
- ✅ Build-time ldflags injection
- ✅ Makefile integration for local builds
- ✅ GoReleaser integration for releases
- ✅ Validation scripts for release safety

**No manual version updates required!** Just tag and release.

---

**Last Updated:** October 2025  
**Status:** ✅ Implemented in v1.0.0  
**Maintainer:** Gopher Development Team

