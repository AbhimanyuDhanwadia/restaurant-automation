# Print Queue UI

The Print Queue workspace refreshes durable jobs every two seconds. It shows queue state, ticket detail, attempt count, errors, and reprint actions. Reprints create a new durable job from the latest stored ticket, so the workflow does not depend on the printer manager's in-memory history. On API startup, jobs left in `queued` state are claimed before dispatch so concurrent API instances do not replay the same ticket. Jobs already in `printing` stay available for operator review to avoid duplicate physical tickets and also appear as printer-anomaly insights.
