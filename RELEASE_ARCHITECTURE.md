# Release Architecture

## 🎯 What You Asked For

✅ **Independent release jobs** - Each distribution channel runs independently  
✅ **Parallel execution** - Homebrew, Chocolatey, Scoop run simultaneously  
✅ **Validation before release** - Config check + dry-run build before publishing  
✅ **Dry-run workflow** - Test everything without creating tags or releases  
✅ **Amended commit** - All changes in one commit

## 📊 New Architecture Diagram

### Production Release (`create-release.yml`)

```
┌─────────────────────────────────────────────────────────────────┐
│  STAGE 1: VALIDATION & SETUP                                    │
└─────────────────────────────────────────────────────────────────┘
                            │
        ┌──────────────────┼──────────────────┐
        ▼                   ▼                   ▼
   [validate]        [create-tag]      [validate-config]
   • Tests           • Create tag      • Check config
   • Linting         • Push tag        • Dry-run build
   • Builds          
        │                   │                   │
        └──────────────────┼───────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  STAGE 2: CRITICAL RELEASE                                      │
└─────────────────────────────────────────────────────────────────┘
                            │
                    [release-github]
                    • Build all platforms
                    • Create archives
                    • Upload to GitHub
                    • Linux packages
                            │
                            │ ✅ MUST SUCCEED
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  STAGE 3: PARALLEL DISTRIBUTION (Non-Critical)                  │
└─────────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
[release-homebrew]  [release-chocolatey]  [release-scoop]
• Update formula    • Publish package     • Update manifest
• molmedoz/tap      • chocolatey.org      • scoop-bucket
        │                   │                   │
        └───────────────────┼───────────────────┘
                            ▼
                   [release-summary]
                   • Aggregate status
                   • Installation guide
```

### Validation Only (`validate-release.yml`)

```
┌─────────────────────────────────────────────────────────────────┐
│  VALIDATION - NO TAGS, NO RELEASES, NO PUBLISHING              │
└─────────────────────────────────────────────────────────────────┘
                            │
                    [validate-code]
                    • Tests
                    • Linting  
                    • Builds
                            │
                            ▼
              [validate-goreleaser-config]
              • Check YAML syntax
              • Validate config
                            │
        ┌───────────────────┼───────────────────┬───────────────┐
        │                   │                   │               │
        ▼                   ▼                   ▼               ▼
[validate-github]   [validate-homebrew] [validate-chocolatey] ...
• Dry-run build     • Check token       • Test package build
• List artifacts    • Check repo        • Check API key
        │                   │                   │               │
        └───────────────────┼───────────────────┴───────────────┘
                            ▼
                       [summary]
                       • Ready to release?
                       • What would be built
                       • What would be published
```

## 🎮 How to Use

### 1. Validate First (Dry Run - No Changes)

```bash
# Test your configuration without making any changes
gh workflow run validate-release.yml -f version=1.0.0

# Watch the validation
gh run watch

# Check the summary
gh run view --log | grep "ready to release"
```

**This will**:
- ✅ Run all tests and linting
- ✅ Validate GoReleaser configuration
- ✅ Build binaries (but not publish)
- ✅ Check all tokens and permissions
- ✅ Show what would be created
- ❌ NOT create any tags
- ❌ NOT publish anything

### 2. Create Real Release (Production)

```bash
# Create a production release
gh workflow run create-release.yml -f version=1.0.0

# Or create as draft first
gh workflow run create-release.yml -f version=1.0.0 -f draft=true

# Monitor in real-time
gh run watch
```

**This will**:
1. **Validate** (Stage 1): Run tests, create tag, validate config
2. **Release to GitHub** (Stage 2 - Critical): Build and publish binaries
3. **Distribute** (Stage 3 - Parallel): Update Homebrew, Chocolatey, Scoop simultaneously
4. **Summarize**: Show status of all channels

## 🔄 Job Dependencies

### Sequential Flow
```
validate → create-tag → validate-config → release-github
```
These MUST run in order. If any fails, the workflow stops.

### Parallel Flow
```
release-github (completes) 
    ↓
    ├→ release-homebrew    ┐
    ├→ release-chocolatey  │ All run at the same time
    └→ release-scoop       ┘
    ↓
release-summary (always runs)
```

**Time saved**: ~5-7 minutes compared to sequential execution

## ⚙️ Job Configurations

### Critical Jobs (Must Succeed)
- `validate` - Code quality
- `create-tag` - Tag creation
- `validate-config` - Config validation  
- `release-github` - Main release

**If any of these fail**: Entire workflow fails, tag is deleted

### Non-Critical Jobs (Can Fail Gracefully)
- `release-homebrew` - Homebrew formula
- `release-chocolatey` - Chocolatey package
- `release-scoop` - Scoop manifest

**If any of these fail**: 
- Workflow continues
- Other channels still publish
- Manual instructions provided in summary

## 📦 What Gets Published Where

