# Release Workflows - Quick Reference

## Available Workflows

| Workflow | Purpose | Creates Tag? | Creates Release? |
|----------|---------|--------------|------------------|
| `validate-release.yml` | Dry-run validation | ❌ No | ❌ No |
| `create-release.yml` | Full production release | ✅ Yes | ✅ Yes |
| `release-channel.yml` | Deploy to specific channel(s) | ❌ No | Depends |

## Common Use Cases

### 1. Test Before Releasing (Dry Run)

**Scenario**: You want to make sure everything works before creating a real release.

```bash
# Validate without creating tags or releases
gh workflow run validate-release.yml -f version=1.0.0

# Watch the validation
gh run watch

# Check if ready
gh run view --log | grep "ready to release"
```

**What it does**:
- ✅ Runs all tests
- ✅ Validates GoReleaser config
- ✅ Builds binaries (dry-run)
- ✅ Checks tokens and permissions
- ❌ Does NOT create tags
- ❌ Does NOT publish anything

---

### 2. Full Production Release (All Channels)

**Scenario**: You're ready to release a new version to all channels.

```bash
# Create full release
gh workflow run create-release.yml -f version=1.0.0

# Or as draft (can edit before publishing)
gh workflow run create-release.yml -f version=1.0.0 -f draft=true

# Monitor
gh run watch
```

**What it does**:
- ✅ Creates git tag
- ✅ Publishes to GitHub (binaries, archives, Linux packages)
- ✅ Updates Homebrew (parallel)
- ✅ Publishes to Chocolatey (parallel)
- ✅ Updates Scoop (parallel)

---

### 3. Deploy Only to GitHub

**Scenario**: You want to publish binaries to GitHub only, without updating package managers.

```bash
# For an existing tag
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=github

# For a new tag, create tag first
git tag v1.0.0
git push origin v1.0.0
gh workflow run release-channel.yml -f tag=v1.0.0 -f channels=github
```

**What it does**:
- ✅ Builds binaries for all platforms
- ✅ Creates archives (.tar.gz, .zip)
- ✅ Generates checksums
- ✅ Builds Linux packages
- ✅ Publishes to GitHub Releases
- ❌ Skips Homebrew
- ❌ Skips Chocolatey
- ❌ Skips Scoop

**Use when**:
- Testing GitHub release process
- You want to publish binaries but update package managers later
- Fixing a failed GitHub release

---

### 4. Deploy Only to Homebrew

**Scenario**: Homebrew failed during full release, or you want to update the formula manually.

```bash
# Update Homebrew formula for existing release
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=homebrew
```

**Requirements**:
- ⚠️ GitHub release must already exist with assets
- ⚠️ `HOMEBREW_TAP_GITHUB_TOKEN` must be set

**What it does**:
- ✅ Updates Formula/gopher.rb in molmedoz/homebrew-tap
- ✅ Calculates SHA256 from GitHub release assets
- ❌ Skips building binaries (uses existing GitHub release)

**Use when**:
- Homebrew publishing failed in full release
- You fixed the Homebrew tap repo and want to retry
- Testing Homebrew formula updates

---

### 5. Deploy Only to Chocolatey

**Scenario**: Chocolatey failed during full release, or you want to republish.

```bash
# Publish to Chocolatey for existing release
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=chocolatey
```

**Requirements**:
- ⚠️ GitHub release must already exist
- ⚠️ `CHOCOLATEY_API_KEY` must be set

**What it does**:
- ✅ Builds .nupkg package
- ✅ Publishes to chocolatey.org
- ❌ Skips building binaries (uses existing GitHub release)

**Use when**:
- Chocolatey publishing failed in full release
- API key expired and you renewed it
- Testing Chocolatey package

---

### 6. Deploy Only to Scoop

**Scenario**: Update Scoop manifest for existing release.

```bash
# Update Scoop manifest
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=scoop
```

**Requirements**:
- ⚠️ GitHub release must already exist
- ⚠️ `molmedoz/scoop-bucket` repository must exist
- ⚠️ `HOMEBREW_TAP_GITHUB_TOKEN` must have access to scoop-bucket

**What it does**:
- ✅ Updates gopher.json in molmedoz/scoop-bucket
- ✅ Calculates SHA256 from GitHub release assets
- ❌ Skips building binaries (uses existing GitHub release)

