# Kitchen Performance UI

The Kitchen Performance workspace refreshes the authenticated analytics overview and durable order list every 15 seconds. It reports current order-state backlog, ready orders, non-cancelled completion, and printer telemetry. Completion is calculated by the backend as ready or delivered orders divided by all non-cancelled orders, with completed and eligible counts returned for coverage.

The status distribution and action queue come from `GET /api/v1/orders`. The automation metrics come from `GET /api/v1/analytics/overview`. Printer utilization reports runtime busy time rather than treating ticket count as a percentage. Cycle-time analytics are deliberately shown as unavailable until kitchen stage timestamps are persisted; the workspace does not infer duration from incomplete data.

No database migration was required because the existing durable order status is the source. The analytics API now returns completion counts and `kitchen_completion_source`; database-free deployments retain the runtime queue-progression fallback.
