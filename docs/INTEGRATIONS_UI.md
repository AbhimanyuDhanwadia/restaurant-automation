# Integrations UI

The Integrations workspace reads the authenticated `GET /api/v1/integrations` registry snapshot every ten seconds. It reports only providers registered at API startup, their current connection health, and each provider's inbound webhook URL.

The screen intentionally has no connect/disconnect control. Provider collectors are created and started as part of API startup; changing an individual connection at runtime would leave the collector lifecycle ambiguous. Configure providers through API environment variables and restart the API.

For Swiggy, set `SWIGGY_WEBHOOK_SECRET` on the API server. The page then reports the `swiggy` provider and its intake URL. Secrets are never returned by the API or displayed in the browser.

The current generic receiver de-duplicates accepted order IDs for 15 minutes. With PostgreSQL configured, the receipt is durable and shared across API instances; otherwise it uses an in-memory local-development fallback. It acknowledges duplicate retries without queueing a second order. Provider-specific timestamp and nonce validation must be added once Swiggy confirms its production webhook contract.
