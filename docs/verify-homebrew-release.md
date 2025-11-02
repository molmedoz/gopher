# 🍺 Verifying Homebrew Release

This guide helps you verify that your Gopher Homebrew release is working correctly.

## Prerequisites

- macOS or Linux with Homebrew installed
- Internet connection
- A release must have been created with the `release-homebrew` job successful

## Current Status Check

First, let's check if the formula exists:

```bash
# Quick status check
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb 2>&1 | head -5

# If you see "404: Not Found", the formula hasn't been created yet.
# This means either:
# 1. No release has been created yet
# 2. The release workflow's `release-homebrew` job failed
# 3. The job hasn't completed yet (check GitHub Actions)
```

## Step-by-Step Verification

### 1. Check if Tap Repository Exists

```bash
# Check if the tap repository exists on GitHub
curl -s https://api.github.com/repos/molmedoz/homebrew-tap | jq '.name'

# Expected output: "homebrew-tap"
```

### 2. Check if Formula File Exists

```bash
# Check if the formula file exists
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | head -20

# Expected: Should show the Ruby formula content
```

### 3. Add the Tap

```bash
# Add your Homebrew tap
brew tap molmedoz/tap

# Expected output:
# ==> Tapping molmedoz/tap
# Cloning into '/opt/homebrew/Library/Taps/molmedoz/homebrew-tap'...
# Tapped 1 formula.
```

### 4. Check Available Versions

```bash
# Check what versions are available
brew info molmedoz/tap/gopher

# Expected output should show:
# - Version information
# - Installation instructions
# - Dependencies (if any)
```

### 5. Test Installation

```bash
# Install gopher from your tap
brew install molmedoz/tap/gopher

# Expected output:
# ==> Installing gopher from molmedoz/tap
# ==> Downloading https://github.com/molmedoz/gopher/releases/download/vX.X.X/...
# ...
# ==> Summary
# 🍺  /opt/homebrew/Cellar/gopher/X.X.X: X files, XX.XMB
```

### 6. Verify Installation Works

```bash
# Check that gopher is installed and working
gopher version

# Expected output: Should show the version number

# Test a basic command
gopher list

# Expected output: Should show available/installed Go versions
```

### 7. Verify Binary Location

```bash
# Check where gopher was installed
which gopher

# Expected: Should show Homebrew path
# Example: /opt/homebrew/bin/gopher (macOS ARM)
# Example: /usr/local/bin/gopher (macOS Intel or Linux)

# Check version matches
gopher version

# Compare with what's in the formula
brew info molmedoz/tap/gopher | grep -E "version|stable"
```

### 8. Test Upgrade Path

```bash
# Check current version
gopher version

# Uninstall and reinstall to test upgrade
brew uninstall gopher
brew install molmedoz/tap/gopher

# Or test upgrade command (if already installed)
brew upgrade molmedoz/tap/gopher
```

## Troubleshooting

### Issue: Tap not found

**Error:** `Error: No available formula with the name "molmedoz/tap/gopher"`

**Solutions:**
1. Check if tap repository exists:
   ```bash
   curl -s https://api.github.com/repos/molmedoz/homebrew-tap
   ```

2. Verify the repository has the `Formula/gopher.rb` file:
   ```bash
   curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb
   ```

3. Check GitHub Actions workflow ran successfully:
   - Go to: https://github.com/molmedoz/gopher/actions
   - Look for `release-homebrew` job in the latest release workflow

### Issue: Formula download fails

**Error:** `Error: Failed to download resource "gopher"`

**Solutions:**
1. Verify GitHub release exists and has assets:
   ```bash
   # Replace vX.X.X with your version
   curl -s https://api.github.com/repos/molmedoz/gopher/releases/tags/v1.0.0 | jq '.assets[] | .name'
   ```

2. Check if the download URL in formula is correct:
   ```bash
   curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | grep url
   ```

