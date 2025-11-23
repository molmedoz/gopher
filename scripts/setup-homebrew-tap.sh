#!/bin/bash
# Setup script for Homebrew tap repository
# This creates the Formula directory that GoReleaser needs

set -e

echo "🍺 Setting up Homebrew tap repository..."
echo ""

# Configuration
TAP_REPO="molmedoz/homebrew-tap"
TEMP_DIR="/tmp/homebrew-tap-setup"

# Check if repository exists
echo "1. Checking if tap repository exists..."
if ! curl -s "https://api.github.com/repos/${TAP_REPO}" | grep -q '"name"'; then
    echo "❌ Repository ${TAP_REPO} does not exist!"
    echo "   Please create it first at: https://github.com/new"
    echo "   Repository name: homebrew-tap"
    echo "   Visibility: Public"
    exit 1
fi
echo "✅ Repository exists: ${TAP_REPO}"

# Check if Formula directory exists
echo ""
echo "2. Checking if Formula directory exists..."
if curl -s "https://api.github.com/repos/${TAP_REPO}/contents/Formula" > /dev/null 2>&1; then
    echo "✅ Formula directory already exists!"
    echo "   You can skip the rest of this script."
        exit 0
    fi
echo "⚠️  Formula directory is missing - creating it now..."

# Clone repository
echo ""
echo "3. Cloning repository..."
rm -rf "${TEMP_DIR}"
git clone "https://github.com/${TAP_REPO}.git" "${TEMP_DIR}"
cd "${TEMP_DIR}"

# Create Formula directory
echo ""
echo "4. Creating Formula directory..."
mkdir -p Formula
touch Formula/.gitkeep

# Check current branch
BRANCH=$(git branch --show-current || echo "main")
echo "   Current branch: ${BRANCH}"

# Commit and push
echo ""
echo "5. Committing changes..."
git add Formula/
git commit -m "Add Formula directory for GoReleaser" || {
    echo "⚠️  No changes to commit (directory might already exist locally)"
    echo "   This is OK - the directory might need to be pushed"
}

# Check if we need to push
if git diff --quiet origin/${BRANCH}..HEAD 2>/dev/null; then
    echo "   No changes to push"
else
    echo "   Pushing to ${BRANCH} branch..."
    git push origin "${BRANCH}"
fi

# Verify
echo ""
echo "6. Verifying Formula directory..."
sleep 2  # Wait for GitHub to update
if curl -s "https://api.github.com/repos/${TAP_REPO}/contents/Formula" | grep -q '"name"'; then
    echo "✅ Formula directory created successfully!"
    echo ""
    echo "📋 Next steps:"
    echo "   1. Verify token: gh secret list | grep HOMEBREW_TAP_GITHUB_TOKEN"
    echo "   2. Test release: Use the validate-release workflow"
    echo "   3. Create release: GitHub Actions → Create Release workflow"
else
    echo "⚠️  Formula directory not found yet (might take a few seconds)"
    echo "   Check manually: https://github.com/${TAP_REPO}/tree/${BRANCH}/Formula"
fi

# Cleanup
cd /
rm -rf "${TEMP_DIR}"

echo ""
echo "✅ Setup complete!"
