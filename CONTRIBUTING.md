# Contributing

Thanks for contributing to `teams-migrator`.

This project is a conservative migration CLI. The bar for new behavior is correctness, reviewability, and predictable operator experience before convenience.

## Before opening a pull request

1. Open or reference an issue for non-trivial changes.
2. Keep the change focused. One user-visible change per PR is preferred.
3. Update docs or examples when the CLI behavior changes.
4. Add or update tests when behavior changes.

## Development setup

Requirements:

- Go toolchain version from `go.mod`

Useful commands:

```bash
go test ./...
go vet -unsafeptr=false ./...
gofmt -w ./cmd ./internal
go build ./cmd/teams-migrator
```

## Style and scope

- Prefer narrowly scoped changes.
- Preserve dry-run-first behavior.
- Avoid broad refactors mixed with behavior changes.
- Keep operator-facing output explicit and reviewable.
- Do not add Jira Cloud support assumptions to Server/Data Center code paths.

## Pull requests

- Use Conventional Commits for branch, commit, and PR titles when practical.
- Include operator impact, risk, and rollback notes in the PR description.
- Mention any Jira permissions, ScriptRunner assumptions, or migration edge cases touched by the change.

## Documentation

Update the relevant files when behavior changes:

- `README.md` for public-facing usage
- `docs/` for operator workflows and support expectations
- `examples/` for sample inputs and outputs

## Reporting bugs

Use the bug report template and include:

- Jira edition and version
- source and target environment type
- exact command used
- dry-run or apply mode
- sanitized snippets from generated reports

For security-sensitive issues, follow [SECURITY.md](SECURITY.md) instead of opening a public issue.
