# Commit Convention for Gopher

This document outlines the commit message convention used in the Gopher project to ensure consistent, readable, and automated changelog generation.

## Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Types

### Primary Types
- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries

### Special Types
- **ci**: Changes to CI configuration files and scripts
- **build**: Changes that affect the build system or external dependencies
- **revert**: Reverts a previous commit

## Scopes (Optional)

- **alias**: Alias management features
- **version**: Version management features
- **install**: Installation and uninstallation
- **config**: Configuration management
- **cli**: Command-line interface
- **api**: API changes
- **docs**: Documentation
- **tests**: Test-related changes
- **ci**: CI/CD related changes
- **docker**: Docker-related changes
- **security**: Security-related changes

## Description

- Use the imperative mood ("add feature" not "added feature")
- Don't capitalize the first letter
- No period at the end
- Keep it concise but descriptive

## Body (Optional)

- Explain the "what" and "why" of the change
- Wrap at 72 characters
- Use blank lines to separate paragraphs

## Footer (Optional)

- Reference issues: `Fixes #123`, `Closes #456`
- Breaking changes: `BREAKING CHANGE: description`
- Co-authors: `Co-authored-by: Name <email>`

## Examples

### Good Examples

```
feat(alias): add support for alias groups and tags
fix(version): resolve version detection on Windows
docs: update installation instructions for macOS
refactor(manager): split large functions into smaller modules
perf(install): optimize download speed with parallel downloads
test(alias): add comprehensive alias management tests
chore: update Go version in CI to 1.23
ci: add security scanning with Trivy
build: update dependencies to latest versions
```

### Breaking Changes

```
feat(alias): redesign alias storage format

BREAKING CHANGE: Alias storage format has changed from JSON to YAML.
Existing aliases will be automatically migrated on first run.
```

### With Scope and Body

```
feat(alias): add bulk alias operations

Add support for creating multiple aliases in a single command.
This improves user experience when setting up multiple aliases
for different environments.

Closes #123
```

## Benefits

1. **Automated Changelog**: Commit messages are automatically categorized
2. **Clear History**: Easy to understand what each commit does
3. **Release Notes**: Automatic generation of release notes
4. **Code Review**: Easier to review changes
5. **Debugging**: Easier to find specific changes

## Tools

- **Commitizen**: Interactive commit message creation
- **Commitlint**: Lint commit messages
- **Conventional Changelog**: Generate changelogs from commits
- **Semantic Release**: Automated versioning and releasing

## Installation

### Commitizen (Recommended)
```bash
npm install -g commitizen cz-conventional-changelog
echo '{ "path": "cz-conventional-changelog" }' > ~/.czrc
```

### Commitlint
```bash
npm install -g @commitlint/cli @commitlint/config-conventional
echo 'module.exports = {extends: ["@commitlint/config-conventional"]};' > commitlint.config.js
```

## Usage

### With Commitizen
```bash
git add .
git cz  # Interactive commit message creation
```

### Manual
```bash
git commit -m "feat(alias): add support for alias groups"
```

## Enforcement

- Pre-commit hooks can be set up to enforce this convention
- CI can validate commit messages
- PR templates can remind contributors of the convention

## Migration

If you have existing commits that don't follow this convention:

1. Use the enhanced release notes generator which handles both formats
2. Gradually adopt the convention for new commits
3. Consider squashing old commits when appropriate

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit)
- [Commitizen](https://github.com/commitizen/cz-cli)
- [Commitlint](https://github.com/conventional-changelog/commitlint)
