# Compatibility

## Supported scope

- Jira Server
- Jira Data Center
- Advanced Roadmaps team and team membership migration
- optional Parent Link and saved filter Team ID corrections

## Not supported

- Jira Cloud
- plan creation
- program creation
- plan-only team creation
- broad Jira content migration

## Source input modes

- direct Jira API access
- exported JSON files for teams, persons, and resources

## Target requirements

- reachable Jira Server/Data Center target instance
- credentials with permission to create teams and memberships where applicable
- existing destination context for non-shared teams

## Filter rewrite prerequisites

One of:

- ScriptRunner custom endpoints from `scripts/`
- a DB-derived CSV matching the expected filter export columns

## Platform support

Published binaries target:

- Linux `amd64`, `arm64`
- macOS `amd64`, `arm64`
- Windows `amd64`, `arm64`

## Support posture

The project is maintained on a best-effort basis. The latest release and the current default branch receive the strongest attention.
