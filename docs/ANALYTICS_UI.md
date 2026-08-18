# Analytics API Connection

The Analytics workspace reads `GET /api/v1/analytics/overview` through the authenticated frontend API client. It refreshes every 15 seconds while the page is open and includes the current Supabase access token when authentication is enabled.

## Data availability

The API intentionally marks metrics as unavailable until their source data exists. At this stage, order events and printer telemetry populate order count, kitchen completion, printer availability, printer utilization, and peak hours. Recorded durable order totals populate sales and average-ticket metrics when eligible totals share one currency. Delivery-provider orders populate active delivery, delivered-order, and intake-to-delivery metrics. Existing shift-task records populate staff task-completion percentage and completed/total counts; the metric remains unavailable until at least one task exists.

The interface does not substitute demo metrics if the service is unavailable. It displays the API error and provides a retry action instead.
