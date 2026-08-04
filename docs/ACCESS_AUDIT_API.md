# Access Audit API

Migration `0016_access_audit_logs.sql` adds an append-only role-assignment trail. Each record stores the actor subject and email, target subject, previous and new role IDs, and timestamp. It does not store JWTs, passwords, or authentication headers.

`PATCH /api/v1/admin/users/{subject}/role` updates `app_users` and inserts its access-audit record in the same PostgreSQL transaction. If either statement fails, neither change is committed.

`GET /api/v1/admin/access-audit` returns the 200 newest role-assignment events with role names resolved from the active role catalog. It requires `audit.view`.
