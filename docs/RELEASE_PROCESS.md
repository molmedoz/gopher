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
# Run local validation
make prepare-release VERSION=1.0.0

# This checks:
# - All tests pass with race detection
# - Linter passes
# - Code is formatted
# - Imports are clean
# - Everything builds
```

### 2. Create Release via GitHub Actions

**⚠️ IMPORTANT: We use GitHub Actions workflow for releases (NOT manual tags)**

```bash
# 1. Commit all changes
git add -A
git commit -m "chore: prepare release $VERSION"

# 2. Push to your branch
git push origin your-branch

# 3. Create and merge PR to master

# 4. After merge, go to GitHub Actions
# https://github.com/molmedoz/gopher/actions/workflows/create-release.yml

# 5. Click "Run workflow"
# 6. Enter version (e.g., 1.0.0)
# 7. Choose options (draft/prerelease)
# 8. Click "Run workflow"
```

### 3. Automated Release Process

The GitHub Actions workflow will automatically:
1. ✅ Validate version format
2. ✅ Check tag doesn't already exist
3. ✅ Verify version exists in CHANGELOG.md
4. ✅ Run all tests with race detection
5. ✅ Run linter (golangci-lint)
6. ✅ Build for all platforms
7. ✅ Create tag (ONLY if all above pass)
8. ✅ Run GoReleaser
9. ✅ Extract release notes from CHANGELOG
10. ✅ Create GitHub Release
11. ✅ Upload all artifacts

**Key Benefit:** Tag is created AFTER validation, not before!

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
- [ ] Final `make prepare-release VERSION=X.Y.Z` check
- [ ] Finalize CHANGELOG.md
- [ ] Finalize RELEASE_NOTES.md
- [ ] Commit with message: `chore: prepare release vX.Y.Z`
- [ ] Push and merge PR to master
- [ ] Go to GitHub Actions → "Create Release" workflow
- [ ] Run workflow with version number
- [ ] Monitor workflow (validation → test → tag → build → release)
- [ ] Verify release artifacts on GitHub
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

## Automation

### Using Makefile

```bash
# Prepare release locally (runs all checks)
make prepare-release VERSION=1.1.0

# This runs 'make ci' and shows next steps
```

### Using GitHub Actions

The `create-release.yml` workflow is fully automated:
- Validates everything
- Creates tag only if validation passes
- Builds and releases automatically
- No manual tag creation needed!

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
make prepare-release VERSION=1.1.0   # Check everything locally
vim CHANGELOG.md                      # Finalize changes
vim docs/RELEASE_NOTES.md            # Finalize notes
git commit -m "chore: prepare release v1.1.0"
git push && create PR → merge to master

# Creating release (via GitHub Actions)
# 1. Go to: https://github.com/molmedoz/gopher/actions/workflows/create-release.yml
# 2. Click "Run workflow"
# 3. Enter version: 1.1.0
# 4. Click "Run workflow"
# GitHub Actions validates → tests → tags → builds → releases!

# Post-release
# Update package managers
# Announce release
# Create new [Unreleased] section
```

---

**Last Updated**: October 2025  
**Version**: 1.0