3. Verify the SHA256 checksum matches:
   ```bash
   # Get checksum from formula
   curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | grep sha256
   
   # Compare with release checksums
   curl -s https://github.com/molmedoz/gopher/releases/download/vX.X.X/checksums.txt
   ```

### Issue: Installation succeeds but binary doesn't work

**Solutions:**
1. Check binary permissions:
   ```bash
   ls -la $(which gopher)
   # Should show: -rwxr-xr-x
   ```

2. Verify binary is executable:
   ```bash
   file $(which gopher)
   # Should show: Mach-O binary (macOS) or ELF binary (Linux)
   ```

3. Check for architecture mismatch:
   ```bash
   uname -m  # Your architecture
   file $(which gopher)  # Binary architecture
   # Should match (arm64 vs x86_64)
   ```

## Quick Verification Script

Save this as `verify-homebrew.sh`:

```bash
#!/bin/bash

set -e

echo "🔍 Verifying Homebrew release..."
echo ""

# Check tap exists
echo "1️⃣ Checking tap repository..."
if curl -s -f https://api.github.com/repos/molmedoz/homebrew-tap > /dev/null; then
    echo "   ✅ Tap repository exists"
else
    echo "   ❌ Tap repository not found"
    exit 1
fi

# Check formula exists
echo ""
echo "2️⃣ Checking formula file..."
if curl -s -f https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb > /dev/null; then
    echo "   ✅ Formula file exists"
    FORMULA_VERSION=$(curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | grep -m1 'version "' | sed 's/.*version "\(.*\)".*/\1/')
    echo "   📦 Formula version: $FORMULA_VERSION"
else
    echo "   ❌ Formula file not found"
    exit 1
fi

# Tap and check info
echo ""
echo "3️⃣ Tapping repository..."
brew tap molmedoz/tap 2>&1 | head -3

# Get info
echo ""
echo "4️⃣ Checking package info..."
brew info molmedoz/tap/gopher 2>&1 | head -15

# Check if already installed
echo ""
echo "5️⃣ Checking installation..."
if command -v gopher > /dev/null 2>&1; then
    INSTALLED_VERSION=$(gopher version 2>&1 | head -1)
    echo "   ✅ Gopher is installed: $INSTALLED_VERSION"
    echo "   📍 Location: $(which gopher)"
else
    echo "   ⚠️  Gopher is not installed"
    echo "   💡 Run: brew install molmedoz/tap/gopher"
fi

echo ""
echo "✨ Verification complete!"
```

Make it executable and run:
```bash
chmod +x verify-homebrew.sh
./verify-homebrew.sh
```

## Expected Workflow Behavior

After creating a GitHub release:

1. **GitHub Actions workflow triggers:**
   - Workflow: `.github/workflows/create-release.yml`
   - Job: `release-homebrew`

2. **GoReleaser processes Homebrew:**
   - Downloads release assets
   - Calculates SHA256 checksums
   - Generates/updates `Formula/gopher.rb`
   - Commits to `molmedoz/homebrew-tap` repository

3. **Verification (5-10 minutes after release):**
   - Formula should be available in tap
   - `brew tap molmedoz/tap` should work
   - `brew install molmedoz/tap/gopher` should work

## Manual Formula Check

If you want to manually inspect the formula:

```bash
# Clone the tap repository
git clone https://github.com/molmedoz/homebrew-tap.git /tmp/homebrew-tap
cd /tmp/homebrew-tap

# View the formula
cat Formula/gopher.rb

# Check recent commits
git log --oneline -10

# Verify formula syntax
brew audit --tap molmedoz/tap/gopher
```

## Next Steps After Verification

Once verified, you can:

1. **Update documentation** with installation instructions
2. **Announce** the Homebrew availability
3. **Test on multiple machines/platforms**
4. **Monitor** for user issues

## Related Files

- `.goreleaser.yml` - GoReleaser configuration for Homebrew
- `.github/workflows/create-release.yml` - Release workflow
- `docs/PACKAGE_MANAGER_SETUP.md` - Setup documentation

