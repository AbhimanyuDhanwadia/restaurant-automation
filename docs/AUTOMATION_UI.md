# Automation Dashboard API Connection

The Automation dashboard reads the authenticated API every five seconds:

- `GET /api/v1/integrations`
- `GET /api/v1/printers`
- `GET /api/v1/automation/events`
- `GET /api/v1/automation/queue`
- `GET /api/v1/insights`

It uses the shared frontend API client, including the current Supabase access token when authentication is enabled.

The dashboard reports only data exposed by those services. Captured order count, queue depth, retries, provider and printer status, events, and intelligence signals are live. Average print time, queue latency, database health, and internet health remain unavailable until dedicated telemetry or health endpoints are introduced.

An API failure is visible and can be retried; the dashboard does not replace unavailable runtime data with mock values.
