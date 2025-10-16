# Package Manager Setup Guide

This guide helps you set up Gopher for distribution through package managers (Homebrew, Chocolatey, Snap).

---

## 📦 Current Status

| Package Manager | Status | Required For v1.0.0 |
|----------------|--------|---------------------|
| **GitHub Releases** | ✅ Ready | ✅ YES - Primary distribution |
| **Linux Packages** | ✅ Ready | ✅ YES - deb, rpm, apk, Arch |
| **Homebrew** | ⚠️ Needs Setup | ⚠️ Recommended |
| **Chocolatey** | ⚠️ Needs Setup | ⚠️ Optional |
| **Snap Store** | ⚠️ Needs Setup | ⚠️ Optional |

**For v1.0.0:** GitHub Releases + Linux packages work immediately. Package managers can be added later.

---

## ✅ What Works NOW (No Setup Needed)

### 1. GitHub Releases
- ✅ Binaries for all platforms (10 combinations)
- ✅ Automatically uploaded on release
- ✅ Users download directly

**Installation:**
```bash
# Download binary for your platform
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_1.0.0_Linux_x86_64.tar.gz
tar -xzf gopher_1.0.0_Linux_x86_64.tar.gz
sudo mv gopher /usr/local/bin/
```

### 2. Linux Packages (deb, rpm, apk, Arch)
- ✅ Created automatically
- ✅ Uploaded to GitHub Releases
- ✅ Users can download and install

**Installation:**
```bash
# Debian/Ubuntu
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_1.0.0_linux_amd64.deb
sudo dpkg -i gopher_1.0.0_linux_amd64.deb

# RedHat/Fedora/CentOS
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_1.0.0_linux_amd64.rpm
sudo rpm -i gopher_1.0.0_linux_amd64.rpm

# Alpine
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_1.0.0_linux_amd64.apk
sudo apk add --allow-untrusted gopher_1.0.0_linux_amd64.apk

# Arch Linux
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_1.0.0_linux_amd64.pkg.tar.zst
sudo pacman -U gopher_1.0.0_linux_amd64.pkg.tar.zst
```

---

## 🍺 Homebrew Setup (macOS/Linux)

### Why Set This Up?
- **Most popular** package manager for macOS
- **Easy user experience**: `brew install gopher`
- **Automatic updates**: `brew upgrade gopher`
- **5 minutes** to set up

### Setup Steps

#### 1. Create Homebrew Tap Repository
```bash
# Using GitHub CLI
gh repo create molmedoz/homebrew-tap --public --description "Homebrew tap for Gopher"

# Initialize the repository
cd /tmp
git clone https://github.com/molmedoz/homebrew-tap.git
cd homebrew-tap
mkdir -p Formula
echo "# Homebrew Tap for Gopher" > README.md
git add .
git commit -m "Initial commit"
git push origin main
```

#### 2. Enable in .goreleaser.yml
```yaml
brews:
  - skip_upload: false  # Change from true to false
```

#### 3. That's it!
GoReleaser will automatically:
- Create/update the Homebrew formula
- Commit to molmedoz/homebrew-tap
- Users can install with: `brew tap molmedoz/tap && brew install gopher`

### Testing
```bash
# After release
brew tap molmedoz/tap
brew install gopher
gopher version
```

---

## 🍫 Chocolatey Setup (Windows)

### Why Set This Up?
- **Most popular** Windows package manager
- **Easy for Windows users**: `choco install gopher`
- **15 minutes** to set up

### Setup Steps

#### 1. Create Chocolatey Account
```
1. Go to: https://community.chocolatey.org/account/register
2. Create account
3. Verify email
```

#### 2. Get API Key
```
1. Go to: https://community.chocolatey.org/account
2. Click "API Keys" tab
3. Copy your API key
```

#### 3. Add API Key to GitHub Secrets
```bash
# Using GitHub CLI
gh secret set CHOCOLATEY_API_KEY

# Paste your API key when prompted
```

