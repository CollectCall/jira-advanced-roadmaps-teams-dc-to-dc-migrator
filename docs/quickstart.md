# Quickstart

This guide is for operators moving Advanced Roadmaps team data between Jira Server or Data Center instances.

## Before you start

Make sure you have:

- admin or equivalent access for the Jira APIs you plan to use
- source and target Jira base URLs
- credentials for both environments
- an identity mapping CSV if users do not match by email
- ScriptRunner endpoints or a DB-derived filter CSV if filter rewrites are in scope

## 1. Install the CLI

Use a release artifact from the GitHub Releases page and verify `checksums.txt`.

## 2. Create a saved profile

```bash
teams-migrator init
```

The wizard stores stable inputs in `config.yaml` and prompts for secrets later at runtime instead of saving them.

## 3. Run a preparation pass

```bash
teams-migrator migrate --profile default --phase pre-migrate
```

Review the generated artifacts under `out/`. This is where you confirm what will be reused, created, skipped, or corrected later.

## 4. Apply the team and membership migration

```bash
teams-migrator migrate --profile default --phase migrate --apply
```

This phase creates teams and memberships only when the mapping is resolvable and safe.

## 5. Apply follow-up corrections

If Parent Link or Team-ID filter rewrites are in scope:

```bash
teams-migrator migrate --profile default --phase post-migrate --apply
```

## 6. Re-render a saved report

```bash
teams-migrator report --input out/migrate-report.json --format csv
```

## Common safety checks

- keep the first run in dry-run mode
- limit `--issue-project-scope` to the projects you actually want to correct
- verify non-shared teams already exist in the destination plan
- inspect skipped items before rerunning with `--apply`
