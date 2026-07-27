# Analytics API Connection

The Analytics workspace reads `GET /api/v1/analytics/overview` through the authenticated frontend API client. It refreshes every 15 seconds while the page is open and includes the current Supabase access token when authentication is enabled.

## Data availability

The API intentionally marks metrics as unavailable until their source data exists. At this stage, order events and printer telemetry can populate order count, kitchen completion, printer availability, printer utilization, and peak hours. Sales, delivery, staff, and average-ticket metrics require later durable data models.

The interface does not substitute demo metrics if the service is unavailable. It displays the API error and provides a retry action instead.
