# 📍 Where Are Releases Published?

## Quick Answer

**GitHub Releases** are published to the **main repository**:
- **Repository**: `molmedoz/gopher` (this repository)
- **URL**: https://github.com/molmedoz/gopher/releases
- **Tab Location**: Go to the repository → Click "Releases" tab (next to "Code", "Issues", etc.)

**Homebrew Formula** is published to a **separate tap repository**:
- **Repository**: `molmedoz/homebrew-tap` (different repository)
- **URL**: https://github.com/molmedoz/homebrew-tap
- **Formula File**: https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb

---

## Where GoReleaser Publishes

### 1. GitHub Releases (Main Repository)

**Location**: `molmedoz/gopher` repository

**What gets published**:
- ✅ Binary archives (`.tar.gz`, `.zip`) for all platforms
- ✅ Checksums (SHA256)
- ✅ Source code archives
- ✅ Linux packages (`.deb`, `.rpm`, `.apk`, `.pkg.tar.zst`)
- ✅ Release notes (auto-generated from CHANGELOG)

**How to view**:
1. Go to: https://github.com/molmedoz/gopher
2. Click the **"Releases"** tab (next to "Code", "Issues", "Pull requests")
3. Or direct link: https://github.com/molmedoz/gopher/releases

**Current releases**:
- `v1.0.1` - Published 2025-11-03
- `v1.0.0` - Published 2025-10-19

**If you don't see the Releases tab**:
- The tab might be disabled in repository settings
- Go to: Settings → General → Features
- Check that "Releases" is enabled
- Or use direct URL: https://github.com/molmedoz/gopher/releases

### 2. Homebrew Formula (Tap Repository)

**Location**: `molmedoz/homebrew-tap` repository (separate repository)

**What gets published**:
- ✅ Ruby formula file: `Formula/gopher.rb`
- ✅ Auto-updated on each release
- ✅ Contains checksums, URLs, install instructions

**How to view**:
1. Go to: https://github.com/molmedoz/homebrew-tap
2. Navigate to: `Formula/gopher.rb`
3. Or direct link: https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb

**Installation**:
```bash
brew tap molmedoz/tap
brew install gopher
```

### 3. Chocolatey (Windows)

**Location**: https://community.chocolatey.org/packages/gopher

**What gets published**:
- ✅ `.nupkg` package file
- ✅ Package metadata
- ✅ Auto-published if `CHOCOLATEY_API_KEY` is set

**Installation**:
```bash
choco install gopher
```

### 4. Scoop (Windows)

**Location**: `molmedoz/scoop-bucket` repository (if it exists)

**What gets published**:
- ✅ Manifest file: `bucket/gopher.json`
- ✅ Auto-published if repository exists

**Installation**:
```bash
scoop bucket add molmedoz https://github.com/molmedoz/scoop-bucket
scoop install gopher
```

---

## Why You Don't See Releases Tab

### Issue 1: Releases Tab is Hidden

**Solution**:
1. Go to: https://github.com/molmedoz/gopher/settings
2. Scroll to "Features" section
3. Make sure "Releases" checkbox is **enabled** ✅
4. Click "Save changes"

### Issue 2: Looking in Wrong Repository

**Releases are in**: `molmedoz/gopher` (main repository)  
**NOT in**: `molmedoz/homebrew-tap` (tap repository has no releases)

**Solution**: Use direct link: https://github.com/molmedoz/gopher/releases

### Issue 3: Releases Tab Doesn't Exist

This can happen if:
- Repository is brand new
- No releases have been created yet
- Releases feature is disabled

**Solution**:
1. Check if releases exist:
   ```bash
   curl -s https://api.github.com/repos/molmedoz/gopher/releases | jq '.[].tag_name'
   ```

2. Enable Releases feature:
   - Go to repository Settings → Features
   - Enable "Releases"

3. Create a release (even if manually):
   - This will create the Releases tab
   - Or use the workflow: Actions → Create Release

---

## How to Verify Releases

### Check GitHub Releases API

```bash
# List all releases
curl -s https://api.github.com/repos/molmedoz/gopher/releases | jq '.[] | "\(.tag_name) - \(.published_at)"'

# Check specific release
curl -s https://api.github.com/repos/molmedoz/gopher/releases/tags/v1.0.1 | jq '{tag_name, published_at, assets: [.assets[].name]}'
```

### Check Release Assets

```bash
# List assets for latest release
curl -s https://api.github.com/repos/molmedoz/gopher/releases/latest | jq '.assets[].name'
```

### Check Homebrew Formula

```bash
# Check if formula exists
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | head -20

# Check formula version
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | grep "version"
```

---

## Release Workflow Overview

```
1. User triggers "Create Release" workflow
   ↓
2. Validation (tests, linter, builds)
   ↓
3. Tag created and pushed
   ↓
4. GoReleaser runs:
   ├─→ GitHub Releases (main repo)
   ├─→ Homebrew Formula (tap repo)
   ├─→ Chocolatey (chocolatey.org)
   └─→ Scoop (bucket repo)
```

**All releases go to different places**:
- **GitHub Releases** → Main repository (`molmedoz/gopher`)
- **Homebrew** → Tap repository (`molmedoz/homebrew-tap`)
- **Chocolatey** → Community repository
- **Scoop** → Bucket repository

---

## Quick Links

- **GitHub Releases**: https://github.com/molmedoz/gopher/releases
- **Latest Release**: https://github.com/molmedoz/gopher/releases/latest
- **Homebrew Tap**: https://github.com/molmedoz/homebrew-tap
- **Formula File**: https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb
- **Chocolatey**: https://community.chocolatey.org/packages/gopher

---

## Troubleshooting

### "I don't see any releases"

1. **Check if releases exist**:
   ```bash
   curl -s https://api.github.com/repos/molmedoz/gopher/releases
   ```

2. **Check Releases tab is enabled**:
   - Go to: Settings → Features → Enable "Releases"

3. **Use direct URL**:
   - https://github.com/molmedoz/gopher/releases

### "Homebrew release didn't work"

This is separate from GitHub Releases. See:
- [HOMEBREW_SETUP_FIX.md](HOMEBREW_SETUP_FIX.md)
- [HOMEBREW_RELEASE_CHECKLIST.md](HOMEBREW_RELEASE_CHECKLIST.md)

### "I see releases but no assets"

This means GoReleaser didn't run successfully. Check:
- GitHub Actions logs
- GoReleaser configuration
- Build errors

---

**Last Updated**: After clarifying release locations  
**Status**: Current releases are at https://github.com/molmedoz/gopher/releases

