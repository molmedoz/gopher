# 🚨 Immediate Actions to Fix Homebrew Release

## Current Status

✅ Tap repository exists: `molmedoz/homebrew-tap`  
❌ Formula directory missing: No `Formula/` directory  
❌ Formula file missing: Returns 404  
⚠️ Latest workflow succeeded, but formula wasn't created

## Immediate Steps (Do These Now)

### Step 1: Create Formula Directory (5 minutes)

The tap repository needs a `Formula` directory. Run this:

```bash
# Clone the tap repository
git clone https://github.com/molmedoz/homebrew-tap.git /tmp/homebrew-tap
cd /tmp/homebrew-tap

# Create Formula directory
mkdir -p Formula
touch Formula/.gitkeep

# Commit and push
git add Formula/
git commit -m "Add Formula directory for gopher"
git push origin main

# Verify
curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula
```

### Step 2: Verify Secret is Set (2 minutes)

1. Go to: https://github.com/molmedoz/gopher/settings/secrets/actions
2. Check if `HOMEBREW_TAP_GITHUB_TOKEN` exists
3. If **missing**, create it (see Step 3)

### Step 3: Create Token Secret (If Missing) (5 minutes)

#### Create a Personal Access Token:

1. Go to: https://github.com/settings/tokens/new
2. **Name**: `homebrew-tap-bot`
3. **Expiration**: 90 days (or No expiration)
4. **Scopes**: Check `repo` (Full control of private repositories)
5. Click "Generate token"
6. **Copy the token** (starts with `ghp_...`)

#### Add to Repository Secrets:

**Using GitHub CLI:**
```bash
gh secret set HOMEBREW_TAP_GITHUB_TOKEN
# Paste the token when prompted
```

**Or using web interface:**
1. Go to: https://github.com/molmedoz/gopher/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Paste your token
5. Click "Add secret"

### Step 4: Check Latest Workflow Run

1. Go to: https://github.com/molmedoz/gopher/actions
2. Find the latest "Create Release" workflow (the one that succeeded)
3. Click on it
4. Check the `release-homebrew` job:
   - If it shows ❌ failure → Check the logs for errors
   - If it shows ⚠️ skipped → It might have failed silently

### Step 5: Re-run or Create New Release

You have two options:

#### Option A: Re-run Failed Job (Quick)

1. Go to the successful workflow run
2. Find the `release-homebrew` job
3. Click "Re-run job"
4. Wait for completion
5. Check if formula appears: `curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb`

#### Option B: Trigger Homebrew Release Only (Recommended)

Use the `release-channel.yml` workflow:

1. Go to: https://github.com/molmedoz/gopher/actions/workflows/release-channel.yml
2. Click "Run workflow"
3. **Tag**: `v1.0.0` (or your latest release tag)
4. **Channels**: `homebrew`
5. **Draft**: `false`
6. Click "Run workflow"
7. Wait for completion

### Step 6: Verify Formula Created (2 minutes)

Wait 5-10 minutes after the workflow completes, then check:

```bash
# Check if formula exists
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | head -20

# If successful, you'll see Ruby code starting with:
# class Gopher < Formula
#   desc "A simple, fast, and dependency-free Go version manager"
```

## Quick Test Script

After completing the steps above, run this to verify:

```bash
#!/bin/bash
echo "🍺 Quick Homebrew Verification"
echo ""

# Check Formula directory
echo "1. Formula directory:"
if curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents/Formula > /dev/null 2>&1; then
    echo "   ✅ Exists"
else
    echo "   ❌ Missing - Run Step 1"
fi

# Check formula file
echo ""
echo "2. Formula file:"
if curl -s -f https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb > /dev/null 2>&1; then
    echo "   ✅ Exists"
    VERSION=$(curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb | grep -m1 'version "' | sed 's/.*version "\(.*\)".*/\1/' || echo "unknown")
    echo "   📦 Version: $VERSION"
else
    echo "   ❌ Missing - Complete Steps 2-5"
fi

echo ""
echo "📋 Next: Run the workflow or check GitHub Actions"
```

## Summary Checklist

- [ ] Formula directory created in tap repo
- [ ] `HOMEBREW_TAP_GITHUB_TOKEN` secret is set
- [ ] Token has `repo` scope permissions
- [ ] Workflow triggered (re-run or new release)
- [ ] Formula file exists (check after 5-10 minutes)
- [ ] Can install: `brew tap molmedoz/tap && brew install gopher`

## Most Likely Issue

Based on the 404 error, the **Formula directory doesn't exist**. This is required for GoReleaser to push the formula. Complete **Step 1** first, then re-run the workflow.

## Need More Help?

See the full troubleshooting guide: `TROUBLESHOOTING_HOMEBREW.md`