#### 4. Enable in .goreleaser.yml
```yaml
chocolateys:
  - skip_publish: false  # Change from true to false
```

#### 5. First Release
The first release will create the package. Chocolatey will review it (24-48 hours).

### Testing
```powershell
# After approval
choco install gopher
gopher version
```

---

## 📦 Snap Store Setup (Linux)

### Why Set This Up?
- **Pre-installed** on Ubuntu and many Linux distros
- **Automatic updates**
- **Sandboxed** installation

### Setup Steps

#### 1. Create Ubuntu One Account
```
1. Go to: https://login.ubuntu.com/
2. Create account
3. Verify email
```

#### 2. Register on Snapcraft
```
1. Go to: https://snapcraft.io/
2. Login with Ubuntu One
3. Click "Publish" → "Register name"
4. Register "gopher"
```

#### 3. Install Snapcraft CLI
```bash
# On Ubuntu/Debian
sudo snap install snapcraft --classic

# Or use Docker
docker run -it snapcore/snapcraft
```

#### 4. Login and Export Credentials
```bash
# Login to Snap Store
snapcraft login

# Export credentials for CI
snapcraft export-login --snaps gopher --channels stable - > snapcraft-credentials.txt

# This creates a base64-encoded credentials file
```

#### 5. Add Credentials to GitHub Secrets
```bash
# Read the credentials file and add to secrets
gh secret set SNAPCRAFT_STORE_CREDENTIALS < snapcraft-credentials.txt

# Clean up
rm snapcraft-credentials.txt
```

#### 6. Enable in .goreleaser.yml
```yaml
snapcrafts:
  - publish: true  # Change from false to true
```

### Testing
```bash
# After release
sudo snap install gopher --classic
gopher version
```

---

## 📋 Setup Verification Checklist

Run this script to check what's configured:

```bash
#!/bin/bash
echo "=== Package Manager Setup Status ==="
echo

echo "1. Homebrew Tap:"
if gh repo view molmedoz/homebrew-tap >/dev/null 2>&1; then
    echo "   ✅ Repository exists"
    echo "   Check .goreleaser.yml: skip_upload should be 'false'"
else
    echo "   ❌ Repository NOT found"
    echo "   Run: gh repo create molmedoz/homebrew-tap --public"
fi
echo

echo "2. Chocolatey:"
if gh secret list | grep -q "CHOCOLATEY_API_KEY"; then
    echo "   ✅ API key configured"
    echo "   Check .goreleaser.yml: skip_publish should be 'false'"
else
    echo "   ❌ API key NOT configured"
    echo "   Get key from: https://community.chocolatey.org/account"
fi
echo

echo "3. Snap Store:"
if gh secret list | grep -q "SNAPCRAFT_STORE_CREDENTIALS"; then
    echo "   ✅ Credentials configured"
    echo "   Check .goreleaser.yml: publish should be 'true'"
else
    echo "   ❌ Credentials NOT configured"
    echo "   Register at: https://snapcraft.io/"
fi
echo

echo "4. GitHub Releases:"
echo "   ✅ Always works (no setup needed)"
echo

echo "5. Linux Packages:"
echo "   ✅ Always works (uploaded to GitHub Releases)"
```

---

## 🎯 Recommended Setup Order

### For v1.0.0 Release

**Minimum (Works Immediately):**
- ✅ GitHub Releases
- ✅ Linux packages (deb, rpm, apk, Arch)

**Recommended (5 minutes):**
1. ✅ Setup Homebrew tap
   - Most users expect `brew install`
   - Easiest to set up
   - Works immediately

**Optional (Post-v1.0.0):**
2. Setup Chocolatey (Windows users)
3. Setup Snap Store (Ubuntu users)

---

## ⚙️ Current .goreleaser.yml Configuration

### What's Enabled:
- ✅ **GitHub Releases**: Always works
- ✅ **Archives**: tar.gz for all platforms
- ✅ **Checksums**: SHA256 for verification
- ✅ **Linux Packages**: deb, rpm, apk, Arch (uploaded to GitHub)
- ✅ **Source Archives**: Complete source code

