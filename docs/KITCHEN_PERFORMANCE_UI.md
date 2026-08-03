# Kitchen Performance UI

The Kitchen Performance workspace refreshes the authenticated analytics overview and durable order list every 15 seconds. It reports current order-state backlog, ready orders, non-cancelled completion, and printer telemetry.

The status distribution and action queue come from `GET /api/v1/orders`. The automation metrics come from `GET /api/v1/analytics/overview`. Cycle-time analytics are deliberately shown as unavailable until kitchen stage timestamps are persisted; the workspace does not infer duration from incomplete data.

No database migration or API contract was required for this milestone.
