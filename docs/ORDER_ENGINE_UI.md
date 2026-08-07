# Order Engine UI

The Order Engine workspace refreshes authenticated `GET /api/v1/automation/events` and `GET /api/v1/automation/queue` endpoints every two seconds. It groups runtime events by order ID and displays the intake stages: received, normalized, validated, stored, and queued. Its event history also shows lifecycle transitions emitted by the Orders API: kitchen accepted, ready, delivered, and cancelled.

This workspace is intentionally observability-only. It does not replay, retry, or create orders from the browser, which prevents an operator from accidentally duplicating a live order. Status changes are made from Orders or Kitchen; repeated status submissions do not create duplicate lifecycle events. Printer processing and kitchen acceptance are separate services with their own screens.

Automation events are held in the active API process and reset when it restarts. Durable operational history remains available through the Audit Logs workspace.
