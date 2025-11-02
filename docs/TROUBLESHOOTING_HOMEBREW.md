# 🔧 Troubleshooting Homebrew Formula Creation

## Current Issue: 404 Not Found

The formula file doesn't exist because the `release-homebrew` job hasn't successfully created it.

## Step-by-Step Troubleshooting

### Step 1: Check GitHub Actions Workflow

1. **Go to Actions**: https://github.com/molmedoz/gopher/actions
2. **Find the latest release workflow** (look for `v1.0.0` or your latest release)
3. **Click on the workflow run**
4. **Find the `release-homebrew` job** in the job list
5. **Click on it** to see the logs

**What to look for:**
- ✅ Job status: Success or Failure
- ❌ Error messages about authentication
- ❌ Error messages about permissions
- ❌ Error messages about missing token

### Step 2: Check if Secret is Set

The workflow requires `HOMEBREW_TAP_GITHUB_TOKEN` to push to `molmedoz/homebrew-tap`.

**Verify secret exists:**
1. Go to: https://github.com/molmedoz/gopher/settings/secrets/actions
2. Look for `HOMEBREW_TAP_GITHUB_TOKEN`
3. If it's **missing**, you need to create it (see Step 4)

### Step 3: Check Token Permissions

The token needs write access to `molmedoz/homebrew-tap` repository.

**Required permissions:**
- ✅ `repo` scope (full repository access)
- OR fine-grained token with:
  - `Contents: Read and write`
  - `Metadata: Read`

### Step 4: Create/Update the Secret (If Missing)

#### Option A: Using GitHub CLI (Recommended)

```bash
# Create a Personal Access Token (Classic) or Fine-grained token
# See: https://github.com/settings/tokens

# Set the secret
gh secret set HOMEBREW_TAP_GITHUB_TOKEN
# Paste your token when prompted
```

#### Option B: Using GitHub Web Interface

