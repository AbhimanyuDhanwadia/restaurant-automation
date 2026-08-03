# Order Engine UI

The Order Engine workspace refreshes authenticated `GET /api/v1/automation/events` and `GET /api/v1/automation/queue` endpoints every two seconds. It groups runtime events by order ID and displays the implemented stages: received, normalized, validated, stored, and queued.

This workspace is intentionally observability-only. It does not replay, retry, or create orders from the browser, which prevents an operator from accidentally duplicating a live order. Printer processing and kitchen acceptance are separate services with their own screens.

Automation events are held in the active API process and reset when it restarts. Durable operational history remains available through the Audit Logs workspace.
