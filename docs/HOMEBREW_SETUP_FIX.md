# 🔧 Fix Homebrew Setup - Formula Directory Missing

## Problem

GoReleaser didn't create the formula because the `Formula/` directory doesn't exist in the tap repository.

**Current status:**
- ✅ Tap repository exists: `molmedoz/homebrew-tap`
- ❌ Formula directory missing
- ❌ GoReleaser can't create formula files

## Quick Fix (5 minutes)

### Option 1: Use the Setup Script (Easiest)

```bash
# Run the setup script
./scripts/setup-homebrew-tap.sh
```

This script will:
1. Check if repository exists
2. Clone the tap repository
3. Create the `Formula/` directory
4. Commit and push the changes
5. Verify the setup

### Option 2: Manual Setup

If you prefer to do it manually:

```bash
# 1. Clone the tap repository
git clone https://github.com/molmedoz/homebrew-tap.git /tmp/homebrew-tap
cd /tmp/homebrew-tap

# 2. Create Formula directory
mkdir -p Formula
touch Formula/.gitkeep

# 3. Commit and push
git add Formula/
git commit -m "Add Formula directory for GoReleaser"
git push origin main

# 4. Verify
curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula
# Should return directory listing, not 404
```

## Complete Setup Checklist

After creating the Formula directory, verify everything is set up:

### 1. ✅ Formula Directory Exists

```bash
curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula | jq '.[].name'
# Should show: ".gitkeep" or similar
```

### 2. ✅ Token is Configured

```bash
# Check if secret exists
gh secret list | grep HOMEBREW_TAP_GITHUB_TOKEN

# Or check via web:
# https://github.com/molmedoz/gopher/settings/secrets/actions
```

### 3. ✅ Token Has Access

```bash
# Test token access (if you have the token value)
export HOMEBREW_TAP_GITHUB_TOKEN="your-token"
curl -s -H "Authorization: token $HOMEBREW_TAP_GITHUB_TOKEN" \
  https://api.github.com/repos/molmedoz/homebrew-tap | jq '.name'
# Expected: "homebrew-tap"
```

### 4. ✅ GoReleaser Config is Correct

Check `.goreleaser.yml`:

```yaml
brews:
  - name: gopher
    repository:
      owner: molmedoz
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
    directory: Formula  # This is correct
    skip_upload: false  # Must be false
```

### 5. ✅ Workflow is Configured

Check `.github/workflows/create-release.yml` has:
- Job: `release-homebrew`
- Env: `HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}`

## Test After Setup

### 1. Run Validation Workflow

1. Go to: https://github.com/molmedoz/gopher/actions/workflows/validate-release.yml
2. Click "Run workflow" → "Run workflow"
3. Check `validate-homebrew` job
4. Should show: ✅ Can access molmedoz/homebrew-tap

### 2. Test Formula Generation (Dry Run)

```bash
# If you have GoReleaser installed locally
export HOMEBREW_TAP_GITHUB_TOKEN="your-token"
goreleaser release --snapshot --clean --skip=validate,archive,chocolatey,nfpm,scoop
```

### 3. Create a Test Release

Once everything is verified:
1. Go to: https://github.com/molmedoz/gopher/actions/workflows/create-release.yml
2. Click "Run workflow"
3. Enter tag: `v1.0.1` (or your version)
4. Monitor the `release-homebrew` job

## Expected Result

After a successful release, you should see:

```bash
# Formula file should exist
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | head -20

# Should show Ruby formula content starting with:
# class Gopher < Formula
#   desc "A simple, fast, and dependency-free Go version manager"
```

## Troubleshooting

### Issue: Script fails with "repository not found"

**Solution:**
1. Verify repository exists: https://github.com/molmedoz/homebrew-tap
2. Make sure it's public (or your token has access)
3. Check repository name is exactly `homebrew-tap`

### Issue: Script fails with "permission denied"

**Solution:**
1. Make sure you have write access to the repository
2. Check if you're authenticated: `gh auth status`
3. Use HTTPS with token or SSH with key

### Issue: Directory created but GoReleaser still fails

**Possible causes:**
1. Token doesn't have write access
2. Token is expired or revoked
3. Wrong branch (should be `main` or `master`)

**Solution:**
1. Verify token permissions (see [CREATE_HOMEBREW_TOKEN.md](CREATE_HOMEBREW_TOKEN.md))
2. Regenerate token if needed
3. Check workflow logs for specific error messages

### Issue: Formula directory exists but formula file doesn't appear

**Possible causes:**
1. `release-homebrew` job failed
2. Token doesn't have write access
3. GitHub release doesn't exist (Homebrew needs it first)

**Solution:**
1. Check GitHub Actions logs for errors
2. Ensure `release-github` job succeeded first
3. Re-run the workflow

## Repository Structure After Setup

Your `molmedoz/homebrew-tap` repository should look like:

```
homebrew-tap/
├── README.md
├── LICENSE
└── Formula/
    └── .gitkeep
```

After first release, you'll also have:

```
homebrew-tap/
├── README.md
├── LICENSE
└── Formula/
    ├── .gitkeep
    └── gopher.rb  # ← Created by GoReleaser
```

## Next Steps

1. ✅ Run setup script: `./scripts/setup-homebrew-tap.sh`
2. ✅ Verify token: `gh secret list | grep HOMEBREW_TAP_GITHUB_TOKEN`
3. ✅ Test validation: Run validate-release workflow
4. ✅ Create release: Use Create Release workflow

## Related Documentation

- [HOMEBREW_RELEASE_CHECKLIST.md](HOMEBREW_RELEASE_CHECKLIST.md) - Complete checklist
- [TROUBLESHOOTING_HOMEBREW.md](TROUBLESHOOTING_HOMEBREW.md) - Detailed troubleshooting
- [FIX_HOMEBREW_NOW.md](FIX_HOMEBREW_NOW.md) - Quick actions
- [HOMEBREW_DRY_RUN.md](HOMEBREW_DRY_RUN.md) - Testing guide

---

**Last Updated**: After fixing Formula directory issue
**Status**: Ready for v1.0.1 release

