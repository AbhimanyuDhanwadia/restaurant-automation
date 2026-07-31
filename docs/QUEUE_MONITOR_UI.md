# Queue Monitor UI

The Queue Monitor refreshes the authenticated `GET /api/v1/automation/queue` endpoint every two seconds. It displays current queue depth and capacity, worker count, automation-event count, retries, and failed jobs from the active API process.

Queue capacity, worker count, and retry policy are startup configuration owned by the API. This screen intentionally observes these values without allowing browser-side changes that could disrupt in-flight order processing.
