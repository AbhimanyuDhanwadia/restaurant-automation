# Event Logs UI

The Event Logs workspace refreshes the authenticated `GET /api/v1/automation/events` endpoint every two seconds. It shows the active automation engine’s captured events and supports filtering by event type or searching by event type and order ID.

The endpoint is a runtime stream and is reset when the API process restarts. PostgreSQL deployments persist operational events through a separate event subscriber, but browsing durable audit history is intentionally deferred to the Administration Audit Logs milestone.
