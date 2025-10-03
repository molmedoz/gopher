# Contributing to Gopher

Thank you for your interest in contributing to Gopher! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Process](#development-process)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

## Code of Conduct

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/) Code of Conduct. By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile commands)

### Setting Up Development Environment

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/your-username/gopher.git
   cd gopher
   ```

2. **Set up upstream remote**
   ```bash
   git remote add upstream https://github.com/molmedoz/gopher.git
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Install development tools**
   ```bash
   make install-tools
   ```

5. **Run tests**
   ```bash
   make test
   ```

## Development Process

### 1. Create a Feature Branch

```bash
# Create and switch to a new branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-description
```

### 2. Make Your Changes

- Write your code following the [coding standards](#coding-standards)
- Add tests for new functionality
- Update documentation as needed
- Ensure all tests pass

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run all checks
make check
```

### 4. Commit Your Changes

Use conventional commit messages:

```bash
git add .
git commit -m "feat: add new feature description"
```

**Commit Message Format:**
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Pull Request Process

### Before Submitting

- [ ] Code follows Go conventions and project style
- [ ] Tests are added/updated and passing
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] Commit messages follow conventional format
- [ ] Branch is up to date with main

### Pull Request Template

When creating a pull request, please include:

1. **Description**: What does this PR do?
2. **Type**: Feature, Bug Fix, Documentation, etc.
3. **Testing**: How was this tested?
4. **Breaking Changes**: Any breaking changes?
5. **Related Issues**: Link to related issues

### Review Process

1. **Automated Checks**: CI will run tests and linting
2. **Code Review**: Maintainers will review your code
3. **Feedback**: Address any feedback or requested changes
4. **Approval**: Once approved, your PR will be merged

## Issue Guidelines

### Before Creating an Issue

1. Check existing issues to avoid duplicates
2. Search closed issues for similar problems
3. Check the documentation for solutions

### Issue Types

- **Bug Report**: Something isn't working
- **Feature Request**: New functionality
- **Documentation**: Documentation improvements
- **Question**: General questions

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Run command '...'
2. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g. macOS, Linux, Windows]
- Go version: [e.g. 1.21.0]
- Gopher version: [e.g. 1.0.0]

**Additional context**
Any other context about the problem.
```

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or workarounds.

**Additional context**
Any other context about the feature request.
```

## Coding Standards

### Go Code Style

Follow standard Go conventions:

```go
// Package comment
package version

import (
    "fmt"
    "time"
)

// Interface comment
type ManagerInterface interface {
    // Method comment
    ListInstalled() ([]Version, error)
}

// Struct comment
type Version struct {
    Version     string    `json:"version"`      // Field comment
    OS          string    `json:"os"`           // Field comment
    InstalledAt time.Time `json:"installed_at"`
}

// Method comment
func (v Version) String() string {
    return fmt.Sprintf("%s (%s/%s)", v.Version, v.OS, v.Arch)
}
```

### Naming Conventions

- **Packages**: lowercase, single word
- **Interfaces**: descriptive name ending with `-er`
- **Structs**: PascalCase
- **Functions**: PascalCase for public, camelCase for private
- **Variables**: camelCase

### Error Handling

Always handle errors explicitly:

```go
// Good
func (m *Manager) Install(version string) error {
    if err := ValidateVersion(version); err != nil {
        return fmt.Errorf("invalid version: %w", err)
    }
    // ... rest of implementation
}

// Bad
func (m *Manager) Install(version string) error {
    ValidateVersion(version) // Error ignored
    // ... rest of implementation
}
```

### Documentation

- Document all public APIs
- Use clear, concise comments
- Include examples for complex functions
- Update documentation when changing APIs

## Testing Guidelines

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    invalidInput,
            expected: nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Coverage

- Aim for high test coverage (>80%)
- Test both success and error cases
- Test edge cases and boundary conditions
- Use table-driven tests for multiple scenarios

### Integration Tests

For integration tests that require external resources:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Integration test code
}
```

Run with: `go test -short ./...` (skips integration tests)

## Documentation

### User Documentation

- Update [USER_GUIDE.md](docs/USER_GUIDE.md) for user-facing changes
- Add examples to [EXAMPLES.md](docs/EXAMPLES.md)
- Update README.md for major features

### API Documentation

- Update [API_REFERENCE.md](docs/API_REFERENCE.md) for API changes
- Include method signatures and examples
- Document any breaking changes

### Developer Documentation

- Update [DEVELOPER_GUIDE.md](docs/DEVELOPER_GUIDE.md) for development process changes
- Document new build requirements or tools

## Release Process

### Version Bumping

- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes (backward compatible)

### Release Checklist

- [ ] Update version in `cmd/gopher/main.go`
- [ ] Update CHANGELOG.md
- [ ] Run full test suite
- [ ] Update documentation
- [ ] Create release tag
- [ ] Build and test release binaries

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Code Review**: Ask questions in pull request comments

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- GitHub contributors list

Thank you for contributing to Gopher! 🎉