1. **Create a Personal Access Token:**
   - Go to: https://github.com/settings/tokens
   - Click "Generate new token" → "Generate new token (classic)"
   - Name: `homebrew-tap-token`
   - Select scopes:
     - ✅ `repo` (Full control of private repositories)
   - Click "Generate token"
   - **Copy the token** (you won't see it again!)

2. **Add Secret to Repository:**
   - Go to: https://github.com/molmedoz/gopher/settings/secrets/actions
   - Click "New repository secret"
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: Paste your token
   - Click "Add secret"

### Step 5: Verify Tap Repository Structure

The tap repository needs a `Formula` directory.

**Check current structure:**
```bash
curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents
```

**If Formula directory is missing**, create it:

```bash
# Clone the tap repository
git clone https://github.com/molmedoz/homebrew-tap.git /tmp/homebrew-tap
cd /tmp/homebrew-tap

# Create Formula directory
mkdir -p Formula
touch Formula/.gitkeep

# Commit and push
git add Formula/
git commit -m "Add Formula directory"
git push origin main
```

### Step 6: Test Locally (Optional)

Before creating a new release, you can test the formula creation locally:

```bash
# Make sure you have GoReleaser installed
go install github.com/goreleaser/goreleaser/v2/cmd/goreleaser@latest

# Set your token
export HOMEBREW_TAP_GITHUB_TOKEN="your-token-here"
export GITHUB_TOKEN="your-github-token-here"

# Test with a snapshot release
goreleaser release --snapshot --clean --skip=validate,archives,source,chocolateys,nfpms,scoops

# Check if formula was created in the tap repo
curl -s https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb
```

### Step 7: Re-run the Workflow (If Needed)

If the secret is now set correctly, you can:

**Option A: Create a new release** (recommended)
- This will trigger the workflow fresh
- Use the "Create Release" workflow in Actions

**Option B: Re-run failed job** (if workflow already ran)
1. Go to the workflow run
2. Click "Re-run failed jobs" or "Re-run all jobs"
3. The `release-homebrew` job should now succeed

**Option C: Manually trigger Homebrew release**

Use the `release-channel.yml` workflow:
1. Go to: https://github.com/molmedoz/gopher/actions/workflows/release-channel.yml
2. Click "Run workflow"
3. Select tag: `v1.0.0` (or your latest release tag)
4. Channels: `homebrew`
5. Click "Run workflow"

## Quick Diagnostic Script

Run this to check your setup:

```bash
#!/bin/bash
echo "🔍 Diagnosing Homebrew setup..."
echo ""

# Check tap repo
echo "1️⃣ Tap repository:"
if curl -s https://api.github.com/repos/molmedoz/homebrew-tap | grep -q '"name"'; then
    echo "   ✅ Repository exists"
else
    echo "   ❌ Repository not found"
fi

# Check formula
echo ""
echo "2️⃣ Formula file:"
if curl -s -f https://raw.githubusercontent.com/molmedoz/homebrew-tap/main/Formula/gopher.rb > /dev/null 2>&1; then
    echo "   ✅ Formula exists"
else
    echo "   ❌ Formula not found"
fi

# Check Formula directory
echo ""
echo "3️⃣ Formula directory:"
DIRS=$(curl -s https://api.github.com/repos/molmedoz/homebrew-tap/contents | grep -o '"name": "[^"]*"' | grep -o '[^"]*$')
if echo "$DIRS" | grep -q "^Formula$"; then
    echo "   ✅ Formula directory exists"
else
    echo "   ❌ Formula directory missing"
fi

# Check latest release
echo ""
echo "4️⃣ Latest release:"
RELEASE=$(curl -s https://api.github.com/repos/molmedoz/gopher/releases/latest | grep -o '"tag_name": "[^"]*"' | head -1)
echo "   $RELEASE"

# Check workflow runs
echo ""
echo "5️⃣ Latest workflow runs:"
echo "   Go to: https://github.com/molmedoz/gopher/actions"
echo "   Look for 'release-homebrew' job status"
```

## Common Errors and Solutions

### Error: "403 Resource not accessible by integration"

**Cause:** Token doesn't have access to `molmedoz/homebrew-tap`

**Solution:**
1. Create a Personal Access Token (not the default GITHUB_TOKEN)
2. Token must have `repo` scope
3. Add token as `HOMEBREW_TAP_GITHUB_TOKEN` secret

### Error: "Repository not found"

**Cause:** The tap repository doesn't exist or token can't access it

**Solution:**
1. Verify repository exists: https://github.com/molmedoz/homebrew-tap
2. Ensure token has access to the repository
3. If repository is private, token needs `repo` scope

### Error: "Formula directory not found"

**Cause:** The `Formula` directory doesn't exist in the tap repo

**Solution:**
```bash
# Clone and create Formula directory
git clone https://github.com/molmedoz/homebrew-tap.git
cd homebrew-tap
mkdir -p Formula
git add Formula
git commit -m "Add Formula directory"
git push
```

### Error: Job succeeds but formula doesn't appear

**Cause:** Formula might be on a different branch or commit failed

**Solution:**
1. Check tap repository commits: https://github.com/molmedoz/homebrew-tap/commits/main
2. Look for recent commits from `goreleaserbot`
3. Check if commit was pushed to `main` branch

## Next Steps After Fixing

Once you've fixed the issues:

1. **Verify secret is set**: https://github.com/molmedoz/gopher/settings/secrets/actions
2. **Check Formula directory exists**: See Step 5
3. **Create a new release** or re-run the workflow
4. **Wait 5-10 minutes** for the formula to appear
5. **Test installation**:
   ```bash
   brew tap molmedoz/tap
   brew install gopher
   gopher version
   ```

## Still Having Issues?

1. Check the full workflow logs in GitHub Actions
2. Verify token permissions on: https://github.com/settings/tokens
3. Ensure tap repository exists and is accessible
4. Check GoReleaser documentation: https://goreleaser.com/customization/homebrew/