### What's Disabled (Until Setup):
- ⏸️ **Homebrew**: `skip_upload: true` (needs tap repo)
- ⏸️ **Chocolatey**: `skip_publish: true` (needs API key)
- ⏸️ **Snap**: `publish: false` (needs credentials)

---

## 🚀 Quick Setup Scripts

### Setup All Three (Complete)

```bash
#!/bin/bash
# setup-package-managers.sh

echo "=== Setting up Package Managers for Gopher ==="

# 1. Homebrew Tap
echo
echo "1. Setting up Homebrew tap..."
gh repo create molmedoz/homebrew-tap --public --description "Homebrew tap for Gopher"
cd /tmp
git clone https://github.com/molmedoz/homebrew-tap.git
cd homebrew-tap
mkdir -p Formula
echo "# Homebrew Tap for Gopher" > README.md
echo "Formulas for Gopher Go version manager." >> README.md
git add .
git commit -m "Initial commit"
git push origin main
echo "✅ Homebrew tap created"

# 2. Chocolatey
echo
echo "2. Setting up Chocolatey..."
echo "Visit: https://community.chocolatey.org/account/register"
echo "After creating account and getting API key:"
read -p "Enter your Chocolatey API key: " CHOCO_KEY
gh secret set CHOCOLATEY_API_KEY --body "$CHOCO_KEY"
echo "✅ Chocolatey API key configured"

# 3. Snap Store
echo
echo "3. Setting up Snap Store..."
echo "Visit: https://snapcraft.io/ and register app 'gopher'"
echo "Then run: snapcraft login"
echo "Then run: snapcraft export-login --snaps gopher --channels stable - > /tmp/snap-creds.txt"
read -p "Press Enter after exporting credentials..."
gh secret set SNAPCRAFT_STORE_CREDENTIALS < /tmp/snap-creds.txt
rm /tmp/snap-creds.txt
echo "✅ Snap Store credentials configured"

echo
echo "=== Setup Complete! ==="
echo "Now update .goreleaser.yml:"
echo "- brews[0].skip_upload: false"
echo "- chocolateys[0].skip_publish: false"
echo "- snapcrafts[0].publish: true"
```

### Setup Just Homebrew (Quickest)

```bash
#!/bin/bash
# setup-homebrew-only.sh

echo "Setting up Homebrew tap..."
gh repo create molmedoz/homebrew-tap --public --description "Homebrew tap for Gopher"

# Initialize repository
cd /tmp
git clone https://github.com/molmedoz/homebrew-tap.git
cd homebrew-tap
mkdir -p Formula
cat > README.md << 'EOFMD'
# Homebrew Tap for Gopher

Homebrew formulas for Gopher - Go Version Manager

## Usage

```bash
brew tap molmedoz/tap
brew install gopher
```

## Available Formulas

- **gopher** - Go version manager
EOFMD

git add .
git commit -m "Initial commit"
git push origin main

echo "✅ Homebrew tap ready!"
echo
echo "Next: Update .goreleaser.yml"
echo "  brews[0].skip_upload: false"
```

---

## 🔍 Verification

### Check GoReleaser Config
```bash
# Validate configuration
goreleaser check

# Test build (doesn't publish)
goreleaser release --snapshot --clean

# Check dist/ folder for packages
ls -la dist/
```

### After Release
```bash
# Check Homebrew
brew tap molmedoz/tap
brew install gopher

# Check Chocolatey (if configured)
choco search gopher

# Check Snap (if configured)
snap find gopher
```

---

## 📝 Updating .goreleaser.yml After Setup

### After Homebrew Tap Created:
```yaml
brews:
  - skip_upload: false  # Changed from true
```

### After Chocolatey API Key Added:
```yaml
chocolateys:
  - skip_publish: false  # Changed from true
```

### After Snap Credentials Added:
```yaml
snapcrafts:
  - publish: true  # Changed from false
```

