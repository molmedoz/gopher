# Contributing to Gopher

Thank you for your interest in contributing to Gopher! This guide will help you set up your development environment and understand our workflow.

## Prerequisites

- Go 1.22 or later
- Git
- Make

## Development Environment Setup

### 1. Clone the Repository

```bash
git clone https://github.com/molmedoz/gopher.git
cd gopher
```

### 2. Install Development Tools

Gopher uses several development tools for code quality and testing. Install them all with:

```bash
make install-tools
```

This will install:
- **goimports** - Automatically format and organize imports
- **golangci-lint** - Comprehensive linting (multiple linters in one)
- **gofumpt** - Strict formatting beyond go fmt

### 3. Verify Installation

The Makefile automatically uses the full path to installed tools, so **no PATH configuration is needed!**

```bash
# Verify tools are installed
ls $(go env GOPATH)/bin/ | grep -E "(goimports|golangci-lint)"

# Run a quick check
make ci
```

**Note:** If you want to use the tools manually (outside of Make), you can optionally add them to your PATH:
```bash
# Optional: Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Development Workflow

### Daily Development Cycle

```bash
# 1. Pull latest changes
git pull origin main

# 2. Create a feature branch
git checkout -b feature/my-feature

# 3. Make your changes...

# 4. Format code and imports (auto-fix)
make format

# 5. Run CI checks locally (matches GitHub Actions)
make ci

# 6. If all checks pass, commit your changes
git add .
git commit -m "feat: add my feature"

# 7. Push and create PR
git push origin feature/my-feature
```

**Pro Tip:** Running `make ci` before pushing ensures your PR will pass all GitHub Actions checks!

### Quick Commands

```bash
# Format code
make fmt            # go fmt only
make imports        # goimports only
make format         # both fmt and imports (auto-fix)

# Check code (CI mode - no modifications)
make check-format   # check code format
make check-imports  # check imports format

# Testing
make test           # run tests with race detection and coverage
make test-verbose   # verbose output
make test-coverage  # detailed coverage report

# Quality checks
make vet            # go vet
make lint           # golangci-lint
make check          # fmt + vet + test (modifies files)
make ci             # check-format + check-imports + vet + lint + test (CI mode)

# Build
make build          # build for current platform
make build-all      # build for all platforms

# Complete cycle
make dev            # fmt + vet + test + build
```

## Code Quality Standards

### Formatting

- Use `goimports` for import organization
- Use `go fmt` for code formatting
- Run `make format` before committing

### Linting

- All code must pass `golangci-lint`
- Run `make lint` before committing
- Fix all warnings and errors

### Testing

- Write tests for new features
- Maintain or improve test coverage
- All tests must pass: `make test`
- Coverage reports: `make test-coverage`

## Dependency Policy

**IMPORTANT**: Gopher uses ONLY Go standard library and official Go extended packages.

### Allowed Dependencies

✅ **Go Standard Library** (all packages)
✅ **golang.org/x/*** - Official Go extended packages only

### NOT Allowed

❌ Third-party libraries (github.com/*, gopkg.in/*, etc.)
❌ External dependencies (except official Go extensions)

### Adding Dependencies

If you need to add a dependency:

1. **Check if it's in stdlib first**
2. **Check if golang.org/x/*** has a solution**
3. **Only then**, discuss with maintainers

Currently, Gopher uses:
- `golang.org/x/term` - Terminal handling (progress bars, TTY detection)
- `golang.org/x/sys` - Low-level system calls (indirect dependency)

## Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
# Open coverage.html in browser
```

### Run E2E Tests

```bash
# Run in Docker (safe, doesn't affect your system)
cd test
bash e2e.sh
```

### Run Specific Tests

```bash
go test -v -run TestClean ./internal/runtime/
go test -v -run TestPurge ./internal/runtime/
```

## Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Performance improvement
- `test`: Adding missing tests
- `chore`: Changes to build process or auxiliary tools

### Examples

```
feat: add clean command to remove download cache

Implement gopher clean command that removes downloaded Go archives
to free disk space. The command shows human-readable size of freed
space.

Closes #123
```

```
fix: correct version detection on Windows

Fix bug where Windows system Go was not detected correctly due to
incorrect path handling.
```

## Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch
3. **Make** your changes
4. **Test** thoroughly (`make check`)
5. **Format** code (`make format`)
6. **Lint** code (`make lint`)
7. **Commit** with conventional commit messages
8. **Push** to your fork
9. **Create** a Pull Request

### PR Checklist

- [ ] Code is formatted (`make format`)
- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] No new dependencies (or discussed with maintainers)
- [ ] Documentation updated (if needed)
- [ ] CHANGELOG.md updated (if needed)
- [ ] Commit messages follow convention
- [ ] PR description explains the change

## Documentation

Update documentation when:

- Adding new commands
- Changing behavior
- Adding configuration options
- Fixing bugs that affect user experience

Documentation files:
- `README.md` - Main documentation
- `QUICK_REFERENCE.md` - Quick reference guide
- `docs/USER_GUIDE.md` - Comprehensive user guide
- `docs/ROADMAP.md` - Future plans
- `CHANGELOG.md` - Version history

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/molmedoz/gopher/issues)
- **Discussions**: [GitHub Discussions](https://github.com/molmedoz/gopher/discussions)
- **Email**: molmedozazo@gmail.com

## Code of Conduct

Be respectful, inclusive, and constructive. We're all here to build something great together!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
