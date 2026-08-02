# System Health API

`GET /api/v1/system/health` is an authenticated Automation endpoint that returns an aggregate health snapshot.

It reports four components:

- `database`: PostgreSQL readiness when configured, otherwise `not_configured`.
- `queue`: the active automation queue’s bounded capacity and depth.
- `integrations`: registered providers and their connection state.
- `printers`: registered print drivers and their readiness state.

The overall result is `degraded` only when a configured component is unavailable. Optional components that are not configured do not degrade the service.