---

## ⚠️ Important Notes

### Homebrew
- **No secrets needed** - Uses `GITHUB_TOKEN` (auto-available in Actions)
- **Easiest to set up** - Just create the repository
- **Instant publishing** - No approval process

### Chocolatey
- **Needs API key** - Get from chocolatey.org account
- **First package reviewed** - 24-48 hours approval time
- **Future updates** - Automatic (no review)
- **Moderation** - Must follow Chocolatey guidelines

### Snap
- **Needs registration** - Register app name first
- **Credentials expire** - Re-export periodically
- **Confinement** - Using `classic` (full system access needed)
- **Review process** - First snap might be reviewed

---

## 🚀 Recommendation for v1.0.0

### Option A: Release Now (Minimum)
```
✅ GitHub Releases (works immediately)
✅ Linux packages (works immediately)
⏭️ Skip package managers for now
```

**Users install with:**
```bash
# Download from GitHub Releases
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_...
```

### Option B: Add Homebrew (Recommended)
```
✅ GitHub Releases
✅ Linux packages
✅ Homebrew (5 min setup)
⏭️ Chocolatey & Snap later
```

**Setup:**
```bash
# 1. Create tap (5 minutes)
./setup-homebrew-only.sh

# 2. Update .goreleaser.yml
skip_upload: false

# 3. Release!
```

**Users install with:**
```bash
brew tap molmedoz/tap
brew install gopher
```

### Option C: Full Setup (All Three)
```
✅ GitHub Releases
✅ Linux packages
✅ Homebrew
✅ Chocolatey  
✅ Snap
```

**Time:** ~30 minutes total
**Users:** Can install via any package manager

---

## 📊 User Distribution Estimate

Based on typical Go developer demographics:

| Platform | % of Users | Preferred Method |
|----------|-----------|------------------|
| macOS | ~40% | Homebrew |
| Linux | ~40% | apt/dnf/pacman or direct download |
| Windows | ~20% | Chocolatey or direct download |

**Priority:**
1. **Homebrew** (covers 40% of users easily)
2. **Linux packages** (already work via GitHub)
3. **Chocolatey** (Windows users)
4. **Snap** (subset of Linux users)

---

## 🎯 My Recommendation

### For v1.0.0:
**Setup Homebrew ONLY (5 minutes)**

Why:
- Covers largest user segment (macOS developers)
- Easiest to set up
- No secrets/credentials needed
- Instant publishing
- Professional installation experience

### Post-v1.0.0:
Add Chocolatey and Snap based on user requests.

---

## 📝 What to Document in README

### If Only GitHub Releases:
```markdown
## Installation

### Binary Download
Download from [GitHub Releases](https://github.com/molmedoz/gopher/releases)
```

### If Homebrew Added:
```markdown
## Installation

### macOS/Linux (Homebrew)
```bash
brew tap molmedoz/tap
brew install gopher
```

### Other Platforms
Download from [GitHub Releases](https://github.com/molmedoz/gopher/releases)
```

### If All Package Managers:
```markdown
## Installation

### macOS/Linux
```bash
brew tap molmedoz/tap
brew install gopher
```

### Windows
```powershell
choco install gopher
```

### Linux (Snap)
```bash
sudo snap install gopher --classic
```

### Linux (Packages)
Download .deb, .rpm, or .apk from [Releases](https://github.com/molmedoz/gopher/releases)
```

---

## ✅ Quick Decision Matrix

**Want to release v1.0.0 TODAY?**
- Use GitHub Releases only (no setup needed)

**Want professional installation experience?**
- Setup Homebrew tap (5 minutes)

**Want Windows users to have easy install?**
- Setup Chocolatey (15 minutes, includes approval wait)

**Want maximum coverage?**
- Setup all three (30 minutes + approval waits)

---

**I recommend: Setup Homebrew tap now (5 min), release v1.0.0, add others in v1.1.0 based on user demand.** 🚀

