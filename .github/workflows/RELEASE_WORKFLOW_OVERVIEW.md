# Release Workflow Overview

## Workflows

### 1. `create-release.yml` - Production Release
**Purpose**: Create and publish a new release with automated distribution to all package managers.

**Trigger**: Manual (`workflow_dispatch`)

**Architecture**:
```
validate (code, tests, linting) 
    ↓
create-tag (git tag)
    ↓
validate-config (GoReleaser config + dry-run build)
    ↓
    ├─→ release-github (binaries, archives, checksums)
    │       ↓
    ├─→ release-homebrew (formula update) ←┐
    ├─→ release-chocolatey (Windows pkg)   ├─ Run in parallel
    └─→ release-scoop (Windows bucket)  ───┘
    ↓
release-summary (aggregate status)
```

**Jobs**:
1. **validate** - Run tests, linting, and build validation
2. **create-tag** - Create and push git tag
3. **validate-config** - Validate GoReleaser config and test build
4. **release-github** - Publish binaries and archives to GitHub Releases
5. **release-homebrew** - Update Homebrew formula (parallel, non-critical)
6. **release-chocolatey** - Publish Chocolatey package (parallel, non-critical)
7. **release-scoop** - Update Scoop manifest (parallel, non-critical)
8. **release-summary** - Aggregate status and provide installation instructions

**Benefits**:
- ✅ Independent jobs - one failure doesn't stop others
- ✅ Parallel execution - Homebrew, Chocolatey, Scoop run simultaneously
- ✅ Non-critical failures - Package manager failures don't fail the release
- ✅ GitHub Release is critical - workflow fails only if GitHub release fails
- ✅ Comprehensive validation before any publishing

### 2. `validate-release.yml` - Dry Run Validation
**Purpose**: Validate release configuration without creating a tag or publishing anything.

**Trigger**: Manual (`workflow_dispatch`)

**Use Cases**:
- Test GoReleaser configuration changes
- Verify builds work for all platforms
- Check token permissions
- Validate before creating actual release

**Jobs**:
1. **validate-code** - Run tests, linting, builds
2. **validate-goreleaser-config** - Check GoReleaser config syntax
3. **validate-github-release** - Dry-run GitHub release build
4. **validate-homebrew** - Check Homebrew token and repo access
5. **validate-chocolatey** - Dry-run Chocolatey package build
6. **validate-linux-packages** - Dry-run Linux package builds
7. **validate-scoop** - Check if Scoop bucket exists
8. **summary** - Aggregate results and show readiness

**Output**:
- Shows what would be built/published
- Lists all artifacts that would be created
- Checks all tokens and permissions
- Reports if ready to release

## Job Dependencies

### create-release.yml
```
validate ─→ create-tag ─→ validate-config ─→ release-github ─┬→ release-homebrew
                                                               ├→ release-chocolatey
                                                               └→ release-scoop
                                                                      ↓
                                                              release-summary
```

- **Sequential**: validate → create-tag → validate-config → release-github
- **Parallel**: release-homebrew, release-chocolatey, release-scoop (all depend on release-github)
- **Always runs**: release-summary (even if some jobs fail)

### validate-release.yml
```
validate-code ─→ validate-goreleaser-config ─┬→ validate-github-release
                                              ├→ validate-homebrew
                                              ├→ validate-chocolatey
                                              ├→ validate-linux-packages
                                              └→ validate-scoop
                                                      ↓
                                                  summary
```

- **Sequential**: validate-code → validate-goreleaser-config
- **Parallel**: All validation jobs run in parallel after config validation
- **Always runs**: summary (even if some validations fail)

## Distribution Channels

| Channel | Job | Critical? | Runs in Parallel? |
|---------|-----|-----------|-------------------|
| GitHub Releases | `release-github` | ✅ Yes | No (runs first) |
| Homebrew | `release-homebrew` | ❌ No | Yes |
| Chocolatey | `release-chocolatey` | ❌ No | Yes |
| Scoop | `release-scoop` | ❌ No | Yes |
| Linux packages | Included in `release-github` | ✅ Yes | No |

### Critical vs Non-Critical

**Critical (Must succeed)**:
- GitHub Release creation
- Binary builds
- Archive generation
- Checksums
- Linux packages

**Non-Critical (Can fail)**:
- Homebrew formula update (manual fallback)
- Chocolatey publishing (can republish manually)
- Scoop manifest update (can update manually)

If a non-critical job fails, the release still succeeds and provides manual instructions.

## Usage

### Create a Release

```bash
# Standard release
gh workflow run create-release.yml -f version=1.0.0

# Draft release (can edit before publishing)
gh workflow run create-release.yml -f version=1.0.0 -f draft=true

# Pre-release (beta, rc, etc)
gh workflow run create-release.yml -f version=1.0.0-beta.1 -f prerelease=true
```

### Validate Before Releasing

```bash
# Dry-run validation
gh workflow run validate-release.yml -f version=1.0.0

# Watch the validation
gh run watch

# Check if ready to release
gh run view --log | grep "ready to release"
```

### Monitor Release

```bash
# Watch the release in real-time
gh run watch

# View release summary
gh release view v1.0.0

# Check release assets
gh release view v1.0.0 --json assets --jq '.assets[].name'
```

## Secrets Required

