# System Health UI

The System Health workspace refreshes the authenticated `GET /api/v1/system/health` endpoint every five seconds. It shows the current aggregate state and the detailed status of database, automation queue, integrations, and printers.

The screen does not test arbitrary internet hosts or expose infrastructure credentials. It only reports health signals owned by this application.
