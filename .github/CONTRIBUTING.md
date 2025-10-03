# Contributing to Gopher

Thank you for your interest in contributing to Gopher! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile commands)

### Development Setup

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/your-username/gopher.git
   cd gopher
   ```

2. **Set up the development environment**
   ```bash
   # Install dependencies
   go mod download
   
   # Run tests to ensure everything works
   go test ./...
   
   # Build the project
   go build ./cmd/gopher
   ```

3. **Create a development branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Using Makefile (Recommended)

```bash
# Check project status
make status

# Start development server
make dev

# Run tests
make test

# Format code
make format

# Build for testing
make build

# Run linting
make lint

# Clean build artifacts
make clean
```

### Manual Setup

```bash
# Run tests
go test ./...

# Format code
go fmt ./...
goimports -w .

# Run linting
go vet ./...
staticcheck ./...

# Build
go build ./cmd/gopher
```

## Making Changes

### Code Style

- Follow Go standard formatting (`gofmt`)
- Use `goimports` for import organization
- Follow Go naming conventions
- Write clear, self-documenting code
- Add comments for public APIs

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Maintenance tasks

Examples:
```
feat(list): add interactive pagination
fix(downloader): resolve 404 error for Go versions
docs: update installation instructions
test: add tests for interactive pagination
```

### File Organization

- Keep related functionality together
- Use meaningful package names
- Separate concerns into different packages
- Follow Go project layout conventions

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestFunctionName ./package
```

### Writing Tests

- Write tests for new functionality
- Aim for good test coverage
- Use table-driven tests when appropriate
- Test edge cases and error conditions
- Mock external dependencies

### Test Structure

```go
func TestFunction(t *testing.T) {
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
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("Function() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. **Create a pull request**
   - Use a descriptive title
   - Reference any related issues
   - Fill out the PR template

2. **Ensure quality**
   - All tests pass
   - Code is properly formatted
   - Documentation is updated
   - No linting errors

3. **Review process**
   - Address review feedback
   - Make requested changes
   - Respond to comments

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No new warnings
- [ ] All tests pass
- [ ] Changes are backwards compatible (if applicable)

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- `MAJOR`: Breaking changes
- `MINOR`: New features (backwards compatible)
- `PATCH`: Bug fixes (backwards compatible)

### Release Steps

1. **Update version**
   - Update version in `main.go`
   - Update changelog
   - Create release notes

2. **Create release**
   - Create and push tag: `git tag v1.0.0`
   - Push tag: `git push origin v1.0.0`
   - GitHub Actions will automatically create the release

3. **Verify release**
   - Check GitHub releases page
   - Verify binaries are available
   - Test installation from release

## Getting Help

- **Documentation**: Check the `docs/` directory
- **Issues**: Search existing issues or create a new one
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Contact [molmedoz@example.com](mailto:molmedoz@example.com)

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to Gopher! 🎉
