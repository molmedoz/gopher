# Package Distribution Quick Start

**Get Gopher on Homebrew, Chocolatey, and Linux package managers in 30 minutes!**

---

## ⚡ **TL;DR - What You Need**

1. ✅ GoReleaser installed
2. ✅ GitHub Personal Access Token  
3. ✅ Homebrew tap repository created
4. ✅ Chocolatey account (optional, can add later)
5. ✅ Snapcraft account (optional, can add later)
6. ✅ Tag your release: `git tag v1.0.0 && git push origin v1.0.0`
7. ✅ **Done!** GitHub Actions publishes everywhere automatically

---

## 🚀 **30-Minute Setup**

### **Minute 0-5: Install GoReleaser**

```bash
# macOS
brew install goreleaser

# Linux
brew install goreleaser
# Or: curl -sfL https://goreleaser.com/static/run | bash

# Windows
choco install goreleaser
# Or: scoop install goreleaser

# Verify
goreleaser --version
```

### **Minute 5-10: Create Homebrew Tap**

1. Go to https://github.com/new
2. Repository name: `homebrew-tap`
3. Public repository
4. Click "Create repository"

**That's it for Homebrew setup!** GoReleaser will populate it automatically.

### **Minute 10-15: Configure GitHub Secrets**

#### **Create Personal Access Token:**

1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Name: `HOMEBREW_TAP_TOKEN`
4. Scopes: Select `repo` (all)
5. Click "Generate token"
6. **Copy the token!**

#### **Add to Repository:**

1. Go to your gopher repository
2. Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `HOMEBREW_TAP_TOKEN`
5. Value: Paste your token
6. Click "Add secret"

### **Minute 15-20: Verify Configuration**

```bash
# Test GoReleaser config
cd /path/to/gopher
goreleaser check

# Test build (doesn't publish)
goreleaser release --snapshot --clean

# Check output
ls dist/
```

### **Minute 20-30: Create First Release!**

```bash
# Make sure everything is committed
git status
git add -A
git commit -m "feat: ready for first release"
git push

# Create and push tag
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0

# Watch GitHub Actions
# Go to: https://github.com/molmedoz/gopher/actions
```

**Done!** In 5-10 minutes:
- ✅ GitHub Release created
- ✅ Binaries built for all platforms
- ✅ Homebrew tap updated
- ✅ Users can `brew install molmedoz/tap/gopher`

---

## 📦 **What's Automated**

After pushing a tag, GitHub Actions automatically:

1. ✅ **Builds** for all platforms (Linux, macOS, Windows, ARM, etc.)
2. ✅ **Creates** GitHub Release with binaries
3. ✅ **Updates** Homebrew formula in your tap
4. ✅ **Generates** checksums and signatures
5. ✅ **Creates** changelog from commits
6. ✅ **Publishes** to Chocolatey (if configured)
7. ✅ **Publishes** to Snap Store (if configured)

**Zero manual work needed!**

---

## 🍫 **Adding Chocolatey (Optional)**

### **Quick Setup:**

```bash
# 1. Create account at https://community.chocolatey.org/
# 2. Get API key from account settings
# 3. Add to GitHub Secrets as: CHOCOLATEY_API_KEY
```

**Users can then install:**
```powershell
choco install gopher
```

**Note:** First submission to Chocolatey requires manual approval (can take 24-48 hours). Subsequent releases are automatic!

---

## 🐧 **Adding Snap (Optional)**

### **Quick Setup:**

```bash
# 1. Create account at https://snapcraft.io/
# 2. Reserve name "gopher"

# 3. Export credentials
snapcraft login
snapcraft export-login snapcraft-credentials

# 4. Add file contents to GitHub Secrets as: SNAPCRAFT_TOKEN
cat snapcraft-credentials
```

**Users can then install:**
```bash
sudo snap install gopher --classic
```

---

## 📋 **Release Checklist**

### **Every Release:**

```bash
# 1. Update version in code (if applicable)
# 2. Update CHANGELOG.md
vim CHANGELOG.md

# 3. Commit changes
git add CHANGELOG.md
git commit -m "docs: update changelog for v1.0.0"
git push

# 4. Create tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 5. Push tag (triggers automation!)
git push origin v1.0.0

# 6. Watch GitHub Actions
# Visit: https://github.com/molmedoz/gopher/actions

# 7. Verify release
# Visit: https://github.com/molmedoz/gopher/releases

# 8. Test installation
brew upgrade gopher  # macOS
choco upgrade gopher  # Windows
```

**That's it!** Everything else is automated.

---

## 🎯 **What Users Get**

### **macOS:**
```bash
# One command (recommended):
brew install molmedoz/tap/gopher

# Or two steps:
# brew tap molmedoz/tap
# brew install gopher
```

### **Windows:**
```powershell
choco install gopher
```

### **Linux (Snap):**
```bash
sudo snap install gopher --classic
```

### **Linux (Direct):**
```bash
# Download from releases
wget https://github.com/molmedoz/gopher/releases/download/v1.0.0/gopher_1.0.0_Linux_x86_64.tar.gz
tar xzf gopher_1.0.0_Linux_x86_64.tar.gz
sudo mv gopher /usr/local/bin/
```

---

## ⚠️ **Important Notes**

### **First Time:**
- Homebrew: Instant (tap updates immediately)
- Chocolatey: **24-48 hour approval** for first package
- Snap: **Manual review** for first publication

### **Subsequent Releases:**
- All are **instant and automatic!**
- Just push a tag, automation handles the rest

### **Versioning:**
- Use semantic versioning: `v1.0.0`, `v1.1.0`, `v1.1.1`
- Always prefix with `v`: `v1.0.0` (not `1.0.0`)

---

## 🔗 **Files Created**

All configuration is ready:

- ✅ `.goreleaser.yml` - GoReleaser configuration
- ✅ `.github/workflows/release.yml` - GitHub Actions workflow
- ✅ `PACKAGE_DISTRIBUTION_GUIDE.md` - Complete guide
- ✅ This quick-start guide

---

## 🎓 **Next Steps**

### **Immediate:**
1. Install GoReleaser: `brew install goreleaser`
2. Create `homebrew-tap` repository on GitHub
3. Create GitHub Personal Access Token
4. Add token to repository secrets
5. Test: `goreleaser release --snapshot --clean`

### **Optional (Can Add Later):**
1. Create Chocolatey account
2. Create Snapcraft account  
3. Add additional package managers

### **When Ready:**
1. Create tag: `git tag -a v0.1.0 -m "First release"`
2. Push tag: `git push origin v0.1.0`
3. Watch the magic happen! ✨

---

**See `PACKAGE_DISTRIBUTION_GUIDE.md` for complete details!**

**Last Updated:** 2025-10-13

