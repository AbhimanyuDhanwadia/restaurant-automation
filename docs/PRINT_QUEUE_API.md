# Print Queue API

Authenticated printer endpoints:

- `GET /api/v1/printers/queue` lists durable print jobs.
- `POST /api/v1/printers/tickets` creates a queued print job and submits it to the printer manager.
- `POST /api/v1/printers/tickets/{orderID}/reprint` recreates the latest durable ticket for that order as a reprint job.

`POST /api/v1/orders` is also an automatic ticket source: each accepted order creates one initial kitchen ticket from its line items. Manual ticket submission remains available for cashier and operational tickets.

Migration `0009_print_jobs.sql` stores ticket lines, status, attempts, errors, and timestamps. Jobs progress through `queued`, `printing`, `printed`, or `failed`; printer lifecycle callbacks update the record without blocking physical printing.

The same callbacks also publish `print.started`, `print.finished`, and `print.failed` automation events with the ticket destination and attempt number. PostgreSQL deployments persist them through the existing `operational_events` subscriber; a failed ticket additionally appears in the automation intelligence feed as a printer anomaly.

At server startup, queued jobs are atomically claimed before they are re-dispatched after printer workers and observers are ready. This prevents multiple API instances from replaying the same durable ticket during concurrent startup. Jobs already marked `printing` are not automatically replayed because their physical outcome is indeterminate after a restart; they are surfaced as printer-anomaly insights for operator review. Failed jobs also require an explicit reprint.

Tickets select their target through the required `destination` field. The current server configuration routes `kitchen` and `cashier` independently to their named drivers.
