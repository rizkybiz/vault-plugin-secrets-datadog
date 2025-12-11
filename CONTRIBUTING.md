# Contributing to vault-plugin-secrets-datadog

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/vault-plugin-secrets-datadog.git`
3. Create a feature branch: `git checkout -b feature/my-new-feature`
4. Make your changes
5. Test your changes: `go test ./...`
6. Commit using conventional commits (see below)
7. Push to your fork and submit a pull request

## Development Setup

### Prerequisites

- Go 1.25.0 or later
- Git
- Make (optional, for convenience)

### Building

```bash
# Build the plugin
go build -o vault/plugins/vault-plugin-secrets-datadog cmd/vault-plugin-secrets-datadog/main.go

# Or use make
make build
```

### Testing

```bash
# Run all tests
go test -v ./...

# Or use make
make test
```

## Commit Message Guidelines

**⚠️ IMPORTANT:** This project enforces [Conventional Commits](https://www.conventionalcommits.org/) format. Pull requests with non-compliant commit messages will fail CI checks.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: A new feature (triggers minor version bump)
- **fix**: A bug fix (triggers patch version bump)
- **docs**: Documentation changes only
- **style**: Code style changes (formatting, semicolons, etc.)
- **refactor**: Code changes that neither fix bugs nor add features
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **build**: Changes to build system or dependencies
- **ci**: Changes to CI configuration
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit

### Scopes (Optional)

- `deps`: Dependency updates
- `security`: Security-related changes
- `auth`: Authentication-related
- `api`: API key generation
- `roles`: Role management
- `config`: Configuration management
- `client`: Datadog client

### Examples

#### Good Commit Messages ✅

```bash
feat: add support for incident_write scope

fix(auth): handle expired credentials gracefully

docs: update README with installation instructions

fix!: remove deprecated configuration fields

feat(api): implement key rotation with zero downtime

BREAKING CHANGE: Configuration format has changed
```

#### Bad Commit Messages ❌

```bash
# Missing type
Added new feature

# Doesn't start with lowercase type
Feat: new feature

# Subject ends with period
fix: resolve bug.

# Not descriptive
fix: fix issue

# Too vague
update code
```

### Breaking Changes

For changes that break backwards compatibility:

```bash
# Option 1: Use ! after type
feat!: change API response format

# Option 2: Add footer
feat: redesign configuration structure

BREAKING CHANGE: Configuration file format has changed from YAML to JSON.
Migration guide: https://...
```

## Pull Request Process

1. **Ensure tests pass**: Run `go test ./...`
2. **Update documentation**: If you've changed APIs or behavior
3. **Use conventional commits**: PR title and commits must follow format
4. **Keep PRs focused**: One feature/fix per PR
5. **Write clear descriptions**: Explain what and why, not just how

### PR Title Format

Your PR title must follow conventional commit format (it becomes the commit message when squash merged):

```
feat: Add support for new Datadog scopes
fix(client): Handle API rate limiting properly
docs: Improve role configuration examples
```

### Automated Checks

Every PR will run:
- ✅ **Test Suite**: All unit tests must pass
- ✅ **Conventional Commits**: PR title and commits validated
- ✅ **Build**: Code must compile successfully

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep functions small and focused
- Add comments for complex logic
- Write tests for new features

## Security

If you discover a security vulnerability:

1. **DO NOT** open a public issue
2. Email the maintainers privately
3. Allow time for a fix before public disclosure

## License

By contributing, you agree that your contributions will be licensed under the same license as this project.

## Questions?

- Check existing issues and PRs
- Read the documentation in `RELEASE.md`
- Open a discussion for questions

## Recognition

Contributors will be recognized in release notes and the project README.

Thank you for contributing! 🎉
