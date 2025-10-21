# Release Distribution Guide

## Overview

Gopher is distributed through multiple package managers and platforms to make installation easy for users across different operating systems.

## Current Distribution Channels

### ✅ Active (Auto-Published on Release)

1. **GitHub Releases** 
   - Binary archives for all platforms (.tar.gz, .zip)
   - Checksums (SHA256)
   - Source code archives
   - 📍 Location: https://github.com/molmedoz/gopher/releases

2. **Homebrew (macOS/Linux)**
   - Formula repository: https://github.com/molmedoz/homebrew-tap
   - Install: `brew tap molmedoz/tap && brew install gopher`
   - Auto-updates via GitHub Actions

3. **Chocolatey (Windows)**
   - Repository: https://community.chocolatey.org/packages/gopher
   - Install: `choco install gopher`
   - Auto-publishes to Chocolatey.org (requires `CHOCOLATEY_API_KEY`)

4. **Linux Packages**
   - **Debian/Ubuntu**: `.deb` packages
   - **Red Hat/Fedora/CentOS**: `.rpm` packages
   - **Alpine Linux**: `.apk` packages
   - **Arch Linux**: `.pkg.tar.zst` packages
   - 📍 Available as GitHub Release assets

5. **Scoop (Windows)**
   - Bucket: https://github.com/molmedoz/scoop-bucket
   - Install: `scoop bucket add molmedoz https://github.com/molmedoz/scoop-bucket && scoop install gopher`
   - Auto-publishes if bucket repository exists (`skip_upload: auto`)

### 🚧 Configured but Disabled

6. **Snapcraft (Linux)**
   - ❌ Disabled: Name "gopher" requires manual approval from Snapcraft team
   - Status: Pending approval
   - Once approved: `sudo snap install gopher --classic`

7. **AUR - Arch User Repository**
   - ❌ Disabled: Requires AUR account setup and SSH key
   - Package name: `gopher-bin`
   - Future: Manual submission to AUR

8. **Winget (Windows Package Manager)**
   - ❌ Disabled: Requires fork of microsoft/winget-pkgs
   - Package: `molmedoz.gopher`
   - Future: Submit to Microsoft's winget-pkgs repository

9. **Docker/OCI Images**
   - ❌ Disabled: Requires `Dockerfile.release`
   - Future location: `ghcr.io/molmedoz/gopher:latest`

## Architecture

### Single Job Release (Simplified)

```
┌─────────────────────────────────────────────────────┐
│  Job: validate (Ubuntu)                            │
│  - Run tests with race detection                  │
│  - Run linter                                      │
│  - Build for all platforms                        │
└─────────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────┐
│  Job: create-tag (Ubuntu)                          │
│  - Create git tag                                  │
│  - Push to repository                              │
└─────────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────┐
│  Job: release (Ubuntu)                             │
│  - Build all binaries (Linux, macOS, Windows)      │
│  - Create archives                                 │
│  - Generate checksums                              │
│  - Publish to GitHub Releases                      │
│  - Publish to Homebrew tap                         │
│  - Publish to Chocolatey                           │
│  - Build Linux packages (.deb, .rpm, .apk, etc)    │
│  - Publish to Scoop (if bucket exists)             │
└─────────────────────────────────────────────────────┘
```

### Why Ubuntu for Everything?

- **GoReleaser on Linux can build Chocolatey packages** (no Windows required)
- **Better token authentication** for cross-repository writes (Homebrew, Scoop)
- **Simpler workflow** with one job instead of multiple
- **Faster execution** (Linux runners start faster than Windows)
- **Cost effective** (Linux runners are cheaper on GitHub Actions)

## Required Secrets