| Distribution | Job | Parallel? | Critical? | Fallback |
|--------------|-----|-----------|-----------|----------|
| **GitHub Release** | `release-github` | No (runs first) | ✅ Yes | None - must work |
| **Binaries** | `release-github` | No | ✅ Yes | None - must work |
| **Archives (.tar.gz, .zip)** | `release-github` | No | ✅ Yes | Manual download |
| **Checksums** | `release-github` | No | ✅ Yes | Manual verification |
| **Linux (.deb, .rpm, .apk)** | `release-github` | No | ✅ Yes | Manual install |
| **Homebrew Formula** | `release-homebrew` | Yes | ❌ No | Manual PR to tap |
| **Chocolatey Package** | `release-chocolatey` | Yes | ❌ No | Manual push |
| **Scoop Manifest** | `release-scoop` | Yes | ❌ No | Manual PR to bucket |

## 🎭 Failure Scenarios

### Scenario 1: GitHub Release Fails
```
Status: ❌ CRITICAL FAILURE
Impact: Entire release fails, tag deleted
Action: Fix issue, retry entire workflow
```

### Scenario 2: Homebrew Fails, Others Succeed
```
Status: ⚠️ PARTIAL SUCCESS
Impact: GitHub ✅, Chocolatey ✅, Scoop ✅, Homebrew ❌
Action: Manual Homebrew formula update (instructions provided)
Result: Release still published, users can install via other methods
```

### Scenario 3: All Package Managers Fail
```
Status: ✅ MINIMAL SUCCESS
Impact: GitHub Release ✅, All package managers ❌
Action: Users can download directly from GitHub
Result: Binaries available, package managers can be fixed post-release
```

### Scenario 4: Validation Fails (Dry Run)
```
Status: 🛑 PREVENTED FAILURE
Impact: No tag created, no release attempted
Action: Fix issue, re-run validation
Result: Saved from a failed release!
```

## 🚀 Performance Comparison

### Old Architecture (Single Job)
```
Total time: ~15-20 minutes
└─ Sequential: Validate → Build → GitHub → Homebrew → Chocolatey → Scoop

Failure impact: One failure = entire release fails
```

### New Architecture (Parallel Jobs)
```
Total time: ~8-12 minutes

Stage 1 (Sequential): ~3-5 min
  └─ Validate → Create Tag → Validate Config

Stage 2 (Critical): ~5-7 min  
  └─ GitHub Release

Stage 3 (Parallel): ~2-3 min
  ├─ Homebrew (2 min) ┐
  ├─ Chocolatey (3 min) │ Running simultaneously  
  └─ Scoop (2 min)     ┘

Failure impact: GitHub failure = release fails
                Package manager failure = release succeeds with warning
```

**Time saved**: ~5-8 minutes per release  
**Reliability**: Much better (isolated failures)

## 📋 Checklist for Release

### Pre-Release
- [ ] Run dry-run validation: `gh workflow run validate-release.yml -f version=X.Y.Z`
- [ ] Check validation passed
- [ ] Update `CHANGELOG.md`
- [ ] Verify all secrets are set: `gh secret list`

### Release
- [ ] Run create-release: `gh workflow run create-release.yml -f version=X.Y.Z`
- [ ] Monitor workflow: `gh run watch`
- [ ] Check GitHub Release created
- [ ] Verify package managers updated (or note failures)

### Post-Release
- [ ] Test GitHub download: `curl -LO https://github.com/molmedoz/gopher/releases/download/vX.Y.Z/...`
- [ ] Test Homebrew: `brew install molmedoz/tap/gopher`
- [ ] Test Chocolatey: `choco install gopher`
- [ ] Test Scoop: `scoop install gopher` (if enabled)
- [ ] Fix any failed package managers manually (if needed)

## 🎓 Learning the New System

### Day 1: Understand the Flow
1. Read `RELEASE_WORKFLOW_OVERVIEW.md` (detailed docs)
2. Review this file (architecture overview)
3. Look at `.github/workflows/create-release.yml` (implementation)

### Day 2: Test It Out
1. Run `validate-release.yml` with a test version
2. Review the validation output
3. Understand what would happen in a real release

### Day 3: First Real Release
1. Follow the pre-release checklist
2. Create a draft release first (`-f draft=true`)
3. Review the draft, then publish
4. Observe parallel job execution

### Day 4: Handle Failures
1. Intentionally break a non-critical job (e.g., Homebrew token)
2. Run release and observe graceful degradation
3. Practice manual recovery from summary instructions

## 🔗 Related Documentation

- `RELEASE_WORKFLOW_OVERVIEW.md` - Complete workflow documentation
- `RELEASE_DISTRIBUTION.md` - Distribution channel setup guide
- `RELEASE_FIXES_SUMMARY.md` - What was fixed and why
- `.github/workflows/create-release.yml` - Production release workflow
- `.github/workflows/validate-release.yml` - Validation/dry-run workflow

## 💡 Key Takeaways

1. **Independent Jobs** = More resilient (one failure doesn't break everything)
2. **Parallel Execution** = Faster releases (5-8 minutes saved)
3. **Validation First** = Prevent failed releases (catch issues early)
4. **Non-Critical Failures** = Graceful degradation (release succeeds, fix later)
5. **Dry-Run Workflow** = Safe testing (no tags, no releases)

**Bottom Line**: The release process is now faster, more reliable, and easier to debug! 🎉

