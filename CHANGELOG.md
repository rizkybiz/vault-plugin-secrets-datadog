# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Security

- Resolved 9 critical security vulnerabilities in dependencies
- Updated github.com/docker/docker from v27.2.1 to v28.3.3 (GO-2025-3829: Moby firewalld network isolation)
- Updated github.com/go-jose/go-jose/v4 from v4.0.2 to v4.0.5 (GO-2025-3485: DoS in parsing)
- Updated golang.org/x/crypto from v0.32.0 to v0.45.0 (GO-2025-4135, 4134, 4116, 3487)
- Updated golang.org/x/net from v0.34.0 to v0.47.0 (GO-2025-3595, 3503)
- Updated golang.org/x/oauth2 from v0.22.0 to v0.27.0 (GO-2025-3488)

### Changed

- Updated github.com/hashicorp/vault/sdk from 0.14.0 to 0.15.2
- Updated github.com/hashicorp/vault/api from 1.15.0 to 1.16.0
- Updated github.com/DataDog/datadog-api-client-go/v2 from 2.34.0 to 2.36.1
- Upgraded Go toolchain from 1.22 to 1.25.0

## [0.1.7] - 2024-11-28

### Changed

- Updated dependencies via dependabot

## [0.1.6] - Previous release

### Changed

- Dependency updates

## [0.1.5] - Previous release

### Changed

- Dependency updates

## [0.1.4] - Previous release

### Changed

- Dependency updates

## [0.1.3] - Previous release

### Changed

- Dependency updates

## [0.1.2] - Initial release

### Added

- Dynamic Datadog API key generation
- Dynamic Datadog Application key generation with scoped permissions
- Role-based key management with TTL support
- Configuration rotation capability
- Support for 26 Datadog permission scopes

[Unreleased]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/compare/v0.1.7...HEAD
[0.1.7]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/compare/v0.1.6...v0.1.7
[0.1.6]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/rizkybiz/vault-plugin-secrets-datadog/releases/tag/v0.1.2
