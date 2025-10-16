# Release Process

This document describes how to create and maintain releases for Gopher.

## Release Workflow

### 1. Prepare Release

#### Update Version Number
```bash
# Decide on version number (semantic versioning)
# MAJOR.MINOR.PATCH
# - MAJOR: Breaking changes
# - MINOR: New features (backward compatible)
# - PATCH: Bug fixes

export VERSION="v1.1.0"
```

#### Update CHANGELOG.md
```bash
# Move content from [Unreleased] to new version section
# 1. Change [Unreleased] header to [X.Y.Z] - YYYY-MM-DD
# 2. Add new empty [Unreleased] section at top
# 3. Review and organize changes
```

**Template:**
```markdown
## [Unreleased]

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
```

#### Update docs/RELEASE_NOTES.md
```bash
# 1. Change "Upcoming" version to actual version with date
# 2. Change status from "In Development" to "Released"
# 3. Add new "Upcoming" section for next version
# 4. Review and polish content
```

#### Run Pre-Release Checks
```bash
# 1. Run all tests
make ci

# 2. Build for all platforms
make build-all

# 3. Run E2E tests
bash test/e2e.sh

# 4. Check documentation
make help | grep -E "(clean|purge|ci)"

# 5. Validate release
make validate-release TAG=$VERSION
```

### 2. Create Release

```bash
# 1. Commit all changes
git add -A
git commit -m "chore: prepare release $VERSION"

# 2. Create and push tag
make create-tag TAG=$VERSION

# 3. Push changes
git push origin main
```

### 3. GitHub Release

The GitHub Actions workflow will automatically:
1. Build binaries for all platforms
2. Run all tests
3. Create checksums
4. Generate release notes
5. Upload artifacts

#### Manual Steps (if needed)
1. Go to https://github.com/molmedoz/gopher/releases
2. Find the auto-created release
3. Edit release notes if needed
4. Publish release

### 4. Post-Release

#### Update Package Managers
```bash
# Homebrew (when available)
# Update formula in molmedoz/homebrew-tap

# Chocolatey (when available)
# Update package in molmedoz/chocolatey-packages

# AUR (when available)
# Update PKGBUILD in molmedoz/aur-gopher
```

#### Announce Release
- Tweet/social media
- Update website
- Notify users
- Update documentation

---

## Maintaining Release Notes

### CHANGELOG.md vs RELEASE_NOTES.md

| File | Purpose | Audience | Format |
|------|---------|----------|--------|
| **CHANGELOG.md** | Technical change log | Developers | Structured, detailed |
| **RELEASE_NOTES.md** | User-facing release info | End users | Polished, highlights |

### CHANGELOG.md Maintenance

**During Development:**
```markdown
## [Unreleased]

### Added
- New feature X
- New command Y

### Changed
- Updated behavior of Z

### Fixed
- Bug in feature W
```

**On Release:**
```markdown
## [Unreleased]

*No unreleased changes yet*

## [1.1.0] - 2025-10-20

### Added
- New feature X
- New command Y

### Changed
- Updated behavior of Z

### Fixed
- Bug in feature W
```

### RELEASE_NOTES.md Maintenance

**During Development:**
```markdown
## v1.1.0 - Feature Name (Upcoming)

**Status**: 🚧 **In Development**

### What's New

#### Feature Category
- **Feature Name**: Description
  - Benefit 1
  - Benefit 2
  - Example: `command output`

### Installation

Same as previous version

### Breaking Changes

None
```

**On Release:**
```markdown
## v1.1.0 - Feature Name (October 2025)

**Status**: ✅ **Released**

### What's New

[Same content as above]

### Installation

[Updated if needed]

### Upgrade from v1.0.0

[Migration guide if needed]
```

---

## Release Checklist

### Pre-Release (1 week before)
- [ ] All planned features completed
- [ ] All tests passing (`make ci`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] RELEASE_NOTES.md drafted
- [ ] Breaking changes documented
- [ ] Migration guide written (if needed)

### Release Day
- [ ] Final `make ci` check
- [ ] Update version numbers
- [ ] Finalize CHANGELOG.md
- [ ] Finalize RELEASE_NOTES.md
- [ ] Commit with message: `chore: prepare release vX.Y.Z`
- [ ] Create and push tag: `make create-tag TAG=vX.Y.Z`
- [ ] Monitor GitHub Actions
- [ ] Verify release artifacts
- [ ] Publish GitHub release
- [ ] Update package managers

