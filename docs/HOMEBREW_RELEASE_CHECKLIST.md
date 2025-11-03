# 🍺 Homebrew Release Checklist

Use this checklist to verify everything is set up correctly for Homebrew releases.

## ✅ Pre-Release Checklist

### 1. Repository Configuration

- [ ] **Tap repository exists**: `molmedoz/homebrew-tap`
  ```bash
  curl -s https://api.github.com/repos/molmedoz/homebrew-tap | jq '.name'
  # Expected: "homebrew-tap"
  ```

- [ ] **Formula directory exists**: `Formula/` directory in tap repository
  ```bash
  curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula
  # Should return directory listing, not 404
  ```

### 2. GitHub Secrets

- [ ] **HOMEBREW_TAP_GITHUB_TOKEN is set**
  - Check: https://github.com/molmedoz/gopher/settings/secrets/actions
  - Must exist: `HOMEBREW_TAP_GITHUB_TOKEN`
  
- [ ] **Token has correct permissions**
  - Required: `repo` scope (Full control) OR fine-grained token with:
    - `Contents: Read and write` on `homebrew-tap` repository
    - `Metadata: Read` on `homebrew-tap` repository

### 3. GoReleaser Configuration

- [ ] **`.goreleaser.yml` has correct Homebrew config**
  ```yaml
  brews:
    - name: gopher
      repository:
        owner: molmedoz
        name: homebrew-tap
        token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
      directory: Formula
      skip_upload: false  # Must be false
  ```

- [ ] **GoReleaser config is valid**
  ```bash
  goreleaser check
  # Should show warnings only (deprecated fields are OK)
  ```

### 4. Workflow Configuration

- [ ] **`.github/workflows/create-release.yml` has `release-homebrew` job**
  - Job name: `release-homebrew`
  - Runs on: `ubuntu-latest`
  - Needs: `validate`, `create-tag`, `validate-config`, `release-github`
  - Uses: `goreleaser/goreleaser-action@v6`
  - Args: `release --skip=validate,archive,chocolatey,nfpm,scoop`
  - Env: `HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}`

### 5. GitHub Release

- [ ] **GitHub release exists** (required before Homebrew can publish)
  - Homebrew job depends on `release-github` job
  - GitHub release must have assets (zip/tar.gz archives)
  - Check: https://github.com/molmedoz/gopher/releases

## 🧪 Testing Before Release

### Test GoReleaser Config Locally

```bash
# Check configuration
goreleaser check

# Test Homebrew formula generation (dry-run)
goreleaser release --snapshot --clean --skip=validate,archive,chocolatey,nfpm,scoop
```

### Verify Token Access

```bash
# Test token access to tap repository
export HOMEBREW_TAP_GITHUB_TOKEN="your-token-here"
curl -s -H "Authorization: token $HOMEBREW_TAP_GITHUB_TOKEN" \
  https://api.github.com/repos/molmedoz/homebrew-tap | jq '.name'
# Expected: "homebrew-tap"
```

## 🚀 Release Process

### 1. Create Release

- Trigger workflow: https://github.com/molmedoz/gopher/actions/workflows/create-release.yml
- Or use: `gh workflow run create-release.yml -f tag=v1.0.1`

### 2. Monitor Workflow

- Watch `release-homebrew` job in GitHub Actions
- Job should run after `release-github` succeeds
- Check for any error messages

### 3. Verify Formula Was Created

```bash
# Wait 2-3 minutes after workflow completes, then check:
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | head -20

# Expected: Ruby formula content (not 404)
```

### 4. Test Installation

```bash
# Add tap
brew tap molmedoz/tap

# Check formula info
brew info molmedoz/tap/gopher

# Install (test)
brew install molmedoz/tap/gopher

# Verify
gopher version
```

## 🔍 Troubleshooting

### Issue: Formula returns 404

**Possible causes:**
1. `release-homebrew` job failed
2. `HOMEBREW_TAP_GITHUB_TOKEN` is missing or invalid
3. Token doesn't have write access to `homebrew-tap`
4. `Formula/` directory doesn't exist in tap repository

**Solution:**
1. Check GitHub Actions logs for `release-homebrew` job
2. Verify token exists and has correct permissions
3. Create `Formula/` directory if missing:
   ```bash
   git clone https://github.com/molmedoz/homebrew-tap.git /tmp/homebrew-tap
   cd /tmp/homebrew-tap
   mkdir -p Formula
   touch Formula/.gitkeep
   git add Formula/
   git commit -m "Add Formula directory"
   git push origin main
   ```

### Issue: Workflow fails with "403 Forbidden"

**Cause:** Token doesn't have write access to `homebrew-tap` repository

**Solution:**
1. Regenerate token with `repo` scope (full access)
2. Or create fine-grained token with:
   - Repository: `molmedoz/homebrew-tap`
   - Permissions: `Contents: Read and write`, `Metadata: Read`
3. Update secret: `gh secret set HOMEBREW_TAP_GITHUB_TOKEN`

### Issue: Job skipped or not running

**Cause:** Dependency job (`release-github`) failed

**Solution:**
1. Check `release-github` job status
2. Fix any issues in GitHub release creation
3. Re-run the workflow

## 📚 Related Documentation

- [TROUBLESHOOTING_HOMEBREW.md](TROUBLESHOOTING_HOMEBREW.md) - Detailed troubleshooting guide
- [verify-homebrew-release.md](verify-homebrew-release.md) - Step-by-step verification
- [FIX_HOMEBREW_NOW.md](FIX_HOMEBREW_NOW.md) - Quick fix actions
- [PACKAGE_MANAGER_SETUP.md](PACKAGE_MANAGER_SETUP.md) - Complete setup guide

## ✅ Quick Verification Commands

```bash
# 1. Check tap repository exists
curl -s https://api.github.com/repos/molmedoz/homebrew-tap | jq '.name'

# 2. Check Formula directory
curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula

# 3. Check if formula file exists
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | head -5

# 4. Test installation (if formula exists)
brew tap molmedoz/tap
brew info molmedoz/tap/gopher
```

---

**Last Updated**: After security improvements and root scoping implementation
**Status**: Ready for v1.0.1 release

