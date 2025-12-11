# Release Process

This project uses [Release Please](https://github.com/googleapis/release-please) for automated releases with semantic versioning.

## How It Works

Release Please automates the entire release process:

1. **Analyzes Commits**: Scans commit messages since the last release
2. **Determines Version**: Calculates the next version based on conventional commits
3. **Updates CHANGELOG**: Automatically maintains CHANGELOG.md
4. **Creates Release PR**: Opens a pull request with version bump and changelog updates
5. **Triggers Release**: When the release PR is merged, creates a GitHub release and tag
6. **Builds Artifacts**: GoReleaser builds and publishes release artifacts

## Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/) for automatic changelog generation.

**⚠️ ENFORCED:** Pull requests are automatically validated to ensure commit messages follow the conventional format. PRs that don't comply will fail status checks and cannot be merged.

### Version Bumps

- `fix:` - Patch version bump (0.1.0 → 0.1.1)
- `feat:` - Minor version bump (0.1.0 → 0.2.0)
- `BREAKING CHANGE:` or `!` - Major version bump (0.1.0 → 1.0.0)

### Changelog Categories

```bash
# Bug Fixes
git commit -m "fix: resolve API key generation timeout"
git commit -m "fix(auth): handle expired credentials properly"

# New Features
git commit -m "feat: add support for new Datadog scopes"
git commit -m "feat(roles): implement role inheritance"

# Breaking Changes
git commit -m "feat!: change configuration format"
git commit -m "fix!: remove deprecated endpoints"

# Documentation
git commit -m "docs: update README with new examples"
git commit -m "docs(api): improve API reference"

# Dependency Updates
git commit -m "build(deps): bump vault SDK to v0.20.0"
git commit -m "chore(deps): update Datadog client"

# Other Changes (won't appear in changelog)
git commit -m "chore: update .gitignore"
git commit -m "ci: fix test workflow"
git commit -m "style: format code"
git commit -m "refactor: simplify client initialization"
```

## Release Workflow

### Automated Release (Recommended)

1. **Merge PRs to main** using conventional commit messages
2. **Release Please opens a PR** automatically with:
   - Version bump in necessary files
   - Updated CHANGELOG.md
   - Release notes
3. **Review the release PR** to verify changes
4. **Merge the release PR** to trigger the release
5. **GitHub release is created** automatically with built artifacts

### Example Release PR

When you merge commits to main, Release Please will open a PR like:

```
Title: chore(main): release 0.2.0

Changes:
- Updates version in files
- Adds entries to CHANGELOG.md
- Includes all changes since last release
```

### Manual Tag Release (Legacy)

If you need to create a release manually:

```bash
# Create and push a tag
git tag v0.2.0
git push origin v0.2.0

# The legacy release.yml workflow will trigger
# Note: This bypasses automatic changelog generation
```

## Version Strategy

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (0.X.0): New features (backward compatible)
- **PATCH** (0.0.X): Bug fixes (backward compatible)

### Pre-1.0.0 Versions

While in 0.x.x versions:
- Minor version changes (0.X.0) may include breaking changes
- Patch versions (0.0.X) are for bug fixes only

## Commit Message Validation

### How It Works

PRs are automatically checked for conventional commit compliance:

1. **PR Title Validation**: Ensures PR title follows conventional format
2. **Commit Message Validation**: Validates all commits in the PR
3. **Status Check**: Required check must pass before merging

### Validation Rules

- **Type**: Must be one of: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`
- **Scope** (optional): e.g., `feat(api):` or `fix(auth):`
- **Subject**: Must start with uppercase letter, no period at end
- **Length**: Header max 100 characters

### Bypass Validation (Emergency Only)

If you need to bypass validation temporarily:

1. Add the label `ignore-semantic-pr` to the PR
2. Merge carefully - this will not generate proper changelog entries

### Fixing Validation Failures

**If PR title is invalid:**
```bash
# Edit the PR title to follow format:
# type(scope): Description starts with capital letter
# Example: feat: Add new authentication method
```

**If commit messages are invalid:**
```bash
# Option 1: Amend the last commit
git commit --amend -m "fix: correct commit message format"
git push --force

# Option 2: Interactive rebase to fix multiple commits
git rebase -i HEAD~3  # Edit last 3 commits
# Change "pick" to "reword" for commits to fix
# Save and edit each commit message

git push --force
```

**For squash merges:**
- Only the PR title matters (it becomes the commit message)
- Ensure PR title follows conventional format

## Troubleshooting

### Release PR not created

**Possible causes:**
- No commits since last release
- Commits don't follow conventional commit format
- Release Please workflow failed (check Actions tab)

**Solution:**
- Ensure commits use `feat:`, `fix:`, etc. prefixes
- Check GitHub Actions logs for errors

### Wrong version bump

**Cause:** Commit message didn't match expected format

**Solution:**
- Use correct conventional commit prefix
- For breaking changes, add `!` or `BREAKING CHANGE:` in commit body

### Need to modify release PR

**Options:**
1. Push more commits to main - Release Please will update the PR
2. Close the release PR and it will be recreated with new commits
3. Manually edit the PR's CHANGELOG.md if needed

## CI/CD Integration

### Workflows

- **`.github/workflows/release-please.yml`**: Automated releases
- **`.github/workflows/release.yml`**: Legacy manual tag releases (deprecated)
- **`.github/workflows/test.yml`**: Runs on all pushes

### Build Configuration

- **`.goreleaser.yaml`**: Defines build targets and artifacts
- Builds for: Linux, macOS, Windows (amd64, arm64, 386, arm)
- Generates SHA256 checksums

## Best Practices

1. **Always use conventional commits** for main branch
2. **Squash merge PRs** to keep a clean commit history
3. **Review release PRs** before merging to verify changelog
4. **Tag releases semantically** if manual tagging is required
5. **Test before releasing** - CI runs automatically on PRs

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Release Please Documentation](https://github.com/googleapis/release-please)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GoReleaser Documentation](https://goreleaser.com/)
