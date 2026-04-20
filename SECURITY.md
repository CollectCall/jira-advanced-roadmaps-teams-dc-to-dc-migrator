# Security Policy

## Supported versions

This project is still in an early public stage. Security fixes are handled on a best-effort basis on the latest released version and the current `master` branch.

## Reporting a vulnerability

Do not open a public GitHub issue for:

- exposed credentials
- unsafe write behavior
- bypasses of the dry-run or confirmation flow
- SSRF, command execution, or path traversal concerns
- release or installer supply-chain concerns

Instead, use GitHub private vulnerability reporting for this repository:

`https://github.com/CollectCall/jira-advanced-roadmaps-teams-dc-to-dc-migrator/security/advisories/new`

If that flow is unavailable, contact the repository owner privately through GitHub with:

- a short description
- impact
- reproduction steps
- affected version or commit
- whether credentials or customer data may be involved

If you do not have a safe way to share a proof of concept, say so in the initial private report.

## Response expectations

- Initial acknowledgement target: 5 business days
- Triage target: 10 business days
- Fix timing: depends on severity, exploitability, and maintainer availability

If the reporting path changes, this file will be updated in the default branch first.
