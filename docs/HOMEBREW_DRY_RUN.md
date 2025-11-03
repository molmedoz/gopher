# 🧪 Homebrew Release Dry Run Guide

This guide shows you how to test the Homebrew release setup without actually publishing.

## Option 1: GitHub Actions Validate Workflow (Recommended)

The easiest way is to use the built-in validation workflow:

### Steps:

1. **Go to GitHub Actions**:
   - https://github.com/molmedoz/gopher/actions/workflows/validate-release.yml

2. **Click "Run workflow"** (right side)

3. **Click "Run workflow"** button
   - No inputs needed - it validates everything

4. **Wait for completion** (2-3 minutes)

5. **Check the results**:
   - Look for `validate-homebrew` job
   - It will check:
     - ✅ Token is configured
     - ✅ Can access homebrew-tap repository
     - ✅ GoReleaser config is valid

### What it validates:

- ✅ `HOMEBREW_TAP_GITHUB_TOKEN` secret exists
- ✅ Token has access to `molmedoz/homebrew-tap` repository
- ✅ GoReleaser configuration is valid
- ✅ Repository structure is correct

**Note**: This doesn't actually generate the formula, but validates all prerequisites.

---

## Option 2: Local Dry Run with GoReleaser

Test locally by generating the formula without publishing:

### Prerequisites:

```bash
# Install GoReleaser if not already installed
go install github.com/goreleaser/goreleaser/v2/cmd/goreleaser@latest

# Or with Homebrew
brew install goreleaser

# Verify installation
goreleaser --version
```

### Step 1: Set Environment Variables

```bash
# Set your tokens (you'll need these)
export GITHUB_TOKEN="your-github-token-here"
export HOMEBREW_TAP_GITHUB_TOKEN="your-homebrew-tap-token-here"

# Note: You can use GITHUB_TOKEN for both, but HOMEBREW_TAP_GITHUB_TOKEN
# specifically needs write access to molmedoz/homebrew-tap
```

### Step 2: Test GoReleaser Config

```bash
# Check configuration is valid
goreleaser check

# Expected: Should show deprecation warnings (OK) but no errors
```

### Step 3: Generate Formula (Dry Run)

```bash
# Generate Homebrew formula without publishing
# This creates the formula file locally but doesn't push to GitHub
goreleaser release --snapshot --clean --skip=validate,archive,chocolatey,nfpm,scoop

# This will:
# - Build binaries (for testing)
# - Generate the Homebrew formula
# - Show what would be published
# - NOT actually publish to homebrew-tap
```

### Step 4: Check Generated Formula

```bash
# The formula would be generated in dist/ directory
# But with --snapshot, it won't actually create the tap commit

# Better: Check the formula template that would be used
goreleaser release --snapshot --clean --skip=validate,archive,chocolatey,nfpm,scoop --skip-publish

# Or inspect the GoReleaser template directly
cat .goreleaser.yml | grep -A 50 "brews:"
```

### Step 5: Test Formula Syntax (If Formula Generated)

If GoReleaser generates a formula file, you can test it:

```bash
# If formula was generated in a test directory
brew install --build-from-source /path/to/gopher.rb

# Or just validate the Ruby syntax
ruby -c /path/to/gopher.rb
```

---

## Option 3: Test with Snapshot Release (Most Realistic)

This creates a test release that you can verify:

```bash
# Set tokens
export GITHUB_TOKEN="your-token"
export HOMEBREW_TAP_GITHUB_TOKEN="your-token"

# Create a snapshot release (uses random version)
goreleaser release --snapshot --clean --skip=validate,archive,chocolatey,nfpm,scoop

# This will:
# - Build all binaries
# - Create GitHub release (draft)
# - Generate Homebrew formula
# - Actually publish to homebrew-tap (if you have token)
```

**⚠️ Warning**: This will actually publish to the tap repository if `HOMEBREW_TAP_GITHUB_TOKEN` is set and valid. Use with caution!

---

## Option 4: Validate Workflow Manually

You can also manually check each component:

### Check Token Access

```bash
# Test if token can access the tap repository
export HOMEBREW_TAP_GITHUB_TOKEN="your-token-here"

curl -s -H "Authorization: token $HOMEBREW_TAP_GITHUB_TOKEN" \
  https://api.github.com/repos/molmedoz/homebrew-tap | jq '.name'

# Expected: "homebrew-tap"
# If you get 404 or 403, token doesn't have access
```

### Check Repository Structure

```bash
# Check if Formula directory exists
curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula

# Should return directory listing, not 404
```

### Check GoReleaser Config

```bash
# Validate .goreleaser.yml syntax
goreleaser check

# Should show only deprecation warnings, no errors
```

---

## What Each Option Does

| Option | Publishes? | Tests Formula? | Validates Config? | Best For |
|--------|-----------|----------------|-------------------|----------|
| Validate Workflow | ❌ No | ❌ No | ✅ Yes | Quick validation |
| Local Dry Run | ❌ No | ✅ Yes | ✅ Yes | Testing formula generation |
| Snapshot Release | ⚠️ Maybe* | ✅ Yes | ✅ Yes | Full end-to-end test |
| Manual Checks | ❌ No | ❌ No | ✅ Yes | Troubleshooting |

\* Snapshot release will publish if token is valid and `--skip-publish` is not used.

---

## Recommended Workflow

1. **First time setup**: Use **Option 1** (Validate Workflow) to check prerequisites
2. **Before release**: Use **Option 2** (Local Dry Run) to test formula generation
3. **If issues**: Use **Option 4** (Manual Checks) to troubleshoot specific problems

---

## Troubleshooting Dry Run

### Issue: "goreleaser: command not found"

**Solution**: Install GoReleaser
```bash
go install github.com/goreleaser/goreleaser/v2/cmd/goreleaser@latest
```

### Issue: "token not found" or "401 Unauthorized"

**Solution**: 
1. Check token is set: `echo $HOMEBREW_TAP_GITHUB_TOKEN`
2. Verify token has correct permissions
3. Token might have expired - regenerate it

### Issue: "repository not found" or "403 Forbidden"

**Solution**:
1. Verify `molmedoz/homebrew-tap` repository exists
2. Check token has write access to the repository
3. For fine-grained tokens, ensure `Contents: Read and write` permission

### Issue: "Formula directory not found"

**Solution**: Create the `Formula/` directory in the tap repository:
```bash
git clone https://github.com/molmedoz/homebrew-tap.git /tmp/homebrew-tap
cd /tmp/homebrew-tap
mkdir -p Formula
touch Formula/.gitkeep
git add Formula/
git commit -m "Add Formula directory"
git push origin main
```

---

## Next Steps After Dry Run

Once the dry run passes:

1. ✅ All validations pass
2. ✅ Token has correct permissions
3. ✅ Repository structure is correct
4. ✅ GoReleaser config is valid

**You're ready to create a real release!**

Use the "Create Release" workflow:
- https://github.com/molmedoz/gopher/actions/workflows/create-release.yml
- Or: `gh workflow run create-release.yml -f tag=v1.0.1`

---

## Related Documentation

- [HOMEBREW_RELEASE_CHECKLIST.md](HOMEBREW_RELEASE_CHECKLIST.md) - Complete checklist
- [TROUBLESHOOTING_HOMEBREW.md](TROUBLESHOOTING_HOMEBREW.md) - Detailed troubleshooting
- [verify-homebrew-release.md](verify-homebrew-release.md) - Post-release verification

