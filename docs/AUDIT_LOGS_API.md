# Audit Logs API

`GET /api/v1/admin/audit-logs` requires `audit.view` and returns up to 200 operational events.

With PostgreSQL configured, it reads records from the existing `operational_events` table created by migration `0002_operational_events.sql` and returns `source: "persistent"`. In database-free development, it returns the active engine stream with `source: "runtime"`.

No new database migration is required for this milestone.
