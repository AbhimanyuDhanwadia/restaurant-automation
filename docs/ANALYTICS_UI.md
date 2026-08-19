# Analytics API Connection

The Analytics workspace reads `GET /api/v1/analytics/overview` through the authenticated frontend API client. It refreshes every 15 seconds while the page is open and includes the current Supabase access token when authentication is enabled.

## Data availability

The API intentionally marks metrics as unavailable until their source data exists. Durable order records populate order count and peak hours when the order service is available, so those metrics survive API restarts. The response identifies this with `order_volume_source: "durable"`; database-free instances fall back to the current automation event stream and return `"runtime"`. Printer telemetry populates kitchen completion, printer availability, and printer utilization. Recorded durable order totals populate sales and average-ticket metrics when eligible totals share one currency. Delivery-provider orders populate active delivery, delivered-order, and intake-to-delivery metrics. Existing shift-task records populate staff task-completion percentage and completed/total counts; the metric remains unavailable until at least one task exists.

The interface does not substitute demo metrics if the service is unavailable. It displays the API error and provides a retry action instead.
