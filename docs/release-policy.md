# Release Policy

## Versioning

Releases are tagged with `v<major>.<minor>.<patch>`.

Until the project matures further, minor releases may still include CLI ergonomics changes or workflow refinements alongside new capabilities.

## Release notes

GitHub Releases are generated from merged commit and PR titles. Keep them readable:

- use Conventional Commits
- keep titles specific
- avoid vague summaries

## Release process

Preferred flow for user-visible work:

1. create a focused branch
2. make one logical change per commit
3. open a PR with a Conventional Commit-style title
4. merge to `master`
5. tag from `master`

## Artifact policy

Each tagged release publishes:

- platform-specific archives
- `checksums.txt`

## Support window

The latest released version is the default supported version for public users.