| Secret | Required For | Used By | Type |
|--------|--------------|---------|------|
| `GITHUB_TOKEN` | GitHub Releases | All jobs | Auto-provided |
| `HOMEBREW_TAP_GITHUB_TOKEN` | Homebrew + Scoop | `release-homebrew`, `release-scoop` | PAT |
| `CHOCOLATEY_API_KEY` | Chocolatey | `release-chocolatey` | API Key |
| `SNAPCRAFT_STORE_CREDENTIALS` | Snapcraft | Future use | Credentials |

### Setting Up Secrets

```bash
# Homebrew/Scoop token (Fine-grained PAT)
gh secret set HOMEBREW_TAP_GITHUB_TOKEN --repo molmedoz/gopher

# Chocolatey API key
gh secret set CHOCOLATEY_API_KEY --repo molmedoz/gopher

# Verify secrets
gh secret list --repo molmedoz/gopher
```

## Troubleshooting

### Release Fails at GitHub Job

**Symptom**: `release-github` job fails

**Impact**: Critical - entire release fails, tag is deleted

**Solution**:
1. Check GoReleaser logs for build errors
2. Verify `.goreleaser.yml` syntax: `goreleaser check`
3. Test locally: `goreleaser build --snapshot --clean`
4. Fix issues and retry

### Homebrew Publishing Fails

**Symptom**: `release-homebrew` job fails with 403 error

**Impact**: Non-critical - main release succeeds

**Solution**:
1. Check `HOMEBREW_TAP_GITHUB_TOKEN` has write access to `homebrew-tap`
2. Regenerate token if needed with `Contents: Read and write` permission
3. Or update formula manually:
   ```bash
   cd /tmp
   git clone https://github.com/molmedoz/homebrew-tap.git
   cd homebrew-tap
   # Update Formula/gopher.rb
   git commit -am "Update gopher to v1.0.0"
   git push
   ```

### Chocolatey Publishing Fails

**Symptom**: `release-chocolatey` job fails

**Impact**: Non-critical - main release succeeds

**Solution**:
1. Check `CHOCOLATEY_API_KEY` is valid
2. Verify package builds: Run `validate-release.yml` workflow
3. Or publish manually:
   ```bash
   choco push dist/gopher.1.0.0.nupkg --key YOUR_API_KEY
   ```

### All Jobs Run But None Publish

**Symptom**: Jobs succeed but nothing is published

**Impact**: Major - release doesn't distribute

**Solution**:
1. Check if using `--skip` flags incorrectly
2. Verify GoReleaser isn't in `skip_upload: true` mode
3. Review workflow logs for "skipping" messages

### Validation Fails Before Release

**Symptom**: `validate-config` or `validate-release.yml` fails

**Impact**: Release prevented (good!)

**Solution**:
1. Review which validation failed
2. Fix the issue (build error, config error, etc)
3. Re-run validation to confirm fix
4. Then proceed with release

## Best Practices

### Before Releasing

1. ✅ Run `validate-release.yml` first
2. ✅ Update `CHANGELOG.md`
3. ✅ Verify version follows semver
4. ✅ Test locally: `goreleaser build --snapshot --clean`
5. ✅ Check all secrets are set: `gh secret list`

### During Release

1. ✅ Monitor workflow: `gh run watch`
2. ✅ Watch for failures in parallel jobs
3. ✅ Check GitHub Release creation
4. ✅ Verify assets are uploaded

### After Release

1. ✅ Test installation from each package manager:
   ```bash
   # Homebrew
   brew tap molmedoz/tap && brew install gopher
   
   # Chocolatey
   choco install gopher
   
   # Direct download
   curl -LO https://github.com/molmedoz/gopher/releases/download/v1.0.0/...
   ```
2. ✅ Verify release notes look correct
3. ✅ Check package manager repositories updated:
   - https://github.com/molmedoz/homebrew-tap
   - https://community.chocolatey.org/packages/gopher
   - https://github.com/molmedoz/scoop-bucket (if exists)

### For Failed Package Manager Jobs

If Homebrew, Chocolatey, or Scoop fail:
1. ✅ Don't panic - main release succeeded
2. ✅ Check job logs for specific error
3. ✅ Follow manual instructions in summary
4. ✅ Or re-run just that job (if supported)

## Workflow Evolution

### v1 (Old - Single Job)
- One job doing everything
- Windows runner for Chocolatey
- Sequential execution
- One failure = entire release fails

### v2 (Current - Parallel Jobs)
- Independent jobs per distribution channel
- Linux runner for everything (GoReleaser supports Chocolatey on Linux)
- Parallel execution after GitHub release
- Non-critical failures don't stop release
- Separate validation workflow
- Better observability and debugging

## Performance

### Typical Execution Times

- **validate**: ~3-5 minutes
- **create-tag**: ~10 seconds
- **validate-config**: ~2-3 minutes
- **release-github**: ~5-10 minutes (builds all platforms)
- **release-homebrew**: ~1-2 minutes (parallel)
- **release-chocolatey**: ~2-3 minutes (parallel)
- **release-scoop**: ~1-2 minutes (parallel)
- **release-summary**: ~5 seconds

**Total**: ~8-15 minutes (depending on parallel job overlap)

### Optimization

- ✅ Parallel jobs save ~5-7 minutes vs sequential
- ✅ Validation upfront prevents wasted release attempts
- ✅ Continue-on-error prevents cascading failures

## Future Enhancements

- [ ] Add AUR (Arch User Repository) support
- [ ] Add Winget (Windows Package Manager) support
- [ ] Add Docker image publishing
- [ ] Add Snap support (when name approved)
- [ ] Add rollback workflow
- [ ] Add release announcement workflow
- [ ] Add metrics/monitoring