### Post-Release
- [ ] Announce release
- [ ] Update website/docs site
- [ ] Monitor for issues
- [ ] Create new [Unreleased] sections
- [ ] Plan next release

---

## Version Numbering Guide

### Semantic Versioning (MAJOR.MINOR.PATCH)

**MAJOR (X.0.0)**: Breaking changes
- API changes that break backward compatibility
- Removed features or commands
- Changed default behavior significantly
- Example: v2.0.0

**MINOR (1.X.0)**: New features (backward compatible)
- New commands
- New features
- Significant improvements
- Example: v1.1.0, v1.2.0

**PATCH (1.0.X)**: Bug fixes
- Bug fixes
- Security patches
- Documentation fixes
- Performance improvements
- Example: v1.0.1, v1.0.2

### Examples

```
v1.0.0 → v1.0.1  (bug fixes)
v1.0.0 → v1.1.0  (new commands: clean, purge)
v1.0.0 → v2.0.0  (breaking: changed CLI interface)
```

---

## Release Notes Template

### For CHANGELOG.md

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- Feature description with technical details
- Links to PRs/issues if relevant

### Changed
- What changed and why
- Impact on users

### Deprecated
- What will be removed in future
- Alternative solutions

### Removed
- What was removed and why
- Migration path

### Fixed
- Bug descriptions
- Issue numbers

### Security
- Security updates
- CVE numbers if applicable
```

### For RELEASE_NOTES.md

```markdown
## vX.Y.Z - Release Title (Month YYYY)

**Status**: ✅ **Released** / 🚧 **In Development**

### What's New

#### Feature Category
- **Feature Name**: User-friendly description
  - Key benefit 1
  - Key benefit 2
  - Usage example with output

### Highlights
- Top 3-5 most important changes
- What users will immediately notice

### Installation

[Platform-specific instructions]

### Upgrade from vX.Y.Z

[Migration steps if needed]

### Breaking Changes

[List any breaking changes and how to handle them]

### Known Issues

[Any known issues and workarounds]
```

---

## Automation Tips

### Use Scripts

```bash
# scripts/prepare-release.sh
VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 v1.1.0"
    exit 1
fi

# Update CHANGELOG.md
sed -i '' "s/## \[Unreleased\]/## [$VERSION] - $(date +%Y-%m-%d)/" CHANGELOG.md

# Update RELEASE_NOTES.md
sed -i '' "s/Upcoming/$(date +"%B %Y")/" docs/RELEASE_NOTES.md
sed -i '' "s/In Development/Released/" docs/RELEASE_NOTES.md

# Commit
git add CHANGELOG.md docs/RELEASE_NOTES.md
git commit -m "chore: prepare release $VERSION"

echo "✓ Release prepared for $VERSION"
echo "Next: make create-tag TAG=$VERSION"
```

### Use Makefile

Already integrated! See:
- `make validate-release TAG=vX.Y.Z`
- `make create-tag TAG=vX.Y.Z`

---

## Communication Plan

### Release Announcement Template

**Subject**: Gopher vX.Y.Z Released - [Feature Highlights]

**Body**:
```
We're excited to announce Gopher vX.Y.Z! 

🎉 What's New:
- Feature 1: Brief description
- Feature 2: Brief description
- Feature 3: Brief description

📦 Installation:
[Quick install command]

📖 Full release notes:
https://github.com/molmedoz/gopher/releases/tag/vX.Y.Z

🙏 Special thanks to our contributors!
```

### Channels
- GitHub Releases
- Twitter/X
- Reddit (r/golang)
- Dev.to
- GitHub Discussions
- Website blog (if applicable)

---

## Maintaining This Process

### After Each Release
1. Create new [Unreleased] section in CHANGELOG.md
2. Create new "Upcoming" section in RELEASE_NOTES.md
3. Update version in next release planning docs
4. Archive release artifacts

### Continuous Improvement
- Review what went well
- Update this process document
- Improve automation scripts
- Gather feedback from users

---

## Quick Reference

```bash
# Daily development
git commit -m "feat: add feature"    # Adds to [Unreleased]

# Preparing release
make validate-release TAG=v1.1.0     # Check everything
vim CHANGELOG.md                      # Finalize changes
vim docs/RELEASE_NOTES.md            # Finalize notes
git commit -m "chore: prepare release v1.1.0"

# Creating release  
make create-tag TAG=v1.1.0           # Create and push tag
# GitHub Actions handles the rest!

# Post-release
# Update package managers
# Announce release
# Create new [Unreleased] section
```

---

**Last Updated**: October 2025  
**Version**: 1.0

