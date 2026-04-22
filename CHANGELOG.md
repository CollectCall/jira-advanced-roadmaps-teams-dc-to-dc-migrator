# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and the versioning model used in this repository's Git tags.

## [Unreleased]

## [1.3.4] - 2026-04-22

### Changed

- post-migrate target issue lookups now use batched issue-key searches to reduce target Jira REST requests
- post-migrate progress now reports specific input loading and target lookup steps
- target ScriptRunner filter lookup pushes name and owner narrowing into SQL for supported databases

### Fixed

- Jira REST requests now retry rate-limited and transient target API responses with bounded backoff

## [1.3.3] - 2026-04-22

### Added

- public project baseline docs and community files
- CI workflow for formatting, tests, `go vet`, and `govulncheck`
- examples for config, identity mapping, issues CSV, and filter CSV
- roadmap and sample output documentation

### Changed

- README rewritten for public-facing onboarding and operator clarity
- smoke test workflow updated to use `config show`
- documented the required ScriptRunner `local` database resource for filter endpoint setup

### Fixed

- ScriptRunner filter endpoints now handle Jira versions that expose `NotClause` children through different property names

## [1.3.2] - 2026-04-15

### Changed

- streamlined migrate workflow and docs

## [1.3.1] - 2026-04-15

### Added

- parent-link query restriction to epic issues during migration

## [1.3.0] - 2026-04-15

### Added

- filter clause scanning for Team ID rewrites

## [1.2.0] - 2026-04-15

### Added

- batch team migration scope support

## [1.1.0] - 2026-04-15

### Added

- issue export progress and team usage context