**Use when**:
- Scoop publishing failed (bucket didn't exist)
- You created the scoop-bucket repo and want to publish
- Testing Scoop manifest updates

---

### 7. Deploy to Multiple Specific Channels

**Scenario**: You want to update only package managers, not GitHub.

```bash
# Update all package managers (but not GitHub)
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=homebrew,chocolatey,scoop

# Or use the 'all' option to deploy everything
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=all
```

**What it does**:
- When using `homebrew,chocolatey,scoop`:
  - ✅ Updates Homebrew formula
  - ✅ Publishes to Chocolatey
  - ✅ Updates Scoop manifest
  - ❌ Skips GitHub release (uses existing)

- When using `all`:
  - ✅ Publishes to GitHub (if not already)
  - ✅ Updates all package managers

**Use when**:
- All package managers failed but GitHub succeeded
- You want to republish to all channels
- Testing all package manager updates

---

### 8. Fix a Failed Release

**Scenario**: Full release ran, GitHub succeeded, but Homebrew failed.

```bash
# Step 1: Check what failed
gh run view --log | grep "❌"

# Step 2: Fix the issue (e.g., regenerate HOMEBREW_TAP_GITHUB_TOKEN)

# Step 3: Re-run just the failed channel
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=homebrew
```

**Use when**:
- Token expired during release
- Repository wasn't accessible
- Network issues during publish
- Want to retry without rebuilding everything

---

### 9. Republish Everything

**Scenario**: You want to republish an existing release to all channels (e.g., after fixing package manager issues).

```bash
# Republish to all channels
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=all
```

**What it does**:
- ✅ Rebuilds and republishes to GitHub
- ✅ Updates Homebrew formula
- ✅ Publishes to Chocolatey
- ✅ Updates Scoop manifest

**Use when**:
- You deleted a release and want to recreate it
- Multiple channels failed
- Starting fresh after fixing issues

---

## Workflow Decision Tree

```
┌─────────────────────────────────────────────┐
│ What do you want to do?                     │
└─────────────────┬───────────────────────────┘
                  │
                  ├─→ Test configuration?
                  │   └─→ validate-release.yml
                  │
                  ├─→ Create new release?
                  │   └─→ create-release.yml
                  │
                  ├─→ Fix failed channel?
                  │   └─→ release-channel.yml (specific channel)
                  │
                  ├─→ Publish to one channel only?
                  │   └─→ release-channel.yml (specific channel)
                  │
                  └─→ Republish existing release?
                      └─→ release-channel.yml (all)
```

## Channel Options in `release-channel.yml`

| Option | Deploys To |
|--------|------------|
| `github` | GitHub Releases only |
| `homebrew` | Homebrew tap only |
| `chocolatey` | Chocolatey.org only |
| `scoop` | Scoop bucket only |
| `homebrew,chocolatey,scoop` | All package managers (not GitHub) |
| `all` | Everything (GitHub + all package managers) |

## Important Notes

### For Package Managers (Homebrew, Chocolatey, Scoop)

⚠️ **GitHub release must exist first!**

Package managers need:
- GitHub release with assets (binaries, archives)
- SHA256 checksums
- Proper version tag

**Correct order**:
```bash
# 1. First, create GitHub release (if not exists)
gh workflow run release-channel.yml -f tag=v1.0.0 -f channels=github

# 2. Wait for completion
gh run watch

# 3. Then deploy to package managers
gh workflow run release-channel.yml -f tag=v1.0.0 -f channels=homebrew
```

### Tag Format

Tags can be with or without `v` prefix:
- ✅ `v1.0.0` - preferred
- ✅ `1.0.0` - will auto-add `v` prefix

### Draft Releases

Only applies to `github` channel:
```bash
# Create draft (can edit before publishing)
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=github \
  -f draft=true

# Edit draft on GitHub, then publish manually
```

## Troubleshooting

### "Tag does not exist"

**Error**: `Tag v1.0.0 does not exist`

**Solution**:
```bash
# Create the tag first
git tag v1.0.0
git push origin v1.0.0

# Then run release-channel
gh workflow run release-channel.yml -f tag=v1.0.0 -f channels=github
```

### "GitHub release needed first"

**Error**: Package manager job skips with "GitHub release needed first"

**Solution**:
```bash
# Create GitHub release first
gh workflow run release-channel.yml -f tag=v1.0.0 -f channels=github

# Wait for it to complete
gh run watch

# Then run package manager
gh workflow run release-channel.yml -f tag=v1.0.0 -f channels=homebrew
```

### Token Errors

**Error**: `403 Resource not accessible by integration`

**Solution**:
```bash
# Check if token exists
gh secret list --repo molmedoz/gopher

# Set/update the token
gh secret set HOMEBREW_TAP_GITHUB_TOKEN --repo molmedoz/gopher
# Or for Chocolatey
gh secret set CHOCOLATEY_API_KEY --repo molmedoz/gopher
```

## Examples

### Example 1: New Release (Full Workflow)

```bash
# Step 1: Validate
gh workflow run validate-release.yml -f version=1.0.0
gh run watch

# Step 2: If validation passed, create release
gh workflow run create-release.yml -f version=1.0.0
gh run watch

# Step 3: Verify
gh release view v1.0.0
```

### Example 2: Fix Failed Homebrew

```bash
# During release, Homebrew failed
# Check the error
gh run view --log

# Fix the issue (e.g., token), then:
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=homebrew

# Verify
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | grep version
```

### Example 3: Test GitHub Release Only

```bash
# Create tag manually
git tag v1.0.0-test
git push origin v1.0.0-test

# Publish to GitHub as draft
gh workflow run release-channel.yml \
  -f tag=v1.0.0-test \
  -f channels=github \
  -f draft=true

# Review draft on GitHub
# Delete if not good
gh release delete v1.0.0-test --yes
git tag -d v1.0.0-test
git push origin :refs/tags/v1.0.0-test
```

### Example 4: Publish After Creating Scoop Bucket

```bash
# You released but scoop-bucket didn't exist
# Now you created it:
gh repo create molmedoz/scoop-bucket --public

# Publish to Scoop
gh workflow run release-channel.yml \
  -f tag=v1.0.0 \
  -f channels=scoop

# Verify
curl -s https://raw.githubusercontent.com/molmedoz/scoop-bucket/main/bucket/gopher.json
```

## Summary

| Workflow | Use Case | Speed |
|----------|----------|-------|
| `validate-release.yml` | Testing, validation | Fast (~5 min) |
| `create-release.yml` | New releases | Medium (~10 min) |
| `release-channel.yml` (one channel) | Fix/test single channel | Fast (~2-3 min) |
| `release-channel.yml` (all) | Republish everything | Medium (~8 min) |

**Pro tip**: Always run `validate-release.yml` before `create-release.yml` to catch issues early!