| Secret | Purpose | Required? | How to Get |
|--------|---------|-----------|------------|
| `GITHUB_TOKEN` | GitHub releases, archives | ✅ Auto | Automatically provided by GitHub Actions |
| `HOMEBREW_TAP_GITHUB_TOKEN` | Homebrew + Scoop publishing | ✅ Manual | [Create PAT](https://github.com/settings/tokens) with `Contents: Read and write` on `homebrew-tap` and `scoop-bucket` repos |
| `CHOCOLATEY_API_KEY` | Chocolatey publishing | Optional | [Get from Chocolatey](https://community.chocolatey.org/account) |
| `SNAPCRAFT_STORE_CREDENTIALS` | Snapcraft publishing | Optional | When enabled: `snapcraft export-login` |
| `AUR_SSH_PRIVATE_KEY` | AUR publishing | Optional | When enabled: Generate SSH key for AUR |

## Setting Up Additional Distribution Channels

### Enable Scoop (Windows Package Manager)

1. **Create the bucket repository:**
   ```bash
   gh repo create molmedoz/scoop-bucket --public --description "Scoop bucket for molmedoz tools"
   cd /tmp && git clone https://github.com/molmedoz/scoop-bucket.git
   cd scoop-bucket
   echo "# Scoop Bucket" > README.md
   git add README.md && git commit -m "Initial commit" && git push
   ```

2. **Grant token access:**
   - Your `HOMEBREW_TAP_GITHUB_TOKEN` should already have access
   - If not, add `scoop-bucket` to the token's repository list

3. **Next release will auto-publish** (no code changes needed)

### Enable Winget (Future)

1. Fork https://github.com/microsoft/winget-pkgs
2. Update `.goreleaser.yml`:
   ```yaml
   winget:
     - repository:
         owner: molmedoz
         name: winget-pkgs  # Your fork
   skip_upload: false
   ```
3. After GoReleaser updates your fork, create PR to microsoft/winget-pkgs

### Enable AUR (Arch User Repository)

1. Create AUR account at https://aur.archlinux.org/register
2. Generate SSH key: `ssh-keygen -t ed25519 -C "your-email@example.com"`
3. Add public key to AUR account
4. Add private key to GitHub Secrets: `AUR_SSH_PRIVATE_KEY`
5. Update `.goreleaser.yml`:
   ```yaml
   aurs:
     - name: gopher-bin
       skip_upload: false
   ```

### Enable Docker Images

1. Create `Dockerfile.release`:
   ```dockerfile
   FROM alpine:latest
   COPY gopher /usr/local/bin/gopher
   ENTRYPOINT ["/usr/local/bin/gopher"]
   ```

2. Update `.goreleaser.yml`:
   ```yaml
   dockers:
     - skip_push: false
   ```

3. The image will be published to `ghcr.io/molmedoz/gopher:latest`

## Testing Releases

### Create a Draft Release

```bash
gh workflow run create-release.yml \
  -f version=1.0.1 \
  -f draft=true
```

### Monitor the Release

```bash
# Watch the workflow
gh run watch

# Or view in browser
gh workflow view create-release.yml --web
```

### Verify Distribution

After release:
- ✅ Check GitHub: https://github.com/molmedoz/gopher/releases
- ✅ Check Homebrew: https://github.com/molmedoz/homebrew-tap/blob/main/Formula/gopher.rb
- ✅ Check Chocolatey: https://community.chocolatey.org/packages/gopher
- ✅ Check Scoop: https://github.com/molmedoz/scoop-bucket (if enabled)

## Troubleshooting

### GitHub Archives Not Appearing

**Symptom**: Release created but no binary archives attached

**Solution**: Archives are built by GoReleaser automatically. Check:
1. GoReleaser didn't fail (check workflow logs)
2. `.goreleaser.yml` has `archives:` section configured
3. `release --clean` is being used (not `--skip=archives`)

### Homebrew 403 Error

**Symptom**: `403 Resource not accessible by integration`

**Solution**: Token permissions issue
1. Go to https://github.com/settings/tokens
2. Find `HOMEBREW_TAP_GITHUB_TOKEN`
3. Ensure it has `Contents: Read and write` on `homebrew-tap` repository
4. Regenerate if needed and update secret:
   ```bash
   gh secret set HOMEBREW_TAP_GITHUB_TOKEN --repo molmedoz/gopher
   ```

### Chocolatey Build Fails

**Symptom**: Chocolatey package build fails on Linux

**Solution**: This is normal if:
- `CHOCOLATEY_API_KEY` is not set (package is built but not published)
- GoReleaser successfully builds the `.nupkg` file even on Linux

If the build actually fails, check:
1. Version format is correct (semver)
2. URLs in `.goreleaser.yml` are accessible

## Release Checklist

- [ ] Update `CHANGELOG.md` with version
- [ ] Verify all tests pass locally
- [ ] Run workflow: `gh workflow run create-release.yml -f version=X.Y.Z`
- [ ] Monitor workflow execution
- [ ] Verify GitHub release created
- [ ] Test installation from each package manager
- [ ] Announce release (optional)

## Future Distribution Channels

- [ ] **Flox** (nix-based package manager)
- [ ] **Nix packages**
- [ ] **MacPorts** (alternative to Homebrew)
- [ ] **APT repository** (self-hosted)
- [ ] **Yum/DNF repository** (self-hosted)
- [ ] **Cargo** (if Rust bindings created)
- [ ] **npm** (if Node.js wrapper created)

## Documentation

- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Homebrew Tap Creation](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [Chocolatey Package Creation](https://docs.chocolatey.org/en-us/create/create-packages)
- [Scoop Bucket Creation](https://github.com/ScoopInstaller/Scoop/wiki/Buckets)

