# Audit Logs API

`GET /api/v1/admin/audit-logs` requires `audit.view` and returns up to 200 operational events.

Access-control changes are exposed separately through `GET /api/v1/admin/access-audit`; see `ACCESS_AUDIT_API.md`.

With PostgreSQL configured, it reads records from the existing `operational_events` table created by migration `0002_operational_events.sql` and returns `source: "persistent"`. In database-free development, it returns the active engine stream with `source: "runtime"`.

Role-assignment history is stored separately by migration `0016_access_audit_logs.sql`.
